package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DogeProtocol/dp/accounts/keystore"
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/console/prompt"
	"github.com/DogeProtocol/dp/conversionutil"
	"github.com/DogeProtocol/dp/crypto/crosssign"
	"github.com/DogeProtocol/dp/crypto/cryptobase"
	"github.com/DogeProtocol/dp/log"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"
)

const READ_API_URL = ""  //"https://scan.dpapi.org"
const WRITE_API_URL = "" //"https://txn.dpapi.org"

func printHelp() {
	fmt.Println("===========")
	fmt.Println(" dputil ")
	fmt.Println("      Set a default environment variables:")
	fmt.Println("           DP_RAW_URL")
	fmt.Println("dputil genesis-sign ETH_ADDRESS DEPOSITOR_QUANTUM_ADDRESS VALIDATOR_QUANTUM_ADDRESS AMOUNT")
	fmt.Println("      Set the following environment variables:")
	fmt.Println("           DP_KEY_FILE_DIR, DP_DEPOSITOR_ACC_PWD, DP_VALIDATOR_ACC_PWD")
	fmt.Println("===========")
	fmt.Println("dputil genesis-verify JSON_FILE_NAME")
	fmt.Println("===========")
	fmt.Println("dputil getconversionmessage ETH_ADDRESS")
	fmt.Println("      Set the following environment variables:")
	fmt.Println("           DP_KEY_FILE")
	fmt.Println("===========")
	fmt.Println("dputil getcoinsfortokens ETH_ADDRESS ETH_SIGNATURE")
	fmt.Println("      Set the following environment variables:")
	fmt.Println("           DP_KEY_FILE")
	fmt.Println("===========")
	fmt.Println("dputil balance ACCOUNT_ADDRESS")
	fmt.Println("      Set the following environment variables:")
	fmt.Println("           DP_RAW_URL")
	fmt.Println("===========")
	fmt.Println("dputil newdeposit DEPOSITOR_ADDRESS VALIDATOR_ADDRESS DEPOSITOR_AMOUNT")
	fmt.Println("      Set the following environment variables:")
	fmt.Println("           DP_KEY_FILE")
	fmt.Println("===========")
	fmt.Println("dputil initiatewithdrawal DEPOSITOR_ADDRESS")
	fmt.Println("      Set the following environment variables:")
	fmt.Println("           DP_KEY_FILE")
	fmt.Println("===========")
	fmt.Println("dputil completewithdrawal DEPOSITOR_ADDRESS")
	fmt.Println("      Set the following environment variables:")
	fmt.Println("           DP_KEY_FILE")
	fmt.Println("===========")
	fmt.Println("dputil stakingbalance DEPOSITOR_ADDRESS")
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
	/*
		if len(rawURL) == 0 {
			os := runtime.GOOS
			if os == "windows" {
				rawURL = "//./pipe/geth.ipc"
			} else {
				rawURL = "~/.ethereum/geth.ipc"
			}
		}
	*/
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
	} else if os.Args[1] == "newdeposit" {
		err := Deposit()
		if err != nil {
			fmt.Println("Error", err)
		}
	} else if os.Args[1] == "initiatewithdrawal" {
		err := InitiateWithdrawal()
		if err != nil {
			fmt.Println("Error", err)
		}
	} else if os.Args[1] == "completewithdrawal" {
		err := CompleteWithdrawal()
		if err != nil {
			fmt.Println("Error", err)
		}
	} else if os.Args[1] == "stakingbalance" {
		err := DepositorBalance()
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

	if len(rawURL) == 0 {
		ethBalance, weiBalance, nonce, err := requestGetBalance(addr)
		if err != nil {
			fmt.Println("Error", err)
		}
		fmt.Println("Address", addr, "coins", ethBalance, "wei", weiBalance, "nonce", nonce)
	} else {
		ethBalance, weiBalance, err := getBalance(addr)
		if err != nil {
			fmt.Println("Error", err)
		}
		fmt.Println("Address", addr, "coins", ethBalance, "wei", weiBalance)
	}
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

	fmt.Println(fmt.Sprintf("Quantum wallet address %s", keyFile))
	accPwd, err := prompt.Stdin.PromptPassword(fmt.Sprintf("Enter the quantum wallet password : "))
	if err != nil {
		return err
	}
	if len(accPwd) == 0 {
		return errors.New("password is not set")
	}

	key, err := GetKeyFromFile(keyFile, accPwd)
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

	_, ok := conversionutil.SnapshotMap[strings.ToLower(ethAddress)]

	if ok == false {
		log.Trace("IsGasExemptTxn address not in snapshot", "ethAddress", ethAddress)
		return errors.New("unidentified eth address")
	}

	ethConfirm, err := prompt.Stdin.PromptConfirm(fmt.Sprintf("Do you confirm that your ETH ADDRESS having the Dogep tokens is %s ?", ethAddress))
	if err != nil {
		return err
	}
	if ethConfirm != true {
		return errors.New("confirmation not made")
	}
	fmt.Println()

	ethSignature := os.Args[3]

	keyFile := os.Getenv("DP_KEY_FILE")
	if len(keyFile) == 0 {
		return errors.New("DP_KEY_FILE environment variable is not set")
	}

	fmt.Println(fmt.Sprintf("Quantum wallet addres %s", keyFile))
	accPwd, err := prompt.Stdin.PromptPassword(fmt.Sprintf("Enter the quantum wallet password : "))
	if err != nil {
		return err
	}
	if len(accPwd) == 0 {
		return errors.New("password is not set")
	}
	fmt.Println()

	backupConfirm, err := prompt.Stdin.PromptConfirm(fmt.Sprintf("Do you confirm that you have backed up your quantum wallet located at %s ?", keyFile))
	if err != nil {
		return err
	}
	if backupConfirm != true {
		return errors.New("confirmation not made")
	}
	fmt.Println()

	passwordConfirm, err := prompt.Stdin.PromptConfirm(fmt.Sprintf("Do you understand that the wallet password will always be required to use the quantum wallet at %s?", keyFile))
	if err != nil {
		return err
	}
	if passwordConfirm != true {
		return errors.New("confirmation not made")
	}
	fmt.Println()

	key, err := GetKeyFromFile(keyFile, accPwd)
	if err != nil {
		return err
	}

	qAddr, err := cryptobase.SigAlg.PublicKeyToAddress(&key.PublicKey)
	if err != nil {
		return err
	}

	quantumAddress := qAddr.Hex()

	time.Sleep(500 * time.Millisecond)

	fmt.Println()
	quantumConfirm, err := prompt.Stdin.PromptConfirm(fmt.Sprintf("Do you confirm that you want the coins deposited to QUANTUM ADDRESS %s ?", quantumAddress))
	if err != nil {
		return err
	}
	if quantumConfirm != true {
		return errors.New("confirmation not made")
	}
	fmt.Println()

	crossSignDetails := &crosssign.ConversionSignDetails{
		EthAddress:        strings.ToLower(ethAddress),
		EthereumSignature: ethSignature,
		QuantumAddress:    strings.ToLower(quantumAddress),
	}

	_, err = crosssign.VerifyConversion(crossSignDetails)
	if err != nil {
		fmt.Println("An error occurred while verifying the ethereum signature.")
		return err
	}

	time.Sleep(3000 * time.Millisecond)
	fmt.Println("Final conformation!!!")
	time.Sleep(3000 * time.Millisecond)
	fmt.Println("Verify your message...")
	time.Sleep(3000 * time.Millisecond)

	message := strings.Replace(crosssign.ConversionMessageTemplate, "[ETH_ADDRESS]", strings.ToLower(ethAddress), 1)
	message = strings.Replace(message, "[QUANTUM_ADDRESS]", strings.ToLower(quantumAddress), 1)

	finalConfirm, err := prompt.Stdin.PromptConfirm(fmt.Sprintf("%s", message))
	if err != nil {
		return err
	}
	if finalConfirm != true {
		return errors.New("confirmation not made")
	}

	if len(rawURL) == 0 {
		return requestConvertCoins(ethAddress, ethSignature, key)
	} else {
		return convertCoins(ethAddress, ethSignature, key)
	}
}

func Deposit() error {
	if len(os.Args) < 5 {
		printHelp()
		return errors.New("incorrect usage")
	}

	if len(os.Getenv("DP_KEY_FILE_DIR")) == 0 {
		return errors.New("set the keyfile directory environment variable DP_KEY_FILE_DIR")
	}

	depositorAddr := os.Args[2]
	validatorAddr := os.Args[3]
	depositorAmount := os.Args[4]

	if common.IsHexAddress(depositorAddr) == false {
		return errors.New("invalid depositor address " + depositorAddr)
	}

	if common.IsHexAddress(validatorAddr) == false {
		return errors.New("invalid validator address " + validatorAddr)
	}

	_, err := ParseBigFloat(depositorAmount)
	if err != nil {
		return err
	}

	depositorKeyFile, err := findKeyFile(depositorAddr)
	if err != nil {
		return errors.New("error finding DEPOSITOR_ADDRESS in DP_KEY_FILE_DIR " + err.Error())
	}

	fmt.Println(fmt.Sprintf("Depositor wallet addres %s", depositorKeyFile))
	depositorPwd, err := prompt.Stdin.PromptPassword(fmt.Sprintf("Enter the depositor wallet password : "))
	if err != nil {
		return err
	}
	if len(depositorPwd) == 0 {
		return errors.New("depositor password is not set")
	}

	fmt.Println()

	depositorPasswordConfirm, err := prompt.Stdin.PromptConfirm(fmt.Sprintf("Do you understand that the depositor password will always be required to use the quantum depositor wallet at %s?", depositorKeyFile))
	if err != nil {
		return err
	}
	if depositorPasswordConfirm != true {
		return errors.New("confirmation not made")
	}
	fmt.Println()

	depKey, err := GetKeyFromFile(depositorKeyFile, depositorPwd)
	if err != nil {
		return errors.New("error decrypting depositor key " + err.Error())
	}

	depAddressFromKey, err := cryptobase.SigAlg.PublicKeyToAddress(&depKey.PublicKey)
	if err != nil {
		return errors.New("depositor public key to address " + err.Error())
	}

	if !depAddressFromKey.IsEqualTo(common.HexToAddress(depositorAddr)) {
		return errors.New("depositor key address check failed " + err.Error())
	}

	validatorKeyFile, err := findKeyFile(validatorAddr)
	if err != nil {
		return errors.New("error finding VALIDATOR_ADDRESS in DP_KEY_FILE_DIR " + err.Error())
	}

	fmt.Println(fmt.Sprintf("Validator wallet addres %s", validatorKeyFile))
	validatorPwd, err := prompt.Stdin.PromptPassword(fmt.Sprintf("Enter the validator wallet password : "))
	if err != nil {
		return err
	}
	if len(validatorPwd) == 0 {
		return errors.New("validator password is not set")
	}
	fmt.Println()

	validatorPasswordConfirm, err := prompt.Stdin.PromptConfirm(fmt.Sprintf("Do you understand that the validator password will always be required to use the quantum validator wallet at %s?", validatorKeyFile))
	if err != nil {
		return err
	}
	if validatorPasswordConfirm != true {
		return errors.New("confirmation not made")
	}
	fmt.Println()

	valKey, err := GetKeyFromFile(validatorKeyFile, validatorPwd)
	if err != nil {
		return errors.New("error decrypting validator key " + err.Error())
	}

	valAddressFromKey, err := cryptobase.SigAlg.PublicKeyToAddress(&valKey.PublicKey)
	if err != nil {
		return errors.New("validator PublicKeyToAddress " + err.Error())
	}

	if !valAddressFromKey.IsEqualTo(common.HexToAddress(validatorAddr)) {
		return errors.New("validator key address check failed " + err.Error())
	}

	if len(rawURL) == 0 {
		return requestNewDeposit(validatorAddr, depositorAmount, depKey)
	} else {
		return newDeposit(validatorAddr, depositorAmount, depKey)
	}
}

func InitiateWithdrawal() error {
	if len(os.Args) < 3 {
		printHelp()
		return errors.New("incorrect usage")
	}

	if len(os.Getenv("DP_KEY_FILE_DIR")) == 0 {
		return errors.New("set the keyfile directory environment variable DP_KEY_FILE_DIR")
	}

	depositorAddr := os.Args[2]

	if common.IsHexAddress(depositorAddr) == false {
		return errors.New("invalid depositor address " + depositorAddr)
	}

	depositorKeyFile, err := findKeyFile(depositorAddr)
	if err != nil {
		return errors.New("error finding DEPOSITOR_ADDRESS in DP_KEY_FILE_DIR " + err.Error())
	}

	fmt.Println(fmt.Sprintf("Depositor wallet addres %s", depositorKeyFile))
	depositorPwd, err := prompt.Stdin.PromptPassword(fmt.Sprintf("Enter the depositor wallet password : "))
	if err != nil {
		return err
	}
	if len(depositorPwd) == 0 {
		return errors.New("depositor password is not set")
	}

	depKey, err := GetKeyFromFile(depositorKeyFile, depositorPwd)
	if err != nil {
		return errors.New("error decrypting depositor key " + err.Error())
	}

	fmt.Println()

	depositorPasswordConfirm, err := prompt.Stdin.PromptConfirm(fmt.Sprintf("Do you understand that the depositor password will always be required to use the quantum depositor wallet at %s?", depositorKeyFile))
	if err != nil {
		return err
	}
	if depositorPasswordConfirm != true {
		return errors.New("confirmation not made")
	}
	fmt.Println()

	depAddressFromKey, err := cryptobase.SigAlg.PublicKeyToAddress(&depKey.PublicKey)
	if err != nil {
		return errors.New("depositor public key to address " + err.Error())
	}

	if !depAddressFromKey.IsEqualTo(common.HexToAddress(depositorAddr)) {
		return errors.New("depositor key address check failed " + err.Error())
	}

	if len(rawURL) == 0 {
		return requestInitiateWithdrawal(depKey)
	} else {
		return initiateWithdrawal(depKey)
	}
}

func CompleteWithdrawal() error {
	if len(os.Args) < 3 {
		printHelp()
		return errors.New("incorrect usage")
	}

	if len(os.Getenv("DP_KEY_FILE_DIR")) == 0 {
		return errors.New("set the keyfile directory environment variable DP_KEY_FILE_DIR")
	}

	depositorAddr := os.Args[2]

	if common.IsHexAddress(depositorAddr) == false {
		return errors.New("invalid depositor address " + depositorAddr)
	}

	depositorKeyFile, err := findKeyFile(depositorAddr)
	if err != nil {
		return errors.New("error finding DEPOSITOR_ADDRESS in DP_KEY_FILE_DIR " + err.Error())
	}

	fmt.Println(fmt.Sprintf("Depositor wallet addres %s", depositorKeyFile))
	depositorPwd, err := prompt.Stdin.PromptPassword(fmt.Sprintf("Enter the depositor wallet password : "))
	if err != nil {
		return err
	}
	if len(depositorPwd) == 0 {
		return errors.New("depositor password is not set")
	}

	depKey, err := GetKeyFromFile(depositorKeyFile, depositorPwd)
	if err != nil {
		return errors.New("error decrypting depositor key " + err.Error())
	}

	fmt.Println()

	depositorPasswordConfirm, err := prompt.Stdin.PromptConfirm(fmt.Sprintf("Do you understand that the depositor password will always be required to use the quantum depositor wallet at %s?", depositorKeyFile))
	if err != nil {
		return err
	}
	if depositorPasswordConfirm != true {
		return errors.New("confirmation not made")
	}
	fmt.Println()

	depAddressFromKey, err := cryptobase.SigAlg.PublicKeyToAddress(&depKey.PublicKey)
	if err != nil {
		return errors.New("depositor public key to address " + err.Error())
	}

	if !depAddressFromKey.IsEqualTo(common.HexToAddress(depositorAddr)) {
		return errors.New("depositor key address check failed " + err.Error())
	}

	if len(rawURL) == 0 {
		return requestCompleteWithdrawal(depKey)
	} else {
		return completeWithdrawal(depKey)
	}
}

func DepositorBalance() error {
	if len(os.Args) < 3 {
		printHelp()
		return errors.New("incorrect usage")
	}

	depositorAddr := os.Args[2]

	if common.IsHexAddress(depositorAddr) == false {
		return errors.New("invalid depositor address " + depositorAddr)
	}

	if len(rawURL) == 0 {
		return nil
	} else {
		return getBalanceOfDepositor(depositorAddr)
	}
}
