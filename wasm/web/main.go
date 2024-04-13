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
	abi "github.com/DogeProtocol/dp/wasm/accounts/abi"
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
	js.Global().Set("Scrypt", js.FuncOf(Scrypt))
	js.Global().Set("PublicKeyToAddress", js.FuncOf(PublicKeyToAddress))
	js.Global().Set("TxnSigningHash", js.FuncOf(TxnSigningHash))
	js.Global().Set("TxnHash", js.FuncOf(TxnHash))
	js.Global().Set("TxnData", js.FuncOf(TxnData))
	js.Global().Set("ContractData", js.FuncOf(ContractData))
	js.Global().Set("KeyPairToWalletJson", js.FuncOf(KeyPairToWalletJson))
	js.Global().Set("JsonToWalletKeyPair", js.FuncOf(JsonToWalletKeyPair))
	js.Global().Set("ParseBigFloat", js.FuncOf(ParseBigFloat))
	js.Global().Set("IsValidAddress", js.FuncOf(IsValidAddress))
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

func IsValidAddress(this js.Value, args []js.Value) interface{} {
	address := args[0].String()
	return common.IsHexAddress(address)
}

func TxnSigningHash(this js.Value, args []js.Value) interface{} {
	ts, err := transactionData(args)
	if err != nil {
		fmt.Println("TxnSigningHash err", err)
		return nil
	}

	tx := wasm.NewDefaultFeeTransaction(ts.Transaction[0].ChainId, ts.Transaction[0].Nonce,
		&ts.Transaction[0].ToAddress, ts.Transaction[0].Value,
		ts.Transaction[0].GasLimit, wasm.GAS_TIER_DEFAULT, ts.Transaction[0].Data)

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

func TxnHash(this js.Value, args []js.Value) interface{} {
	ts, err := transactionData(args)
	if err != nil {
		fmt.Println("txnhash err", err)
		return nil
	}

	tx := wasm.NewDefaultFeeTransaction(ts.Transaction[0].ChainId, ts.Transaction[0].Nonce,
		&ts.Transaction[0].ToAddress, ts.Transaction[0].Value,
		ts.Transaction[0].GasLimit, wasm.GAS_TIER_DEFAULT, ts.Transaction[0].Data)

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

func TxnData(this js.Value, args []js.Value) interface{} {
	ts, err := transactionData(args)
	if err != nil {
		fmt.Println("TxnData err", err)
		return nil
	}

	tx := wasm.NewDefaultFeeTransaction(ts.Transaction[0].ChainId, ts.Transaction[0].Nonce,
		&ts.Transaction[0].ToAddress, ts.Transaction[0].Value,
		ts.Transaction[0].GasLimit, wasm.GAS_TIER_DEFAULT, ts.Transaction[0].Data)

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

func ContractData(this js.Value, args []js.Value) interface{} {
	method := args[0].String()

	abiData, err := abi.JSON(strings.NewReader((args[1].String())))

	if err != nil {
		return nil
	}

	arguments := make([]interface{}, 0, len(args)-2)
	for _, i := range args[2:] {
		arguments = append(arguments, i.String())
	}

	data, err := abiData.Pack(method, arguments...)
	if err != nil {
		return nil
	}

	var d strings.Builder
	for i := 0; i < len(data); i++ {
		sh := data[i]
		d.WriteString(string(sh))
	}

	return d.String()
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

func ParseBigFloatInner(value string) (*big.Float, error) {
	f := new(big.Float)
	f.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	f.SetMode(big.ToNearestEven)
	_, err := fmt.Sscan(value, f)
	return f, err
}

func transactionData(args []js.Value) (transaction Transaction, err error) {
	fromAddress := common.HexToAddress(args[0].String())

	var nonceString string
	var nonceUint64 uint64
	fmt.Sscan(args[1].String(), &nonceString, &nonceUint64)
	nonce := nonceUint64

	toAddress := common.HexToAddress(args[2].String())

	var ethVal *big.Float
	var weiVal *big.Int
	ethVal, err = ParseBigFloatInner(args[3].String())
	if err != nil {
		fmt.Println("ParseBigFloatInner", args[3].String(), "err", err)
		return Transaction{}, err
	}
	weiVal = etherToWeiFloat(ethVal)

	var gasString string
	var gasUint64 uint64
	fmt.Sscan(args[4].String(), &gasString, &gasUint64)
	gasLimit := gasUint64

	var chainIdString string
	var chainIdInt64 int64
	fmt.Sscan(args[5].String(), &chainIdString, &chainIdInt64)
	chainId := big.NewInt(chainIdInt64)

	dataString := js.Global().Get("Uint8Array").New(args[6])
	data := make([]byte, dataString.Get("length").Int())
	js.CopyBytesToGo(data, dataString)

	transactionDetails := TransactionDetails{
		FromAddress: fromAddress, ToAddress: toAddress, Nonce: nonce, GasLimit: gasLimit,
		Value: weiVal, Data: data, ChainId: chainId}

	var t Transaction
	t.Transaction = append(t.Transaction, transactionDetails)

	return t, nil
}

func signTxHash(tx *wasm.Transaction, signer wasm.Signer, pubBytes, sigBytes []byte) (*wasm.Transaction, error) {
	sig := common.CombineTwoParts(sigBytes, pubBytes)
	return tx.WithSignature(signer, sig)
}

func weiToEther(val *big.Int) *big.Int {
	return new(big.Int).Div(val, big.NewInt(params.Ether))
}

func etherToWeiFloat(eth *big.Float) *big.Int {
	truncInt, _ := eth.Int(nil)
	truncInt = new(big.Int).Mul(truncInt, big.NewInt(params.Ether))
	fracStr := strings.Split(fmt.Sprintf("%.18f", eth), ".")[1]
	fracStr += strings.Repeat("0", 18-len(fracStr))
	fracInt, _ := new(big.Int).SetString(fracStr, 10)
	wei := new(big.Int).Add(truncInt, fracInt)
	return wei
}
