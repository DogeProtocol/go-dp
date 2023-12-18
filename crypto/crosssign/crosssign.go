package crosssign

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DogeProtocol/dp/accounts"
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/common/hexutil"
	"github.com/DogeProtocol/dp/crypto"
	"github.com/DogeProtocol/dp/crypto/cryptobase"
	"github.com/DogeProtocol/dp/crypto/secp256k1"
	"github.com/DogeProtocol/dp/crypto/signaturealgorithm"
	"github.com/status-im/keycard-go/hexutils"
	"strings"
)

const (
	ERC20AddressLength = 20
	MessageTemplate    = "I AGREE TO BECOME A GENESIS VALIDATOR FOR MAINNET. MY ETH ADDRESS IS [ETH_ADDRESS]. MY CORRESPONDING DEPOSITOR QUANTUM ADDRESS IS [DEPOSITOR_ADDRESS] AND VALIDATOR QUANTUM ADDRESS IS [VALIDATOR_ADDRESS]. VALIDATOR AMOUNT IS [AMOUNT] DOGEP."
)

type SignDetails struct {
	Address string `json:"address"`
	Msg     string `json:"msg"`
	Sig     string `json:"sig"`
	Version string `json:"version"`
}

type GenesisCrossSignDetails struct {
	EthAddress        string `json:"ethAddress"`
	DepositorAddress  string `json:"depositorAddress"`
	ValidatorAddress  string `json:"validatorAddress"`
	Amount            string `json:"amount"`
	Message           string `json:"message"`
	QuantumSignature  string `json:"quantumSignature"`
	EthereumSignature string `json:"ethereumSignature"`
}

//signJsonData := "{\r\n  \"address\": \"0xF422Ec881E87B934A165DB64132a87fbd1753daD\",\r\n  \"msg\": \"Test message waller\",\r\n  \"sig\": \"0x5c73e35d19d6656f826c82513a4523a8c789762bacfd1ce5127f24c1e61cd59f7779132c3a390294db158735e398c4e87b726b87bef44ad840a47ac6ca06ef8d1b\",\r\n  \"version\": \"2\"\r\n}"

func SignGenesis(depKey *signaturealgorithm.PrivateKey, valKey *signaturealgorithm.PrivateKey,
	ethAddr string, amount string) (*GenesisCrossSignDetails, error) {
	depositorAddr := cryptobase.SigAlg.PublicKeyToAddressNoError(&depKey.PublicKey).Hex()
	validatorAddr := cryptobase.SigAlg.PublicKeyToAddressNoError(&valKey.PublicKey).Hex()

	message := strings.Replace(MessageTemplate, "[ETH_ADDRESS]", ethAddr, 1)
	message = strings.Replace(message, "[DEPOSITOR_ADDRESS]", depositorAddr, 1)
	message = strings.Replace(message, "[VALIDATOR_ADDRESS]", validatorAddr, 1)
	message = strings.Replace(message, "[AMOUNT]", amount, 1)

	messageDigest, _ := accounts.TextAndHash([]byte(message))

	depositorSignature, err := cryptobase.SigAlg.Sign(messageDigest, depKey)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Error signing using depositor key")
	}

	validatorSignature, err := cryptobase.SigAlg.Sign(messageDigest, valKey)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Error signing using validator key")
	}

	combined := common.CombineTwoParts(depositorSignature, validatorSignature)
	hexSigCombined := hexutils.BytesToHex(combined)

	details := &GenesisCrossSignDetails{
		EthAddress:        ethAddr,
		DepositorAddress:  depositorAddr,
		ValidatorAddress:  validatorAddr,
		Amount:            amount,
		Message:           message,
		QuantumSignature:  hexSigCombined,
		EthereumSignature: "", //unavailable, to be done via https://app.mycrypto.com/sign-message
	}

	return details, nil
}

func VerifyGenesis(details *GenesisCrossSignDetails) ([]byte, error) {
	if len(details.EthAddress) == 0 || len(details.DepositorAddress) == 0 || len(details.ValidatorAddress) == 0 || len(details.Amount) == 0 || len(details.Message) == 0 || len(details.QuantumSignature) == 0 || len(details.EthereumSignature) == 0 {
		return nil, errors.New("malformed json")
	}

	if common.IsLegacyEthereumHexAddress(details.EthAddress) == false {
		return nil, errors.New("invalid EthAddress")
	}

	if common.IsHexAddress(details.DepositorAddress) == false {
		return nil, errors.New("invalid DepositorAddress")
	}

	if common.IsHexAddress(details.ValidatorAddress) == false {
		return nil, errors.New("invalid ValidatorAddress")
	}

	ethSig, err := hexutil.MustDecodeWithError(details.EthereumSignature)
	if err != nil {
		return nil, err
	}

	//todo: verify other fields to avoid panic and deeper input validations

	message := strings.Replace(MessageTemplate, "[ETH_ADDRESS]", details.EthAddress, 1)
	message = strings.Replace(message, "[DEPOSITOR_ADDRESS]", details.DepositorAddress, 1)
	message = strings.Replace(message, "[VALIDATOR_ADDRESS]", details.ValidatorAddress, 1)
	message = strings.Replace(message, "[AMOUNT]", details.Amount, 1)

	if details.Message != message {
		return nil, errors.New("invalid message")
	}

	messageDigest, _ := accounts.TextAndHash([]byte(message))
	sigBytes := hexutils.HexToBytes(details.QuantumSignature)

	depSig, valSig, err := common.ExtractTwoParts(sigBytes)
	if err != nil {
		return nil, err
	}

	depPubKey, err := cryptobase.SigAlg.PublicKeyFromSignature(messageDigest, depSig)
	if err != nil {
		return nil, err
	}

	if cryptobase.SigAlg.Verify(depPubKey.PubData, messageDigest, depSig) == false {
		return nil, errors.New("depositor signature verify failed")
	}

	valPubKey, err := cryptobase.SigAlg.PublicKeyFromSignature(messageDigest, valSig)
	if err != nil {
		return nil, err
	}

	if cryptobase.SigAlg.Verify(valPubKey.PubData, messageDigest, valSig) == false {
		return nil, errors.New("validator signature verify failed")
	}

	depositorAddr2 := cryptobase.SigAlg.PublicKeyToAddressNoError(depPubKey).Hex()
	if strings.Compare(details.DepositorAddress, depositorAddr2) != 0 {
		return nil, errors.New("depositor address verify failed")
	}

	validatorAddr2 := cryptobase.SigAlg.PublicKeyToAddressNoError(valPubKey).Hex()
	if strings.Compare(details.ValidatorAddress, validatorAddr2) != 0 {
		return nil, errors.New("validator address verify failed")
	}

	err = VerifyEthereumAddressAndMessage(details.EthAddress, messageDigest, ethSig)
	if err != nil {
		fmt.Println("VerifyEthereumAddressAndMessage failed", err)
		return nil, err
	}

	return messageDigest, nil
}

func CrossSignVerification(signJsonData string) error {

	var signDetails SignDetails

	err := json.Unmarshal([]byte(signJsonData), &signDetails)
	if err != nil {
		return fmt.Errorf("error 1 : " + err.Error())
	}

	if len(signDetails.Msg) == 0 || len(signDetails.Sig) == 0 || len(signDetails.Address) == 0 {
		return fmt.Errorf("error 1-1 : Some data is empty")
	}

	msgData := []byte(signDetails.Msg)
	msgHash, _ := accounts.TextAndHash(msgData)
	sig := hexutil.MustDecode(signDetails.Sig)
	addressBytes := hexToAddress(signDetails.Address)

	if len(sig) != 65 {
		return fmt.Errorf("error 2 : mismatch sign length")
	}
	if sig[64] != 27 && sig[64] != 28 {
		return fmt.Errorf("error 3 : Sign last byte mismatch")
	}

	sig[64] -= 27 // Transform yellow paper V from 27/28 to 0/1
	sign := make([]byte, 65)
	copy(sign, sig)

	recovered, err := sigToPub(msgHash, sign)
	if err != nil {
		return fmt.Errorf("error 4 : " + err.Error())
	}

	recoveredAddressBytes := pubkeyToAddress(*recovered)

	if len(addressBytes) != ERC20AddressLength {
		return fmt.Errorf("error 5 : mismatch length addressBytes")
	}

	if len(recoveredAddressBytes) != ERC20AddressLength {
		return fmt.Errorf("error 6 : mismatch length recoveredAddressBytes")
	}

	if bytes.Compare(recoveredAddressBytes, addressBytes) != 0 {
		return fmt.Errorf("error 7 : mismatch address bytes (recoveredAddressBytes, addressBytes) ")
	}

	//fmt.Println("recoveredAddress ", recoveredAddressBytes)
	//fmt.Println("address: ", addressBytes)
	//fmt.Println("Success...")

	return nil
}

func VerifyEthereumAddressAndMessage(ethAddress string, messageDigest []byte, signature []byte) error {
	if len(signature) != 65 {
		return fmt.Errorf("error 2 : mismatch sign length")
	}
	if signature[64] != 27 && signature[64] != 28 {
		return fmt.Errorf("error 3 : Sign last byte mismatch")
	}
	signature[64] -= 27 // Transform yellow paper V from 27/28 to 0/1
	sign := make([]byte, 65)
	copy(sign, signature)

	recovered, err := sigToPub(messageDigest, sign)
	if err != nil {
		return fmt.Errorf("error : " + err.Error())
	}

	recoveredAddressBytes := pubkeyToAddress(*recovered)
	addressBytes := hexToAddress(ethAddress)

	if bytes.Compare(recoveredAddressBytes, addressBytes) != 0 {
		return fmt.Errorf("error : mismatch address bytes (recoveredAddressBytes, addressBytes) ")
	}

	return nil
}

func sigToPub(hash, sig []byte) (*ecdsa.PublicKey, error) {
	s, err := ecrecover(hash, sig)
	if err != nil {
		return nil, err
	}

	x, y := elliptic.Unmarshal(S256(), s)
	return &ecdsa.PublicKey{Curve: S256(), X: x, Y: y}, nil
}

func ecrecover(hash, sig []byte) ([]byte, error) {
	return secp256k1.RecoverPubkey(hash, sig)
}

func pubkeyToAddress(p ecdsa.PublicKey) []byte {
	pubBytes := fromECDSAPub(&p)
	return bytesToAddress(crypto.Keccak256(pubBytes[1:])[12:])
}

func fromECDSAPub(pub *ecdsa.PublicKey) []byte {
	if pub == nil || pub.X == nil || pub.Y == nil {
		return nil
	}
	return elliptic.Marshal(S256(), pub.X, pub.Y)
}

func S256() elliptic.Curve {
	return secp256k1.S256()
}

func bytesToAddress(b []byte) []byte {
	a := make([]byte, 20)
	if len(b) > len(a) {
		b = b[len(b)-ERC20AddressLength:]
	}
	copy(a[ERC20AddressLength-len(b):], b)
	return a
}

func hexToAddress(s string) []byte {
	return bytesToAddress(common.FromHex(s))
}
