package main

import (
	"context"
	"fmt"
	"github.com/DogeProtocol/dp/accounts/abi/bind"
	"github.com/DogeProtocol/dp/accounts/keystore"
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/crypto/cryptobase"
	"github.com/DogeProtocol/dp/crypto/signaturealgorithm"
	"github.com/DogeProtocol/dp/ethclient"
	"log"
	"math/big"
	"os"
	"time"
)

var contractAddress = os.Getenv("DP_STAKING_CONTRACT_ADDRESS")

func AccountPasswordCheckManual() {
	password := "dummy"
	ks := SetUpKeyStore("./data/keystore")
	accounts := ks.GetAllKeys()

	for _, account := range accounts {
		acc := account.Address.String()
		fmt.Println("address : ", acc)
		secretKey, _ := ReadDataFile(account.URL.Path)
		key, _ := keystore.DecryptKey(secretKey, password)
		pubKey, err := cryptobase.SigAlg.SerializePublicKey(&key.PrivateKey.PublicKey)
		if err != nil {
			panic(err)
		}
		fmt.Println("pubKey : ", pubKey)
	}
	fmt.Println("accounts : ", accounts)
	return
}

func Deposit() {
	//Payment transaction
	account3 := "486b2f8ae61c013da9b9d4c311026802572c9f86"
	account3Slave := "0x4643635a54Ca29C1E803B9c0Eca489426757c4C2"
	deposit(account3, account3Slave)
	time.Sleep(10 * time.Second)

	account4 := "66d6e1400a300516f832cdbd7499901b986a798c"
	account4Slave := "0xCaceA5D099136Ea3bD10ADcE285D454Cda0b67b1"
	deposit(account4, account4Slave)
	time.Sleep(10 * time.Second)

	account5 := "99a0a81da594bccdb4820f27639839bae4273e5d"
	account5Slave := "0x9F01AB712dB8F438DB59B0B506fd2E764E841517"
	deposit(account5, account5Slave)
	time.Sleep(10 * time.Second)

}

func Withdraw() {
	account1 := "43BAc60aA592706fb962C333D2515CFd788F9A24"
	withdraw(account1)
	time.Sleep(10 * time.Second)

	account2 := "c82127a093e99fbbc3bfbaa8056bd4bdf45396a0"
	withdraw(account2)
	time.Sleep(10 * time.Second)
}

func withdraw(address string) {
	transFile := address
	transAddress := "0x" + address
	transSecretKey, _ := ReadDataFile(transFile)
	transKey, _ := keystore.DecryptKey(transSecretKey, "dummy")

	tx, err := withdrawContract(transAddress, contractAddress, transKey.PrivateKey)
	fmt.Println("tx : ", tx, " err : ", err, " transAddress : ", transAddress)

	time.Sleep(100 * time.Millisecond)
}

func ValidatorList(address string) {
	transFile := address
	transAddress := "0x" + address
	transSecretKey, _ := ReadDataFile(transFile)
	transKey, _ := keystore.DecryptKey(transSecretKey, "dummy")

	_, err := validatorList(transAddress, contractAddress, transKey.PrivateKey)
	if err != nil {
		log.Println(err.Error())
		return
	}
	time.Sleep(100 * time.Millisecond)
	return
}

func deposit(address string, accountSlave string) {

	transFile := address
	transAddress := "0x" + address
	transSecretKey, _ := ReadDataFile(transFile)
	transKey, _ := keystore.DecryptKey(transSecretKey, "dummy")

	password := "dummy"

	ks := SetUpKeyStore("./data/keystore")

	accounts := ks.GetAllKeys()

	for _, account := range accounts {
		acc := account.Address.String()
		if accountSlave == acc {
			fmt.Println("address, accountSlave", address, accountSlave)
			secretKey, _ := ReadDataFile(account.URL.Path)
			key, _ := keystore.DecryptKey(secretKey, password)
			pubKey, err := cryptobase.SigAlg.SerializePublicKey(&key.PrivateKey.PublicKey)
			if err != nil {
				panic(err)
			}

			tx, err := depositContract(transAddress, contractAddress, pubKey, transKey.PrivateKey)
			fmt.Println("tx : ", tx, " err : ", err, " transAddress : ", transAddress)

			time.Sleep(100 * time.Millisecond)
			return
		}
	}
	fmt.Println("accounts : ", accounts)
	return
}

func depositContract(fromaddress string, contractaddress string, pubKey []byte,
	key *signaturealgorithm.PrivateKey) (string, error) {

	client, err := ethclient.Dial(rawURL)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	fromAddress := common.HexToAddress(fromaddress)
	contractAddress := common.HexToAddress(contractaddress)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Println(err.Error())
		return "1", err
	}
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Println(err.Error())
		return "2", err
	}
	fmt.Println("chainID", chainID)
	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		log.Println(err.Error())
		return "3", err
	}
	auth.From = fromAddress
	auth.Nonce = big.NewInt(int64(nonce))
	p, _ := ParseBigFloat("5") //ether
	auth.Value = etherToWei(p)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Println(err.Error())
		return "4", err
	}
	auth.GasPrice = gasPrice

	contract, err := NewStakingContractAddress1(contractAddress, client)
	if err != nil {
		log.Println(err.Error())
		return "5", err
	}

	tx, err := contract.NewDeposit(auth, pubKey)
	if err != nil {

		log.Println(err.Error())
		return "6", err
	}

	// Don't even wait, check its presence in the local pending state
	time.Sleep(250 * time.Millisecond) // Allow it to be processed by the local node :P

	return "Tx sent: " + tx.Hash().Hex(), nil
}

func withdrawContract(address string, contractaddress string,
	key *signaturealgorithm.PrivateKey) (string, error) {

	client, err := ethclient.Dial(rawURL)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	toAddress := common.HexToAddress(address)
	contractAddress := common.HexToAddress(contractaddress)

	nonce, err := client.PendingNonceAt(context.Background(), toAddress)
	if err != nil {
		log.Println(err.Error())
		return "1", err
	}
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Println(err.Error())
		return "2", err
	}
	fmt.Println("chainID", chainID)
	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		log.Println(err.Error())
		return "3", err
	}
	auth.From = toAddress
	auth.Nonce = big.NewInt(int64(nonce))
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Println(err.Error())
		return "4", err
	}
	auth.GasPrice = gasPrice

	contract, err := NewStakingContractAddress1(contractAddress, client)
	if err != nil {
		log.Println(err.Error())
		return "5", err
	}

	p, _ := ParseBigFloat("1") //ether
	value := etherToWei(p)
	tx, err := contract.Withdraw(auth, value)
	if err != nil {

		log.Println(err.Error())
		return "6", err
	}

	// Don't even wait, check its presence in the local pending state
	time.Sleep(250 * time.Millisecond) // Allow it to be processed by the local node :P

	return "Tx sent: " + tx.Hash().Hex(), nil
}

func validatorList(fromaddress string, contractaddress string,
	key *signaturealgorithm.PrivateKey) (string, error) {

	client, err := ethclient.Dial(rawURL)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	fromAddress := common.HexToAddress(fromaddress)
	contractAddress := common.HexToAddress(contractaddress)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Println(err.Error())
		return "1", err
	}
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Println(err.Error())
		return "2", err
	}
	fmt.Println("chainID", chainID)
	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		log.Println(err.Error())
		return "3", err
	}
	auth.From = fromAddress
	auth.Nonce = big.NewInt(int64(nonce))
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Println(err.Error())
		return "4", err
	}
	auth.GasPrice = gasPrice

	contract, err := NewStakingContract(contractAddress, client)

	if err != nil {
		log.Println(err.Error())
		return "5", err
	}

	fmt.Println("contract.....:", contract)
	validator, err := contract.ListValidator(nil)
	if err != nil {
		log.Println(err.Error())
		return "666....", err
	}
	fmt.Println("Validator address.....:", validator)

	// Don't even wait, check its presence in the local pending state
	time.Sleep(250 * time.Millisecond) // Allow it to be processed by the local node :P

	return "", nil
}
