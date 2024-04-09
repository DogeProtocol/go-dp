package proofofstake

import (
	"fmt"
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/core/state"
	"github.com/DogeProtocol/dp/params"
	"math"
	"math/big"
	"strings"
)

var (
	totalCoin            = big.NewInt(100000000000000)
	percentageChangeYear = big.NewInt(4)
	percentageDefault    = big.NewInt(20)
	percentageDivided    = big.NewInt(2)

	blockSecond = 6
	blockYearly = big.NewInt(int64((((60 * 60) * 24) / blockSecond) * 365))

	rewardStartBlock = big.NewInt(int64(rewardStartBlockNumber))
)

func GetReward(blockNumber *big.Int) *big.Int {

	blockReward := big.NewInt(0)

	if rewardStartBlock.Int64() <= blockNumber.Int64() {
		//Step 0
		block := common.SafeSubBigInt(blockNumber, rewardStartBlock)
		s := common.SafeDivBigInt(block, blockYearly)
		s0 := common.SafeDivBigInt(s, percentageChangeYear)

		//Step 1
		s1 := common.SafeAddBigInt(s0, big.NewInt(1))

		//Step 2
		s2 := MathPow(int(percentageDivided.Int64()), int(s1.Int64()))

		//Step 3
		s3 := (float64(percentageDefault.Int64()) / s2) / float64(percentageDivided.Int64())

		//Step 4 (1 Year Reward)
		totalReward := (float64(totalCoin.Int64()) * s3) / 100

		//Step 5 (Block reward)
		perBlock := big.NewFloat(totalReward / float64(blockYearly.Int64()))
		blockReward = etherToWeiFloat(perBlock)

	}

	return blockReward
}

// MathPow calculates n to the mth power with the math.Pow() function
func MathPow(n, m int) float64 {
	return math.Pow(float64(n), float64(m))
}

func etherToWeiFloat(eth *big.Float) *big.Int {
	truncInt, _ := eth.Int(nil)
	truncInt = new(big.Int).Mul(truncInt, big.NewInt(params.Ether))
	fracStr := strings.Split(fmt.Sprintf("%.18f", eth), ".")[1]
	fracStr += strings.Repeat("0", 18-len(fracStr))
	fracInt, _ := new(big.Int).SetString(fracStr, 10)
	wei := new(big.Int).Add(truncInt, fracInt)
	return wei
}

func (c *ProofOfStake) accumulateBalance(state *state.StateDB, amount *big.Int, addr common.Address) error {
	state.AddBalance(addr, amount)
	return nil
}
