//go:build js && wasm
// +build js,wasm

package main

import (
	"encoding/base64"
	"fmt"
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/common/hexutil"
	"github.com/DogeProtocol/dp/crypto"
	"github.com/DogeProtocol/dp/params"
	ks "github.com/DogeProtocol/dp/wasm/accounts/keystore"
	wasm "github.com/DogeProtocol/dp/wasm/core/types"
	"github.com/google/uuid"
	"golang.org/x/crypto/scrypt"
	"math/big"
	"strings"
	"syscall/js"
)

type Transaction struct {
	Transaction []TransactionDetails `json:"transaction"`
}

type TransactionDetails struct {
	FromAddress common.Address `json:"fromAddress"`
	ToAddress   common.Address `json:"toAddress"`
	Nonce       uint64         `json:"nonce"`
	GasLimit    uint64         `json:"gasLimit"`
	Value       *big.Int       `json:"value"`
	Data        []byte         `json:"data"`
	ChainId     *big.Int       `json:"chainId"`
}

func main() {
	done := make(chan struct{}, 0)
	js.Global().Set("PublicKeyToAddress", js.FuncOf(PublicKeyToAddress))
	js.Global().Set("Scrypt", js.FuncOf(Scrypt))
	js.Global().Set("TxMessage", js.FuncOf(TxMessage))
	js.Global().Set("TxHash", js.FuncOf(TxHash))
	js.Global().Set("TxData", js.FuncOf(TxData))
	js.Global().Set("KeyPairToWalletJson", js.FuncOf(KeyPairToWalletJson))
	js.Global().Set("JsonToWalletKeyPair", js.FuncOf(JsonToWalletKeyPair))
	js.Global().Set("DogeProtocolToWei", js.FuncOf(DogeProtocolToWei))
	js.Global().Set("WeiToDogeProtocol", js.FuncOf(WeiToDogeProtocol))
	js.Global().Set("ParseBigFloat", js.FuncOf(ParseBigFloat))
	<-done
}

func Scrypt(this js.Value, args []js.Value) interface{} {
	secret := args[0].String()

	salt, err := base64.StdEncoding.DecodeString(args[1].String())
	if err != nil {
		return nil
	}

	derivedKey, err := scrypt.Key([]byte(secret), salt, 262144, 8, 1, 32)
	if err != nil {
		return nil
	}

	return base64.StdEncoding.EncodeToString(derivedKey)
}

func PublicKeyToAddress(this js.Value, args []js.Value) interface{} {
	pubData := js.Global().Get("Uint8Array").New(args[0])
	pubBytes := make([]byte, pubData.Get("length").Int())
	js.CopyBytesToGo(pubBytes, pubData)
	return common.BytesToAddress(crypto.Keccak256(pubBytes[:])[common.AddressTruncateBytes:]).String()
}

func TxMessage(this js.Value, args []js.Value) interface{} {
	ts := transaction(args)

	tx := wasm.NewTransaction(ts.Transaction[0].Nonce,
		ts.Transaction[0].ToAddress, ts.Transaction[0].Value,
		ts.Transaction[0].GasLimit, ts.Transaction[0].Data)

	signer := wasm.NewLondonSigner(ts.Transaction[0].ChainId)

	signerHash, err := signer.Hash(tx)
	if err != nil {
		return nil
	}

	var message strings.Builder
	for i := 0; i < len(signerHash); i++ {
		sh := signerHash[i]
		message.WriteString(string(sh))
	}

	return message.String()
}

func TxHash(this js.Value, args []js.Value) interface{} {
	ts := transaction(args)

	tx := wasm.NewTransaction(ts.Transaction[0].Nonce,
		ts.Transaction[0].ToAddress, ts.Transaction[0].Value,
		ts.Transaction[0].GasLimit, ts.Transaction[0].Data)

	signer := wasm.NewLondonSigner(ts.Transaction[0].ChainId)

	pubData := js.Global().Get("Uint8Array").New(args[7])
	pubBytes := make([]byte, pubData.Get("length").Int())
	js.CopyBytesToGo(pubBytes, pubData)

	sigData := js.Global().Get("Uint8Array").New(args[8])
	sigBytes := make([]byte, sigData.Get("length").Int())
	js.CopyBytesToGo(sigBytes, sigData)

	signTx, err := signTxHash(tx, signer, pubBytes, sigBytes)
	if err != nil {
		return nil
	}

	return signTx.Hash().String()
}

func TxData(this js.Value, args []js.Value) interface{} {
	ts := transaction(args)

	tx := wasm.NewTransaction(ts.Transaction[0].Nonce,
		ts.Transaction[0].ToAddress, ts.Transaction[0].Value,
		ts.Transaction[0].GasLimit, ts.Transaction[0].Data)

	signer := wasm.NewLondonSigner(ts.Transaction[0].ChainId)

	pubData := js.Global().Get("Uint8Array").New(args[7])
	pubBytes := make([]byte, pubData.Get("length").Int())
	js.CopyBytesToGo(pubBytes, pubData)

	sigData := js.Global().Get("Uint8Array").New(args[8])
	sigBytes := make([]byte, sigData.Get("length").Int())
	js.CopyBytesToGo(sigBytes, sigData)

	signTx, err := signTxHash(tx, signer, pubBytes, sigBytes)
	if err != nil {
		return nil
	}

	signTxBinary, err := signTx.MarshalBinary()
	if err != nil {
		return nil
	}

	signTxEncode := hexutil.Encode(signTxBinary)
	return signTxEncode
}

func KeyPairToWalletJson(this js.Value, args []js.Value) interface{} {
	privData := js.Global().Get("Uint8Array").New(args[0])
	privBytes := make([]byte, privData.Get("length").Int())
	js.CopyBytesToGo(privBytes, privData)

	pubData := js.Global().Get("Uint8Array").New(args[1])
	pubBytes := make([]byte, pubData.Get("length").Int())
	js.CopyBytesToGo(pubBytes, pubData)

	passphrase := args[2].String()

	var pubKeyAddress = crypto.PublicKeyBytesToAddress(pubBytes)

	id, err := uuid.NewRandom()
	if err != nil {
		panic(fmt.Sprintf("Could not create random uuid: %v", err))
	}

	publicKey := ks.PublicKey{
		PubData: pubBytes,
	}

	privateKey := &ks.PrivateKey{
		PublicKey: publicKey,
		PriData:   privBytes,
	}

	key := &ks.Key{
		Id:         id,
		Address:    pubKeyAddress,
		PrivateKey: privateKey,
	}

	keyJson, err := ks.EncryptKey(key, pubKeyAddress.Bytes(), passphrase, ks.StandardScryptN, ks.StandardScryptP)
	if err != nil {
		return nil
	}

	return string(keyJson[:])
}

func JsonToWalletKeyPair(this js.Value, args []js.Value) interface{} {
	keyJson := []byte(args[0].String())
	passphrase := args[1].String()

	key, err := ks.DecryptKey(keyJson, passphrase)
	if err != nil {
		return nil
	}
	return base64.StdEncoding.EncodeToString(key.PrivateKey.PriData) + "," + base64.StdEncoding.EncodeToString(key.PrivateKey.PubData)
}

func DogeProtocolToWei(this js.Value, args []js.Value) interface{} {
	dp := new(big.Float)
	_, err := fmt.Sscan(args[0].String(), dp)
	if err != nil {
		return nil
	}
	truncInt, _ := dp.Int(nil)
	truncInt = new(big.Int).Mul(truncInt, big.NewInt(params.Ether))
	fracStr := strings.Split(fmt.Sprintf("%.18f", dp), ".")[1]
	fracStr += strings.Repeat("0", 18-len(fracStr))
	fracInt, _ := new(big.Int).SetString(fracStr, 10)
	wei := new(big.Int).Add(truncInt, fracInt)
	return wei.String()
}

// ParseBigFloat parse string value to big.Float
func ParseBigFloat(this js.Value, args []js.Value) interface{} {
	var value string
	value = args[0].String()
	f := new(big.Float)
	f.SetPrec(236)
	f.SetMode(big.ToNearestEven)
	_, err := fmt.Sscan(value, f)
	if err != nil {
		return nil
	}
	return f.String()
}

func WeiToDogeProtocol(this js.Value, args []js.Value) interface{} {
	wei := new(big.Int)
	_, err := fmt.Sscan(args[0].String(), wei)
	if err != nil {
		return nil
	}
	f := new(big.Float)
	f.SetPrec(236)
	f.SetMode(big.ToNearestEven)
	fWei := new(big.Float)
	fWei.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	fWei.SetMode(big.ToNearestEven)
	dp := f.Quo(fWei.SetInt(wei), big.NewFloat(params.Ether))
	return dp.String()
}

func transaction(args []js.Value) (transaction Transaction) {

	fromAddress := common.HexToAddress(args[0].String())

	var nTitle string
	var n uint64
	fmt.Sscan(args[1].String(), &nTitle, &n)
	nonce := n

	toAddress := common.HexToAddress(args[2].String())

	value := new(big.Int)
	_, err := fmt.Sscan(args[3].String(), value)
	if err != nil {
		panic(err)
	}

	var lTitle string
	var l uint64
	fmt.Sscan(args[4].String(), &lTitle, &l)
	gasLimit := l

	var data []byte //args[5].String()

	var cTitle string
	var c int64
	fmt.Sscan(args[6].String(), &cTitle, &c)
	chainId := big.NewInt(c)

	transactionDetails := TransactionDetails{
		FromAddress: fromAddress, ToAddress: toAddress, Nonce: nonce, GasLimit: gasLimit,
		Value: value, Data: data, ChainId: chainId}

	var t Transaction
	t.Transaction = append(t.Transaction, transactionDetails)
	return t
}

func signTxHash(tx *wasm.Transaction, signer wasm.Signer, pubBytes, sigBytes []byte) (*wasm.Transaction, error) {
	sig := common.CombineTwoParts(sigBytes, pubBytes)
	return tx.WithSignature(signer, sig)
}
