package main

import (
	"fmt"
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/common/hexutil"
	"github.com/DogeProtocol/dp/crypto"
	"github.com/DogeProtocol/dp/params"
	wasm "github.com/DogeProtocol/dp/wasm/core/types"
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
	GasPrice    *big.Int       `json:"gasPrice"`
	Value       *big.Int       `json:"value"`
	Data        []byte         `json:"data"`
	ChainId     *big.Int       `json:"chainId"`
}

func main() {
	done := make(chan struct{}, 0)
	js.Global().Set("PublicKeyToAddress", js.FuncOf(PublicKeyToAddress))
	js.Global().Set("TxMessage", js.FuncOf(TxMessage))
	js.Global().Set("TxHash", js.FuncOf(TxHash))
	js.Global().Set("TxData", js.FuncOf(TxData))
	js.Global().Set("DogeProtocolToWei", js.FuncOf(DogeProtocolToWei))
	js.Global().Set("WeiToDogeProtocol", js.FuncOf(WeiToDogeProtocol))
	js.Global().Set("ParseBigFloat", js.FuncOf(ParseBigFloat))
	<-done
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
		ts.Transaction[0].GasLimit, ts.Transaction[0].GasPrice,
		ts.Transaction[0].Data)

	signer := wasm.NewLondonSigner(ts.Transaction[0].ChainId)

	signerHash, err := signer.Hash(tx)
	if err != nil {
		return err
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
		ts.Transaction[0].GasLimit, ts.Transaction[0].GasPrice,
		ts.Transaction[0].Data)

	signer := wasm.NewLondonSigner(ts.Transaction[0].ChainId)

	pubData := js.Global().Get("Uint8Array").New(args[8])
	pubBytes := make([]byte, pubData.Get("length").Int())
	js.CopyBytesToGo(pubBytes, pubData)

	sigData := js.Global().Get("Uint8Array").New(args[9])
	sigBytes := make([]byte, sigData.Get("length").Int())
	js.CopyBytesToGo(sigBytes, sigData)

	signTx, err := signTxHash(tx, signer, pubBytes, sigBytes)
	if err != nil {
		return err
	}

	return signTx.Hash().String()
}

func TxData(this js.Value, args []js.Value) interface{} {
	ts := transaction(args)

	tx := wasm.NewTransaction(ts.Transaction[0].Nonce,
		ts.Transaction[0].ToAddress, ts.Transaction[0].Value,
		ts.Transaction[0].GasLimit, ts.Transaction[0].GasPrice,
		ts.Transaction[0].Data)

	signer := wasm.NewLondonSigner(ts.Transaction[0].ChainId)

	pubData := js.Global().Get("Uint8Array").New(args[8])
	pubBytes := make([]byte, pubData.Get("length").Int())
	js.CopyBytesToGo(pubBytes, pubData)

	sigData := js.Global().Get("Uint8Array").New(args[8])
	sigBytes := make([]byte, sigData.Get("length").Int())
	js.CopyBytesToGo(sigBytes, sigData)

	signTx, err := signTxHash(tx, signer, pubBytes, sigBytes)
	if err != nil {
		return err
	}

	signTxBinary, err := signTx.MarshalBinary()
	if err != nil {
		return err
	}

	signTxEncode := hexutil.Encode(signTxBinary)
	return signTxEncode
}

func DogeProtocolToWei(this js.Value, args []js.Value) interface{} {
	dp := new(big.Float)
	_, err := fmt.Sscan(args[0].String(), dp)
	if err != nil {
		return err
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
		return err
	}
	return f.String()
}

func WeiToDogeProtocol(this js.Value, args []js.Value) interface{} {
	wei := new(big.Int)
	_, err := fmt.Sscan(args[0].String(), wei)
	if err != nil {
		return err
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

	gasPrice := new(big.Int)
	_, err = fmt.Sscan(args[5].String(), gasPrice)
	if err != nil {
		panic(err)
	}

	var data []byte //args[6].String()

	var cTitle string
	var c int64
	fmt.Sscan(args[7].String(), &cTitle, &c)
	chainId := big.NewInt(c)

	transactionDetails := TransactionDetails{
		FromAddress: fromAddress, ToAddress: toAddress, Nonce: nonce, GasLimit: gasLimit,
		GasPrice: gasPrice, Value: value, Data: data, ChainId: chainId}

	var t Transaction
	t.Transaction = append(t.Transaction, transactionDetails)
	return t
}

func signTxHash(tx *wasm.Transaction, signer wasm.Signer, pubBytes, sigBytes []byte) (*wasm.Transaction, error) {
	sig := common.CombineTwoParts(sigBytes, pubBytes)
	return tx.WithSignature(signer, sig)
}
