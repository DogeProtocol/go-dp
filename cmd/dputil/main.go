package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DogeProtocol/dp/common"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

func printHelp() {
	fmt.Println("Set the environment variable DP_RAW_URL")
	fmt.Println("dputil balance ACCOUNT_ADDRESS")
	fmt.Println("dputil send FROM_ADDRESS TO_ADDRESS QUANTITY")
	fmt.Println("dputil bulksend CSV_FILE")
	fmt.Println("dputil bulksendsingle FROM_ADDRESS QUANTITY")
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
	} else {
		printHelp()
	}
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
		fmt.Errorf("GetConnectionContext error occurred", "error", err)
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
		fmt.Errorf("GetConnectionContext error occurred", "error", err)
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
