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

package runtime

import (
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/DogeProtocol/dp/accounts/abi"
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/consensus"
	"github.com/DogeProtocol/dp/core"
	"github.com/DogeProtocol/dp/core/asm"
	"github.com/DogeProtocol/dp/core/rawdb"
	"github.com/DogeProtocol/dp/core/state"
	"github.com/DogeProtocol/dp/core/types"
	"github.com/DogeProtocol/dp/core/vm"
	"github.com/DogeProtocol/dp/params"
)

func TestDefaults(t *testing.T) {
	cfg := new(Config)
	setDefaults(cfg)

	if cfg.Difficulty == nil {
		t.Error("expected difficulty to be non nil")
	}

	if cfg.Time == nil {
		t.Error("expected time to be non nil")
	}
	if cfg.GasLimit == 0 {
		t.Error("didn't expect gaslimit to be zero")
	}
	if cfg.GasPrice == nil {
		t.Error("expected time to be non nil")
	}
	if cfg.Value == nil {
		t.Error("expected time to be non nil")
	}
	if cfg.GetHashFn == nil {
		t.Error("expected time to be non nil")
	}
	if cfg.BlockNumber == nil {
		t.Error("expected block number to be non nil")
	}
}

func TestEVM(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("crashed with: %v", r)
		}
	}()

	Execute([]byte{
		byte(vm.DIFFICULTY),
		byte(vm.TIMESTAMP),
		byte(vm.GASLIMIT),
		byte(vm.PUSH1),
		byte(vm.ORIGIN),
		byte(vm.BLOCKHASH),
		byte(vm.COINBASE),
	}, nil, nil)
}

func TestExecute(t *testing.T) {
	ret, _, err := Execute([]byte{
		byte(vm.PUSH1), 10,
		byte(vm.PUSH1), 0,
		byte(vm.MSTORE),
		byte(vm.PUSH1), 32,
		byte(vm.PUSH1), 0,
		byte(vm.RETURN),
	}, nil, nil)
	if err != nil {
		t.Fatal("didn't expect error", err)
	}

	num := new(big.Int).SetBytes(ret)
	if num.Cmp(big.NewInt(10)) != 0 {
		t.Error("Expected 10, got", num)
	}
}

func TestCall(t *testing.T) {
	state, _ := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
	address := common.HexToAddress("0x0a")
	state.SetCode(address, []byte{
		byte(vm.PUSH1), 10,
		byte(vm.PUSH1), 0,
		byte(vm.MSTORE),
		byte(vm.PUSH1), 32,
		byte(vm.PUSH1), 0,
		byte(vm.RETURN),
	})

	ret, _, err := Call(address, nil, &Config{State: state})
	if err != nil {
		t.Fatal("didn't expect error", err)
	}

	num := new(big.Int).SetBytes(ret)
	if num.Cmp(big.NewInt(10)) != 0 {
		t.Error("Expected 10, got", num)
	}
}

func BenchmarkCall(b *testing.B) {
	var definition = `[{"constant":true,"inputs":[],"name":"seller","outputs":[{"name":"","type":"address"}],"type":"function"},{"constant":false,"inputs":[],"name":"abort","outputs":[],"type":"function"},{"constant":true,"inputs":[],"name":"value","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":false,"inputs":[],"name":"refund","outputs":[],"type":"function"},{"constant":true,"inputs":[],"name":"buyer","outputs":[{"name":"","type":"address"}],"type":"function"},{"constant":false,"inputs":[],"name":"confirmReceived","outputs":[],"type":"function"},{"constant":true,"inputs":[],"name":"state","outputs":[{"name":"","type":"uint8"}],"type":"function"},{"constant":false,"inputs":[],"name":"confirmPurchase","outputs":[],"type":"function"},{"inputs":[],"type":"constructor"},{"anonymous":false,"inputs":[],"name":"Aborted","type":"event"},{"anonymous":false,"inputs":[],"name":"PurchaseConfirmed","type":"event"},{"anonymous":false,"inputs":[],"name":"ItemReceived","type":"event"},{"anonymous":false,"inputs":[],"name":"Refunded","type":"event"}]`

	var code = common.Hex2Bytes("6060604052361561006c5760e060020a600035046308551a53811461007457806335a063b4146100865780633fa4f245146100a6578063590e1ae3146100af5780637150d8ae146100cf57806373fac6f0146100e1578063c19d93fb146100fe578063d696069714610112575b610131610002565b610133600154600160a060020a031681565b610131600154600160a060020a0390811633919091161461015057610002565b61014660005481565b610131600154600160a060020a039081163391909116146102d557610002565b610133600254600160a060020a031681565b610131600254600160a060020a0333811691161461023757610002565b61014660025460ff60a060020a9091041681565b61013160025460009060ff60a060020a9091041681146101cc57610002565b005b600160a060020a03166060908152602090f35b6060908152602090f35b60025460009060a060020a900460ff16811461016b57610002565b600154600160a060020a03908116908290301631606082818181858883f150506002805460a060020a60ff02191660a160020a179055506040517f72c874aeff0b183a56e2b79c71b46e1aed4dee5e09862134b8821ba2fddbf8bf9250a150565b80546002023414806101dd57610002565b6002805460a060020a60ff021973ffffffffffffffffffffffffffffffffffffffff1990911633171660a060020a1790557fd5d55c8a68912e9a110618df8d5e2e83b8d83211c57a8ddd1203df92885dc881826060a15050565b60025460019060a060020a900460ff16811461025257610002565b60025460008054600160a060020a0390921691606082818181858883f150508354604051600160a060020a0391821694503090911631915082818181858883f150506002805460a060020a60ff02191660a160020a179055506040517fe89152acd703c9d8c7d28829d443260b411454d45394e7995815140c8cbcbcf79250a150565b60025460019060a060020a900460ff1681146102f057610002565b6002805460008054600160a060020a0390921692909102606082818181858883f150508354604051600160a060020a0391821694503090911631915082818181858883f150506002805460a060020a60ff02191660a160020a179055506040517f8616bbbbad963e4e65b1366f1d75dfb63f9e9704bbbf91fb01bec70849906cf79250a15056")

	abi, err := abi.JSON(strings.NewReader(definition))
	if err != nil {
		b.Fatal(err)
	}

	cpurchase, err := abi.Pack("confirmPurchase")
	if err != nil {
		b.Fatal(err)
	}
	creceived, err := abi.Pack("confirmReceived")
	if err != nil {
		b.Fatal(err)
	}
	refund, err := abi.Pack("refund")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 400; j++ {
			Execute(code, cpurchase, nil)
			Execute(code, creceived, nil)
			Execute(code, refund, nil)
		}
	}
}
func benchmarkEVM_Create(bench *testing.B, code string) {
	var (
		statedb, _ = state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
		sender     = common.BytesToAddress([]byte("sender"))
		receiver   = common.BytesToAddress([]byte("receiver"))
	)

	statedb.CreateAccount(sender)
	statedb.SetCode(receiver, common.FromHex(code))
	runtimeConfig := Config{
		Origin:      sender,
		State:       statedb,
		GasLimit:    10000000,
		Difficulty:  big.NewInt(0x200000),
		Time:        new(big.Int).SetUint64(0),
		Coinbase:    common.Address{},
		BlockNumber: new(big.Int).SetUint64(1),
		ChainConfig: &params.ChainConfig{
			ChainID:             big.NewInt(1),
			HomesteadBlock:      new(big.Int),
			ByzantiumBlock:      new(big.Int),
			ConstantinopleBlock: new(big.Int),
			DAOForkBlock:        new(big.Int),
			DAOForkSupport:      false,
			EIP150Block:         new(big.Int),
			EIP155Block:         new(big.Int),
			EIP158Block:         new(big.Int),
		},
		EVMConfig: vm.Config{},
	}
	// Warm up the intpools and stuff
	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		Call(receiver, []byte{}, &runtimeConfig)
	}
	bench.StopTimer()
}

func BenchmarkEVM_CREATE_500(bench *testing.B) {
	// initcode size 500K, repeatedly calls CREATE and then modifies the mem contents
	benchmarkEVM_Create(bench, "5b6207a120600080f0600152600056")
}
func BenchmarkEVM_CREATE2_500(bench *testing.B) {
	// initcode size 500K, repeatedly calls CREATE2 and then modifies the mem contents
	benchmarkEVM_Create(bench, "5b586207a120600080f5600152600056")
}
func BenchmarkEVM_CREATE_1200(bench *testing.B) {
	// initcode size 1200K, repeatedly calls CREATE and then modifies the mem contents
	benchmarkEVM_Create(bench, "5b62124f80600080f0600152600056")
}
func BenchmarkEVM_CREATE2_1200(bench *testing.B) {
	// initcode size 1200K, repeatedly calls CREATE2 and then modifies the mem contents
	benchmarkEVM_Create(bench, "5b5862124f80600080f5600152600056")
}

func fakeHeader(n uint64, parentHash common.Hash) *types.Header {
	header := types.Header{
		Coinbase:   common.HexToAddress("0x00000000000000000000000000000000deadbeef"),
		Number:     big.NewInt(int64(n)),
		ParentHash: parentHash,
		Time:       1000,
		Nonce:      types.BlockNonce{0x1},
		Extra:      []byte{},
		Difficulty: big.NewInt(0),
		GasLimit:   100000,
	}
	return &header
}

type dummyChain struct {
	counter int
}

// Engine retrieves the chain's consensus engine.
func (d *dummyChain) Engine() consensus.Engine {
	return nil
}

// GetHeader returns the hash corresponding to their hash.
func (d *dummyChain) GetHeader(h common.Hash, n uint64) *types.Header {
	d.counter++
	parentHash := common.Hash{}
	s := common.LeftPadBytes(big.NewInt(int64(n-1)).Bytes(), 32)
	copy(parentHash[:], s)

	//parentHash := common.Hash{byte(n - 1)}
	//fmt.Printf("GetHeader(%x, %d) => header with parent %x\n", h, n, parentHash)
	return fakeHeader(n, parentHash)
}

// TestBlockhash tests the blockhash operation. It's a bit special, since it internally
// requires access to a chain reader.
func TestBlockhash(t *testing.T) {
	// Current head
	n := uint64(1000)
	parentHash := common.Hash{}
	s := common.LeftPadBytes(big.NewInt(int64(n-1)).Bytes(), 32)
	copy(parentHash[:], s)
	header := fakeHeader(n, parentHash)

	// This is the contract we're using. It requests the blockhash for current num (should be all zeroes),
	// then iteratively fetches all blockhashes back to n-260.
	// It returns
	// 1. the first (should be zero)
	// 2. the second (should be the parent hash)
	// 3. the last non-zero hash
	// By making the chain reader return hashes which correlate to the number, we can
	// verify that it obtained the right hashes where it should

	/*

		pragma solidity ^0.5.3;
		contract Hasher{

			function test() public view returns (bytes32, bytes32, bytes32){
				uint256 x = block.number;
				bytes32 first;
				bytes32 last;
				bytes32 zero;
				zero = blockhash(x); // Should be zeroes
				first = blockhash(x-1);
				for(uint256 i = 2 ; i < 260; i++){
					bytes32 hash = blockhash(x - i);
					if (uint256(hash) != 0){
						last = hash;
					}
				}
				return (zero, first, last);
			}
		}

	*/
	// The contract above
	//data := common.Hex2Bytes("6080604052348015600f57600080fd5b50600436106045576000357c010000000000000000000000000000000000000000000000000000000090048063f8a8fd6d14604a575b600080fd5b60506074565b60405180848152602001838152602001828152602001935050505060405180910390f35b600080600080439050600080600083409050600184034092506000600290505b61010481101560c35760008186034090506000816001900414151560b6578093505b5080806001019150506094565b508083839650965096505050505090919256fea165627a7a72305820462d71b510c1725ff35946c20b415b0d50b468ea157c8c77dff9466c9cb85f560029")
	data := common.Hex2Bytes("6080604052600436106100e75760003560e01c8063adeb73d91161008a578063e1c2f06711610059578063e1c2f067146102c7578063e324684714610304578063e9ba669214610341578063f4125c2a1461037e576100e7565b8063adeb73d9146101f7578063b8ffa78714610222578063bae4e56a1461024d578063d0c7a7911461028a576100e7565b8063461d8863116100c6578063461d88631461016b578063699a0be1146101825780639a0a200d1461019e5780639f9347ec146101ba576100e7565b80624cb7ae146100ec5780631a8ab4681461010357806341941b181461012e575b600080fd5b3480156100f857600080fd5b506101016103bb565b005b34801561010f57600080fd5b5061011861054b565b6040516101259190611b66565b60405180910390f35b34801561013a57600080fd5b506101556004803603810190610150919061146d565b6105a3565b6040516101629190611b4b565b60405180910390f35b34801561017757600080fd5b506101806105c5565b005b61019c6004803603810190610197919061146d565b610823565b005b6101b860048036038101906101b3919061146d565b610c7a565b005b3480156101c657600080fd5b506101e160048036038101906101dc919061146d565b6110da565b6040516101ee9190611de8565b60405180910390f35b34801561020357600080fd5b5061020c6110f7565b6040516102199190611de8565b60405180910390f35b34801561022e57600080fd5b50610237611101565b6040516102449190611de8565b60405180910390f35b34801561025957600080fd5b50610274600480360381019061026f919061146d565b61110b565b6040516102819190611de8565b60405180910390f35b34801561029657600080fd5b506102b160048036038101906102ac919061146d565b611128565b6040516102be9190611b4b565b60405180910390f35b3480156102d357600080fd5b506102ee60048036038101906102e99190611496565b61114a565b6040516102fb9190611de8565b60405180910390f35b34801561031057600080fd5b5061032b6004803603810190610326919061146d565b61121e565b6040516103389190611de8565b60405180910390f35b34801561034d57600080fd5b5061036860048036038101906103639190611496565b61123b565b6040516103759190611de8565b60405180910390f35b34801561038a57600080fd5b506103a560048036038101906103a0919061146d565b61130f565b6040516103b29190611de8565b60405180910390f35b6000339050600115156005600083815260200190815260200160002060009054906101000a900460ff16151514610426576040517e31cb2800000000000000000000000000000000000000000000000000000000815260040161041d90611d28565b60405180910390fd5b6000600c6000838152602001908152602001600020541461047b576040517e31cb2800000000000000000000000000000000000000000000000000000000815260040161047290611d88565b60405180910390fd5b60006001600083815260200190815260200160002054116104d0576040517e31cb280000000000000000000000000000000000000000000000000000000081526004016104c790611dc8565b60405180910390fd5b6203d0904301600c600083815260200190815260200160002081905550600015156005600083815260200190815260200160002060009054906101000a905050507f7dc4fb4a5597fc9380af932772d767d2eb12bc71f0c2855cbe5716f81351aa49816040516105409190611b4b565b60405180910390a150565b6060600080548060200260200160405190810160405280929190818152602001828054801561059957602002820191906000526020600020905b815481526020019060010190808311610585575b5050505050905090565b6000806008600084815260200190815260200160002054905080915050919050565b60003390506000600c6000838152602001908152602001600020541161061f576040517e31cb2800000000000000000000000000000000000000000000000000000000815260040161061690611d48565b60405180910390fd5b600c6000828152602001908152602001600020544311610673576040517e31cb2800000000000000000000000000000000000000000000000000000000815260040161066a90611ba8565b60405180910390fd5b60003063f4125c2a836040518263ffffffff1660e01b81526004016106989190611b4b565b60206040518083038186803b1580156106b057600080fd5b505afa1580156106c4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106e891906114d2565b90506001600083815260200190815260200160002060009055600b600083815260200190815260200160002060009055600a6000838152602001908152602001600020600090556005600083815260200190815260200160002060006101000a81549060ff02191690556000828260405161076290611b36565b60006040518083038185875af1925050503d806000811461079f576040519150601f19603f3d011682016040523d82523d6000602084013e6107a4565b606091505b50509050806107e7576040517e31cb280000000000000000000000000000000000000000000000000000000081526004016107de90611ce8565b60405180910390fd5b7f0840e10e73b227adb5caa9bb15f74177aaf76925099437c14f96d58e7e3eef74836040516108169190611b4b565b60405180910390a1505050565b600033905060003490506acecb8f27f4200f3a00000081101561087a576040517e31cb2800000000000000000000000000000000000000000000000000000000815260040161087190611be8565b60405180910390fd5b828214156108bc576040517e31cb280000000000000000000000000000000000000000000000000000000081526004016108b390611c68565b60405180910390fd5b60008314156108ff576040517e31cb280000000000000000000000000000000000000000000000000000000081526004016108f690611c08565b60405180910390fd5b600015156004600085815260200190815260200160002060009054906101000a900460ff16151514610965576040517e31cb2800000000000000000000000000000000000000000000000000000000815260040161095c90611da8565b60405180910390fd5b600015156006600085815260200190815260200160002060009054906101000a900460ff161515146109cb576040517e31cb280000000000000000000000000000000000000000000000000000000081526004016109c290611d68565b60405180910390fd5b60008331905060008114610a13576040517e31cb28000000000000000000000000000000000000000000000000000000008152600401610a0a90611cc8565b60405180910390fd5b600015156005600085815260200190815260200160002060009054906101000a900460ff16151514610a79576040517e31cb28000000000000000000000000000000000000000000000000000000008152600401610a7090611d08565b60405180910390fd5b600015156007600085815260200190815260200160002060009054906101000a900460ff16151514610adf576040517e31cb28000000000000000000000000000000000000000000000000000000008152600401610ad690611bc8565b60405180910390fd5b6000849080600181540180825580915050600190039060005260206000200160009091909190915055610b1d826002546113fb90919063ffffffff16565b600281905550610b3960016003546113fb90919063ffffffff16565b60038190555081600160008581526020019081526020016000208190555060016004600086815260200190815260200160002060006101000a81548160ff02191690831515021790555060016005600085815260200190815260200160002060006101000a81548160ff02191690831515021790555060016006600086815260200190815260200160002060006101000a81548160ff02191690831515021790555060016007600085815260200190815260200160002060006101000a81548160ff02191690831515021790555082600860008681526020019081526020016000208190555083600960008581526020019081526020016000208190555083837ffd869cc875fcc5dc89490a820f180ac603fd56f6a12f20a2d681a76947ef898d844342604051610c6c93929190611e03565b60405180910390a350505050565b600015156004600083815260200190815260200160002060009054906101000a900460ff16151514610ce0576040517e31cb28000000000000000000000000000000000000000000000000000000008152600401610cd790611da8565b60405180910390fd5b600015156005600083815260200190815260200160002060009054906101000a900460ff16151514610d46576040517e31cb28000000000000000000000000000000000000000000000000000000008152600401610d3d90611ca8565b60405180910390fd5b600015156006600083815260200190815260200160002060009054906101000a900460ff16151514610dac576040517e31cb28000000000000000000000000000000000000000000000000000000008152600401610da390611b88565b60405180910390fd5b600015156007600083815260200190815260200160002060009054906101000a900460ff16151514610e12576040517e31cb28000000000000000000000000000000000000000000000000000000008152600401610e0990611c48565b60405180910390fd5b6000813114610e55576040517e31cb28000000000000000000000000000000000000000000000000000000008152600401610e4c90611cc8565b60405180910390fd5b6000811415610e98576040517e31cb28000000000000000000000000000000000000000000000000000000008152600401610e8f90611c08565b60405180910390fd5b600033905081811415610edf576040517e31cb28000000000000000000000000000000000000000000000000000000008152600401610ed690611c68565b60405180910390fd5b600115156005600083815260200190815260200160002060009054906101000a900460ff16151514610f45576040517e31cb28000000000000000000000000000000000000000000000000000000008152600401610f3c90611d28565b60405180910390fd5b6000600c60008381526020019081526020016000205414610f9a576040517e31cb28000000000000000000000000000000000000000000000000000000008152600401610f9190611c88565b60405180910390fd5b60016004600084815260200190815260200160002060006101000a81548160ff02191690831515021790555060016006600084815260200190815260200160002060006101000a81548160ff021916908315150217905550806008600084815260200190815260200160002081905550816009600083815260200190815260200160002081905550600082908060018154018082558091505060019003906000526020600020016000909190919091505560006009600083815260200190815260200160002054905060006004600083815260200190815260200160002060006101000a81548160ff02191690831515021790555060086000828152602001908152602001600020600090558281837f416faa49381864da90ecda6efc4e4c051fcd583aae4fbf48b48cf3538db78be760405160405180910390a4505050565b600060016000838152602001908152602001600020549050919050565b6000600354905090565b6000600254905090565b6000600a6000838152602001908152602001600020549050919050565b6000806008600084815260200190815260200160002054905080915050919050565b600080331461118d576040517e31cb2800000000000000000000000000000000000000000000000000000000815260040161118490611c28565b60405180910390fd5b6111b382600b6000868152602001908152602001600020546113fb90919063ffffffff16565b600b600085815260200190815260200160002081905550827fdff0b7fe8242a4caacdf182c8ac9188d4c318cd6fb7dfd9e5dbb28bcbf223198836040516111fa9190611de8565b60405180910390a2600b600084815260200190815260200160002054905092915050565b6000600b6000838152602001908152602001600020549050919050565b600080331461127e576040517e31cb2800000000000000000000000000000000000000000000000000000000815260040161127590611c28565b60405180910390fd5b6112a482600a6000868152602001908152602001600020546113fb90919063ffffffff16565b600a600085815260200190815260200160002081905550827f2b24b52ec200ef483de9b526f769c3f8197432e6028e84c83801789055f18030836040516112eb9190611de8565b60405180910390a2600a600084815260200190815260200160002054905092915050565b60008015156005600084815260200190815260200160002060009054906101000a900460ff161515141561134657600090506113f6565b6000600c600084815260200190815260200160002054111561136b57600090506113f6565b60006113a6600b60008581526020019081526020016000205460016000868152602001908152602001600020546113fb90919063ffffffff16565b9050600a60008481526020019081526020016000205481116113cc5760009150506113f6565b6113f2600a6000858152602001908152602001600020548261141790919063ffffffff16565b9150505b919050565b60008082840190508381101561140d57fe5b8091505092915050565b60008282111561142357fe5b818303905092915050565b60008135905061143d81611eab565b92915050565b60008135905061145281611ec2565b92915050565b60008151905061146781611ec2565b92915050565b60006020828403121561147f57600080fd5b600061148d8482850161142e565b91505092915050565b600080604083850312156114a957600080fd5b60006114b78582860161142e565b92505060206114c885828601611443565b9150509250929050565b6000602082840312156114e457600080fd5b60006114f284828501611458565b91505092915050565b60006115078383611513565b60208301905092915050565b61151c81611e8f565b82525050565b61152b81611e8f565b82525050565b600061153c82611e4a565b6115468185611e62565b935061155183611e3a565b8060005b8381101561158257815161156988826114fb565b975061157483611e55565b925050600181019050611555565b5085935050505092915050565b600061159c601983611e7e565b91507f56616c696461746f7220616c72656164792065786973746564000000000000006000830152602082019050919050565b60006115dc602483611e7e565b91507f4465706f7369746f72207769746864726177616c20726571756573742070656e60008301527f64696e67000000000000000000000000000000000000000000000000000000006020830152604082019050919050565b6000611642601683611e7e565b91507f4465706f7369746f722065786973746564206f6e6365000000000000000000006000830152602082019050919050565b6000611682600083611e73565b9150600082019050919050565b600061169c602b83611e7e565b91507f4465706f73697420616d6f756e742062656c6f77206d696e696d756d2064657060008301527f6f73697420616d6f756e740000000000000000000000000000000000000000006020830152604082019050919050565b6000611702601183611e7e565b91507f496e76616c69642076616c696461746f720000000000000000000000000000006000830152602082019050919050565b6000611742601983611e7e565b91507f4f6e6c7920564d2063616c6c732061726520616c6c6f776564000000000000006000830152602082019050919050565b6000611782601983611e7e565b91507f4465706f7369746f7220616c72656164792065786973746564000000000000006000830152602082019050919050565b60006117c2603583611e7e565b91507f4465706f7369746f7220616464726573732063616e6e6f742062652073616d6560008301527f2061732056616c696461746f72206164647265737300000000000000000000006020830152604082019050919050565b6000611828601583611e7e565b91507f5769746864726177616c2069732070656e64696e6700000000000000000000006000830152602082019050919050565b6000611868601883611e7e565b91507f56616c696461746f722069732061206465706f7369746f7200000000000000006000830152602082019050919050565b60006118a8602083611e7e565b91507f76616c696461746f722062616c616e63652073686f756c64206265207a65726f6000830152602082019050919050565b60006118e8600f83611e7e565b91507f5769746864726177206661696c656400000000000000000000000000000000006000830152602082019050919050565b6000611928601883611e7e565b91507f4465706f7369746f7220616c72656164792065786973747300000000000000006000830152602082019050919050565b6000611968601883611e7e565b91507f4465706f7369746f7220646f6573206e6f7420657869737400000000000000006000830152602082019050919050565b60006119a8602b83611e7e565b91507f4465706f7369746f72207769746864726177616c207265717565737420646f6560008301527f73206e6f742065786973740000000000000000000000000000000000000000006020830152604082019050919050565b6000611a0e601683611e7e565b91507f56616c696461746f722065786973746564206f6e6365000000000000000000006000830152602082019050919050565b6000611a4e602383611e7e565b91507f4465706f7369746f72207769746864726177616c20726571756573742065786960008301527f73747300000000000000000000000000000000000000000000000000000000006020830152604082019050919050565b6000611ab4601883611e7e565b91507f56616c696461746f7220616c72656164792065786973747300000000000000006000830152602082019050919050565b6000611af4601983611e7e565b91507f4465706f7369746f722062616c616e6365206973207a65726f000000000000006000830152602082019050919050565b611b3081611ea1565b82525050565b6000611b4182611675565b9150819050919050565b6000602082019050611b606000830184611522565b92915050565b60006020820190508181036000830152611b808184611531565b905092915050565b60006020820190508181036000830152611ba18161158f565b9050919050565b60006020820190508181036000830152611bc1816115cf565b9050919050565b60006020820190508181036000830152611be181611635565b9050919050565b60006020820190508181036000830152611c018161168f565b9050919050565b60006020820190508181036000830152611c21816116f5565b9050919050565b60006020820190508181036000830152611c4181611735565b9050919050565b60006020820190508181036000830152611c6181611775565b9050919050565b60006020820190508181036000830152611c81816117b5565b9050919050565b60006020820190508181036000830152611ca18161181b565b9050919050565b60006020820190508181036000830152611cc18161185b565b9050919050565b60006020820190508181036000830152611ce18161189b565b9050919050565b60006020820190508181036000830152611d01816118db565b9050919050565b60006020820190508181036000830152611d218161191b565b9050919050565b60006020820190508181036000830152611d418161195b565b9050919050565b60006020820190508181036000830152611d618161199b565b9050919050565b60006020820190508181036000830152611d8181611a01565b9050919050565b60006020820190508181036000830152611da181611a41565b9050919050565b60006020820190508181036000830152611dc181611aa7565b9050919050565b60006020820190508181036000830152611de181611ae7565b9050919050565b6000602082019050611dfd6000830184611b27565b92915050565b6000606082019050611e186000830186611b27565b611e256020830185611b27565b611e326040830184611b27565b949350505050565b6000819050602082019050919050565b600081519050919050565b6000602082019050919050565b600082825260208201905092915050565b600081905092915050565b600082825260208201905092915050565b6000611e9a82611ea1565b9050919050565b6000819050919050565b611eb481611e8f565b8114611ebf57600080fd5b50565b611ecb81611ea1565b8114611ed657600080fd5b5056fea2646970667358221220008b18e5e2e5b96a822c44f75aa8f8342beba90434d8967dc455078db4e2afbc64736f6c63782b302e372e362d646576656c6f702e323032332e31312e392b636f6d6d69742e33313838663336632e6d6f64005c")
	// The method call to 'test()'
	input := common.Hex2Bytes("1a8ab468")
	chain := &dummyChain{}
	ret, _, err := Execute(data, input, &Config{
		GetHashFn:   core.GetHashFn(header, chain),
		BlockNumber: new(big.Int).Set(header.Number),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(ret) != 96 {
		t.Fatalf("expected returndata to be 96 bytes, got %d", len(ret))
	}

	zero := new(big.Int).SetBytes(ret[0:32])
	first := new(big.Int).SetBytes(ret[32:64])
	last := new(big.Int).SetBytes(ret[64:96])
	if zero.BitLen() != 0 {
		t.Fatalf("expected zeroes, got %x", ret[0:32])
	}
	if first.Uint64() != 999 {
		t.Fatalf("second block should be 999, got %d (%x)", first, ret[32:64])
	}
	if last.Uint64() != 744 {
		t.Fatalf("last block should be 744, got %d (%x)", last, ret[64:96])
	}
	if exp, got := 255, chain.counter; exp != got {
		t.Errorf("suboptimal; too much chain iteration, expected %d, got %d", exp, got)
	}
}

type stepCounter struct {
	inner *vm.JSONLogger
	steps int
}

func (s *stepCounter) CaptureStart(env *vm.EVM, from common.Address, to common.Address, create bool, input []byte, gas uint64, value *big.Int) {
}

func (s *stepCounter) CaptureFault(env *vm.EVM, pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, depth int, err error) {
}

func (s *stepCounter) CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) {}

func (s *stepCounter) CaptureState(env *vm.EVM, pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, rData []byte, depth int, err error) {
	s.steps++
	// Enable this for more output
	//s.inner.CaptureState(env, pc, op, gas, cost, memory, stack, rStack, contract, depth, err)
}

// benchmarkNonModifyingCode benchmarks code, but if the code modifies the
// state, this should not be used, since it does not reset the state between runs.
func benchmarkNonModifyingCode(gas uint64, code []byte, name string, b *testing.B) {
	cfg := new(Config)
	setDefaults(cfg)
	cfg.State, _ = state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
	cfg.GasLimit = gas
	var (
		destination = common.BytesToAddress([]byte("contract"))
		vmenv       = NewEnv(cfg)
		sender      = vm.AccountRef(cfg.Origin)
	)
	cfg.State.CreateAccount(destination)
	eoa := common.HexToAddress("E0")
	{
		cfg.State.CreateAccount(eoa)
		cfg.State.SetNonce(eoa, 100)
	}
	reverting := common.HexToAddress("EE")
	{
		cfg.State.CreateAccount(reverting)
		cfg.State.SetCode(reverting, []byte{
			byte(vm.PUSH1), 0x00,
			byte(vm.PUSH1), 0x00,
			byte(vm.REVERT),
		})
	}

	//cfg.State.CreateAccount(cfg.Origin)
	// set the receiver's (the executing contract) code for execution.
	cfg.State.SetCode(destination, code)
	vmenv.Call(sender, destination, nil, gas, cfg.Value)

	b.Run(name, func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			vmenv.Call(sender, destination, nil, gas, cfg.Value)
		}
	})
}

// BenchmarkSimpleLoop test a pretty simple loop which loops until OOG
// 55 ms
func BenchmarkSimpleLoop(b *testing.B) {

	staticCallIdentity := []byte{
		byte(vm.JUMPDEST), //  [ count ]
		// push args for the call
		byte(vm.PUSH1), 0, // out size
		byte(vm.DUP1),       // out offset
		byte(vm.DUP1),       // out insize
		byte(vm.DUP1),       // in offset
		byte(vm.PUSH1), 0x4, // address of identity
		byte(vm.GAS), // gas
		byte(vm.STATICCALL),
		byte(vm.POP),      // pop return value
		byte(vm.PUSH1), 0, // jumpdestination
		byte(vm.JUMP),
	}

	callIdentity := []byte{
		byte(vm.JUMPDEST), //  [ count ]
		// push args for the call
		byte(vm.PUSH1), 0, // out size
		byte(vm.DUP1),       // out offset
		byte(vm.DUP1),       // out insize
		byte(vm.DUP1),       // in offset
		byte(vm.DUP1),       // value
		byte(vm.PUSH1), 0x4, // address of identity
		byte(vm.GAS), // gas
		byte(vm.CALL),
		byte(vm.POP),      // pop return value
		byte(vm.PUSH1), 0, // jumpdestination
		byte(vm.JUMP),
	}

	callInexistant := []byte{
		byte(vm.JUMPDEST), //  [ count ]
		// push args for the call
		byte(vm.PUSH1), 0, // out size
		byte(vm.DUP1),        // out offset
		byte(vm.DUP1),        // out insize
		byte(vm.DUP1),        // in offset
		byte(vm.DUP1),        // value
		byte(vm.PUSH1), 0xff, // address of existing contract
		byte(vm.GAS), // gas
		byte(vm.CALL),
		byte(vm.POP),      // pop return value
		byte(vm.PUSH1), 0, // jumpdestination
		byte(vm.JUMP),
	}

	callEOA := []byte{
		byte(vm.JUMPDEST), //  [ count ]
		// push args for the call
		byte(vm.PUSH1), 0, // out size
		byte(vm.DUP1),        // out offset
		byte(vm.DUP1),        // out insize
		byte(vm.DUP1),        // in offset
		byte(vm.DUP1),        // value
		byte(vm.PUSH1), 0xE0, // address of EOA
		byte(vm.GAS), // gas
		byte(vm.CALL),
		byte(vm.POP),      // pop return value
		byte(vm.PUSH1), 0, // jumpdestination
		byte(vm.JUMP),
	}

	loopingCode := []byte{
		byte(vm.JUMPDEST), //  [ count ]
		// push args for the call
		byte(vm.PUSH1), 0, // out size
		byte(vm.DUP1),       // out offset
		byte(vm.DUP1),       // out insize
		byte(vm.DUP1),       // in offset
		byte(vm.PUSH1), 0x4, // address of identity
		byte(vm.GAS), // gas

		byte(vm.POP), byte(vm.POP), byte(vm.POP), byte(vm.POP), byte(vm.POP), byte(vm.POP),
		byte(vm.PUSH1), 0, // jumpdestination
		byte(vm.JUMP),
	}

	calllRevertingContractWithInput := []byte{
		byte(vm.JUMPDEST), //
		// push args for the call
		byte(vm.PUSH1), 0, // out size
		byte(vm.DUP1),        // out offset
		byte(vm.PUSH1), 0x20, // in size
		byte(vm.PUSH1), 0x00, // in offset
		byte(vm.PUSH1), 0x00, // value
		byte(vm.PUSH1), 0xEE, // address of reverting contract
		byte(vm.GAS), // gas
		byte(vm.CALL),
		byte(vm.POP),      // pop return value
		byte(vm.PUSH1), 0, // jumpdestination
		byte(vm.JUMP),
	}

	//tracer := vm.NewJSONLogger(nil, os.Stdout)
	//Execute(loopingCode, nil, &Config{
	//	EVMConfig: vm.Config{
	//		Debug:  true,
	//		Tracer: tracer,
	//	}})
	// 100M gas
	benchmarkNonModifyingCode(100000000, staticCallIdentity, "staticcall-identity-100M", b)
	benchmarkNonModifyingCode(100000000, callIdentity, "call-identity-100M", b)
	benchmarkNonModifyingCode(100000000, loopingCode, "loop-100M", b)
	benchmarkNonModifyingCode(100000000, callInexistant, "call-nonexist-100M", b)
	benchmarkNonModifyingCode(100000000, callEOA, "call-EOA-100M", b)
	benchmarkNonModifyingCode(100000000, calllRevertingContractWithInput, "call-reverting-100M", b)

	//benchmarkNonModifyingCode(10000000, staticCallIdentity, "staticcall-identity-10M", b)
	//benchmarkNonModifyingCode(10000000, loopingCode, "loop-10M", b)
}

// TestEip2929Cases contains various testcases that are used for
// EIP-2929 about gas repricings
func TestEip2929Cases(t *testing.T) {

	id := 1
	prettyPrint := func(comment string, code []byte) {

		instrs := make([]string, 0)
		it := asm.NewInstructionIterator(code)
		for it.Next() {
			if it.Arg() != nil && 0 < len(it.Arg()) {
				instrs = append(instrs, fmt.Sprintf("%v 0x%x", it.Op(), it.Arg()))
			} else {
				instrs = append(instrs, fmt.Sprintf("%v", it.Op()))
			}
		}
		ops := strings.Join(instrs, ", ")
		fmt.Printf("### Case %d\n\n", id)
		id++
		fmt.Printf("%v\n\nBytecode: \n```\n0x%x\n```\nOperations: \n```\n%v\n```\n\n",
			comment,
			code, ops)
		Execute(code, nil, &Config{
			EVMConfig: vm.Config{
				Debug:     true,
				Tracer:    vm.NewMarkdownLogger(nil, os.Stdout),
				ExtraEips: []int{2929},
			},
		})
	}

	{ // First eip testcase
		code := []byte{
			// Three checks against a precompile
			byte(vm.PUSH1), 1, byte(vm.EXTCODEHASH), byte(vm.POP),
			byte(vm.PUSH1), 2, byte(vm.EXTCODESIZE), byte(vm.POP),
			byte(vm.PUSH1), 3, byte(vm.BALANCE), byte(vm.POP),
			// Three checks against a non-precompile
			byte(vm.PUSH1), 0xf1, byte(vm.EXTCODEHASH), byte(vm.POP),
			byte(vm.PUSH1), 0xf2, byte(vm.EXTCODESIZE), byte(vm.POP),
			byte(vm.PUSH1), 0xf3, byte(vm.BALANCE), byte(vm.POP),
			// Same three checks (should be cheaper)
			byte(vm.PUSH1), 0xf2, byte(vm.EXTCODEHASH), byte(vm.POP),
			byte(vm.PUSH1), 0xf3, byte(vm.EXTCODESIZE), byte(vm.POP),
			byte(vm.PUSH1), 0xf1, byte(vm.BALANCE), byte(vm.POP),
			// Check the origin, and the 'this'
			byte(vm.ORIGIN), byte(vm.BALANCE), byte(vm.POP),
			byte(vm.ADDRESS), byte(vm.BALANCE), byte(vm.POP),

			byte(vm.STOP),
		}
		prettyPrint("This checks `EXT`(codehash,codesize,balance) of precompiles, which should be `100`, "+
			"and later checks the same operations twice against some non-precompiles. "+
			"Those are cheaper second time they are accessed. Lastly, it checks the `BALANCE` of `origin` and `this`.", code)
	}

	{ // EXTCODECOPY
		code := []byte{
			// extcodecopy( 0xff,0,0,0,0)
			byte(vm.PUSH1), 0x00, byte(vm.PUSH1), 0x00, byte(vm.PUSH1), 0x00, //length, codeoffset, memoffset
			byte(vm.PUSH1), 0xff, byte(vm.EXTCODECOPY),
			// extcodecopy( 0xff,0,0,0,0)
			byte(vm.PUSH1), 0x00, byte(vm.PUSH1), 0x00, byte(vm.PUSH1), 0x00, //length, codeoffset, memoffset
			byte(vm.PUSH1), 0xff, byte(vm.EXTCODECOPY),
			// extcodecopy( this,0,0,0,0)
			byte(vm.PUSH1), 0x00, byte(vm.PUSH1), 0x00, byte(vm.PUSH1), 0x00, //length, codeoffset, memoffset
			byte(vm.ADDRESS), byte(vm.EXTCODECOPY),

			byte(vm.STOP),
		}
		prettyPrint("This checks `extcodecopy( 0xff,0,0,0,0)` twice, (should be expensive first time), "+
			"and then does `extcodecopy( this,0,0,0,0)`.", code)
	}

	{ // SLOAD + SSTORE
		code := []byte{

			// Add slot `0x1` to access list
			byte(vm.PUSH1), 0x01, byte(vm.SLOAD), byte(vm.POP), // SLOAD( 0x1) (add to access list)
			// Write to `0x1` which is already in access list
			byte(vm.PUSH1), 0x11, byte(vm.PUSH1), 0x01, byte(vm.SSTORE), // SSTORE( loc: 0x01, val: 0x11)
			// Write to `0x2` which is not in access list
			byte(vm.PUSH1), 0x11, byte(vm.PUSH1), 0x02, byte(vm.SSTORE), // SSTORE( loc: 0x02, val: 0x11)
			// Write again to `0x2`
			byte(vm.PUSH1), 0x11, byte(vm.PUSH1), 0x02, byte(vm.SSTORE), // SSTORE( loc: 0x02, val: 0x11)
			// Read slot in access list (0x2)
			byte(vm.PUSH1), 0x02, byte(vm.SLOAD), // SLOAD( 0x2)
			// Read slot in access list (0x1)
			byte(vm.PUSH1), 0x01, byte(vm.SLOAD), // SLOAD( 0x1)
		}
		prettyPrint("This checks `sload( 0x1)` followed by `sstore(loc: 0x01, val:0x11)`, then 'naked' sstore:"+
			"`sstore(loc: 0x02, val:0x11)` twice, and `sload(0x2)`, `sload(0x1)`. ", code)
	}
	{ // Call variants
		code := []byte{
			// identity precompile
			byte(vm.PUSH1), 0x0, byte(vm.DUP1), byte(vm.DUP1), byte(vm.DUP1), byte(vm.DUP1),
			byte(vm.PUSH1), 0x04, byte(vm.PUSH1), 0x0, byte(vm.CALL), byte(vm.POP),

			// random account - call 1
			byte(vm.PUSH1), 0x0, byte(vm.DUP1), byte(vm.DUP1), byte(vm.DUP1), byte(vm.DUP1),
			byte(vm.PUSH1), 0xff, byte(vm.PUSH1), 0x0, byte(vm.CALL), byte(vm.POP),

			// random account - call 2
			byte(vm.PUSH1), 0x0, byte(vm.DUP1), byte(vm.DUP1), byte(vm.DUP1), byte(vm.DUP1),
			byte(vm.PUSH1), 0xff, byte(vm.PUSH1), 0x0, byte(vm.STATICCALL), byte(vm.POP),
		}
		prettyPrint("This calls the `identity`-precompile (cheap), then calls an account (expensive) and `staticcall`s the same"+
			"account (cheap)", code)
	}
}

// TestColdAccountAccessCost test that the cold account access cost is reported
// correctly
// see: https://github.com/ethereum/go-ethereum/issues/22649
func TestColdAccountAccessCost(t *testing.T) {
	for i, tc := range []struct {
		code []byte
		step int
		want uint64
	}{
		{ // EXTCODEHASH(0xff)
			code: []byte{byte(vm.PUSH1), 0xFF, byte(vm.EXTCODEHASH), byte(vm.POP)},
			step: 1,
			want: 2600,
		},
		{ // BALANCE(0xff)
			code: []byte{byte(vm.PUSH1), 0xFF, byte(vm.BALANCE), byte(vm.POP)},
			step: 1,
			want: 2600,
		},
		{ // CALL(0xff)
			code: []byte{
				byte(vm.PUSH1), 0x0,
				byte(vm.DUP1), byte(vm.DUP1), byte(vm.DUP1), byte(vm.DUP1),
				byte(vm.PUSH1), 0xff, byte(vm.DUP1), byte(vm.CALL), byte(vm.POP),
			},
			step: 7,
			want: 2855,
		},
		{ // CALLCODE(0xff)
			code: []byte{
				byte(vm.PUSH1), 0x0,
				byte(vm.DUP1), byte(vm.DUP1), byte(vm.DUP1), byte(vm.DUP1),
				byte(vm.PUSH1), 0xff, byte(vm.DUP1), byte(vm.CALLCODE), byte(vm.POP),
			},
			step: 7,
			want: 2855,
		},
		{ // DELEGATECALL(0xff)
			code: []byte{
				byte(vm.PUSH1), 0x0,
				byte(vm.DUP1), byte(vm.DUP1), byte(vm.DUP1),
				byte(vm.PUSH1), 0xff, byte(vm.DUP1), byte(vm.DELEGATECALL), byte(vm.POP),
			},
			step: 6,
			want: 2855,
		},
		{ // STATICCALL(0xff)
			code: []byte{
				byte(vm.PUSH1), 0x0,
				byte(vm.DUP1), byte(vm.DUP1), byte(vm.DUP1),
				byte(vm.PUSH1), 0xff, byte(vm.DUP1), byte(vm.STATICCALL), byte(vm.POP),
			},
			step: 6,
			want: 2855,
		},
		{ // SELFDESTRUCT(0xff)
			code: []byte{
				byte(vm.PUSH1), 0xff, byte(vm.SELFDESTRUCT),
			},
			step: 1,
			want: 7600,
		},
	} {
		tracer := vm.NewStructLogger(nil)
		Execute(tc.code, nil, &Config{
			EVMConfig: vm.Config{
				Debug:  true,
				Tracer: tracer,
			},
		})
		have := tracer.StructLogs()[tc.step].GasCost
		if want := tc.want; have != want {
			for ii, op := range tracer.StructLogs() {
				t.Logf("%d: %v %d", ii, op.OpName(), op.GasCost)
			}
			t.Fatalf("tescase %d, gas report wrong, step %d, have %d want %d", i, tc.step, have, want)
		}
	}
}
