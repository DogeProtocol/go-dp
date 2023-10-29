// Copyright 2017 The go-ethereum Authors
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

package proofofstake

import (
	"encoding/json"
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/common/hexutil"
	"github.com/DogeProtocol/dp/consensus"
	"github.com/DogeProtocol/dp/rpc"
	"math/big"
)

// API is a user facing RPC API to allow controlling the signer and voting
// mechanisms of the proof-of-authority scheme.
type API struct {
	chain        consensus.ChainHeaderReader
	proofofstake *ProofOfStake
}

type blockNumberOrHashOrRLP struct {
	*rpc.BlockNumberOrHash
	RLP hexutil.Bytes `json:"rlp,omitempty"`
}

func (sb *blockNumberOrHashOrRLP) UnmarshalJSON(data []byte) error {
	bnOrHash := new(rpc.BlockNumberOrHash)
	// Try to unmarshal bNrOrHash
	if err := bnOrHash.UnmarshalJSON(data); err == nil {
		sb.BlockNumberOrHash = bnOrHash
		return nil
	}
	// Try to unmarshal RLP
	var input string
	if err := json.Unmarshal(data, &input); err != nil {
		return err
	}
	sb.RLP = hexutil.MustDecode(input)
	return nil
}

// GetValidators retrieves the list of authorized signers at the specified block.
func (api *API) GetValidators() (map[common.Address]*big.Int, error) {
	// Retrieve the requested block number (or current if none requested)
	var header = api.chain.CurrentHeader()
	if header == nil {
		return nil, errUnknownBlock
	}
	validators, err := api.proofofstake.GetValidators(header.ParentHash)
	if err != nil {
		return nil, err
	}
	return validators, nil
}
