package systemcontracts

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/core/state"
	"github.com/DogeProtocol/dp/log"
	"github.com/DogeProtocol/dp/params"
)

type UpgradeConfig struct {
	BeforeUpgrade upgradeHook
	AfterUpgrade  upgradeHook
	ContractAddr  common.Address
	CommitUrl     string
	Code          string
}

type Upgrade struct {
	UpgradeName string
	Configs     []*UpgradeConfig
}

type upgradeHook func(blockNumber *big.Int, contractAddr common.Address, statedb *state.StateDB) error

var (
// GenesisHash common.Hash
)

func init() {
	// reserved for future use to instantiate Upgrade vars
}

func UpgradeBuildInSystemContract(config *params.ChainConfig, blockNumber *big.Int, statedb *state.StateDB) error {
	if config == nil || blockNumber == nil || statedb == nil {
		return nil
	}
	return nil
}

func collectContracts(config *params.ChainConfig) ([]*UpgradeConfig, error) {
	contracts := GetContracts()
	if len(contracts) == 0 {
		return nil, errors.New("Missing systemContracts in  config for Bombay fork")
	}

	upgrades := make([]*UpgradeConfig, len(contracts))
	for i, contract := range contracts {
		c := GetContract_Data(contract)
		upgrades[i] = &UpgradeConfig{
			ContractAddr: c.ContractAddress,
			Code:         c.BIN,
		}
	}

	return upgrades, nil
}

func applySystemContractUpgrade(upgrade *Upgrade, blockNumber *big.Int, statedb *state.StateDB, logger log.Logger) {
	if upgrade == nil {
		logger.Info("Empty upgrade config", "height", blockNumber.String())
		return
	}

	logger.Info(fmt.Sprintf("Applying upgrade %s at height %d", upgrade.UpgradeName, blockNumber.Int64()))
	for _, cfg := range upgrade.Configs {

		logger.Info(fmt.Sprintf("Upgrade contract %s to commit %s", cfg.ContractAddr.String(), cfg.CommitUrl))

		if cfg.BeforeUpgrade != nil {
			err := cfg.BeforeUpgrade(blockNumber, cfg.ContractAddr, statedb)
			if err != nil {
				panic(fmt.Errorf("contract address: %s, execute beforeUpgrade error: %s", cfg.ContractAddr.String(), err.Error()))
			}
		}

		newContractCode, err := hex.DecodeString(strings.TrimPrefix(cfg.Code, "0x"))
		fmt.Println("applySystemContractUpgrade : ", cfg.ContractAddr.String(), newContractCode)
		if err != nil {
			panic(fmt.Errorf("failed to decode new contract code: %s", err.Error()))
		}
		statedb.SetCode(cfg.ContractAddr, newContractCode)

		if cfg.AfterUpgrade != nil {
			err := cfg.AfterUpgrade(blockNumber, cfg.ContractAddr, statedb)
			if err != nil {
				panic(fmt.Errorf("contract address: %s, execute afterUpgrade error: %s", cfg.ContractAddr.String(), err.Error()))
			}
		}
	}
}
