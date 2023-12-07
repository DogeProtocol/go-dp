// Copyright 2014 The go-ethereum Authors
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

package common

import (
	"math/big"
)

// Common big integers often used
var (
	Big1   = big.NewInt(1)
	Big2   = big.NewInt(2)
	Big3   = big.NewInt(3)
	Big0   = big.NewInt(0)
	Big32  = big.NewInt(32)
	Big256 = big.NewInt(256)
	Big257 = big.NewInt(257)
)

func SafeAddBigInt(x, y *big.Int) *big.Int {
	result := big.NewInt(0)
	result.Add(x, y)
	return result
}

func SafeSubBigInt(x, y *big.Int) *big.Int {
	result := big.NewInt(0)
	result.Sub(x, y)
	return result
}

func SafeMulBigInt(x, y *big.Int) *big.Int {
	result := big.NewInt(0)
	result.Mul(x, y)
	return result
}

func SafeDivBigInt(x, y *big.Int) *big.Int {
	result := big.NewInt(0)
	result.Div(x, y)
	return result
}

func SafeDivBigFloat(x, y *big.Float) *big.Float {
	result := big.NewFloat(0)
	result.Quo(x, y)
	return result
}

func SafePercentageOfBigInt(x, y *big.Int) *big.Int {
	hundred := big.NewInt(100)
	return SafeDivBigInt(SafeMulBigInt(hundred, x), y)
}

func SafeRelativePercentageBigInt(total, percentage *big.Int) *big.Int {
	hundred := big.NewInt(100)
	return SafeDivBigInt(SafeMulBigInt(total, percentage), hundred)
}
