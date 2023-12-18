package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"

	"fmt"
	"github.com/DogeProtocol/dp/accounts/keystore"
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/crypto/crosssign"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

func printHelp() {
	fmt.Println("===========")
	fmt.Println("dputil genesis-sign ETH_ADDRESS DEPOSITOR_QUANTUM_ADDRESS VALIDATOR_QUANTUM_ADDRESS AMOUNT")
	fmt.Println("      Set the following environment variables:")
	fmt.Println("           DP_KEY_FILE_DIR, DP_DEPOSITOR_ACC_PWD, DP_VALIDATOR_ACC_PWD")
	fmt.Println("===========")
	fmt.Println("dputil genesis-verify JSON_FILE_NAME")
	fmt.Println("===========")
	fmt.Println("Set the environment variable DP_RAW_URL")
	fmt.Println("dputil balance ACCOUNT_ADDRESS")
	fmt.Println("===========")
	fmt.Println("dputil send FROM_ADDRESS TO_ADDRESS QUANTITY")
	fmt.Println("===========")
	fmt.Println("dputil bulksend CSV_FILE")
	fmt.Println("===========")
	fmt.Println("dputil bulksendsingle FROM_ADDRESS QUANTITY")
	fmt.Println("===========")
	fmt.Println("dputil bulksendreverse TO_ADDRESS QUANTITY COUNT TXN_PER_BATCH")
}

var rawURL string
var wg sync.WaitGroup

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}
	rawURL = os.Getenv("DP_RAW_URL")
	if os.Args[1] == "balance" {
		balance()
	} else if os.Args[1] == "send" {
		sendTxn()
	} else if os.Args[1] == "txn" {
		getTxn()
	} else if os.Args[1] == "bulksend" {
		sendTxnBulk()
	} else if os.Args[1] == "bulksendsingle" {
		sendTxnBulkFromSingleAddress()
	} else if os.Args[1] == "bulksendreverse" {
		sendTxnBulkToSingleAddress()
	} else if os.Args[1] == "genesis-sign" {
		GenesisSign()
	} else if os.Args[1] == "genesis-verify" {
		GenesisVerify()
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
	valPassword := os.Getenv("DP_DEPOSITOR_ACC_PWD")
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

func sendTxnBulkFromSingleAddress() {
	if len(os.Args) < 4 {
		printHelp()
		return
	}

	from := os.Args[2]
	quantity := os.Args[3]

	addresses, err := findAllAddresses()
	if err != nil {
		fmt.Println("findAllAddresses error", err)
		return
	}

	fmt.Println("addresses", len(addresses), "from", from)

	connectionContext, err := GetConnectionContext(from)
	if err != nil {
		fmt.Println("GetConnectionContext error occurred", "error", err)
		return
	}

	for i := 0; i < len(addresses); i++ {
		sendVia(connectionContext, addresses[i], quantity, 0)
	}
}

func sendTxnBulkToSingleAddress() {
	if len(os.Args) < 6 {
		printHelp()
		return
	}

	to := os.Args[2]
	quantity := os.Args[3]
	count, err := strconv.Atoi(os.Args[4])
	if err != nil {
		panic("conversion error")
	}

	txnPerBatch, err := strconv.Atoi(os.Args[5])
	if err != nil {
		panic("conversion error")
	}

	addresses, err := findAllAddresses()
	if err != nil {
		fmt.Println("findAllAddresses error", err)
		return
	}

	ctr := 0
	for i := 0; i < len(addresses); i++ {
		txn := Txn{
			FromAddress: addresses[i],
			ToAddress:   to,
			Quantity:    quantity,
			Count:       count,
		}

		wg.Add(1)
		go sendTxnSingleSender(txn)

		ctr = ctr + 1
		if ctr >= txnPerBatch {
			wg.Wait()
			ctr = 0
		}
	}

	fmt.Println("Waiting for all transactions")
	wg.Wait()
	fmt.Println("Done waiting for all transactions")
}

func sendTxnBulk() {
	if len(os.Args) < 3 {
		printHelp()
		return
	}

	csv := os.Args[2]

	txnMap := make(map[string][]Txn)

	//read to addresses
	file, err := os.Open(csv)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		columns := strings.Split(scanner.Text(), ",")
		count, err := strconv.Atoi(columns[3])
		if err != nil {
			panic("conversion error")
		}
		txn := Txn{
			FromAddress: columns[0],
			ToAddress:   columns[1],
			Quantity:    columns[2],
			Count:       count,
		}
		_, ok := txnMap[txn.FromAddress]
		if ok == false {
			txnMap[txn.FromAddress] = make([]Txn, 0)
		}
		txnMap[txn.FromAddress] = append(txnMap[txn.FromAddress], txn)
	}

	for from, txns := range txnMap {
		if common.IsHexAddress(from) == false {
			fmt.Println("Invalid from address", from)
			return
		}

		for i := 0; i < len(txns); i++ {
			txn := txns[i]
			wg.Add(1)
			go sendTxnSingleSender(txn)
		}
	}

	fmt.Println("Waiting for all transactions")
	wg.Wait()
	fmt.Println("Done waiting for all transactions")
}

func sendTxnSingleSender(txn Txn) {
	defer wg.Done()
	connectionContext, err := GetConnectionContext(txn.FromAddress)
	if err != nil {
		fmt.Println("GetConnectionContext error occurred", "error", err)
		return
	}

	var nonce uint64
	nonce = 0
	for j := 0; j < txn.Count; j++ {
		if common.IsHexAddress(txn.ToAddress) == false {
			fmt.Println("Invalid to address", txn.ToAddress)
			return
		}

		flt, err := ParseBigFloat(txn.Quantity)
		if err != nil {
			fmt.Println(err)
			return
		}

		wei := etherToWeiFloat(flt)
		ether := weiToEther(wei)

		fmt.Println("Send", "from", txn.FromAddress, "to", txn.ToAddress, "quantity", txn.Quantity, "ether", ether)

		txHash, nonceTmp, err := sendVia(connectionContext, txn.ToAddress, txn.Quantity, nonce)
		if err != nil {
			fmt.Println(err)
			return
		}
		nonce = nonceTmp + 1

		fmt.Println("TxnHash", txHash)
	}
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
