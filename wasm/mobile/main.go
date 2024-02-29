package main

import "C"
import (
	"fmt"
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/common/hexutil"
	"github.com/DogeProtocol/dp/crypto"
	"github.com/DogeProtocol/dp/params"
	ks "github.com/DogeProtocol/dp/wasm/accounts/keystore"
	wasm "github.com/DogeProtocol/dp/wasm/core/types"
	"github.com/google/uuid"
	"math/big"
	"strconv"
	"strings"
	"unsafe"
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

}

//export PublicKeyToAddress
func PublicKeyToAddress(pKey_str *C.char, pk_count int) (*C.char, *C.char) {
	pubBytes := C.GoBytes(unsafe.Pointer(pKey_str), C.int(pk_count))
	address := common.BytesToAddress(crypto.Keccak256(pubBytes[:])[common.AddressTruncateBytes:]).String()
	return C.CString(address), nil
}

//export TxMessage
func TxMessage(from, nonce, to, value, gasLimit, data, chainId *C.char) (*C.char, *C.char) {
	ts := transaction(C.GoString(from), C.GoString(nonce), C.GoString(to),
		C.GoString(value), C.GoString(gasLimit), C.GoString(data), C.GoString(chainId))

	tx := wasm.NewTransaction(ts.Transaction[0].Nonce,
		ts.Transaction[0].ToAddress, ts.Transaction[0].Value,
		ts.Transaction[0].GasLimit, ts.Transaction[0].Data)

	signer := wasm.NewLondonSigner(ts.Transaction[0].ChainId)

	signerHash, err := signer.Hash(tx)
	if err != nil {
		return nil, C.CString(err.Error())
	}

	var message strings.Builder
	for i := 0; i < len(signerHash); i++ {
		sh := signerHash[i]
		message.WriteString(string(sh))
	}

	return C.CString(message.String()), nil
}

//export TxHash
func TxHash(from, nonce, to, value, gasLimit, data, chainId,
	pKey_str, sig_str *C.char, pk_count int, sig_count int) (*C.char, *C.char) {
	ts := transaction(C.GoString(from), C.GoString(nonce), C.GoString(to),
		C.GoString(value), C.GoString(gasLimit), C.GoString(data), C.GoString(chainId))

	tx := wasm.NewTransaction(ts.Transaction[0].Nonce,
		ts.Transaction[0].ToAddress, ts.Transaction[0].Value,
		ts.Transaction[0].GasLimit, ts.Transaction[0].Data)

	signer := wasm.NewLondonSigner(ts.Transaction[0].ChainId)

	pubBytes := C.GoBytes(unsafe.Pointer(pKey_str), C.int(pk_count))
	sigBytes := C.GoBytes(unsafe.Pointer(sig_str), C.int(sig_count))

	signTx, err := signTxHash(tx, signer, pubBytes, sigBytes)
	if err != nil {
		return nil, C.CString(err.Error())
	}

	return C.CString(signTx.Hash().String()), nil
}

//export TxData
func TxData(from, nonce, to, value, gasLimit, data, chainId,
	pKey_str, sig_str *C.char, pk_count int, sig_count int) (*C.char, *C.char) {

	ts := transaction(C.GoString(from), C.GoString(nonce), C.GoString(to),
		C.GoString(value), C.GoString(gasLimit), C.GoString(data), C.GoString(chainId))

	tx := wasm.NewTransaction(ts.Transaction[0].Nonce,
		ts.Transaction[0].ToAddress, ts.Transaction[0].Value,
		ts.Transaction[0].GasLimit, ts.Transaction[0].Data)

	signer := wasm.NewLondonSigner(ts.Transaction[0].ChainId)

	pubBytes := C.GoBytes(unsafe.Pointer(pKey_str), C.int(pk_count))
	sigBytes := C.GoBytes(unsafe.Pointer(sig_str), C.int(sig_count))

	signTx, err := signTxHash(tx, signer, pubBytes, sigBytes)
	if err != nil {
		return nil, C.CString(err.Error())
	}

	signTxBinary, err := signTx.MarshalBinary()
	if err != nil {
		return nil, C.CString(err.Error())
	}

	signTxEncode := hexutil.Encode(signTxBinary)
	return C.CString(signTxEncode), nil
}

//export ExportKey
func ExportKey(skKeyStr, pkKeyStr, authentication *C.char, skCount int, pkCount int) (*C.char, *C.char) {

	privateBytes := C.GoBytes(unsafe.Pointer(skKeyStr), C.int(skCount))
	publicBytes := C.GoBytes(unsafe.Pointer(pkKeyStr), C.int(pkCount))

	auth := C.GoString(authentication)

	var pubKeyAddress = common.BytesToAddress(crypto.Keccak256(publicBytes[:])[common.AddressTruncateBytes:])

	id, err := uuid.NewRandom()
	if err != nil {
		panic(fmt.Sprintf("Could not create random uuid: %v", err))
	}

	publicKey := ks.PublicKey{
		PubData: publicBytes,
	}

	privateKey := &ks.PrivateKey{
		PublicKey: publicKey,
		PriData:   privateBytes,
	}

	key := &ks.Key{
		Id:         id,
		Address:    pubKeyAddress,
		PrivateKey: privateKey,
	}

	keyJson, err := ks.EncryptKey(key, pubKeyAddress.Bytes(), auth, ks.StandardScryptN, ks.StandardScryptP)
	if err != nil {
		return nil, C.CString(err.Error())
	}

	return C.CString(string(keyJson)), nil
}

//export ImportKey
func ImportKey(skKeyStr, authentication *C.char, skCount int) (*C.char, *C.char) {

	keyJson := C.GoBytes(unsafe.Pointer(skKeyStr), C.int(skCount))

	auth := C.GoString(authentication)

	key, err := ks.DecryptKey(keyJson, auth)
	if err != nil {
		return nil, C.CString(err.Error())
	}
	return C.CString(string(key.PrivateKey.PriData)), nil
}

//export DogeProtocolToWei
func DogeProtocolToWei(quantity *C.char) (*C.char, *C.char) {
	dp := new(big.Float)
	_, err := fmt.Sscan(C.GoString(quantity), dp)
	if err != nil {
		return nil, C.CString(err.Error())
	}
	truncInt, _ := dp.Int(nil)
	truncInt = new(big.Int).Mul(truncInt, big.NewInt(params.Ether))
	fracStr := strings.Split(fmt.Sprintf("%.18f", dp), ".")[1]
	fracStr += strings.Repeat("0", 18-len(fracStr))
	fracInt, _ := new(big.Int).SetString(fracStr, 10)
	wei := new(big.Int).Add(truncInt, fracInt)
	return C.CString(wei.String()), nil
}

//export ParseBigFloat
func ParseBigFloat(quantity *C.char) (*C.char, *C.char) {
	var value string
	value = C.GoString(quantity)
	f := new(big.Float)
	f.SetPrec(236)
	f.SetMode(big.ToNearestEven)
	_, err := fmt.Sscan(value, f)
	if err != nil {
		return nil, C.CString(err.Error())
	}
	return C.CString(f.String()), nil
}

//export WeiToDogeProtocol
func WeiToDogeProtocol(quantity *C.char) (*C.char, *C.char) {
	wei := new(big.Int)
	_, err := fmt.Sscan(C.GoString(quantity), wei)
	if err != nil {
		return nil, C.CString(err.Error())
	}
	f := new(big.Float)
	f.SetPrec(236)
	f.SetMode(big.ToNearestEven)
	fWei := new(big.Float)
	fWei.SetPrec(236)
	fWei.SetMode(big.ToNearestEven)
	dp := f.Quo(fWei.SetInt(wei), big.NewFloat(params.Ether))
	return C.CString(dp.String()), nil
}

func transaction(args0, args1, args2, args3, args4, args5, args6 string) (transaction Transaction) {

	var fromAddress = common.HexToAddress(args0)
	n, _ := strconv.Atoi(args1)
	var nonce = uint64(n)
	var toAddress = common.HexToAddress(args2)
	var value, _ = new(big.Int).SetString(args3, 0)
	g, _ := strconv.Atoi(args4)
	var gasLimit = uint64(g)
	var data []byte //args5.String()
	var chainId, _ = new(big.Int).SetString(args6, 0)

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
