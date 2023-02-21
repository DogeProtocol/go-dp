package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/signaturealgorithm"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	//"strconv"
)

// Accounts struct which contains
// an array of users
type Accounts struct {
	Accounts []Account `json:"accounts"`
}

type Tokendps struct {
	Tokendps []Tokendp `json:"tokendps"`
}

type Contractdps struct {
	Contractdps []Contractdp `json:"contractdps"`
}

type Tokens struct {
	Tokens []TokenList `json:"tokens"`
}

// User struct which contains a name
// a type and a list of social links
type Account struct {
	Address   string `json:"address"`
	SecretKey string `json:"secretkey"`
	Password  string `json:"password"`
	Status    int64  `json:"status"`
	Amount    string `json:"amount"`
}

type Tokendp struct {
	Address         string `json:"address"`
	ContractAddress string `json:"contractaddress"`
	SecretKey       string `json:"secretkey"`
	Password        string `json:"password"`
	Amount          string `json:"amount"`
}

type Contractdp struct {
	Address         string `json:"address"`
	ContractAddress string `json:"contractaddress"`
	SecretKey       string `json:"secretkey"`
	Password        string `json:"password"`
}

type TokenList struct {
	Address  string `json:"address"`
	ChainId  string `json:"chainId"`
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals int64  `json:"decimals"`
	logoURI  string `json:"logoURI"`
}

type KeyStore struct {
	Handle *keystore.KeyStore
}

var rawURL = "" //"http://127.0.0.1:8545" //"http://172.31.34.126:8545" "\\\\.\\pipe\\geth.ipc"
var accountPrimary = ""
var tokenInfoPath = ""

var accountPrimaryPassword = ""
var dataFilePath = ""
var accountPassword = ""

const letterBytes = "abcdefghijklmnopqrstuvwxyz"

var reduceLength = 30

var accountHistory Accounts

func main() {
	//AccountPasswordCheckManual()
	//Deposit()
	//Withdraw()
	//ValidatorList("46f8c16c50b122a568c96fb5e97e44ca9cd205ce")

	rawURL = os.Getenv("GETH_URL")
	accountPrimary = os.Getenv("GETH_ALLOC_ACCOUNT")
	accountPrimaryPassword = os.Getenv("GETH_ALLOC_ACCOUNT_PASSWORD")
	tokenInfoPath = os.Getenv("TOKENS_INFO")
	dataFilePath = os.Getenv("GETH_DATA_PATH")
	accountPassword = os.Getenv("GETH_ACCOUNT_PASSWORD")

	account_history()

	fmt.Println("Start  new account coin transaction...")
	go startTestCoinByNewAccount()

	/*
		fmt.Println("Start account to account coin transaction...")
		go startTestCoinAccountByAccount()

		fmt.Println("Start new token...")
		go startTestAccountByNewToken()

		fmt.Println("Start new token account transaction...")
		go startTestNewTokenAccountByAccount()

		fmt.Println("Start  account to account token transaction")
		go startTestTokenByAccount()

		fmt.Println("Start new Contract to dynamic account")
		go startTestAccountByContract()

		fmt.Println("Start new contract account ")
		go startTestNewContractByAccount()

		fmt.Println("Start  load account history...")
		go load_account_history()
	*/
	fmt.Println("Waiting indefinitely...")
	<-make(chan int)
}

/*
	func startTestCoinByNewAccountManual(address string) {
		transFile := accountPrimary
		transAddress := "0x" + accountPrimary

		transSecretKey, _ := ReadDataFile(transFile)
		transKey, _ := keystore.DecryptKey(transSecretKey, "TestNet123$$&#$%$!#s%")

		//Transfer fund
		newAccountAddress := strings.Split(address, ",")
		fmt.Println("newAccountAddress : ", newAccountAddress[0])
		fmt.Println("newAccountAddress : ", newAccountAddress[1])

		for _, element := range newAccountAddress {
			fmt.Println("element : ", element)
			if len(element) > 40 {
				for i := 0; i < 5; i++ {
					//amount randam
					rand.Seed(time.Now().UnixNano())
					min := int64(7000000000000000000)
					max := int64(9000000000000000000)
					amount := rand.Int63n(max-min) + min
					fmt.Println("transAddress element : ", transAddress, element)
					trans, _ := transferCoinToAccount(transAddress,
						element, amount, transKey.PrivateKey)

					fmt.Println(trans)
					fmt.Println("newAccountAddress : ", element, amount)
					time.Sleep(9 * time.Second)
				}
			}
		}

}
*/

func startTestCoinByNewAccount() {
	////oldTime := time.Now()
	////timeDiff := float64(0.60)

	password := accountPassword
	transFile := accountPrimary
	transAddress := "0x" + accountPrimary

	transSecretKey, _ := ReadDataFile(transFile)
	transKey, _ := keystore.DecryptKey(transSecretKey, accountPrimaryPassword)
	ks := SetUpKeyStore("./" + dataFilePath)

	////accounts := accountHistory
	////accLen := len(accounts.Accounts)

	////i := 0
	////j := 1
	loop := 1

	for loop != 0 {

		//val, err := dogep_getBalance(transAddress)
		//if err != nil {
		//	log.Println("Error occurred. dogep_getBalance : " + transAddress +
		//		" : " + err.Error())
		////return
		//}

		/****
				if accLen > 10 {
					rand.Seed(time.Now().UnixNano())
					min := int64(1)

					fmt.Println("i", i)
					fmt.Println("min", min)
					fmt.Println("len(accountHistory.Accounts)", len(accountHistory.Accounts))
					fmt.Println("len(accounts.Accounts)", len(accounts.Accounts))
					fmt.Println("accounts.Accounts[i].Amount", accounts.Accounts[i].Amount)

					var max int64
					fmt.Sscan(accounts.Accounts[i].Amount, &max)

					if max > 5 {
						max := int64(5)
						amount := rand.Int63n(max-min) + min

						//to Account selection
						k := i + j
						if k >= accLen {
							k = k - accLen
						}
						//Coin
						secretKey, _ := ReadDataFile(accounts.Accounts[i].SecretKey)
						key, _ := keystore.DecryptKey(secretKey, password)
						trans, err := transferCoinToAccount(accounts.Accounts[i].Address,
							accounts.Accounts[k].Address, amount, key.PrivateKey)
						if err != nil {
							log.Println("Error occurred. transferCoinToAccount : " + trans + " : " + err.Error())
							//return
						}
						fmt.Println("Start Transaction :")
						fmt.Println("From", accounts.Accounts[i].Address)
						fmt.Println("To", accounts.Accounts[k].Address)
						fmt.Println("amount", amount)
						fmt.Println("Tx hash :", trans)

						var accountAmount int64

						fmt.Sscan(accounts.Accounts[i].Amount, &accountAmount)
						accountbalance := accountAmount - amount
						accounts.Accounts[i].Amount = fmt.Sprint(accountbalance)

						fmt.Sscan(accounts.Accounts[k].Amount, &accountAmount)
						accountbalance = accountAmount + amount
						accounts.Accounts[k].Amount = fmt.Sprint(accountbalance)
					}
				}

				time.Sleep(500 * time.Millisecond)

				i++

				if i >= accLen {
					accounts = accountHistory
					accLen = len(accounts.Accounts)

					rand.Seed(time.Now().UnixNano())
					min := 1
					max := 4
					j = rand.Intn(max-min) + min

					rand.Seed(time.Now().UnixNano())
					min = 1
					max = accLen - j

					i = rand.Intn(max-min) + min
					fmt.Println("To address start and skip:", i, j)
				}
		****/
		////diff := time.Now().Sub(oldTime).Minutes()

		//mainBalance, _ := ParseBigFloat(val)
		//var mbalance int64
		//fmt.Sscan(mainBalance.String(), &mbalance)

		//if mbalance >= 1 {
		//Create account
		newAccount := ks.CreateNewKeys(password)
		fmt.Println("newAccount : ", newAccount.Address.String())

		//Transfer fund
		rand.Seed(time.Now().UnixNano())
		min := int64(7)
		max := int64(9)
		amount := rand.Int63n(max-min) + min

		trans, err := transferCoinToAccount(transAddress,
			newAccount.Address.String(), amount, transKey.PrivateKey)
		if err != nil {
			log.Println("Error occurred. transferCoinToAccount : " + trans + " : " + err.Error())
		}
		fmt.Println("Tx hash :", trans)

		////acc := Account{Address: newAccount.Address.String(), SecretKey: newAccount.URL.Path,
		////	Password: password, Status: 0, Amount: fmt.Sprint(amount)}
		////accountHistory.Accounts = append(accountHistory.Accounts, acc)
		////oldTime = time.Now()
		//}
		time.Sleep(185 * time.Second)
		////time.Sleep(500 * time.Millisecond)
	}
}

func startTestCoinAccountByAccount() {
	password := accountPassword

	accounts := accountHistory
	accLen := len(accounts.Accounts)

	i := 0
	j := 1
	loop := 1

	for loop != 0 {

		if accLen > 10 {
			rand.Seed(time.Now().UnixNano())
			min := int64(1)

			var max int64
			fmt.Sscan(accounts.Accounts[i].Amount, &max)

			if max > 5 {
				max := int64(5)
				amount := rand.Int63n(max-min) + min

				//to Account selection
				k := i + j
				if k >= accLen {
					k = k - accLen
				}

				//Coin
				secretKey, _ := ReadDataFile(accounts.Accounts[i].SecretKey)
				key, _ := keystore.DecryptKey(secretKey, password)
				trans, err := transferCoinToAccount(accounts.Accounts[i].Address,
					accounts.Accounts[k].Address, amount, key.PrivateKey)
				if err != nil {
					log.Println("Error occurred. transferCoinToAccount : " + trans + " : " + err.Error())
				}
				fmt.Println("Start Transaction :")
				fmt.Println("From", accounts.Accounts[i].Address)
				fmt.Println("To", accounts.Accounts[k].Address)
				fmt.Println("amount", amount)
				fmt.Println("Tx hash :", trans)

				var accountAmount int64
				fmt.Sscan(accounts.Accounts[i].Amount, &accountAmount)
				accountbalance := accountAmount - amount
				accounts.Accounts[i].Amount = fmt.Sprint(accountbalance)

				fmt.Sscan(accounts.Accounts[k].Amount, &accountAmount)
				accountbalance = accountAmount + amount
				accounts.Accounts[k].Amount = fmt.Sprint(accountbalance)
			}
		}

		i++
		if i >= accLen {
			accounts = accountHistory
			accLen = len(accounts.Accounts)

			rand.Seed(time.Now().UnixNano())
			min := 1
			max := 4
			j = rand.Intn(max-min) + min

			rand.Seed(time.Now().UnixNano())
			min = 1
			max = accLen - j
			i = rand.Intn(max-min) + min
			fmt.Println("To address start and skip:", i, j)
		}
		time.Sleep(9 * time.Second)
	}
}

func transferCoinToAccount(fromaddress string, toaddress string, amount int64,
	key *signaturealgorithm.PrivateKey) (string, error) {

	client, err := ethclient.Dial(rawURL)
	if err != nil {
		return "0", err
	}

	fromAddress := common.HexToAddress(fromaddress)
	toAddress := common.HexToAddress(toaddress)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "1", err
	}
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return "2", err
	}
	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "3", err
	}
	v, _ := ParseBigFloat(strconv.Itoa(int(amount)))
	value := etherToWei(v)
	fmt.Println("toAddress", toAddress)
	fmt.Println("value", value)

	var data []byte
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), key)
	if err != nil {
		return "4", err
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "5", err
	}
	return signedTx.Hash().Hex(), nil
}

func startTestAccountByNewToken() {
	var tokendps Tokendps

	//oldTime := time.Now()
	//timeDiff := float64(0.3)

	accounts := accountHistory
	password := accountPassword

	//load old contract
	contract := ReadNewTokenJsonDataFile("contract.json")
	tokendps.Tokendps = contract
	accLen := len(accounts.Accounts) - reduceLength

	i := 0
	loop := 1
	for loop != 0 {
		if accLen > reduceLength {
			////diff := time.Now().Sub(oldTime).Minutes()
			var max int64
			fmt.Sscan(accounts.Accounts[i].Amount, &max)

			//if diff >= timeDiff && max >= 1 {
			if max >= 1 {
				secretKey, _ := ReadDataFile(accounts.Accounts[i].SecretKey)
				key, _ := keystore.DecryptKey(secretKey, password)
				tokenContract, tokenValue, err := DeployTestTokenContractDynamicAccount(
					accounts.Accounts[i].Address, key.PrivateKey)

				if err == nil {
					tok := Tokendp{Address: accounts.Accounts[i].Address, ContractAddress: tokenContract,
						SecretKey: accounts.Accounts[i].SecretKey, Password: password, Amount: tokenValue}
					tokendps.Tokendps = append(tokendps.Tokendps, tok)
					content, _ := json.Marshal(tokendps.Tokendps)
					_ = ioutil.WriteFile("contract.json", content, 0644)
				}

				fmt.Println("Start transCount :")
				fmt.Println("Address", accounts.Accounts[i].Address)
				fmt.Println("Token contract", tokenContract, err)
				fmt.Println("Token amount", tokenValue)
				////oldTime = time.Now()
			}
		}
		i++
		if i >= accLen {
			accounts = accountHistory
			accLen = len(accounts.Accounts) - reduceLength
			if i < accLen {
				rand.Seed(time.Now().UnixNano())
				min := i
				max := accLen
				i = rand.Intn(max-min) + min
			} else {
				i = 0
			}
		}
		time.Sleep(130 * time.Second)
	}
}

func startTestNewTokenAccountByAccount() {
	var tokendps Tokendps

	password := accountPassword
	transCount := 0
	transrandNum := 0

	accounts := accountHistory
	accLen := len(accounts.Accounts)

	j := 1
	i := 0

	c := 0
	p := 1
	contractLen := 0

	fromAddress := ""
	contractAdd := ""
	fromAddressSecretKeyPath := ""
	totalToken := int64(0)
	tokenAmount := int64(0)

	loop := 1
	for loop != 0 {
		//Block transaction count
		if transrandNum <= transCount {
			rand.Seed(time.Now().UnixNano())
			min := 1
			max := 4
			j = rand.Intn(max-min) + min

			contract := ReadNewTokenJsonDataFile("contract.json")
			contractLen = len(contract)
			contractLen = contractLen - reduceLength
			totalToken = 0
			tokenAmount = 0

			if contractLen > 1 {
				if c > contractLen {
					c = 0
				}

				rand.Seed(time.Now().UnixNano())
				min := 1
				max := 2
				transrandNum = rand.Intn(max-min) + min
				transCount = 0
				fmt.Println("Planed transaction Count :", transrandNum)

				fromAddress = contract[c].Address
				contractAdd = contract[c].ContractAddress
				fromAddressSecretKeyPath = contract[c].SecretKey

				fmt.Sscan(contract[c].Amount, &totalToken)
				tokenAmount = int64(totalToken) * 10 / 100

				c = c + 1
			}
		}
		fmt.Println("contractAddress :", contractAdd)
		fmt.Println("startTestNewTokenAccountByAccount: tokenAmount, contractLen, c, transrandNum, transCount", tokenAmount, contractLen, c, transrandNum, transCount)

		//to Account selection
		if tokenAmount > 0 && contractLen > 1 && c <= contractLen && accLen > 10 {
			k := i + j
			if k >= accLen {
				k = k - accLen
			}
			fmt.Println("k, fromAddress,  accounts.Accounts[k].Address ", k, fromAddress, accounts.Accounts[k].Address)
			if fromAddress != accounts.Accounts[k].Address {
				secretKey, _ := ReadDataFile(fromAddressSecretKeyPath)
				key, _ := keystore.DecryptKey(secretKey, password)
				t, err := transferTokenToAccount(fromAddress,
					accounts.Accounts[k].Address, contractAdd, tokenAmount, key.PrivateKey)

				fmt.Println("Start Token transCount:", transCount, err)
				if err == nil {
					fmt.Println("Start Token transCount:", transCount)
					fmt.Println("From", fromAddress)
					fmt.Println("To", accounts.Accounts[k].Address)
					fmt.Println("Token contract", contractAdd)
					fmt.Println("TokenAmount", tokenAmount)
					fmt.Println(t)

					tok := Tokendp{Address: accounts.Accounts[k].Address, ContractAddress: contractAdd,
						SecretKey: accounts.Accounts[k].SecretKey, Password: password, Amount: fmt.Sprint(tokenAmount)}

					tokendps.Tokendps = append(tokendps.Tokendps, tok)
					content, _ := json.Marshal(tokendps.Tokendps)
					_ = ioutil.WriteFile("clientcontract.json", content, 0644)

					rand.Seed(time.Now().UnixNano())
					min := 1
					max := 10
					p = rand.Intn(max-min) + min

					totalToken = totalToken - tokenAmount
					tokenAmount = int64(totalToken) * int64(p) / 100

					if transCount == transrandNum {
						tok := Tokendp{Address: fromAddress, ContractAddress: contractAdd,
							SecretKey: fromAddressSecretKeyPath, Password: password, Amount: fmt.Sprint(totalToken)}
						tokendps.Tokendps = append(tokendps.Tokendps, tok)
						content, _ := json.Marshal(tokendps.Tokendps)
						_ = ioutil.WriteFile("clientcontract.json", content, 0644)
					}
				}
				//transCount++
			}
		}

		transCount++
		i++

		if i >= accLen {
			accounts = accountHistory
			accLen = len(accounts.Accounts)
			i = 0
			rand.Seed(time.Now().UnixNano())
			min := 1
			max := accLen - j
			i = rand.Intn(max-min) + min
			fmt.Println("To address start and skip:", i, j)
		}

		if totalToken < 1 {
			transCount = 0
			transrandNum = 0
		}
		time.Sleep(10 * time.Second)
	}
}

func startTestTokenByAccount() {
	var tokendps Tokendps

	accounts := accountHistory
	password := accountPassword

	transCount := 0
	transrandNum := 0
	accLen := len(accounts.Accounts)

	j := 1
	i := 0
	c := 0
	p := 1
	contractLen := 1

	fromAddress := ""
	contractAddress := ""
	fromAddressSecretKeyPath := ""
	totalToken := int64(0)
	tokenAmount := int64(0)

	loop := 1
	for loop != 0 {

		//Block transaction count
		if transrandNum <= transCount {
			rand.Seed(time.Now().UnixNano())
			min := 1
			max := 4
			j = rand.Intn(max-min) + min

			contract := ReadNewTokenJsonDataFile("clientcontract.json")
			contract1 := ReadNewTokenJsonDataFile("clientcontract-1.json")

			contract = append(contract, contract1...)

			contractLen = len(contract)
			contractLen = contractLen - reduceLength
			totalToken = 0
			tokenAmount = 0

			if contractLen > 1 {
				if c > contractLen {
					c = 0
				}

				rand.Seed(time.Now().UnixNano())
				min := 1
				max := 3
				transrandNum = rand.Intn(max-min) + min
				transCount = 0

				fmt.Println("Planed transaction Count :", transrandNum)

				fromAddress = contract[c].Address
				contractAddress = contract[c].ContractAddress
				fromAddressSecretKeyPath = contract[c].SecretKey

				fmt.Sscan(contract[c].Amount, &totalToken)
				tokenAmount = int64(totalToken) * 10 / 100

				c = c + 1
			}
		}

		fmt.Println("startTestTokenByAccount: totalToken tokenAmount, contractLen, c, transrandNum, transCount", totalToken, tokenAmount, contractLen, c, transrandNum, transCount)

		//to Account selection
		if tokenAmount > 0 && contractLen > 1 && c <= contractLen && accLen > 10 {
			k := i + j
			if k >= accLen {
				k = k - accLen
			}
			fmt.Println("k, fromAddress,  accounts.Accounts[k].Address ", k, fromAddress, accounts.Accounts[k].Address)
			if fromAddress != accounts.Accounts[k].Address {
				secretKey, _ := ReadDataFile(fromAddressSecretKeyPath)
				key, _ := keystore.DecryptKey(secretKey, password)
				t, err := transferTokenToAccount(fromAddress,
					accounts.Accounts[k].Address, contractAddress, tokenAmount, key.PrivateKey)
				fmt.Println("Start Token transCount:", transCount, err)
				if err == nil {
					fmt.Println("Start Token transCount:", transCount)
					fmt.Println("From", fromAddress)
					fmt.Println("To", accounts.Accounts[k].Address)
					fmt.Println("Token contract", contractAddress)
					fmt.Println("TokenAmount", tokenAmount)
					fmt.Println(t)

					tok := Tokendp{Address: accounts.Accounts[k].Address, ContractAddress: contractAddress,
						SecretKey: accounts.Accounts[k].SecretKey, Password: password, Amount: fmt.Sprint(tokenAmount)}

					tokendps.Tokendps = append(tokendps.Tokendps, tok)
					content, _ := json.Marshal(tokendps.Tokendps)
					_ = ioutil.WriteFile("clientcontract-1.json", content, 0644)

					rand.Seed(time.Now().UnixNano())
					min := 1
					max := 10
					p = rand.Intn(max-min) + min

					totalToken = totalToken - tokenAmount
					tokenAmount = int64(totalToken) * int64(p) / 100
				}
				//transCount++
			}
		}

		transCount++
		i++

		if i >= accLen {
			accLen = len(accounts.Accounts)
			rand.Seed(time.Now().UnixNano())
			min := 1
			max := accLen - j
			i = rand.Intn(max-min) + min
			fmt.Println("To address start and skip:", i, j)
		}

		if totalToken < 1 {
			transCount = 0
			transrandNum = 0
		}
		time.Sleep(10 * time.Second)
	}
}

func startTestAccountByContract() {
	var contractdps Contractdps

	oldTime := time.Now()
	timeDiff := float64(0.5)

	password := accountPassword

	//load old contract
	contract := ReadContractJsonDataFile("newcontract.json")
	contractdps.Contractdps = contract

	accounts := accountHistory
	accLen := len(accounts.Accounts) - reduceLength

	i := 0
	loop := 1
	for loop != 0 {
		if accLen > reduceLength {
			//Coin
			diff := time.Now().Sub(oldTime).Minutes()
			if diff >= timeDiff {
				secretKey, _ := ReadDataFile(accounts.Accounts[i].SecretKey)
				key, _ := keystore.DecryptKey(secretKey, password)
				contract, value, err := DeployTestOtherContractDynamicAccount(
					accounts.Accounts[i].Address, key.PrivateKey)

				if err == nil {
					con := Contractdp{Address: accounts.Accounts[i].Address, ContractAddress: contract,
						SecretKey: accounts.Accounts[i].SecretKey, Password: password}
					contractdps.Contractdps = append(contractdps.Contractdps, con)
					content, _ := json.Marshal(contractdps.Contractdps)
					_ = ioutil.WriteFile("newcontract.json", content, 0644)
				}

				fmt.Println("Start transCount :")
				fmt.Println("Address", accounts.Accounts[i].Address)
				fmt.Println("New contract", contract, err)
				fmt.Println("New value", value)

				oldTime = time.Now()
			}
		}

		i++

		if i >= accLen {
			accounts = accountHistory
			accLen = len(accounts.Accounts) - reduceLength

			rand.Seed(time.Now().UnixNano())
			min := 1
			max := accLen
			i = rand.Intn(max-min) + min
		}

		time.Sleep(5 * time.Second)
	}
}

func startTestNewContractByAccount() {

	password := accountPassword

	c := 0
	contractLen := 0
	fromAddress := ""
	contractAddress := ""
	fromAddressSecretKeyPath := ""

	//accounts := accountHistory
	//accLen := len(accounts.Accounts) - reduceLength

	//i := 0
	loop := 1
	for loop != 0 {
		//if accLen > reduceLength {
		//Block transaction count
		contract := ReadContractJsonDataFile("newcontract.json")
		contractLen = len(contract)
		contractLen = contractLen - reduceLength
		if contractLen > 1 && c < contractLen {

			fromAddress = contract[c].Address
			contractAddress = contract[c].ContractAddress
			fromAddressSecretKeyPath = contract[c].SecretKey

			secretKey, _ := ReadDataFile(fromAddressSecretKeyPath)
			key, _ := keystore.DecryptKey(secretKey, password)
			t, value, err := setContractDataToAccount(fromAddress, contractAddress, key.PrivateKey)

			if err == nil {
				fmt.Println("Start Contract set function:")
				fmt.Println("Address ", fromAddress)
				fmt.Println("Contract ", contractAddress)
				fmt.Println("Value ", value)
				fmt.Println(t)
			}
			c = c + 1
		}
		//}

		/*
			i++
			if i > accLen {
				accounts = accountHistory
				accLen = len(accounts.Accounts) - reduceLength
				rand.Seed(time.Now().UnixNano())
				min := 1
				max := accLen
				i = rand.Intn(max-min) + min
			}
		*/

		time.Sleep(15 * time.Second)
	}
}

func DeployTestTokenContractDynamicAccount(address string,
	key *signaturealgorithm.PrivateKey) (string, string, error) {

	client, err := ethclient.Dial(rawURL)

	if err != nil {
		return "", "0", err
	}
	addr := common.HexToAddress(address)

	nonce, err := client.PendingNonceAt(context.Background(), addr)
	if err != nil {
		return "", "0", err
	}
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return "", "0", err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		return "", "0", err
	}

	auth.From = addr
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Println(err.Error())
		return "", "0", err
	}
	auth.GasPrice = gasPrice

	var tokens Tokens
	tk := ReadTokenJsonDataFile(tokenInfoPath)
	tokenLength := len(tk)
	tokens.Tokens = tk

	rand.Seed(time.Now().UnixNano())
	min := int64(1)
	max := int64(tokenLength)
	dynamictoken := int(rand.Int63n(max-min) + min)

	var sb strings.Builder
	sb.WriteString(tokens.Tokens[dynamictoken-1].Name)
	sb.WriteString(" (" + EncodeToString(6) + ") ")

	tokenName := sb.String()
	tokenSymbol := tokens.Tokens[dynamictoken-1].Symbol

	fmt.Println("Name : ", tokenName)
	fmt.Println("Symbol : ", tokenSymbol)

	rand.Seed(time.Now().UnixNano())
	min = int64(0)
	max = int64(18)
	decimalUnits := uint8(rand.Int63n(max-min) + min)

	rand.Seed(time.Now().UnixNano())
	min = int64(1)
	max = int64(9)
	t := int(rand.Int63n(max-min) + min)

	rand.Seed(time.Now().UnixNano())
	min = int64(12)
	max = int64(15)
	p := int(rand.Int63n(max-min) + min)
	number, _ := strconv.Atoi(RandPowTotal(t, p))
	totalSupply := big.NewInt(int64(number))

	fmt.Println("totalSupply", totalSupply)

	contractAddress, tx, _, err := DeployToken(auth, client, tokenName, tokenSymbol, decimalUnits,
		totalSupply)

	if err != nil {
		log.Println(err.Error())
		return "", "0", err
	}
	fmt.Println("Tx hash:" + tx.Hash().Hex())

	// Don't even wait, check its presence in the local pending state
	time.Sleep(250 * time.Millisecond) // Allow it to be processed by the local node :P
	fmt.Println("Contract address: " + contractAddress.String())

	return contractAddress.String(), totalSupply.String(), nil
}

func DeployTestOtherContractDynamicAccount(address string,
	key *signaturealgorithm.PrivateKey) (string, string, error) {

	client, err := ethclient.Dial(rawURL)

	if err != nil {
		log.Println(err.Error())
		return "", "0", err
	}
	addr := common.HexToAddress(address)

	nonce, err := client.PendingNonceAt(context.Background(), addr)
	if err != nil {
		log.Println(err.Error())
		return "", "0", err
	}
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Println(err.Error())
		return "", "0", err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		log.Println(err.Error())
		return "", "0", err
	}

	auth.From = addr
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Println(err.Error())
		return "", "0", err
	}
	auth.GasPrice = gasPrice

	rand.Seed(time.Now().UnixNano())
	min := int64(100)
	max := int64(10000)
	amount := uint8(rand.Int63n(max-min) + min)

	contractAddress, tx, _, err := DeployGreeter(auth, client, big.NewInt(int64(amount)))

	if err != nil {
		log.Println(err.Error())
		return "", "0", err
	}

	//fmt.Printf("Contract pending deploy: 0x%x\n", contractAddress)

	fmt.Println("Tx hash:" + tx.Hash().Hex())

	// Don't even wait, check its presence in the local pending state
	time.Sleep(250 * time.Millisecond) // Allow it to be processed by the local node :P
	fmt.Println("Contract address: " + contractAddress.String())

	return contractAddress.String(), strconv.Itoa(int(amount)), nil
}

func transferTokenToAccount(fromaddress string, toaddress string, tokenaddress string,
	amount int64, key *signaturealgorithm.PrivateKey) (string, error) {

	client, err := ethclient.Dial(rawURL)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	fromAddress := common.HexToAddress(fromaddress)
	toAddress := common.HexToAddress(toaddress)
	tokenAddress := common.HexToAddress(tokenaddress)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	auth.From = fromAddress
	auth.Nonce = big.NewInt(int64(nonce))
	//auth.Value = big.N

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	auth.GasPrice = gasPrice

	token, err := NewToken(tokenAddress, client)
	if err != nil {
		//log.Fatalf("Failed to instantiate a Token contract: %v", err)
		log.Println(err.Error())
		return "", err
	}

	name, err := token.Name(nil)
	if err != nil {
		//log.Fatalf("Failed to retrieve token name: %v", err)
		log.Println(err.Error())
		return "", err
	}
	fmt.Println("Token name.....:", name)

	tx, err := token.Transfer(auth, toAddress, big.NewInt(amount))
	if err != nil {
		//log.Fatalf("Failed to request token transfer: %v", err)
		log.Println(err.Error())
		return "", err
	}
	//fmt.Printf("Transfer pending: 0x%x\n", tx.Hash())

	// Don't even wait, check its presence in the local pending state
	time.Sleep(250 * time.Millisecond) // Allow it to be processed by the local node :P

	return tx.Hash().Hex(), nil
}

func setContractDataToAccount(fromaddress string, contractaddress string,
	key *signaturealgorithm.PrivateKey) (string, string, error) {

	client, err := ethclient.Dial(rawURL)
	if err != nil {
		log.Println(err.Error())
		return "", "0", err
	}

	fromAddress := common.HexToAddress(fromaddress)
	contractAddress := common.HexToAddress(contractaddress)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Println(err.Error())
		return "", "0", err
	}
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Println(err.Error())
		return "", "0", err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		log.Println(err.Error())
		return "", "0", err
	}

	auth.From = fromAddress
	auth.Nonce = big.NewInt(int64(nonce))
	//auth.Value = big.N

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Println(err.Error())
		return "", "0", err
	}
	auth.GasPrice = gasPrice

	contract, err := NewGreeter(contractAddress, client)
	if err != nil {
		log.Println(err.Error())
		return "", "0", err
	}

	rand.Seed(time.Now().UnixNano())
	min := int64(100)
	max := int64(10000)
	amount := uint8(rand.Int63n(max-min) + min)

	tx, err := contract.Set(auth, big.NewInt(int64(amount)))
	if err != nil {
		//log.Fatalf("Failed to request token transfer: %v", err)
		log.Println(err.Error())
		return "", "0", err
	}
	//fmt.Printf("Transfer pending: 0x%x\n", tx.Hash())

	// Don't even wait, check its presence in the local pending state
	time.Sleep(250 * time.Millisecond) // Allow it to be processed by the local node :P

	return "Tx sent: " + tx.Hash().Hex(), strconv.Itoa(int(amount)), nil
}

func SetUpKeyStore(kp string) *KeyStore {
	ks := &KeyStore{}
	ks.Handle = keystore.NewKeyStore(kp, keystore.LightScryptN, keystore.LightScryptP)
	return ks
}

func (ks *KeyStore) CreateNewKeys(password string) accounts.Account {
	account, err := ks.Handle.NewAccount(password)
	if err != nil {
		log.Println(err.Error())
	}
	return account
}

func (ks *KeyStore) GetKeysByAddress(address string) accounts.Account {

	var account accounts.Account
	var err error
	if ks.Handle.HasAddress(common.HexToAddress(address)) {
		if account, err = ks.Handle.Find(accounts.Account{Address: common.HexToAddress(address)}); err != nil {
			log.Println(err.Error())
		}
	}
	return account
}

func (ks *KeyStore) GetAllKeys() []accounts.Account {
	return ks.Handle.Accounts()
}

//func createKs() {
//	ks := keystore.NewKeyStore("./" + dataFilePath, keystore.StandardScryptN, keystore.StandardScryptP)
//	password := accountPassword
//	account, err := ks.NewAccount(password)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(account.Address.Hex()) // 0x20F8D42FB0F667F2E53930fed426f225752453b3
//}

//func RandTokenName(n int) string {
//	b := make([]byte, n)
//	for i := range b {
//		b[i] = letterBytes[rand.Intn(len(letterBytes))]
//	}
//	return string(b)
//}

//func RandTokenSymbol(n int, s string) string {
//	b := make([]byte, n)
//	for i := range b {
//		b[i] = s[rand.Intn(len(s))]
//	}
//	return string(b)
//}

func RandPowTotal(t int, p int) string {
	b := strconv.Itoa(t)
	for i := 0; i < p; i++ {
		b = b + "0"
	}
	return b
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

func ReadNewTokenJsonDataFile(filename string) (tok []Tokendp) {

	fileContent, err := os.Open(filename)

	if err != nil {
		log.Println(err.Error())
	}
	//fmt.Println("The File is opened successfully...", filename)
	defer fileContent.Close()

	byteResult, _ := ioutil.ReadAll(fileContent)

	var tokendps []Tokendp
	json.Unmarshal(byteResult, &tokendps)

	return tokendps
}

func ReadTokenJsonDataFile(filename string) (tok []TokenList) {
	fileContent, err := os.Open(filename)
	if err != nil {
		log.Println(err.Error())
	}

	defer fileContent.Close()
	byteResult, _ := ioutil.ReadAll(fileContent)

	var tokens []TokenList
	json.Unmarshal(byteResult, &tokens)

	return tokens
}

func ReadContractJsonDataFile(filename string) (contract []Contractdp) {

	fileContent, err := os.Open(filename)

	if err != nil {
		log.Println(err.Error())
	}
	//fmt.Println("The File is opened successfully...", filename)
	defer fileContent.Close()

	byteResult, _ := ioutil.ReadAll(fileContent)

	var contractdps []Contractdp
	json.Unmarshal(byteResult, &contractdps)

	return contractdps
}

func EncodeToString(max int) string {
	b := make([]int, max)
	for i := 0; i < max; i++ {
		rand.Seed(time.Now().UnixNano())
		min := 0
		max := 9
		b[i] = min + rand.Intn(max-min)
		time.Sleep(250 * time.Millisecond)
	}
	//return string(b)
	//b := rand.Perm(max)
	//s, _ := json.Marshal(b)
	st := strings.Trim(strings.Replace(fmt.Sprint(b), " ", "", -1), "[]")
	return st
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

func dogep_getBalance(address string) (string, error) {
	client, err := ethclient.Dial(rawURL)
	if err != nil {
		return "", err
	}
	balance, err := client.BalanceAt(context.Background(), common.HexToAddress(address), nil)
	if err != nil {
		return "", err
	}
	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	accountbalance := fmt.Sprint(ethValue)
	return accountbalance, nil
}

func load_account_history() {
	loop := 1
	for loop != 0 {
		time.Sleep(1 * time.Hour)
		account_history()
	}
}

func account_history() {
	var accounts Accounts
	//Load account information
	password := accountPassword
	ks := SetUpKeyStore("./" + dataFilePath)
	//List accounts
	account := ks.GetAllKeys()
	for i := 0; i < len(account); i++ {
		val, err := dogep_getBalance(account[i].Address.String())
		if err != nil {
			log.Println("Error occurred. dogep_getBalance : " + account[i].Address.String() +
				" : " + err.Error())
		}
		mainBalance, _ := ParseBigFloat(val)
		var mbalance float64
		fmt.Sscan(mainBalance.String(), &mbalance)
		if mbalance >= 1 {
			acc := Account{Address: account[i].Address.String(), SecretKey: account[i].URL.Path,
				Password: password, Status: 0, Amount: val}
			accounts.Accounts = append(accounts.Accounts, acc)
		}
		if mbalance <= 0.000000000000000000 {
			dir, file := filepath.Split(account[i].URL.Path)
			moveFile(account[i].URL.Path, dir+"/temp/"+file)
		}
		time.Sleep(100 * time.Microsecond)
	}
	accountHistory = accounts
}

func moveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %s", err)
	}
	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("Couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("Writing to output file failed: %s", err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("Failed removing original file: %s", err)
	}
	return nil
}
