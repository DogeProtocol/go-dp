package rlpx

import (
	"bytes"
	cipher2 "crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/DogeProtocol/dp/crypto"
	"github.com/DogeProtocol/dp/crypto/cryptobase"
	"github.com/DogeProtocol/dp/crypto/keyestablishmentalgorithm"
	"github.com/DogeProtocol/dp/crypto/oqs"
	"github.com/DogeProtocol/dp/crypto/signaturealgorithm"
	"github.com/DogeProtocol/dp/log"
	"github.com/DogeProtocol/dp/rlp"
	"io"
	"sync"
)

type clientHelloMessage struct {
	ClientKemPublicKey    []byte //kemPublicKeyLen
	ClientHelloRandomData [shaLen]byte
	Version               uint
	Rest                  []rlp.RawValue `rlp:"tail"`
}

type clientVerifyMessage struct {
	Signature    []byte //SignPublicKeyLen
	SignatureLen uint
	Rest         []rlp.RawValue `rlp:"tail"`
}

type Client struct {
	ephemeralKemPrivateKey  *keyestablishmentalgorithm.PrivateKey
	kem                     *oqs.KeyEncapsulation
	kemCipherText           []byte //kemCipherTextLength
	kemSharedSecret         []byte //kemSecretLength
	Nonce                   uint
	clientSigningPrivateKey *signaturealgorithm.PrivateKey
	serverSigningPublicKey  *signaturealgorithm.PublicKey

	clientHelloMessage  *clientHelloMessage
	serverHelloMessage  *serverHelloMessage
	serverVerifyMessage *serverVerifyMessage
	clientVerifyMessage *clientVerifyMessage

	rbuf        ReadBuffer
	wbuf        WriteBuffer
	RecordCount uint

	secret SessionSecret

	conn io.ReadWriter

	serializer Serializer

	serverSeqNumHandshake uint
	clientSeqNumHandshake uint

	serverSeqNumApplication uint
	clientSeqNumApplication uint

	transcript []byte

	server *Server

	handshakeDone bool
	mutex         sync.Mutex

	context string
}

func (c *Client) SetServer(server *Server) {
	c.server = server
}

func NewClient(conn io.ReadWriter, clientSigningPrivateKey *signaturealgorithm.PrivateKey, serverSigningPublicKey *signaturealgorithm.PublicKey, context string) *Client {
	client := Client{
		conn:                    conn,
		clientSigningPrivateKey: clientSigningPrivateKey,
		serverSigningPublicKey:  serverSigningPublicKey,
	}

	client.serializer = NewRlpxSerializer()
	client.serverSeqNumHandshake = 1
	client.clientSeqNumHandshake = 1

	client.serverSeqNumApplication = 1
	client.clientSeqNumApplication = 1
	client.serializer.SetContext("client " + context)
	client.context = context

	return &client
}

func (c *Client) SetClientSigningPrivateKey(clientSigningPrivateKey *signaturealgorithm.PrivateKey) {
	c.clientSigningPrivateKey = clientSigningPrivateKey
}

func (c *Client) SetServerSigningPublicKey(serverSigningPublicKey *signaturealgorithm.PublicKey) {
	c.serverSigningPublicKey = serverSigningPublicKey
}

func (c *Client) PerformHandshake() error {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.handshakeDone == true {
		return errors.New("Handshake already done")
	}

	//Initialize KEM
	kem := oqs.KeyEncapsulation{}

	err := kem.Init(oqs.KemName, nil)
	if err != nil {
		return err
	}
	c.kem = &kem

	//Make client hello message
	err = c.makeClientHello()
	if err != nil {
		return err
	}

	clientHelloPacket, err := c.serializer.Serialize(c.clientHelloMessage)
	if err != nil {
		return err
	}

	//Write client hello message
	if _, err = c.conn.Write(clientHelloPacket); err != nil {
		return err
	}

	//Receive server hello message
	serverHelloMessage := new(serverHelloMessage)
	_, err = c.serializer.Deserialize(serverHelloMessage, c.conn)
	if err != nil {
		return err
	}

	//Handle server hello message
	c.serverHelloMessage = serverHelloMessage
	err = c.handleServerHello()
	if err != nil {
		return err
	}

	//Find the transcript of the session
	clientHelloTranscript, err := c.serializer.SerializeDeterministic(c.clientHelloMessage, 0)
	if err != nil {
		return err
	}

	serverHelloTranscript, err := c.serializer.SerializeDeterministic(c.serverHelloMessage, 0)
	if err != nil {
		return err
	}
	transcript := append(clientHelloTranscript, serverHelloTranscript...)
	transcriptHash := crypto.Keccak256(transcript)
	c.transcript = transcript

	//Create the secrets
	secret, err := NewSessionSecret(transcriptHash, c.kemSharedSecret[:])
	c.secret = *secret

	//Receive the server verify message
	serverVerifyMessage := new(serverVerifyMessage)
	err = c.ReadAndDecryptMessage(serverVerifyMessage, PacketTypeHandshake)

	if err != nil {
		return err
	}

	//Verify the signature to make sure the server is what it is claiming to be
	serverPubKeyDataLocal, err := cryptobase.SigAlg.SerializePublicKey(c.serverSigningPublicKey)
	if err != nil {
		return err
	}

	//Recover the public key from the signature
	serverPubKeyDataRemote, err := cryptobase.SigAlg.PublicKeyBytesFromSignature(transcriptHash, serverVerifyMessage.Signature[:serverVerifyMessage.SignatureLen])
	if err != nil {
		return err
	}

	//Validate that expected public key and remote public key are the same (additional sanity check)
	if !bytes.Equal(serverPubKeyDataLocal, serverPubKeyDataRemote) {
		log.Error("Public Key mismatch",
			"serverSigningPublicKey", base64.StdEncoding.EncodeToString(c.serverSigningPublicKey.PubData),
			"signature", base64.StdEncoding.EncodeToString(serverVerifyMessage.Signature[:serverVerifyMessage.SignatureLen]),
			"serverPubKeyDataRemote", base64.StdEncoding.EncodeToString(serverPubKeyDataRemote))
		return errors.New("Public key mismatch")
	}

	if !cryptobase.SigAlg.Verify(serverPubKeyDataLocal, transcriptHash, serverVerifyMessage.Signature[:serverVerifyMessage.SignatureLen]) {
		return errors.New("server's signature verification failed")
	}

	//Create the transcript
	serverVerifyTranscript, err := c.serializer.SerializeDeterministic(serverVerifyMessage, 0)
	if err != nil {
		return err
	}

	transcript = append(transcript, serverVerifyTranscript...)
	transcriptHash = crypto.Keccak256(transcript)
	c.transcript = transcript
	c.serverVerifyMessage = serverVerifyMessage

	//Sign the transcript hash
	signature, err := cryptobase.SigAlg.Sign(transcriptHash, c.clientSigningPrivateKey)
	if err != nil {
		return err
	}

	//Serialize the server verify message
	clientVerifyMessage := new(clientVerifyMessage)
	clientVerifyMessage.Signature = make([]byte, cryptobase.SigAlg.SignatureWithPublicKeyLength())
	copy(clientVerifyMessage.Signature[:], signature)
	clientVerifyMessage.SignatureLen = uint(len(signature))
	c.clientVerifyMessage = clientVerifyMessage

	clientVerifyPacket, err := c.serializer.Serialize(clientVerifyMessage)
	if err != nil {
		return err
	}

	clientVerifyTranscript, err := c.serializer.SerializeDeterministic(clientVerifyMessage, 0)
	if err != nil {
		return err
	}

	err = c.WriteEncrypted(clientVerifyPacket, 0, PacketTypeHandshake)
	if err != nil {
		return err
	}

	transcript = append(transcript, clientVerifyTranscript...)
	c.transcript = transcript

	tHash := crypto.Keccak256(c.transcript)
	err = c.secret.CreateApplicationSecrets(tHash)
	if err != nil {
		return err
	}

	c.handshakeDone = true

	return nil
}

func (c *Client) makeClientHello() error {
	clientHelloMessage := new(clientHelloMessage)
	clientHelloMessage.Version = 1

	//Generate an ephemeral kem keypair
	kemPrivateKey, err := c.kem.GenerateKemKeyPair()
	if err != nil {
		return err
	}
	c.ephemeralKemPrivateKey = kemPrivateKey
	clientHelloMessage.ClientKemPublicKey = make([]byte, c.kem.AlgDetails.LengthPublicKey)
	copy(clientHelloMessage.ClientKemPublicKey[:], c.ephemeralKemPrivateKey.N.Bytes())

	// Generate ClientRandomData
	randomData := make([]byte, shaLength)
	_, err = rand.Read(randomData)
	if err != nil {
		return err
	}
	copy(clientHelloMessage.ClientHelloRandomData[:], randomData)
	c.Nonce = 1
	c.clientHelloMessage = clientHelloMessage

	return nil
}

func (c *Client) Cleanup() {
	if c.kem != nil {
		c.kem.Clean()
	}
}

func (c *Client) handleServerHello() error {

	sharedSecret, err := c.kem.DecapsulateSecret(c.serverHelloMessage.CipherText[:])
	if err != nil {
		return err
	}

	c.kemSharedSecret = make([]byte, c.kem.AlgDetails.LengthSharedSecret)
	copy(c.kemSharedSecret[:], sharedSecret[:])

	return nil
}

func (c *Client) ReadAndDecryptMessage(msg interface{}, packetType PacketType) error {
	dataPacket, err := c.ReadAndDecrypt(packetType)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(dataPacket.fragment)
	_, err = c.serializer.Deserialize(msg, reader)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) WriteEncrypted(data []byte, context uint64, packetType PacketType) error {
	if packetType == PacketTypeApplicationData {
		if c.handshakeDone != true {
		}
	}

	additionalData := make([]byte, shaLength)
	_, err := rand.Read(additionalData)
	if err != nil {
		return err
	}

	var cipher cipher2.AEAD
	var seqNum uint
	var clientIv []byte
	if packetType == PacketTypeHandshake {
		cipher = c.secret.ClientHandshakeCipher
		seqNum = c.clientSeqNumHandshake
		clientIv = c.secret.ClientHandshakeIv
	} else {
		cipher = c.secret.ClientApplicationCipher
		seqNum = c.clientSeqNumApplication
		clientIv = c.secret.ClientApplicationIv
	}

	encryptedData, err := Encrypt(cipher, data, additionalData, packetType, clientIv, seqNum)
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

	headerPacket, err := c.serializer.Serialize(header)
	if err != nil {
		return err
	}

	_, err = c.conn.Write(headerPacket)
	if err != nil {
		return err
	}

	_, err = c.conn.Write(encryptedData)

	if err != nil {
		return err
	}

	if packetType == PacketTypeHandshake {
		c.clientSeqNumHandshake = c.clientSeqNumHandshake + 1 //important
	} else {
		c.clientSeqNumApplication = c.clientSeqNumApplication + 1 //important
	}

	return nil
}

func (c *Client) ReadAndDecrypt(packetType PacketType) (*DataPacket, error) {
	if packetType == PacketTypeApplicationData {

	}

	header := new(Header)
	_, err := c.serializer.Deserialize(header, c.conn)
	if err != nil {
		return nil, err
	}

	// Read the encrypted data
	recLen := int(header.RecordLength)

	if packetType == PacketTypeApplicationData {

	}

	encryptedData := make([]byte, int(recLen))
	bytesRead, err := io.ReadAtLeast(c.conn, encryptedData, int(recLen))
	if err != nil {

		return nil, err
	}

	if bytesRead != int(recLen) {

		return nil, errors.New("prefix size less")
	}

	var cipher cipher2.AEAD
	var seqNum uint
	var serverIv []byte
	if packetType == PacketTypeHandshake {
		cipher = c.secret.ServerHandshakeCipher
		seqNum = c.serverSeqNumHandshake
		serverIv = c.secret.ServerHandshakeIv
	} else {
		cipher = c.secret.ServerApplicationCipher
		seqNum = c.serverSeqNumApplication
		serverIv = c.secret.ServerApplicationIv
	}

	dataPacket, err := Decrypt(cipher, encryptedData, header.AdditionalData[:], packetType, serverIv, seqNum)
	if err != nil {
		return nil, err
	}

	if dataPacket.packetType != packetType {
		return nil, errors.New("packetType mismatch")
	}
	dataPacket.context = header.Context

	if packetType == PacketTypeHandshake {
		c.serverSeqNumHandshake = c.serverSeqNumHandshake + 1
	} else {
		c.serverSeqNumApplication = c.serverSeqNumApplication + 1
	}

	return dataPacket, nil
}

func (c *Client) InitWithSecrets(secret SessionSecret) {
	c.secret = secret
}
