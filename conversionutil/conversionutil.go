package conversionutil

import (
	"bytes"
	"errors"
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/core/types"
	"github.com/DogeProtocol/dp/crypto/crosssign"
	"github.com/DogeProtocol/dp/log"
	"github.com/DogeProtocol/dp/systemcontracts/conversion"
	"strings"
)

var snapshotMap = map[string]bool{
	strings.ToLower("0xda02553C0D68A251F58024c23E76c99e48315FcC"): true,
	strings.ToLower("0x4609545aA34Ad61d5b19dEB1f019Ba8674c6d8De"): true,
	strings.ToLower("0x62De16972E1C779e9EEA3C53A6c0115C9686032f"): true,
	strings.ToLower("0xd65A4B7526Df5756356009DDf0E89F325506c6d4"): true,
}

var FirstPart = []byte{25, 71, 244, 125, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 160, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 42}
var SecondPart = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 132}

func IsGasExemptTxn(tx *types.Transaction, signer types.Signer) (bool, error) {
	msg, err := tx.AsMessage(signer)

	if err != nil {
		log.Trace("IsGasExemptTxn")
		return false, err
	}

	if msg.To().IsEqualTo(conversion.CONVERSION_CONTRACT_ADDRESS) == false {
		return false, nil
	}

	ethAddress, err := VerifyDataAndGetEthereumAddress(msg.From(), tx.Data())
	if err != nil {
		log.Trace("IsGasExemptTxn VerifyDataAndGetEthereumAddress failed", "err", err)
		return false, err
	}

	_, ok := snapshotMap[strings.ToLower(ethAddress)]

	// If the key exists
	if ok == false {
		log.Trace("IsGasExemptTxn address not in snapshot", "ethAddress", ethAddress)
		return false, errors.New("unidentified eth address")
	}

	log.Trace("Is a GasExemptTxn", "ethAddress", ethAddress)
	//todo: verify if already converted, to prevented flood

	return true, nil
}

func VerifyDataAndGetEthereumAddress(quantumAddress common.Address, data []byte) (string, error) {
	if data == nil {
		return "", errors.New("data is nil")
	}

	if len(data) != 356 {
		return "", errors.New("unexpected data length")
	}

	if bytes.Compare(data[0:100], FirstPart) != 0 {
		return "", errors.New("error parsing data a")
	}

	ethAddress := string(data[100:142])

	second := data[143:196]
	if bytes.Compare(second, SecondPart) != 0 {
		return "", errors.New("error parsing data b")
	}

	ethSignature := string(data[196:328])

	crossSignDetails := &crosssign.ConversionSignDetails{
		EthAddress:        strings.ToLower(ethAddress),
		EthereumSignature: ethSignature,
		QuantumAddress:    strings.ToLower(quantumAddress.Hex()),
	}
	_, err := crosssign.VerifyConversion(crossSignDetails)
	if err != nil {
		return "", err
	}

	if allZero(data[328:356]) == false {
		return "", errors.New("error parsing data c")
	}

	return ethAddress, nil
}

func allZero(b []byte) bool {
	for _, byte := range b {
		if byte != 0 {
			return false
		}
	}
	return true
}
