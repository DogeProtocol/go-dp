// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"errors"
	"fmt"
	"github.com/DogeProtocol/dp/backupmanager"
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/consensus"
	"github.com/DogeProtocol/dp/conversionutil"
	"github.com/DogeProtocol/dp/core/state"
	"github.com/DogeProtocol/dp/core/types"
	"github.com/DogeProtocol/dp/core/vm"
	"github.com/DogeProtocol/dp/crypto"
	"github.com/DogeProtocol/dp/log"
	"github.com/DogeProtocol/dp/params"
	"github.com/DogeProtocol/dp/trie"
	"math/big"
)

// StateProcessor is a basic Processor, which takes care of transitioning
// state from one point to another.
//
// StateProcessor implements Processor.
type StateProcessor struct {
	config *params.ChainConfig // Chain configuration options
	bc     *BlockChain         // Canonical block chain
	engine consensus.Engine    // Consensus engine used for block rewards
}

// NewStateProcessor initialises a new StateProcessor.
func NewStateProcessor(config *params.ChainConfig, bc *BlockChain, engine consensus.Engine) *StateProcessor {
	return &StateProcessor{
		config: config,
		bc:     bc,
		engine: engine,
	}
}

// Process processes the state changes according to the Ethereum rules by running
// the transaction messages using the statedb and applying any rewards to both
// the processor (coinbase) and any included uncles.
//
// Process returns the receipts and logs accumulated during the process and
// returns the amount of gas that was used in the process. If any of the
// transactions failed to execute due to insufficient gas it will return an error.
func (p *StateProcessor) Process(block *types.Block, statedb *state.StateDB, cfg vm.Config) (types.Receipts, []*types.Log, uint64, error) {
	var (
		receipts types.Receipts
		usedGas  = new(uint64)
		header   = block.Header()
		//blockHash   = block.Hash()
		//blockNumber = block.Number()
		allLogs []*types.Log
		gp      = new(GasPool).AddGas(block.GasLimit())
	)
	// Mutate the block and state according to any hard-fork specs
	/*if p.config.DAOForkSupport && p.config.DAOForkBlock != nil && p.config.DAOForkBlock.Cmp(block.Number()) == 0 {
		log.Trace("Process ApplyDAOHardFork")
		misc.ApplyDAOHardFork(statedb)
	}*/
	//blockContext := NewEVMBlockContext(header, p.bc, nil)

	// Iterate over and process the individual transactions
	for i, tx := range block.Transactions() {
		signer := types.MakeSigner(p.config, header.Number)
		msg, err := tx.AsMessage(signer)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("could not apply tx %d [%v]: %w", i, tx.Hash().Hex(), err)
		}

		vmConfig := cfg
		isGasExemptTxn, err := conversionutil.IsGasExemptTxn(tx, signer)
		if err == nil && isGasExemptTxn {
			vmConfig = *cfg.DeepCopy()
			vmConfig.OverrideGasFailure = true
			msg.OverrideGasPrice(big.NewInt(0))
			log.Trace("Process OverrideGasPrice", "txn", tx.Hash(), "price", msg.GasPrice())
		}

		//vmenv := vm.NewEVM(blockContext, vm.TxContext{}, statedb, p.config, vmConfig)

		statedb.Prepare(tx.Hash(), i)

		zeroAddress := common.ZERO_ADDRESS
		receipt, err := ApplyTransaction(p.config, p.bc, &zeroAddress, gp, statedb, block.Header(), tx, usedGas, cfg, isGasExemptTxn)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("could not apply tx %d [%v]: %w", i, tx.Hash().Hex(), err)
		}

		/*receipt, err := applyTransaction(msg, p.config, p.bc, nil, gp, statedb, blockNumber, blockHash, tx, usedGas, vmenv)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("could not apply tx %d [%v]: %w", i, tx.Hash().Hex(), err)
		}*/
		receipts = append(receipts, receipt)
		allLogs = append(allLogs, receipt.Logs...)

		log.Info("Tx", "hash", tx.Hash(), "receipt poststate", common.Bytes2Hex(receipt.PostState), "tx data", common.Bytes2Hex(tx.Data()), "usedGas", &usedGas, "from", msg.From().Hex(), "nonce", tx.Nonce(), "GasUsed", receipt.GasUsed,
			"CumulativeGasUsed", receipt.CumulativeGasUsed, "Type", receipt.Type,
			"status", receipt.Status, "Bloom", receipt.Bloom.Bytes(), "ContractAddress", receipt.ContractAddress, "TxHash", receipt.TxHash.Bytes(),
			"to", tx.To(), "value", tx.Value().String(), "msg", msg.AccessList(), "receipt", receipt.Logs, "timestamp", block.Time(), "Difficulty", block.Difficulty(), "NumberU64", block.NumberU64())

		printTransactionReceipt(*block, receipt, &signer, tx, uint64(i))
	}
	// Finalize the block, applying any consensus engine specific extras (e.g. block rewards)
	p.engine.Finalize(p.bc, header, statedb, block.Transactions())

	backupManager := backupmanager.GetInstance()
	if backupManager != nil {
		err := backupManager.BackupBlock(block)
		if err != nil {
			return nil, nil, 0, err
		}
	}

	receiptHashTestLocal := types.DeriveSha(receipts, trie.NewStackTrie(nil))
	log.Info("Process", "receiptHashTestLocal", common.Bytes2Hex(receiptHashTestLocal.Bytes()), "len receipts", len(receipts))

	panic("done")

	return receipts, allLogs, *usedGas, nil
}

func applyTransaction(msg types.Message, config *params.ChainConfig, bc ChainContext, author *common.Address, gp *GasPool, statedb *state.StateDB, blockNumber *big.Int,
	blockHash common.Hash, tx *types.Transaction, usedGas *uint64, evm *vm.EVM) (*types.Receipt, error) {
	// Create a new context to be used in the EVM environment.
	txContext := NewEVMTxContext(msg)
	evm.Reset(txContext, statedb)

	// Apply the transaction to the current state (included in the env).
	result, err := ApplyMessage(evm, msg, gp)
	if err != nil {
		return nil, err
	}

	// Update the state with pending changes.
	var root []byte
	if config.IsByzantium(blockNumber) {
		log.Info("applyTransaction IsByzantium yes", "blockNumber", blockNumber)
		statedb.Finalise(true)
	} else {
		root = statedb.IntermediateRoot(config.IsEIP158(blockNumber)).Bytes()
		log.Info("applyTransaction IsByzantium no", "stateroot", common.Bytes2Hex(root), "from", msg.From().Hex())
	}

	*usedGas += result.UsedGas

	// Create a new receipt for the transaction, storing the intermediate root and gas used
	// by the tx.
	receipt := &types.Receipt{Type: tx.Type(), PostState: root, CumulativeGasUsed: *usedGas}
	if result.Failed() {
		receipt.Status = types.ReceiptStatusFailed
	} else {
		receipt.Status = types.ReceiptStatusSuccessful
	}
	receipt.TxHash = tx.Hash()
	receipt.GasUsed = result.UsedGas

	// If the transaction created a contract, store the creation address in the receipt.
	if msg.To() == nil {
		receipt.ContractAddress = crypto.CreateAddress(evm.TxContext.Origin, tx.Nonce())
	}

	// Set the receipt logs and create the bloom filter.
	receipt.Logs = statedb.GetLogs(tx.Hash(), blockHash)
	receipt.Bloom = types.CreateBloom(types.Receipts{receipt})
	receipt.BlockHash = blockHash
	receipt.BlockNumber = blockNumber
	receipt.TransactionIndex = uint(statedb.TxIndex())

	l := receipt.Logs[0]
	l.Data = common.Hex2Bytes("0x00000000000000000000000000000000000000000005ca4ec2a79a7f6700000000000000000000000000000000000000000000000000000000000000000676c0000000000000000000000000000000000000000000000000000000006650c7a8")
	receipt.Logs[0] = l

	return receipt, err
}

// ApplyTransaction attempts to apply a transaction to the given state database
// and uses the input parameters for its environment. It returns the receipt
// for the transaction, gas used and an error if the transaction failed,
// indicating the block was invalid.
func ApplyTransaction(config *params.ChainConfig, bc ChainContext, author *common.Address, gp *GasPool, statedb *state.StateDB, header *types.Header, tx *types.Transaction,
	usedGas *uint64, cfg vm.Config, isGasExemptTxn bool) (*types.Receipt, error) {
	if bc == nil {
		return nil, errors.New("ChainContext is nil")
	}

	msg, err := tx.AsMessage(types.MakeSigner(config, header.Number))
	if err != nil {
		return nil, err
	}

	vmConfig := cfg
	if isGasExemptTxn {
		vmConfig = *cfg.DeepCopy()
		vmConfig.OverrideGasFailure = true
		msg.OverrideGasPrice(big.NewInt(0))
		log.Trace("ApplyTransaction OverrideGasPrice", "txn", tx.Hash(), "price", msg.GasPrice())
	}

	// Create a new context to be used in the EVM environment
	blockContext := NewEVMBlockContext(header, bc, nil)
	vmenv := vm.NewEVM(blockContext, vm.TxContext{}, statedb, config, vmConfig)
	return applyTransaction(msg, config, bc, nil, gp, statedb, header.Number, header.Hash(), tx, usedGas, vmenv)
}

func printTransactionReceipt(block types.Block, receipt *types.Receipt, signer *types.Signer, tx *types.Transaction, txIndex uint64) {
	from, _ := types.Sender(*signer, tx)

	receiptStr := ""

	receiptStr = receiptStr + "\nblockHash: " + block.Hash().Hex()
	receiptStr = receiptStr + "\nblockNumber: " + string(block.NumberU64())
	receiptStr = receiptStr + "\ntransactionHash: " + tx.Hash().Hex()
	receiptStr = receiptStr + "\ntransactionIndex: " + string(txIndex)
	receiptStr = receiptStr + "\nfrom: " + from.Hex()
	receiptStr = receiptStr + "\nto: " + tx.To().Hex()
	receiptStr = receiptStr + "\ngasUsed: " + string(receipt.GasUsed)
	receiptStr = receiptStr + "\ncumulativeGasUsed: " + string(receipt.CumulativeGasUsed)
	if receipt.Logs != nil {
		logStr := ""
		for _, log := range receipt.Logs {
			logStr = logStr + "\nData: " + common.Bytes2Hex(log.Data)
			logStr = logStr + "\nAddress: " + log.Address.Hex()
			logStr = logStr + "\nTxHash: " + log.TxHash.Hex()
			logStr = logStr + "\nTopics: "
			for _, topic := range log.Topics {
				logStr = logStr + topic.Hex() + ", "
			}
			logStr = logStr + "\nBlockHash: " + log.BlockHash.Hex()
			logStr = logStr + "\nBlockNumber: " + string(log.BlockNumber)
			logStr = logStr + "\nIndex: " + string(log.Index)
			if log.Removed {
				logStr = logStr + "\nRemoved: true"
			} else {
				logStr = logStr + "\nRemoved: false"
			}
		}
		receiptStr = receiptStr + "\nlogs: " + logStr
	}
	receiptStr = receiptStr + "\nBloom: " + common.Bytes2Hex(receipt.Bloom.Bytes())
	receiptStr = receiptStr + "\nStatus: " + string(receipt.Status)

	receiptStr = receiptStr + "\nreceipt.PostState: " + common.Bytes2Hex(receipt.PostState)

	log.Info("printTransactionReceipt", "receipt", receiptStr)
}
