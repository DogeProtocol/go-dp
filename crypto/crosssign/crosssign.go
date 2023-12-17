package crosssign

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/json"
	"fmt"
	"github.com/DogeProtocol/dp/accounts"
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/common/hexutil"
	"github.com/DogeProtocol/dp/crypto"
	"github.com/DogeProtocol/dp/crypto/secp256k1"
)

const (
	ERC20AddressLength = 20
)

type SignDetails struct {
	Address string `json:"address"`
	Msg     string `json:"msg"`
	Sig     string `json:"sig"`
	Version string `json:"version"`
}

//signJsonData := "{\r\n  \"address\": \"0xF422Ec881E87B934A165DB64132a87fbd1753daD\",\r\n  \"msg\": \"Test message waller\",\r\n  \"sig\": \"0x5c73e35d19d6656f826c82513a4523a8c789762bacfd1ce5127f24c1e61cd59f7779132c3a390294db158735e398c4e87b726b87bef44ad840a47ac6ca06ef8d1b\",\r\n  \"version\": \"2\"\r\n}"

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
	recovered, err := sigToPub(messageDigest, signature)
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
