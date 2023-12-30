package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DogeProtocol/dp/accounts/abi/bind"
	"github.com/DogeProtocol/dp/accounts/keystore"
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/crypto/crosssign"
	"github.com/DogeProtocol/dp/crypto/cryptobase"
	"github.com/DogeProtocol/dp/crypto/signaturealgorithm"
	"github.com/DogeProtocol/dp/ethclient"
	"github.com/DogeProtocol/dp/systemcontracts/conversion"
	"io/ioutil"
	"math/big"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

func printHelp() {
	fmt.Println("===========")
	fmt.Println("dputil genesis-sign ETH_ADDRESS DEPOSITOR_QUANTUM_ADDRESS VALIDATOR_QUANTUM_ADDRESS AMOUNT")
	fmt.Println("      Set the following environment variables:")
	fmt.Println("           DP_KEY_FILE_DIR, DP_DEPOSITOR_ACC_PWD, DP_VALIDATOR_ACC_PWD")
	fmt.Println("===========")
	fmt.Println("dputil genesis-verify JSON_FILE_NAME")
	fmt.Println("===========")
	fmt.Println("dputil getconversionmessage ETH_ADDRESS")
	fmt.Println("      Set the following environment variables:")
	fmt.Println("           DP_KEY_FILE, DP_ACC_PWD")
	fmt.Println("===========")
	fmt.Println("dputil getcoinsfortokens ETH_ADDRESS")
	fmt.Println("      Set the following environment variables:")
	fmt.Println("           DP_KEY_FILE, DP_ACC_PWD")
	fmt.Println("===========")
	fmt.Println("dputil balance ACCOUNT_ADDRESS")
	fmt.Println("      Set the following environment variables:")
	fmt.Println("           DP_RAW_URL")
	fmt.Println("===========")
	fmt.Println("dputil send FROM_ADDRESS TO_ADDRESS QUANTITY")
	fmt.Println("===========")
}

var rawURL string
var wg sync.WaitGroup

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}
	rawURL = os.Getenv("DP_RAW_URL")

	if len(rawURL) == 0 {
		os := runtime.GOOS
		if os == "windows" {
			rawURL = "//./pipe/geth.ipc"
		} else {
			rawURL = "~/.ethereum/geth.ipc"
		}
	}

	if os.Args[1] == "balance" {
		balance()
	} else if os.Args[1] == "send" {
		sendTxn()
	} else if os.Args[1] == "txn" {
		getTxn()
	} else if os.Args[1] == "genesis-sign" {
		GenesisSign()
	} else if os.Args[1] == "genesis-verify" {
		GenesisVerify()
	} else if os.Args[1] == "getconversionmessage" {
		err := GetConversionMessage()
		if err != nil {
			fmt.Println("Error", err)
		}
	} else if os.Args[1] == "getcoinsfortokens" {
		err := ConvertToCoins()
		if err != nil {
			fmt.Println("Error", err)
		}
	} else {
		printHelp()
	}
}

func GenesisSign() {
	if len(os.Args) < 6 {
		printHelp()
		return
	}
	if len(os.Getenv("DP_KEY_FILE_DIR")) == 0 {
		fmt.Println("Set the keyfile directory environment variable DP_KEY_FILE_DIR")
		return
	}
	if len(os.Getenv("DP_DEPOSITOR_ACC_PWD")) == 0 {
		fmt.Println("Set the depositor password environment variable DP_DEPOSITOR_ACC_PWD")
		return
	}
	if len(os.Getenv("DP_VALIDATOR_ACC_PWD")) == 0 {
		fmt.Println("Set the validator password environment variable DP_VALIDATOR_ACC_PWD")
		return
	}

	ethAddr := os.Args[2]
	depositorAddr := os.Args[3]
	validatorAddr := os.Args[4]
	amount := os.Args[5]

	if common.IsLegacyEthereumHexAddress(ethAddr) == false {
		fmt.Println("Invalid eth address", ethAddr)
		return
	}

	if common.IsHexAddress(depositorAddr) == false {
		fmt.Println("Invalid depositor address", depositorAddr)
		return
	}

	if common.IsHexAddress(validatorAddr) == false {
		fmt.Println("Invalid validator address", validatorAddr)
		return
	}

	_, err := ParseBigFloat(amount)
	if err != nil {
		fmt.Println(err)
		return
	}

	depositorKeyFile, err := findKeyFile(depositorAddr)
	if err != nil {
		fmt.Println("Error finding DEPOSITOR_ADDRESS in DP_KEY_FILE_DIR", err)
		return
	}
	depositorKey, err := ReadDataFile(depositorKeyFile)
	if err != nil {
		fmt.Println("Error loading depositor key file", err)
		return
	}
	depPassword := os.Getenv("DP_DEPOSITOR_ACC_PWD")
	depKey, err := keystore.DecryptKey(depositorKey, depPassword)
	if err != nil {
		fmt.Println("Error decrypting depositor key using DP_DEPOSITOR_ACC_PWD", err)
		return
	}

	validatorKeyFile, err := findKeyFile(validatorAddr)
	if err != nil {
		fmt.Println("Error finding VALIDATOR_ADDRESS in DP_KEY_FILE_DIR", err)
		return
	}
	validatorKey, err := ReadDataFile(validatorKeyFile)
	if err != nil {
		fmt.Println("Error loading validator key file", err)
		return
	}
	valPassword := os.Getenv("DP_VALIDATOR_ACC_PWD")
	valKey, err := keystore.DecryptKey(validatorKey, valPassword)
	if err != nil {
		fmt.Println("Error decrypting depositor key using DP_VALIDATOR_ACC_PWD", err)
		return
	}

	details, err := crosssign.SignGenesis(depKey.PrivateKey, valKey.PrivateKey, ethAddr, amount)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Signed the genesis validator message!")

	marshalled, err := json.Marshal(details)
	if err != nil {
		fmt.Println(err)
		return
	}

	fileName := "cross-sign-" + depositorAddr + ".json"
	err = ioutil.WriteFile(fileName, marshalled, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Successfully created cross-sign file", fileName)

	return
}

func GenesisVerify() {
	if len(os.Args) < 3 {
		printHelp()
		return
	}

	jsonFile := os.Args[2]

	jsonString, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		fmt.Println("error opening json file", jsonFile, err)
		return
	}

	jsonBytes := []byte(jsonString)

	details := crosssign.GenesisCrossSignDetails{}
	err = json.Unmarshal(jsonBytes, &details)
	if err != nil {
		fmt.Println("error reading json", jsonFile, err)
		return
	}

	_, err = crosssign.VerifyGenesis(&details)
	if err != nil {
		fmt.Println("verify failed", err)
		return
	}

	fmt.Println("Verify succeeded!")
}

func balance() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	addr := os.Args[2]

	if common.IsHexAddress(addr) == false {
		fmt.Println("Invalid address", addr)
		return
	}

	ethBalance, weiBalance, err := getBalance(addr)
	if err != nil {
		fmt.Println("Error", err)
	}
	fmt.Println("Address", addr, "eth", ethBalance, "wei", weiBalance)

}

type Txn struct {
	FromAddress string
	ToAddress   string
	Quantity    string
	Count       int
}

func sendTxn() {
	if len(os.Args) < 5 {
		printHelp()
		return
	}

	from := os.Args[2]
	to := os.Args[3]
	quantity := os.Args[4]

	if common.IsHexAddress(from) == false {
		fmt.Println("Invalid address", from)
		return
	}

	if common.IsHexAddress(to) == false {
		fmt.Println("Invalid address", to)
		return
	}

	flt, err := ParseBigFloat(quantity)
	if err != nil {
		fmt.Println(err)
		return
	}

	wei := etherToWeiFloat(flt)
	ether := weiToEther(wei)

	fmt.Println("Send", "from", from, "to", to, "quantity", quantity, "ether", ether)

	txHash, err := send(from, to, quantity)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("TxnHash", txHash)
}

func getTxn() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	hash := os.Args[2]

	txnJson, err := GetTransaction(hash)
	if err != nil {
		fmt.Println("GetTransaction Error", err)
		return
	}
	json, err := Prettify(txnJson)
	if err != nil {
		fmt.Println(txnJson)
		fmt.Println(err)
	}
	fmt.Println(json)
}

func Prettify(str string) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}

func GetConversionMessage() error {
	if len(os.Args) < 3 {
		printHelp()
		return errors.New("incorrect usage")
	}

	ethAddress := os.Args[2]
	if common.IsLegacyEthereumHexAddress(ethAddress) == false {
		return errors.New("invalid EthAddress")
	}

	keyFile := os.Getenv("DP_KEY_FILE")
	if len(keyFile) == 0 {
		return errors.New("DP_KEY_FILE environment variable is not set")
	}

	accPwd := os.Getenv("DP_ACC_PWD")
	if len(accPwd) == 0 {
		return errors.New("DP_ACC_PWD environment variable is not set")
	}

	key, err := GetKeyFromFile(keyFile)
	if err != nil {
		return err
	}

	qAddr, err := cryptobase.SigAlg.PublicKeyToAddress(&key.PublicKey)
	if err != nil {
		return err
	}

	quantumAddress := qAddr.Hex()

	message := strings.Replace(crosssign.ConversionMessageTemplate, "[ETH_ADDRESS]", strings.ToLower(ethAddress), 1)
	message = strings.Replace(message, "[QUANTUM_ADDRESS]", strings.ToLower(quantumAddress), 1)

	fmt.Println("Message is: ")
	fmt.Println(message)

	return nil
}

func ConvertToCoins() error {
	if len(os.Args) < 4 {
		printHelp()
		return errors.New("incorrect usage")
	}

	ethAddress := os.Args[2]
	if common.IsLegacyEthereumHexAddress(ethAddress) == false {
		return errors.New("invalid EthAddress")
	}

	ethSignature := os.Args[3]
	fmt.Println("ethSignature len", len(ethSignature))

	keyFile := os.Getenv("DP_KEY_FILE")
	if len(keyFile) == 0 {
		return errors.New("DP_KEY_FILE environment variable is not set")
	}

	accPwd := os.Getenv("DP_ACC_PWD")
	if len(accPwd) == 0 {
		return errors.New("DP_ACC_PWD environment variable is not set")
	}

	key, err := GetKeyFromFile(keyFile)
	if err != nil {
		return err
	}

	return ConvertCoins(ethAddress, ethSignature, key)
}

func ConvertCoins(ethAddress string, ethSignature string, key *signaturealgorithm.PrivateKey) error {
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

	fmt.Println("Txn Hash", tx.Hash())

	//fmt.Println("data", common.Bytes2Hex(tx.Data()), tx.Data(), len(tx.Data()))

	time.Sleep(1000 * time.Millisecond)

	return nil
}
