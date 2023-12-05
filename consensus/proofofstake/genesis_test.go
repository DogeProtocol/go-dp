package proofofstake

import (
	"encoding/hex"
	"fmt"
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/crypto"
	"math/big"
	"strings"
	"testing"
)

func PadHexToHashSize(hex string) string {
	pad := ""
	for i := len(hex); i < common.HashLength*2; i++ {
		pad = pad + "0"
	}
	return pad + hex
}

func TestGenesisJson(t *testing.T) {
	validatorList := []string{
		"0x57079c7528d6c322918ef03dff4a74bfd67be0c78cc2e67c2afc85cd806538e9",
		"0x81f3416867cd2ce0f4fda114454a5ada656e26e3427e32edf451ab090aadbf0e",
		"0x726803cf94b3c990d69dab11f36ac8bbcd1afedf86a33484105e77195365f966",
		"0xc7a5c9d61560f919f39bb2c6e9e9a0aeb9bf5ee4e34372967f74d83c7749f361"}

	depositorList := []string{
		"0x150c94b2ae124fe0edfb4b487dd3c815d81f29f7c207bb311c4793dfa6afe933",
		"0x21f707eee6f7a7377d7d02f7adeb5fd23b0617a0385f4bae466d8a03e9ae078d",
		"0x866f06de82780554fafd6677a1e6b7b6fd90efb7abeee04cdd7e7f3ec4c9724d",
		"0x895fdb851f1ba60ec72088b5a0a46a468a6cb7c5c80b5e471b159b6b3bf6fbbe"}

	if len(validatorList) != len(depositorList) {
		t.Fatalf("failed")
	}

	storageIndexHex1 := "0000000000000000000000000000000000000000000000000000000000000000"

	valIndexKey := crypto.Keccak256Hash(common.Hex2Bytes(storageIndexHex1))

	fmt.Println("Validator List")
	for i := int64(0); i < int64(len(validatorList)); i++ {
		startIndexBigInt := new(big.Int)
		startIndexBigInt.SetString(strings.Replace(valIndexKey.String(), "0x", "", 1), 16)
		startIndexBigInt = common.SafeAddBigInt(startIndexBigInt, big.NewInt(i))
		storageIndexHex1 = PadHexToHashSize(fmt.Sprintf("%x", startIndexBigInt))
		fmt.Println("storageIndexHex", storageIndexHex1, "validator", validatorList[i])
	}

	fmt.Println("Depositor Balance")
	storageIndexHex2 := "0000000000000000000000000000000000000000000000000000000000000001"
	for i := int64(0); i < int64(len(depositorList)); i++ {
		key := strings.Replace(depositorList[i], "0x", "", 1) + storageIndexHex2
		hexKey, err := hex.DecodeString(key)
		if err != nil {
			fmt.Println("error", err)
			t.Fatalf("failed")
		}
		fmt.Println("storageIndexHex", crypto.Keccak256Hash(hexKey))
	}

	fmt.Println("Validator To Depositor Mapping")
	storageIndexHex3 := "0000000000000000000000000000000000000000000000000000000000000008"
	for i := int64(0); i < int64(len(validatorList)); i++ {
		key := strings.Replace(validatorList[i], "0x", "", 1) + storageIndexHex3
		hexKey, err := hex.DecodeString(key)
		if err != nil {
			fmt.Println("error", err)
			t.Fatalf("failed")
		}
		fmt.Println("storageIndexHex", crypto.Keccak256Hash(hexKey), "depositor", depositorList[i])
	}
}
