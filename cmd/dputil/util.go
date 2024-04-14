package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DogeProtocol/dp/accounts"
	"github.com/DogeProtocol/dp/accounts/abi/bind"
	"github.com/DogeProtocol/dp/accounts/keystore"
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/common/hexutil"
	"github.com/DogeProtocol/dp/core/types"
	"github.com/DogeProtocol/dp/crypto/cryptobase"
	"github.com/DogeProtocol/dp/crypto/signaturealgorithm"
	"github.com/DogeProtocol/dp/ethclient"
	"github.com/DogeProtocol/dp/params"
	"github.com/DogeProtocol/dp/systemcontracts/conversion"
	"github.com/DogeProtocol/dp/systemcontracts/staking"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type KeyStore struct {
	Handle *keystore.KeyStore
}

type BalanceData struct {
	Result struct {
		Balance string `json:"_balance"`
		Nonce   string `json:"nonce"`
	}
}

func etherToWei(val *big.Int) *big.Int {
	return new(big.Int).Mul(val, big.NewInt(params.Ether))
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

func getBalance(address string) (ethBalance string, weiBalance string, err error) {
	client, err := ethclient.Dial(rawURL)
	if err != nil {
		return "", "", err
	}
	balance, err := client.BalanceAt(context.Background(), common.HexToAddress(address), nil)
	if err != nil {
		return "", "", err
	}
	return weiToEther(balance).String(), balance.String(), nil
}

func requestGetBalance(address string) (ethBalance string, weiBalance string, nonce string, err error) {
	request, err := http.NewRequest("GET", READ_API_URL+"/api/accounts/"+address+"/balance", nil)
	if err != nil {
		return "", "", "", err
	}
	request.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		return "", "", "", err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", "", "", err
	}

	var balanceData BalanceData
	err = json.Unmarshal(body, &balanceData)
	if err != nil {
		return "", "", "", err
	}

	if len(balanceData.Result.Balance) == 0 {
		balanceData.Result.Balance = "0"
	}
	if len(balanceData.Result.Nonce) == 0 {
		balanceData.Result.Nonce = "0"
	}

	balance := new(big.Int)
	_, err = fmt.Sscan(balanceData.Result.Balance, balance)
	if err != nil {
		return "", "", "", err
	}

	return weiToEther(balance).String(), balanceData.Result.Balance, balanceData.Result.Nonce, nil
}

func findAllAddresses() ([]string, error) {
	keyfileDir := os.Getenv("DP_KEY_FILE_DIR")
	if len(keyfileDir) == 0 {
		return nil, errors.New("Both DP_KEY_FILE and DP_KEY_FILE_DIR environment variables not set")
	}

	files, err := ioutil.ReadDir(keyfileDir)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var addresses []string
	addresses = make([]string, 0)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		columns := strings.Split(file.Name(), "--")
		if len(columns) != 3 {
			continue
		}
		addresses = append(addresses, columns[2])
	}

	return addresses, nil
}

func findKeyFile(keyAddress string) (string, error) {
	keyfile := os.Getenv("DP_KEY_FILE")
	if len(keyfile) > 0 {
		return keyfile, nil
	}

	keyfileDir := os.Getenv("DP_KEY_FILE_DIR")
	if len(keyfileDir) == 0 {
		return "", errors.New("Both DP_KEY_FILE and DP_KEY_FILE_DIR environment variables not set")
	}

	files, err := ioutil.ReadDir(keyfileDir)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	addr := strings.ToLower(strings.Replace(keyAddress, "0x", "", 1))
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.Contains(strings.ToLower(file.Name()), addr) {
			return filepath.Join(keyfileDir, file.Name()), nil
		}
	}

	return "", errors.New("could not find key file")
}

type ConnectionContext struct {
	From   string
	Client *ethclient.Client
	Key    *keystore.Key
}

func GetKeyFromFile(keyFile string, accPwd string) (*signaturealgorithm.PrivateKey, error) {
	secretKey, err := ReadDataFile(keyFile)
	if err != nil {
		return nil, err
	}

	password := accPwd
	key, err := keystore.DecryptKey(secretKey, password)
	if err != nil {
		return nil, err
	}

	return key.PrivateKey, nil
}

func GetConnectionContext(from string) (*ConnectionContext, error) {
	keyFile, err := findKeyFile(from)
	if err != nil {
		return nil, err
	}

	secretKey, err := ReadDataFile(keyFile)
	if err != nil {
		return nil, err
	}

	password := os.Getenv("DP_ACC_PWD")
	key, err := keystore.DecryptKey(secretKey, password)
	if err != nil {
		return nil, err
	}

	client, err := ethclient.Dial(rawURL)
	if err != nil {
		return nil, err
	}

	return &ConnectionContext{
		From:   from,
		Client: client,
		Key:    key,
	}, nil
}

func sendVia(connectionContext *ConnectionContext, to string, quantity string, nonce uint64) (string, uint64, error) {
	if connectionContext == nil {
		return "", 0, errors.New("nil")
	}
	fromAddress := common.HexToAddress(connectionContext.From)
	toAddress := common.HexToAddress(to)

	if nonce == 0 {
		nonceTmp, err := connectionContext.Client.PendingNonceAt(context.Background(), fromAddress)
		if err != nil {
			return "", 0, err
		}
		nonce = nonceTmp
	}

	chainID, err := connectionContext.Client.NetworkID(context.Background())
	if err != nil {
		return "", 0, err
	}
	gasLimit := uint64(21000)
	gasPrice, err := connectionContext.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", 0, err
	}

	v, err := ParseBigFloat(quantity)
	if err != nil {
		return "", 0, err
	}

	value := etherToWeiFloat(v)

	var data []byte
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)
	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainID), connectionContext.Key.PrivateKey)
	if err != nil {
		return "", 0, err
	}
	err = connectionContext.Client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", 0, err
	}

	fmt.Println("Sent Transaction", "from", fromAddress, "to", toAddress, "quantity", quantity, "Transaction", signedTx.Hash().Hex())
	return signedTx.Hash().Hex(), nonce, nil
}

func send(from string, to string, quantity string) (string, error) {
	keyFile, err := findKeyFile(from)
	if err != nil {
		return "", err
	}

	fmt.Println("keyFile", keyFile)
	secretKey, err := ReadDataFile(keyFile)
	if err != nil {
		return "", err
	}
	password := os.Getenv("DP_ACC_PWD")
	key, err := keystore.DecryptKey(secretKey, password)
	if err != nil {
		return "", err
	}

	client, err := ethclient.Dial(rawURL)
	if err != nil {
		return "", err
	}

	fromAddress := common.HexToAddress(from)
	toAddress := common.HexToAddress(to)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", err
	}
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return "", err
	}
	gasLimit := uint64(21000)

	v, err := ParseBigFloat(quantity)
	if err != nil {
		return "", err
	}

	value := etherToWeiFloat(v)

	var data []byte
	tx := types.NewDefaultFeeTransaction(chainID, nonce, &toAddress, value, gasLimit, types.GAS_TIER_DEFAULT, data)
	fmt.Println("chainID", chainID)

	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainID), key.PrivateKey)
	if err != nil {
		fmt.Println("signedTx err", err)
		return "", err
	}
	fmt.Println("signedTx ok")
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}

	fmt.Println("Sent Transaction", "from", fromAddress, "to", toAddress, "quantity", quantity, "Transaction", signedTx.Hash().Hex())
	return signedTx.Hash().Hex(), nil
}

func GetTransaction(txnHash string) (string, error) {
	client, err := ethclient.Dial(rawURL)
	if err != nil {
		return "", err
	}
	hash := common.HexToHash(txnHash)
	fmt.Println("hash", hash)
	return client.RawTransactionByHash(context.Background(), hash)
}

// ParseBigFloat parse string value to big.Float
func ParseBigFloat(value string) (*big.Float, error) {
	f := new(big.Float)
	f.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	f.SetMode(big.ToNearestEven)
	_, err := fmt.Sscan(value, f)
	return f, err
}

func ReadDataFile(filename string) ([]byte, error) {
	// Open our jsonFile
	jsonFile, err := os.Open(filename)
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	fmt.Println("Successfully Opened ", filename)
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	return byteValue, nil
}

func (ks *KeyStore) CreateNewKeys(password string) accounts.Account {
	account, err := ks.Handle.NewAccount(password)
	if err != nil {
		log.Println(err.Error())
	}
	return account
}

func (ks *KeyStore) GetKeysByAddress(address string) (accounts.Account, error) {
	var account accounts.Account
	var err error
	if ks.Handle.HasAddress(common.HexToAddress(address)) {
		if account, err = ks.Handle.Find(accounts.Account{Address: common.HexToAddress(address)}); err != nil {
			return accounts.Account{}, err
		}
	}
	return account, nil
}

func (ks *KeyStore) GetAllKeys() []accounts.Account {
	return ks.Handle.Accounts()
}

func SetUpKeyStore() *KeyStore {
	dataDir := os.Getenv("DP_DATA_PATH")
	if dataDir == "" {
		dataDir = "data"
	}

	ks := &KeyStore{}
	ks.Handle = keystore.NewKeyStore(dataDir, keystore.LightScryptN, keystore.LightScryptP)
	return ks
}

func convertCoins(ethAddress string, ethSignature string, key *signaturealgorithm.PrivateKey) error {
	client, err := ethclient.Dial(rawURL)
	if err != nil {
		return err
	}

	fromAddress, err := cryptobase.SigAlg.PublicKeyToAddress(&key.PublicKey)
	if err != nil {
		return err
	}
	contractAddress := common.HexToAddress(conversion.CONVERSION_CONTRACT)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return err
	}
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return err
	}
	txnOpts, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		return err
	}
	txnOpts.From = fromAddress
	txnOpts.Nonce = big.NewInt(int64(nonce))
	txnOpts.GasLimit = uint64(210000)

	contract, err := conversion.NewConversion(contractAddress, client)
	if err != nil {
		return err
	}

	tx, err := contract.RequestConversion(txnOpts, ethAddress, ethSignature)
	if err != nil {
		return err
	}

	fmt.Println("Your request to get the quantum dp coins has been added to the queue for processing. Please check your account balance after 10 minutes.")
	fmt.Println("The transaction hash for tracking this request is: ", tx.Hash())
	fmt.Println("Your can you use the following command to check your account balance: ")
	fmt.Println("dputil balance [YOUR_QUANTUM_ADDRESS]")
	fmt.Println("Do double check that you have backed up your quantum wallet safely in multiple devices and offline backups. And remember your password!")
	fmt.Println()

	time.Sleep(1000 * time.Millisecond)

	return nil
}

func requestConvertCoins(ethAddress string, ethSignature string, key *signaturealgorithm.PrivateKey) error {

	fromAddress, err := cryptobase.SigAlg.PublicKeyToAddress(&key.PublicKey)

	if err != nil {
		return err
	}
	_, _, n, err := requestGetBalance(fromAddress.String())
	if err != nil {
		return err
	}

	var nonce uint64
	fmt.Sscan(n, &nonce)

	contractAddress := common.HexToAddress(conversion.CONVERSION_CONTRACT)

	txnOpts, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(123123))

	if err != nil {
		return err
	}

	txnOpts.From = fromAddress
	txnOpts.Nonce = big.NewInt(int64(nonce))
	txnOpts.GasLimit = uint64(210000)

	method := conversion.GetContract_Method_requestConversion()
	abiData, err := conversion.GetConversionContract_ABI()
	if err != nil {
		return err
	}

	input, err := abiData.Pack(method, ethAddress, ethSignature)
	if err != nil {
		return err
	}

	baseTx := types.NewDefaultFeeTransactionSimple(nonce, &contractAddress, txnOpts.Value,
		txnOpts.GasLimit, input)

	var rawTx *types.Transaction
	rawTx = types.NewTx(baseTx)

	if txnOpts.Signer == nil {
		return errors.New("no signer to authorize the transaction with")
	}

	signTx, err := txnOpts.Signer(txnOpts.From, rawTx)
	if err != nil {
		return err
	}

	signTxBinary, err := signTx.MarshalBinary()
	if err != nil {
		return err
	}

	tx := signTx
	txData := hexutil.Encode(signTxBinary)

	var jsonStr = []byte(`{"txnData" : "` + txData + `"}`)

	request, err := http.NewRequest("POST", WRITE_API_URL+"/api/transactions", bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	fmt.Println("Your request to get the quantum dp coins has been added to the queue for processing. Please check your account balance after 10 minutes.")
	fmt.Println("The transaction hash for tracking this request is: ", tx.Hash())
	fmt.Println("Your can you use the following command to check your account balance: ")
	fmt.Println("dputil balance [YOUR_QUANTUM_ADDRESS]")
	fmt.Println("Do double check that you have backed up your quantum wallet safely in multiple devices and offline backups. And remember your password!")
	fmt.Println()

	time.Sleep(1000 * time.Millisecond)

	return nil
}

func newDeposit(validatorAddress string, depositAmount string, key *signaturealgorithm.PrivateKey) error {

	client, err := ethclient.Dial(rawURL)
	if err != nil {
		return err
	}

	fromAddress, err := cryptobase.SigAlg.PublicKeyToAddress(&key.PublicKey)

	if err != nil {
		return err
	}

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return err
	}

	contractAddress := common.HexToAddress(staking.STAKING_CONTRACT)
	txnOpts, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(123123))

	if err != nil {
		return err
	}

	txnOpts.From = fromAddress
	txnOpts.Nonce = big.NewInt(int64(nonce))
	txnOpts.GasLimit = uint64(250000)

	val, _ := ParseBigFloat(depositAmount)
	txnOpts.Value = etherToWeiFloat(val)

	contract, err := staking.NewStaking(contractAddress, client)
	if err != nil {
		return err
	}

	tx, err := contract.NewDeposit(txnOpts, common.HexToAddress(validatorAddress))
	if err != nil {
		return err
	}

	fmt.Println("Your request to deposit has been added to the queue for processing. Please check your account balance after 10 minutes.")
	fmt.Println("The transaction hash for tracking this request is: ", tx.Hash())
	fmt.Println()

	time.Sleep(1000 * time.Millisecond)

	return nil
}

func requestNewDeposit(validatorAddress string, depositAmount string, key *signaturealgorithm.PrivateKey) error {

	fromAddress, err := cryptobase.SigAlg.PublicKeyToAddress(&key.PublicKey)

	if err != nil {
		return err
	}
	_, _, n, err := requestGetBalance(fromAddress.String())
	if err != nil {
		return err
	}

	var nonce uint64
	fmt.Sscan(n, &nonce)

	contractAddress := common.HexToAddress(staking.STAKING_CONTRACT)
	txnOpts, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(123123))

	if err != nil {
		return err
	}

	txnOpts.From = fromAddress
	txnOpts.Nonce = big.NewInt(int64(nonce))
	txnOpts.GasLimit = uint64(250000)

	val, _ := ParseBigFloat(depositAmount)
	txnOpts.Value = etherToWeiFloat(val)

	method := staking.GetContract_Method_NewDeposit()
	abiData, err := staking.GetStakingContract_ABI()
	if err != nil {
		return err
	}

	input, err := abiData.Pack(method, common.HexToAddress(validatorAddress))
	if err != nil {
		return err
	}

	baseTx := types.NewDefaultFeeTransactionSimple(nonce, &contractAddress, txnOpts.Value,
		txnOpts.GasLimit, input)

	var rawTx *types.Transaction
	rawTx = types.NewTx(baseTx)

	if txnOpts.Signer == nil {
		return errors.New("no signer to authorize the transaction with")
	}

	signTx, err := txnOpts.Signer(txnOpts.From, rawTx)
	if err != nil {
		return err
	}

	signTxBinary, err := signTx.MarshalBinary()
	if err != nil {
		return err
	}

	tx := signTx
	txData := hexutil.Encode(signTxBinary)

	var jsonStr = []byte(`{"txnData" : "` + txData + `"}`)

	request, err := http.NewRequest("POST", WRITE_API_URL+"/api/transactions", bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	fmt.Println("Your request to deposit has been added to the queue for processing. Please check your account balance after 10 minutes.")
	fmt.Println("The transaction hash for tracking this request is: ", tx.Hash())
	fmt.Println()

	time.Sleep(1000 * time.Millisecond)

	return nil
}

func initiateWithdrawal(key *signaturealgorithm.PrivateKey) error {

	client, err := ethclient.Dial(rawURL)
	if err != nil {
		return err
	}

	fromAddress, err := cryptobase.SigAlg.PublicKeyToAddress(&key.PublicKey)
	if err != nil {
		return err
	}

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return err
	}

	contractAddress := common.HexToAddress(staking.STAKING_CONTRACT)
	txnOpts, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(123123))

	if err != nil {
		return err
	}

	txnOpts.From = fromAddress
	txnOpts.Nonce = big.NewInt(int64(nonce))
	txnOpts.GasLimit = uint64(210000)

	val, _ := ParseBigFloat("0")
	txnOpts.Value = etherToWeiFloat(val)

	contract, err := staking.NewStaking(contractAddress, client)
	if err != nil {
		return err
	}

	tx, err := contract.InitiateWithdrawal(txnOpts)
	if err != nil {
		return err
	}

	fmt.Println("Your request to initial withdrawal has been added to the queue for processing.")
	fmt.Println("The transaction hash for tracking this request is: ", tx.Hash())
	fmt.Println()

	time.Sleep(1000 * time.Millisecond)

	return nil
}

func requestInitiateWithdrawal(key *signaturealgorithm.PrivateKey) error {

	fromAddress, err := cryptobase.SigAlg.PublicKeyToAddress(&key.PublicKey)

	if err != nil {
		return err
	}
	_, _, n, err := requestGetBalance(fromAddress.String())
	if err != nil {
		return err
	}

	var nonce uint64
	fmt.Sscan(n, &nonce)

	contractAddress := common.HexToAddress(staking.STAKING_CONTRACT)
	txnOpts, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(123123))

	if err != nil {
		return err
	}

	txnOpts.From = fromAddress
	txnOpts.Nonce = big.NewInt(int64(nonce))
	txnOpts.GasLimit = uint64(210000)

	val, _ := ParseBigFloat("0")
	txnOpts.Value = etherToWeiFloat(val)

	method := staking.GetContract_Method_InitiateWithdrawal()
	abiData, err := staking.GetStakingContract_ABI()
	if err != nil {
		return err
	}

	input, err := abiData.Pack(method)
	if err != nil {
		return err
	}

	baseTx := types.NewDefaultFeeTransactionSimple(nonce, &contractAddress, txnOpts.Value,
		txnOpts.GasLimit, input)

	var rawTx *types.Transaction
	rawTx = types.NewTx(baseTx)

	if txnOpts.Signer == nil {
		return errors.New("no signer to authorize the transaction with")
	}

	signTx, err := txnOpts.Signer(txnOpts.From, rawTx)
	if err != nil {
		return err
	}

	signTxBinary, err := signTx.MarshalBinary()
	if err != nil {
		return err
	}

	tx := signTx
	txData := hexutil.Encode(signTxBinary)

	var jsonStr = []byte(`{"txnData" : "` + txData + `"}`)

	request, err := http.NewRequest("POST", WRITE_API_URL+"/api/transactions", bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	fmt.Println("Your request to initial withdrawal has been added to the queue for processing.")
	fmt.Println("The transaction hash for tracking this request is: ", tx.Hash())
	fmt.Println()

	time.Sleep(1000 * time.Millisecond)

	return nil
}

func completeWithdrawal(key *signaturealgorithm.PrivateKey) error {

	client, err := ethclient.Dial(rawURL)
	if err != nil {
		return err
	}

	fromAddress, err := cryptobase.SigAlg.PublicKeyToAddress(&key.PublicKey)

	if err != nil {
		return err
	}

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return err
	}

	contractAddress := common.HexToAddress(staking.STAKING_CONTRACT)
	txnOpts, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(123123))

	if err != nil {
		return err
	}

	txnOpts.From = fromAddress
	txnOpts.Nonce = big.NewInt(int64(nonce))
	txnOpts.GasLimit = uint64(210000)

	val, _ := ParseBigFloat("0")
	txnOpts.Value = etherToWeiFloat(val)

	contract, err := staking.NewStaking(contractAddress, client)
	if err != nil {
		return err
	}

	tx, err := contract.CompleteWithdrawal(txnOpts)
	if err != nil {
		return err
	}

	fmt.Println("Your request to complete withdrawal has been added to the queue for processing.")
	fmt.Println("The transaction hash for tracking this request is: ", tx.Hash())
	fmt.Println()

	time.Sleep(1000 * time.Millisecond)

	return nil
}

func requestCompleteWithdrawal(key *signaturealgorithm.PrivateKey) error {

	fromAddress, err := cryptobase.SigAlg.PublicKeyToAddress(&key.PublicKey)

	if err != nil {
		return err
	}
	_, _, n, err := requestGetBalance(fromAddress.String())
	if err != nil {
		return err
	}

	var nonce uint64
	fmt.Sscan(n, &nonce)

	contractAddress := common.HexToAddress(staking.STAKING_CONTRACT)
	txnOpts, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(123123))

	if err != nil {
		return err
	}

	txnOpts.From = fromAddress
	txnOpts.Nonce = big.NewInt(int64(nonce))
	txnOpts.GasLimit = uint64(210000)

	val, _ := ParseBigFloat("0")
	txnOpts.Value = etherToWeiFloat(val)

	method := staking.GetContract_Method_CompleteWithdrawal()
	abiData, err := staking.GetStakingContract_ABI()
	if err != nil {
		return err
	}

	input, err := abiData.Pack(method)
	if err != nil {
		return err
	}

	baseTx := types.NewDefaultFeeTransactionSimple(nonce, &contractAddress, txnOpts.Value,
		txnOpts.GasLimit, input)

	var rawTx *types.Transaction
	rawTx = types.NewTx(baseTx)

	if txnOpts.Signer == nil {
		return errors.New("no signer to authorize the transaction with")
	}

	signTx, err := txnOpts.Signer(txnOpts.From, rawTx)
	if err != nil {
		return err
	}

	signTxBinary, err := signTx.MarshalBinary()
	if err != nil {
		return err
	}

	tx := signTx
	txData := hexutil.Encode(signTxBinary)

	var jsonStr = []byte(`{"txnData" : "` + txData + `"}`)

	request, err := http.NewRequest("POST", WRITE_API_URL+"/api/transactions", bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	fmt.Println("Your request to complete withdrawal has been added to the queue for processing.")
	fmt.Println("The transaction hash for tracking this request is: ", tx.Hash())
	fmt.Println()

	time.Sleep(1000 * time.Millisecond)

	return nil
}

func getBalanceOfDepositor(dep string) (*big.Int, error) {

	client, err := ethclient.Dial(rawURL)
	if err != nil {
		return nil, err
	}

	contractAddress := common.HexToAddress(staking.STAKING_CONTRACT)
	instance, err := staking.NewStaking(contractAddress, client)
	if err != nil {
		return nil, err
	}

	depositor := common.HexToAddress(dep)
	depositorBalance, err := instance.GetBalanceOfDepositor(nil, depositor)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("StakingBalance", "Address", dep, "coins", weiToEther(depositorBalance).String(), "wei", depositorBalance)

	fmt.Println()

	time.Sleep(1000 * time.Millisecond)

	return depositorBalance, nil
}

func getNetBalanceOfDepositor(dep string) (*big.Int, error) {
	client, err := ethclient.Dial(rawURL)
	if err != nil {
		return nil, err
	}

	contractAddress := common.HexToAddress(staking.STAKING_CONTRACT)
	instance, err := staking.NewStaking(contractAddress, client)
	if err != nil {
		return nil, err
	}

	depositor := common.HexToAddress(dep)
	depositorBalance, err := instance.GetNetBalanceOfDepositor(nil, depositor)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("StakingNetBalance", "Address", dep, "coins", weiToEther(depositorBalance).String(), "wei", depositorBalance)

	fmt.Println()

	time.Sleep(1000 * time.Millisecond)

	return depositorBalance, nil
}

func getDepositorOfValidator(val string) (common.Address, error) {

	client, err := ethclient.Dial(rawURL)
	if err != nil {
		return common.ZERO_ADDRESS, err
	}

	contractAddress := common.HexToAddress(staking.STAKING_CONTRACT)
	instance, err := staking.NewStaking(contractAddress, client)
	if err != nil {
		return common.ZERO_ADDRESS, err
	}

	validator := common.HexToAddress(val)
	depositor, err := instance.GetDepositorOfValidator(nil, validator)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Depositor", depositor, "validator", validator)

	fmt.Println()

	time.Sleep(1000 * time.Millisecond)

	return depositor, err
}

func getDepositorBlockRewards(dep string) (*big.Int, error) {

	client, err := ethclient.Dial(rawURL)
	if err != nil {
		return nil, err
	}

	contractAddress := common.HexToAddress(staking.STAKING_CONTRACT)
	instance, err := staking.NewStaking(contractAddress, client)
	if err != nil {
		return nil, err
	}

	depositor := common.HexToAddress(dep)
	depositorBalance, err := instance.GetDepositorRewards(nil, depositor)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("BlockRewards", "Depositor", dep, "coins", weiToEther(depositorBalance).String(), "wei", depositorBalance)

	fmt.Println()

	time.Sleep(1000 * time.Millisecond)

	return depositorBalance, nil
}

type ValidatorDetails struct {
	depositor    common.Address
	validator    common.Address
	balance      *big.Int
	netBalance   *big.Int
	blockRewards *big.Int
}

func listValidators() error {
	if len(rawURL) == 0 {
		return errors.New("DP_RAW_URL environment variable not specified")
	}

	client, err := ethclient.Dial(rawURL)
	if err != nil {
		return err
	}

	contractAddress := common.HexToAddress(staking.STAKING_CONTRACT)
	instance, err := staking.NewStaking(contractAddress, client)
	if err != nil {
		return err
	}

	validatorList, err := instance.ListValidators(nil)
	if err != nil {
		log.Fatal(err)
	}

	totalDepositedBalance := big.NewInt(int64(0))
	ValidatorDetailsList := make([]*ValidatorDetails, len(validatorList))
	for i := 0; i < len(validatorList); i++ {
		depositor, err := getDepositorOfValidator(validatorList[i].String())
		if err != nil {
			return err
		}

		balanceVal, err := getBalanceOfDepositor(depositor.String())
		if err != nil {
			return err
		}

		netBalance, err := getNetBalanceOfDepositor(depositor.String())
		if err != nil {
			return err
		}

		blockrewards, err := getDepositorBlockRewards(depositor.String())
		if err != nil {
			return err
		}

		ValidatorDetailsList[i] = &ValidatorDetails{
			depositor:    depositor,
			validator:    validatorList[i],
			balance:      balanceVal,
			netBalance:   netBalance,
			blockRewards: blockrewards,
		}

		totalDepositedBalance = totalDepositedBalance.Add(totalDepositedBalance, balanceVal)
	}

	for i := 0; i < len(ValidatorDetailsList); i++ {
		validatorDetails := ValidatorDetailsList[i]
		fmt.Println("Depositor", validatorDetails.depositor, "Validator", validatorDetails.validator, "balance coins", weiToEther(validatorDetails.balance).String(),
			"netBalance coins", weiToEther(validatorDetails.netBalance).String(), "blockrewards coins", weiToEther(validatorDetails.blockRewards).String())
	}

	fmt.Println("Total validators", len(validatorList), "totalDepositedBalance", weiToEther(totalDepositedBalance).String())

	fmt.Println()

	time.Sleep(1000 * time.Millisecond)

	return nil
}
