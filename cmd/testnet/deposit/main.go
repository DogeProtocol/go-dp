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
	"github.com/DogeProtocol/dp/systemcontracts/staking"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"strings"
	"time"
)

var rawURL = os.Getenv("DP_RAW_URL")
var contractAddress = os.Getenv("DP_STAKING_CONTRACT_ADDRESS")
var depositorAddress = os.Getenv("DP_STAKING_DEPOSITER")
var validatorAddress = os.Getenv("DP_STAKING_VALIDATOR")
var depositorPassword = os.Getenv("DP_STAKING_DEPOSITER_PASS")
var validatorPassword = os.Getenv("DP_STAKING_VALIDATOR_PASS")
var depositAmount = os.Getenv("DP_STAKING_DEPOSIT_AMOUNT")

var depositorPath = os.Getenv("DP_DEPOSITER_PATH")
var validatorPath = os.Getenv("DP_VALIDATOR_PATH")
var e = "Error occurred. Please ensure that geth is running, is connected to the blockchain and all required environment variables have been set correctly."

type KeyStore struct {
	Handle *keystore.KeyStore
}

func main() {
	if len(rawURL) < 5 {
		log.Println(e + " DP_RAW_URL")
		return
	}
	if len(contractAddress) < 20 {
		log.Println(e + " DP_STAKING_CONTRACT_ADDRESS")
		return
	}
	if len(depositorAddress) < 20 {
		log.Println(e + " DP_STAKING_DEPOSITER")
		return
	}
	if len(validatorAddress) < 20 {
		log.Println(e + " DP_STAKING_VALIDATOR")
		return
	}
	if len(depositAmount) <= 0 {
		log.Println(e + " DP_STAKING_DEPOSIT_AMOUNT")
		return
	}
	fmt.Println("Deposit ...")
	deposit(contractAddress, depositorAddress, common.HexToAddress(validatorAddress),
		depositorPassword, validatorPassword, depositAmount)
}

func deposit(contractAddress string,
	depositorAddress string, validatorAddress common.Address,
	depositorPassword string, validatorPassword string, depositAmount string) {

	var depositorKey *keystore.Key
	var validatorKey *keystore.Key

	depositorSecretKey, err := ReadDataFile(depositorPath)
	if err != nil {
		fmt.Println("Depositor Error", err)
		return
	}
	depositorKey, err = keystore.DecryptKey(depositorSecretKey, depositorPassword)
	if err != nil {
		fmt.Println("Depositor Error", err)
		return
	}
	if depositorKey == nil {
		log.Println(e + " DP_STAKING_DEPOSITER_PASS")
		return
	}

	if depositorKey == nil {
		log.Println(e + " DP_STAKING_DEPOSITER DP_DEPOSITER_PATH")
		return
	}

	validatorSecretKey, err := ReadDataFile(validatorPath)
	if err != nil {
		fmt.Println("Depositor Error", err)
		return
	}

	validatorKey, err = keystore.DecryptKey(validatorSecretKey, validatorPassword)
	if err != nil {
		fmt.Println("Depositor Error", err)
		return
	}

	if validatorKey == nil {
		log.Println(e + " DP_STAKING_VALIDATOR_PASS")
		return
	}
	_, err = cryptobase.SigAlg.SerializePublicKey(&validatorKey.PrivateKey.PublicKey)
	if err != nil {
		log.Println("validator SerializePublicKey", err)
		return
	}
	valAddressFromKey, err := cryptobase.SigAlg.PublicKeyToAddress(&validatorKey.PrivateKey.PublicKey)
	if err != nil {
		log.Println("validator PublicKeyToAddress", err)
		return
	}

	if valAddressFromKey != validatorAddress {
		log.Println("validator key address check failed", err)
		return
	}

	if len(depositorKey.PrivateKey.PublicKey.PubData) >= cryptobase.SigAlg.PublicKeyLength() {
		priBytes, err := cryptobase.SigAlg.SerializePrivateKey(depositorKey.PrivateKey)
		if err != nil {
			panic(err)
		}

		if depositorKey != nil && len(priBytes) >= cryptobase.SigAlg.PrivateKeyLength() {
			tx, err := depositContract(depositorAddress, contractAddress, depositorKey.PrivateKey.PublicKey.PubData,
				depositorKey.PrivateKey, depositAmount, validatorAddress)
			if err != nil {
				log.Println("Error occurred." + tx + " : " + err.Error())
				return
			}
			fmt.Println("Tx hash: ", tx)
			fmt.Println(" Successfully deposited ...")
		} else {
			log.Println(e + " DP_DEPOSITER_PATH")
		}
	} else {
		log.Println(e + " DP_VALIDATOR_PATH")
	}
}

func depositContract(fromaddress string, contractaddress string, validatorPubKey []byte,
	key *signaturealgorithm.PrivateKey, depositAmount string, validatorAddress common.Address) (string, error) {

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

	auth.GasLimit = uint64(21000)
	//gasPrice, err := client.SuggestGasPrice(context.Background())
	//if err != nil {
	//	return "4", err
	//}
	//auth.GasPrice = big.NewInt(1000000)
	//fmt.Println("gasPrice", gasPrice)

	contract, err := staking.NewStaking(contractAddress, client)
	if err != nil {
		return "5", err
	}

	tx, err := contract.NewDeposit(auth, validatorAddress)
	if err != nil {
		return "NewDeposit failed", err
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
