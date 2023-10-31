package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/DogeProtocol/dp/accounts"
	"github.com/DogeProtocol/dp/accounts/keystore"
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/core/types"
	"github.com/DogeProtocol/dp/ethclient"
	"github.com/DogeProtocol/dp/params"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strings"
)

type KeyStore struct {
	Handle *keystore.KeyStore
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

func findAllAddresses() ([]string, error) {
	keyfileDir := os.Getenv("GETH_KEY_FILE_DIR")
	if len(keyfileDir) == 0 {
		return nil, errors.New("Both GETH_KEY_FILE and GETH_KEY_FILE_DIR environment variables not set")
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

func findKeyFile(from string) (string, error) {
	keyfile := os.Getenv("GETH_KEY_FILE")
	if len(keyfile) > 0 {
		return keyfile, nil
	}

	keyfileDir := os.Getenv("GETH_KEY_FILE_DIR")
	if len(keyfileDir) == 0 {
		return "", errors.New("Both GETH_KEY_FILE and GETH_KEY_FILE_DIR environment variables not set")
	}

	files, err := ioutil.ReadDir(keyfileDir)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	fromAddress := strings.ToLower(strings.Replace(from, "0x", "", 0))
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.Contains(strings.ToLower(file.Name()), fromAddress) {
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

func GetConnectionContext(from string) (*ConnectionContext, error) {
	keyFile, err := findKeyFile(from)
	if err != nil {
		return nil, err
	}

	secretKey, err := ReadDataFile(keyFile)
	if err != nil {
		return nil, err
	}

	password := os.Getenv("GETH_ACC_PWD")
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
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), connectionContext.Key.PrivateKey)
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
	password := os.Getenv("GETH_ACC_PWD")
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
	//gasPrice, err := client.SuggestGasPrice(context.Background())
	//if err != nil {
	//	return "", err
	//}

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
	dataDir := os.Getenv("GETH_DATA_PATH")
	if dataDir == "" {
		dataDir = "data"
	}

	ks := &KeyStore{}
	ks.Handle = keystore.NewKeyStore(dataDir, keystore.LightScryptN, keystore.LightScryptP)
	return ks
}
