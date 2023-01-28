package rlpx

import (
	"bytes"
	cipher2 "crypto/cipher"
	"crypto/rand"
	"errors"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/cryptobase"
	"github.com/ethereum/go-ethereum/crypto/keyestablishmentalgorithm"
	"github.com/ethereum/go-ethereum/crypto/oqs"
	"github.com/ethereum/go-ethereum/crypto/signaturealgorithm"
	"github.com/ethereum/go-ethereum/rlp"
	"io"
	"sync"
)

type serverHelloMessage struct {
	CipherText            []byte //kemCipherTextLength
	ServerHelloRandomData [shaLen]byte
	Version               uint
	Rest                  []rlp.RawValue `rlp:"tail"`
}

type serverVerifyMessage struct {
	Signature    []byte //SignPublicKeyLen
	SignatureLen uint
	Rest         []rlp.RawValue `rlp:"tail"`
}

type Server struct {
	ephemeralKemPrivateKey  *keyestablishmentalgorithm.PrivateKey
	kem                     *oqs.KeyEncapsulation
	serverSigningPrivateKey *signaturealgorithm.PrivateKey
	clientSigningPublicKey  *signaturealgorithm.PublicKey

	rbuf ReadBuffer
	wbuf WriteBuffer

	clientHelloMessage  *clientHelloMessage
	serverHelloMessage  *serverHelloMessage
	serverVerifyMessage *serverVerifyMessage
	clientVerifyMessage *clientVerifyMessage

	kemCipherText   []byte //kemCipherTextLength
	kemSharedSecret []byte //kemSecretLength

	serverSeqNumHandshake uint
	clientSeqNumHandshake uint

	serverSeqNumApplication uint
	clientSeqNumApplication uint

	secret SessionSecret

	conn io.ReadWriter

	serializer Serializer

	transcript []byte

	client *Client

	handshakeDone bool
	mutex         sync.Mutex

	context string
}

func NewServer(conn io.ReadWriter, serverSigningPrivateKey *signaturealgorithm.PrivateKey, context string) *Server {
	server := Server{
		conn:                    conn,
		serverSigningPrivateKey: serverSigningPrivateKey,
		context:                 context,
	}

	server.serializer = NewRlpxSerializer()
	server.serverSeqNumHandshake = 1
	server.clientSeqNumHandshake = 1

	server.serverSeqNumApplication = 1
	server.clientSeqNumApplication = 1
	server.serializer.SetContext("server " + context)

	return &server
}

func (s *Server) SetClient(client *Client) {
	s.client = client
}

func (s *Server) SetServerSigningPrivateKey(serverSigningPrivateKey *signaturealgorithm.PrivateKey) {
	s.serverSigningPrivateKey = serverSigningPrivateKey
}

func (s *Server) PerformHandshake() error {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.handshakeDone == true {
		return errors.New("Handshake already done")
	}

	//Initialize KEM
	kem := oqs.KeyEncapsulation{}

	err := kem.Init(oqs.KemName, nil)
	if err != nil {
		return err
	}
	s.kem = &kem

	//Receive client hello message
	clientHelloMessage := new(clientHelloMessage)
	_, err = s.serializer.Deserialize(clientHelloMessage, s.conn)
	if err != nil {
		return err
	}

	//Handle client hello message
	s.clientHelloMessage = clientHelloMessage
	err = s.handleClientHello()
	if err != nil {
		return err
	}

	//Make server hello message
	err = s.makeServerHello()
	if err != nil {
		return err
	}

	serverHelloPacket, err := s.serializer.Serialize(s.serverHelloMessage)
	if err != nil {
		return err
	}

	//Write server hello message
	if _, err = s.conn.Write(serverHelloPacket); err != nil {
		return err
	}

	//Find the transcript of the session
	clientHelloTranscript, err := s.serializer.SerializeDeterministic(s.clientHelloMessage, 0)
	if err != nil {
		return err
	}

	serverHelloTranscript, err := s.serializer.SerializeDeterministic(s.serverHelloMessage, 0)
	if err != nil {
		return err
	}
	s.transcript = append(clientHelloTranscript, serverHelloTranscript...)
	transcriptHash := crypto.Keccak256(s.transcript)

	//Create the secrets
	secret, err := NewSessionSecret(transcriptHash, s.kemSharedSecret[:])
	s.secret = *secret

	//Sign the transcript hash
	signature, err := cryptobase.SigAlg.Sign(transcriptHash, s.serverSigningPrivateKey)
	if err != nil {
		return err
	}

	//Serialize the server verify message
	serverVerifyMessage := new(serverVerifyMessage)
	serverVerifyMessage.Signature = make([]byte, cryptobase.SigAlg.SignatureWithPublicKeyLength())
	copy(serverVerifyMessage.Signature[:], signature)
	serverVerifyMessage.SignatureLen = uint(len(signature))
	s.serverVerifyMessage = serverVerifyMessage

	serverVerifyPacket, err := s.serializer.Serialize(serverVerifyMessage)
	if err != nil {
		return err
	}

	err = s.WriteEncrypted(serverVerifyPacket, 0, PacketTypeHandshake)
	if err != nil {
		return err
	}

	//Create the transcript
	serverVerifyTranscript, err := s.serializer.SerializeDeterministic(s.serverVerifyMessage, 0)
	if err != nil {
		return err
	}

	s.transcript = append(s.transcript, serverVerifyTranscript...)

	err = s.handleClientVerify()
	if err != nil {
		return err
	}

	s.handshakeDone = true

	return nil
}

func (s *Server) Read() error {
	//Receive the server verify message header
	additionalData := new(AdditionalData)
	_, err := s.serializer.Deserialize(additionalData, s.conn)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) makeServerHello() error {
	serverHelloMessage := new(serverHelloMessage)
	serverHelloMessage.Version = 1

	// Generate ServerRandomData
	randomData := make([]byte, shaLength)
	_, err := rand.Read(randomData)
	if err != nil {
		return err
	}
	copy(serverHelloMessage.ServerHelloRandomData[:], randomData)

	serverHelloMessage.CipherText = make([]byte, s.kem.AlgDetails.LengthCiphertext)
	copy(serverHelloMessage.CipherText[:], s.kemCipherText[:])
	s.serverHelloMessage = serverHelloMessage

	return nil
}

func (s *Server) handleClientHello() error {

	ciphertext, sharedSecret, err := s.kem.EncapsulateSecret(s.clientHelloMessage.ClientKemPublicKey[:])
	if err != nil {
		return err
	}

	s.kemCipherText = make([]byte, s.kem.AlgDetails.LengthCiphertext)
	copy(s.kemCipherText[:], ciphertext[:])

	s.kemSharedSecret = make([]byte, s.kem.AlgDetails.LengthSharedSecret)
	copy(s.kemSharedSecret[:], sharedSecret[:])

	return nil
}

func (s *Server) handleClientVerify() error {

	//Receive the client verify message
	clientVerifyMessage := new(clientVerifyMessage)
	err := s.ReadAndDecryptMessage(clientVerifyMessage, PacketTypeHandshake)

	if err != nil {
		return err
	}

	s.clientVerifyMessage = clientVerifyMessage

	//Find the transcript of the session
	clientVerifyTranscript, err := s.serializer.SerializeDeterministic(s.clientVerifyMessage, 0)
	if err != nil {
		return err
	}

	transcriptHash := crypto.Keccak256(s.transcript)

	//Recover the public key from the signature
	clientPubKeyDataRemote, err := cryptobase.SigAlg.PublicKeyBytesFromSignature(transcriptHash, clientVerifyMessage.Signature[:clientVerifyMessage.SignatureLen])
	if err != nil {

		return err
	}

	if !cryptobase.SigAlg.Verify(clientPubKeyDataRemote, transcriptHash, clientVerifyMessage.Signature[:clientVerifyMessage.SignatureLen]) {
		return errors.New("client's signature verification failed")
	}

	s.clientSigningPublicKey, err = cryptobase.SigAlg.DeserializePublicKey(clientPubKeyDataRemote)
	if err != nil {
		return err
	}

	s.transcript = append(s.transcript, clientVerifyTranscript...)
	transcriptHash = crypto.Keccak256(s.transcript)

	err = s.secret.CreateApplicationSecrets(transcriptHash)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) createApplicationSecrets() error {
	return nil
}

func (s *Server) WriteEncrypted(data []byte, context uint64, packetType PacketType) error {
	additionalData := make([]byte, shaLength)
	_, err := rand.Read(additionalData)
	if err != nil {
		return err
	}

	var cipher cipher2.AEAD
	var seqNum uint
	var serverIv []byte
	if packetType == PacketTypeHandshake {
		cipher = s.secret.ServerHandshakeCipher
		seqNum = s.serverSeqNumHandshake
		serverIv = s.secret.ServerHandshakeIv
	} else {
		cipher = s.secret.ServerApplicationCipher
		seqNum = s.serverSeqNumApplication
		serverIv = s.secret.ServerApplicationIv
	}

	encryptedData, err := Encrypt(cipher, data, additionalData, packetType, serverIv, seqNum)
	if err != nil {
		return err
	}

	header := new(Header)
	header.PacketType = uint(packetType)
	header.MajorVersion = majorVersion
	header.MinorVersion = minorVersion
	header.RecordLength = uint(len(encryptedData))
	header.Context = context
	copy(header.AdditionalData[:], additionalData)

	headerPacket, err := s.serializer.Serialize(header)
	if err != nil {
		return err
	}

	_, err = s.conn.Write(headerPacket)

	if err != nil {
		return err
	}

	_, err = s.conn.Write(encryptedData)

	if err != nil {
		return err
	}

	if packetType == PacketTypeHandshake {
		s.serverSeqNumHandshake = s.serverSeqNumHandshake + 1
	} else {
		s.serverSeqNumApplication = s.serverSeqNumApplication + 1
	}

	return nil
}

func (s *Server) ReadAndDecrypt(packetType PacketType) (*DataPacket, error) {

	if packetType == PacketTypeApplicationData {
	}

	header := new(Header)
	_, err := s.serializer.Deserialize(header, s.conn)
	if err != nil {
		return nil, err
	}

	// Read the encrypted data
	recLen := int(header.RecordLength)

	if packetType == PacketTypeApplicationData {
	}

	encryptedData := make([]byte, int(recLen))
	bytesRead, err := io.ReadAtLeast(s.conn, encryptedData, int(recLen))
	if err != nil {
		return nil, err
	}

	if bytesRead != int(recLen) {
		return nil, errors.New("prefix size less")
	}

	var cipher cipher2.AEAD
	var seqNum uint
	var clientIv []byte
	if packetType == PacketTypeHandshake {
		cipher = s.secret.ClientHandshakeCipher
		seqNum = s.clientSeqNumHandshake
		clientIv = s.secret.ClientHandshakeIv
	} else {
		cipher = s.secret.ClientApplicationCipher
		seqNum = s.clientSeqNumApplication
		clientIv = s.secret.ClientApplicationIv
	}

	dataPacket, err := Decrypt(cipher, encryptedData, header.AdditionalData[:], packetType, clientIv, seqNum)
	if err != nil {
		return nil, err
	}

	if dataPacket.packetType != packetType {
		return nil, errors.New("packetType mismatch")
	}
	dataPacket.context = header.Context

	if packetType == PacketTypeHandshake {
		s.clientSeqNumHandshake = s.clientSeqNumHandshake + 1
	} else {
		s.clientSeqNumApplication = s.clientSeqNumApplication + 1
	}

	return dataPacket, nil
}

func (s *Server) ReadAndDecryptMessage(msg interface{}, packetType PacketType) error {
	dataPacket, err := s.ReadAndDecrypt(packetType)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(dataPacket.fragment)
	_, err = s.serializer.Deserialize(msg, reader)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Cleanup() {
	if s.kem != nil {
		s.kem.Clean()
	}
}

func (s *Server) InitWithSecrets(secret SessionSecret) {
	s.secret = secret
}
