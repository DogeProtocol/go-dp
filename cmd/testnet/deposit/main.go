package main

import (
	"context"
	"fmt"
	"github.com/DogeProtocol/dp/accounts"
	"github.com/DogeProtocol/dp/accounts/abi/bind"
	"github.com/DogeProtocol/dp/accounts/keystore"
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/crypto/cryptobase"
	"github.com/DogeProtocol/dp/crypto/signaturealgorithm"
	"github.com/DogeProtocol/dp/ethclient"
	"github.com/DogeProtocol/dp/params"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"strings"
	"time"
)

var rawURL = os.Getenv("GETH_RAW_URL")
var contractAddress = os.Getenv("GETH_STAKING_CONTRACT_ADDRESS")
var depositorAddress = os.Getenv("GETH_STAKING_DEPOSITER")
var validatorAddress = os.Getenv("GETH_STAKING_VALIDATOR")
var depositorPassword = os.Getenv("GETH_STAKING_DEPOSITER_PASS")
var validatorPassword = os.Getenv("GETH_STAKING_VALIDATOR_PASS")
var depositAmount = os.Getenv("GETH_STAKING_DEPOSIT_AMOUNT")

var depositorPath = os.Getenv("GETH_DEPOSITER_PATH")
var validatorPath = os.Getenv("GETH_VALIDATOR_PATH")
var e = "Error occurred. Please ensure that geth is running, is connected to the blockchain and all required environment variables have been set correctly."

type KeyStore struct {
	Handle *keystore.KeyStore
}

func main() {
	if len(rawURL) < 5 {
		log.Println(e + " GETH_RAW_URL")
		return
	}
	if len(contractAddress) < 20 {
		log.Println(e + " GETH_STAKING_CONTRACT_ADDRESS")
		return
	}
	if len(depositorAddress) < 20 {
		log.Println(e + " GETH_STAKING_DEPOSITER")
		return
	}
	if len(validatorAddress) < 20 {
		log.Println(e + " GETH_STAKING_VALIDATOR")
		return
	}
	if len(depositAmount) <= 0 {
		log.Println(e + " GETH_STAKING_DEPOSIT_AMOUNT")
		return
	}
	fmt.Println("Deposit ...")
	deposit(contractAddress, depositorAddress, validatorAddress,
		depositorPassword, validatorPassword, depositAmount)
}

func deposit(contractAddress string,
	depositorAddress string, validatorAddress string,
	depositorPassword string, validatorPassword string, depositAmount string) {

	var depositorSecretKey []byte
	var depositorKey *keystore.Key
	var validatorSecretKey []byte
	var validatorKey *keystore.Key

	var pubKey []byte

	ks := SetUpKeyStore("./data/keystore")
	if len(depositorPath) > 0 {
		path := strings.ReplaceAll(depositorPath, "\\", "/")
		ks = SetUpKeyStore(path)
	}
	accounts := ks.GetAllKeys()
	for _, account := range accounts {
		acc := account.Address.String()
		d := strings.EqualFold(depositorAddress, acc)
		if d == true {
			depositorSecretKey, _ = ReadDataFile(account.URL.Path)
			depositorKey, _ = keystore.DecryptKey(depositorSecretKey, depositorPassword)
			if depositorKey == nil {
				log.Println(e + " GETH_STAKING_DEPOSITER_PASS")
				return
			}
		}
	}

	if depositorKey == nil {
		log.Println(e + " GETH_STAKING_DEPOSITER GETH_DEPOSITER_PATH")
		return
	}

	ks = SetUpKeyStore("./data/keystore")
	if len(validatorPath) > 0 {
		path := strings.ReplaceAll(validatorPath, "\\", "/")
		ks = SetUpKeyStore(path)
	}
	accounts = ks.GetAllKeys()
	for _, account := range accounts {
		acc := account.Address.String()
		v := strings.EqualFold(validatorAddress, acc)
		if v == true {
			validatorSecretKey, _ = ReadDataFile(account.URL.Path)
			validatorKey, _ = keystore.DecryptKey(validatorSecretKey, validatorPassword)
			if validatorKey == nil {
				log.Println(e + " GETH_STAKING_VALIDATOR_PASS")
				return
			}
			_, err := cryptobase.SigAlg.SerializePublicKey(&validatorKey.PrivateKey.PublicKey)
			if err != nil {
				panic(err)
			}
		}
	}

	if validatorKey == nil {
		log.Println(e + " GETH_STAKING_VALIDATOR  GETH_VALIDATOR_PATH")
		return
	}

	if len(pubKey) >= cryptobase.SigAlg.PublicKeyLength() {
		priBytes, err := cryptobase.SigAlg.SerializePrivateKey(depositorKey.PrivateKey)
		if err != nil {
			panic(err)
		}

		if depositorKey != nil && len(priBytes) >= cryptobase.SigAlg.PrivateKeyLength() {
			tx, err := depositContract(depositorAddress, contractAddress, pubKey,
				depositorKey.PrivateKey, depositAmount)
			if err != nil {
				log.Println("Error occurred." + tx + " : " + err.Error())
				return
			}
			fmt.Println("Tx hash: ", tx)
			fmt.Println(" Successfully deposited ...")
		} else {
			log.Println(e + " GETH_DEPOSITER_PATH")
		}
	} else {
		log.Println(e + " GETH_VALIDATOR_PATH")
	}
}

func depositContract(fromaddress string, contractaddress string, pubKey []byte,
	key *signaturealgorithm.PrivateKey, depositAmount string) (string, error) {

	client, err := ethclient.Dial(rawURL)
	if err != nil {
		//log.Println(err.Error())
		return "0", err
	}

	fromAddress := common.HexToAddress(fromaddress)
	contractAddress := common.HexToAddress(contractaddress)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		//log.Println(err.Error())
		return "1", err
	}
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		//log.Println(err.Error())
		return "2", err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		//log.Println(err.Error())
		return "3", err
	}
	auth.From = fromAddress
	auth.Nonce = big.NewInt(int64(nonce))
	p, _ := ParseBigFloat(depositAmount) //ether
	auth.Value = etherToWei(p)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		//log.Println(err.Error())
		return "4", err
	}
	auth.GasPrice = gasPrice
	contract, err := NewStakingContractAddress1(contractAddress, client)
	if err != nil {
		//log.Println(err.Error())
		return "5", err
	}

	tx, err := contract.NewDeposit(auth, pubKey)
	if err != nil {
		//log.Println(err.Error())
		return "6", err
	}

	// Don't even wait, check its presence in the local pending state
	time.Sleep(250 * time.Millisecond) // Allow it to be processed by the local node :P

	return tx.Hash().Hex(), nil
}

func etherToWei(eth *big.Float) *big.Int {
	truncInt, _ := eth.Int(nil)
	truncInt = new(big.Int).Mul(truncInt, big.NewInt(params.Ether))
	fracStr := strings.Split(fmt.Sprintf("%.18f", eth), ".")[1]
	fracStr += strings.Repeat("0", 18-len(fracStr))
	fracInt, _ := new(big.Int).SetString(fracStr, 10)
	wei := new(big.Int).Add(truncInt, fracInt)
	return wei
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

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	return byteValue, nil
}

func SetUpKeyStore(kp string) *KeyStore {
	ks := &KeyStore{}
	ks.Handle = keystore.NewKeyStore(kp, keystore.LightScryptN, keystore.LightScryptP)
	return ks
}

func (ks *KeyStore) GetAllKeys() []accounts.Account {
	return ks.Handle.Accounts()
}
