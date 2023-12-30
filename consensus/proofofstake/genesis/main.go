package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/crypto"
	"github.com/DogeProtocol/dp/systemcontracts/conversion"
	"github.com/DogeProtocol/dp/systemcontracts/staking"
	"io/ioutil"
	"math/big"
	"os"
	"strconv"
	"strings"
)

type GenesisData struct {
	Depositors                []DepositorList `json:"depositors"`
	Alloc                     []AllocList     `json:"alloc"`
	TotalDepositBalance       string          `json:"totalDepositBalance"`
	Snapshot                  []Snapshot      `json:"snapshot"`
	ConversionContractBalance string          `json:"conversionContractBalance"`
}

type DepositorList struct {
	Address          string `json:"address"`
	ValidatorAddress string `json:"validatorAddress"`
	Balance          string `json:"balance"`
}

type ValidatorList struct {
	Address string `json:"address"`
}

type AllocList struct {
	Address string `json:"address"`
	Balance string `json:"balance"`
}

type Snapshot struct {
	EthAddress string `json:"ethAddress"`
	Balance    string `json:"balance"`
}

var genesisPath string

const ALLOC_TEMPLATE string = `    "[ADDRESS]": {
      "balance": "[BALANCE]"
    },
`

const STORAGE_KEY string = "[STORAGE_KEY]"
const STORAGE_VALUE string = "[STORAGE_VALUE]"

const STORAGE_TEMPLATE string = `       	"[STORAGE_KEY]": "[STORAGE_VALUE]",
`

const GENESIS_TEMPLATE string = `{
  "config": {
    "chainId": 123123,
    "homesteadBlock": 0,
    "eip150Block": 0,
    "eip155Block": 0,
    "eip158Block": 0,
    "byzantiumBlock": 0,
    "constantinopleBlock": 0,
    "petersburgBlock": 0,
    "istanbulBlock": 0,
    "berlinBlock": 0,
    "londonBlock": 0,
    "proofofstake": {
      "period": 4,
      "epoch": 30000
    }
  },
  "nonce": "0x0",
  "timestamp": "0x62884ECD",
  "gasLimit": "30000000",
  "difficulty": "0x1",
  "mixHash":  "[MIX_HASH]",
  "coinbase": "[COIN_BASE]",
  "alloc": {
[ALLOC_DATA]
    "[STAKING_CONTRACT_ADDRESS]": {
      "balance": "[STAKING_CONTRACT_BALANCE]",
      "code": "0x[STAKING_CONTRACT_CODE]",
      "storage": {
[STAKING_STORAGE_DATA]
      }
    },
    "[CONVERSION_CONTRACT_ADDRESS]": {
      "balance": "[CONVERSION_CONTRACT_BALANCE]",
      "code": "0x[CONVERSION_CONTRACT_CODE]",
      "storage": {
[CONVERSION_STORAGE_DATA]
      }
    }    
  }
}
`

func main() {
	gTemplate := GENESIS_TEMPLATE

	fmt.Println("Reading genesis data file genesis-template.json")
	genesisData, err := ReadGenesisJsonDataFile(genesisPath + "genesis-data.json")
	if err != nil {
		fmt.Println("error", err)
		return
	}
	fmt.Println("Creating genesis file genesis.json")
	var validatorList []string
	var alloc string
	var mixHash = "0x"
	var coinbase = "0x"
	var storageData = ""
	var snapshotStorageData = ""

	mixHash = mixHash + PadHexToHashSize("0")
	coinbase = coinbase + PadHexToAddressSize("0")

	//Validators
	for _, value := range genesisData.Depositors {
		validatorList = append(validatorList, value.ValidatorAddress)
	}

	//alloc list
	for j := 0; j < len(genesisData.Alloc); j++ {
		var allocRow string
		if genesisData.Alloc[j].Address[0:2] != "0x" {
			allocRow = strings.Replace(ALLOC_TEMPLATE, "0x[ADDRESS]", genesisData.Alloc[j].Address, 1)
		} else {
			allocRow = strings.Replace(ALLOC_TEMPLATE, "[ADDRESS]", genesisData.Alloc[j].Address, 1)
		}
		allocRow = strings.Replace(allocRow, "[BALANCE]", genesisData.Alloc[j].Balance, 1)

		alloc = alloc + allocRow
	}

	storageIndexHex := PadHexToHashSize("0")
	valIndexKey := crypto.Keccak256Hash(common.Hex2Bytes(storageIndexHex))

	//Validator Count
	var vIntlen = int64(len(validatorList))
	var storageIndexHexValue = "0x" + PadHexToHashSize(strconv.FormatInt(vIntlen, 10))
	var storageRow string
	storageRow = strings.Replace(STORAGE_TEMPLATE, STORAGE_KEY, "0x"+storageIndexHex, 1)
	storageRow = strings.Replace(storageRow, STORAGE_VALUE, storageIndexHexValue, 1)
	storageData = storageData + storageRow

	//Validator List
	for i := int64(0); i < int64(len(validatorList)); i++ {
		startIndexBigInt := new(big.Int)
		startIndexBigInt.SetString(strings.Replace(valIndexKey.String(), "0x", "", 1), 16)
		startIndexBigInt = common.SafeAddBigInt(startIndexBigInt, big.NewInt(i))
		storageIndexHex = PadHexToHashSize(fmt.Sprintf("%x", startIndexBigInt))

		var s string
		s = storageIndexHex
		if storageIndexHex[0:2] != "0x" {
			s = "0x" + storageIndexHex
		}

		var v string
		if validatorList[i][0:2] != "0x" {
			v = "0x" + PadHexToHashSize(validatorList[i])
		} else {
			v = "0x" + PadHexToHashSize(validatorList[i][2:])
		}

		storageRow = strings.Replace(STORAGE_TEMPLATE, STORAGE_KEY, s, 1)
		storageRow = strings.Replace(storageRow, STORAGE_VALUE, v, 1)
		storageData = storageData + storageRow
	}

	//Depositor Balance
	storageIndexHex = PadHexToHashSize("1")
	for i := int64(0); i < int64(len(genesisData.Depositors)); i++ {
		var storageRow string

		key := strings.Replace(genesisData.Depositors[i].Address, "0x", "", 1) + storageIndexHex
		hexKey, err := hex.DecodeString(key)
		if err != nil {
			fmt.Println("error", err)
			return
		}

		s := crypto.Keccak256Hash(hexKey).String()
		if s[0:2] != "0x" {
			s = "0x" + s
		}

		var b string
		if genesisData.Depositors[i].Balance[0:2] != "0x" {
			b = "0x" + PadHexToHashSize(genesisData.Depositors[i].Balance)
		} else {
			b = "0x" + PadHexToHashSize(genesisData.Depositors[i].Balance[2:])
		}

		storageRow = strings.Replace(STORAGE_TEMPLATE, STORAGE_KEY, s, 1)
		storageRow = strings.Replace(storageRow, STORAGE_VALUE, b, 1)
		storageData = storageData + storageRow
	}

	//Total Deposited Amount
	t := PadHexToHashSize("2")
	var td string
	if genesisData.TotalDepositBalance[0:2] != "0x" {
		td = "0x" + PadHexToHashSize(genesisData.TotalDepositBalance)
	} else {
		td = "0x" + PadHexToHashSize(genesisData.TotalDepositBalance[2:])
	}
	storageRow = strings.Replace(STORAGE_TEMPLATE, STORAGE_KEY, "0x"+t, 1)
	storageRow = strings.Replace(storageRow, STORAGE_VALUE, td, 1)
	storageData = storageData + storageRow

	//Depositor Count
	d := PadHexToHashSize("3")
	var dc string
	dc = "0x" + PadHexToHashSize(strconv.FormatInt(int64(len(genesisData.Depositors)), 10))
	storageRow = strings.Replace(STORAGE_TEMPLATE, STORAGE_KEY, "0x"+d, 1)
	storageRow = strings.Replace(storageRow, STORAGE_VALUE, dc, 1)
	storageData = storageData + storageRow

	fmt.Println("Validator To Depositor Mapping")
	storageIndexHex = PadHexToHashSize("8")
	for i := int64(0); i < int64(len(genesisData.Depositors)); i++ {
		key := strings.Replace(genesisData.Depositors[i].ValidatorAddress, "0x", "", 1) + storageIndexHex
		hexKey, err := hex.DecodeString(key)
		if err != nil {
			fmt.Println("error", err)
			return
		}

		s := crypto.Keccak256Hash(hexKey).String()
		if s[0:2] != "0x" {
			s = "0x" + s
		}

		var b string
		if genesisData.Depositors[i].Address[0:2] != "0x" {
			b = "0x" + PadHexToHashSize(genesisData.Depositors[i].Address)
		} else {
			b = "0x" + PadHexToHashSize(genesisData.Depositors[i].Address[2:])
		}

		storageRow = strings.Replace(STORAGE_TEMPLATE, STORAGE_KEY, s, 1)
		storageRow = strings.Replace(storageRow, STORAGE_VALUE, b, 1)
		storageData = storageData + storageRow
	}
	//Remove trailing comma
	commaIndex := strings.LastIndex(storageData, ",")
	storageData = storageData[0:commaIndex]

	//Snapshot Balance
	storageIndexHex = PadHexToHashSize("0")
	for i := int64(0); i < int64(len(genesisData.Snapshot)); i++ {
		var storageRow string

		key := strings.Replace(genesisData.Snapshot[i].EthAddress, "0x", "", 1) + storageIndexHex
		hexKey, err := hex.DecodeString(key)
		if err != nil {
			fmt.Println("error", err)
			return
		}

		s := crypto.Keccak256Hash(hexKey).String()
		if s[0:2] != "0x" {
			s = "0x" + s
		}

		var b string
		if genesisData.Snapshot[i].Balance[0:2] != "0x" {
			b = "0x" + PadHexToHashSize(genesisData.Snapshot[i].Balance)
		} else {
			b = "0x" + PadHexToHashSize(genesisData.Snapshot[i].Balance[2:])
		}

		storageRow = strings.Replace(STORAGE_TEMPLATE, STORAGE_KEY, s, 1)
		storageRow = strings.Replace(storageRow, STORAGE_VALUE, b, 1)
		snapshotStorageData = snapshotStorageData + storageRow
	}
	//Remove trailing comma
	commaIndex = strings.LastIndex(snapshotStorageData, ",")
	snapshotStorageData = snapshotStorageData[0:commaIndex]

	//template update
	gTemplate = strings.Replace(gTemplate, "[MIX_HASH]", mixHash, 1)
	gTemplate = strings.Replace(gTemplate, "[COIN_BASE]", coinbase, 1)
	gTemplate = strings.Replace(gTemplate, "[ALLOC_DATA]", alloc, 1)
	gTemplate = strings.Replace(gTemplate, "[STAKING_CONTRACT_ADDRESS]", staking.STAKING_CONTRACT, 1)
	gTemplate = strings.Replace(gTemplate, "[STAKING_CONTRACT_BALANCE]", genesisData.TotalDepositBalance, 1)
	gTemplate = strings.Replace(gTemplate, "[STAKING_CONTRACT_CODE]", staking.STAKING_RUNTIME_BIN, 1)
	gTemplate = strings.Replace(gTemplate, "[STAKING_STORAGE_DATA]", storageData, 1)
	gTemplate = strings.Replace(gTemplate, "[CONVERSION_CONTRACT_ADDRESS]", conversion.CONVERSION_CONTRACT, 1)
	gTemplate = strings.Replace(gTemplate, "[CONVERSION_CONTRACT_BALANCE]", genesisData.ConversionContractBalance, 1)
	gTemplate = strings.Replace(gTemplate, "[CONVERSION_CONTRACT_CODE]", conversion.CONVERSION_RUNTIME_BIN, 1)
	gTemplate = strings.Replace(gTemplate, "[CONVERSION_STORAGE_DATA]", snapshotStorageData, 1)

	// write the whole body at once
	err = ioutil.WriteFile(genesisPath+"genesis.json", []byte(gTemplate), 0644)
	if err != nil {
		panic(err)
	}

	fmt.Print(gTemplate)
	fmt.Println("Wrote to file successfully")
}

func PadHexToHashSize(hex string) string {
	pad := ""
	for i := len(hex); i < common.HashLength*2; i++ {
		pad = pad + "0"
	}
	return pad + hex
}

func PadHexToAddressSize(hex string) string {
	pad := ""
	for i := len(hex); i < common.AddressLength*2; i++ {
		pad = pad + "0"
	}
	return pad + hex
}

func ReadGenesisJsonDataFile(filename string) (genesisData GenesisData, err error) {

	fileContent, err := os.Open(filename)

	if err != nil {
		return GenesisData{}, err
	}

	defer fileContent.Close()

	byteResult, err := ioutil.ReadAll(fileContent)
	if err != nil {
		return GenesisData{}, err
	}
	json.Unmarshal(byteResult, &genesisData)

	return genesisData, nil
}
