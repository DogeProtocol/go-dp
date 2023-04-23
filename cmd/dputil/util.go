package main

import (
	"context"
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

func send(from string, to string, quantity string) (string, error) {
	secretKey, err := ReadDataFile(os.Getenv("GETH_KEY_FILE"))
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
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}

	v, err := ParseBigFloat(quantity)
	if err != nil {
		return "", err
	}

	value := etherToWeiFloat(v)

	var data []byte
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), key.PrivateKey)
	if err != nil {
		return "", err
	}
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
