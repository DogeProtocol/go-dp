package proofofstake

import (
	"fmt"
	"github.com/DogeProtocol/dp/common"
	"github.com/DogeProtocol/dp/params"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

var blockYears = 4

var blockRewardTotal = []float64{951293.759512938, 475646.879756469, 237823.439878234, 118911.719939117, 59455.8599695586, 29727.9299847793, 14863.9649923897, 7431.98249619483,
	3715.99124809741, 1857.99562404871, 928.997812024353, 464.498906012177, 232.249453006088, 116.124726503044, 58.0623632515221, 29.031181625761,
	14.5155908128805, 7.25779540644026, 3.62889770322013, 1.81444885161006, 0.907224425805032, 0.453612212902516, 0.226806106451258, 0.113403053225629,
	0.0567015266128145, 0.0283507633064073, 0.0141753816532036, 0.00708769082660182, 0.00354384541330091, 0.00177192270665045, 0.000885961353325227,
	0.000442980676662613, 0.000221490338331307, 0.000110745169165653, 0.0000553725845828267, 0.0000276862922914133, 0.0000138431461457067,
	0.00000692157307285334, 0.00000346078653642667, 0.00000173039326821333, 0.000000865196634106667, 0.000000432598317053333, 0.000000216299158526667,
	0.000000108149579263333, 0.0000000540747896316667, 0.0000000270373948158333, 0.0000000135186974079167, 0.00000000675934870395834, 0.00000000337967435197917,
	0.00000000168983717598958, 0.000000000844918587994792, 0.000000000422459293997396, 0.000000000211229646998698, 0.000000000105614823499349,
	5.28074117496745e-11, 2.64037058748373e-11, 1.32018529374186e-11, 6.60092646870931e-12, 3.30046323435466e-12, 1.65023161717733e-12,
	8.25115808588664e-13, 4.12557904294332e-13, 2.06278952147166e-13, 1.03139476073583e-13, 5.15697380367915e-14, 2.57848690183957e-14,
	1.28924345091979e-14, 6.44621725459894e-15, 3.22310862729947e-15, 1.61155431364973e-15, 8.05777156824867e-16, 4.02888578412434e-16,
	2.01444289206217e-16, 1.00722144603108e-16, 5.03610723015542e-17, 2.51805361507771e-17}

var blockStartRang = int64(1497600)

var blockEndRange = []int64{22521600, 43545600, 64569600, 85593600, 106617600, 127641600, 148665600, 169689600, 190713600, 211737600, 232761600, 253785600,
	274809600, 295833600, 316857600, 337881600, 358905600, 379929600, 400953600, 421977600, 443001600, 464025600, 485049600, 506073600,
	527097600, 548121600, 569145600, 590169600, 611193600, 632217600, 653241600, 674265600, 695289600, 716313600, 737337600, 758361600,
	779385600, 800409600, 821433600, 842457600, 863481600, 884505600, 905529600, 926553600, 947577600, 968601600, 989625600, 1010649600,
	1031673600, 1052697600, 1073721600, 1094745600, 1115769600, 1136793600, 1157817600, 1178841600, 1199865600, 1220889600, 1241913600,
	1262937600, 1283961600, 1304985600, 1326009600, 1347033600, 1368057600, 1389081600, 1410105600, 1431129600, 1452153600, 1473177600,
	1494201600, 1515225600, 1536249600, 1557273600, 1578297600, 1599321600}

func TestRewardGenerateYearly(t *testing.T) {
	for i := 1; i <= 350; i++ {
		blockNumber := rewardStartBlock.Int64() + (blockYearly.Int64() * int64(i))
		startBlockNumber := big.NewInt(blockNumber - blockYearly.Int64())
		startReward := new(big.Int).Set(GetReward(startBlockNumber))

		endBlockNumber := big.NewInt(blockNumber - 1)
		endReward := new(big.Int).Set(GetReward(endBlockNumber))

		fmt.Println("Year : ", i,
			" Block Range : ", startBlockNumber, " - ", endBlockNumber,
			" Block reward range : ", startReward, " - ", endReward)
	}
}

func TestRewardGenerateBlocks(t *testing.T) {
	startBlockNumber := big.NewInt(22338000 - 1000)
	endBlockNumber := big.NewInt(22338000)
	incrementBlock := big.NewInt(1)

	for startBlockNumber.Int64() <= endBlockNumber.Int64() {
		reward := new(big.Int).Set(GetReward(startBlockNumber))
		fmt.Println("Block Number : ", startBlockNumber, " reward : ", reward)
		startBlockNumber = common.SafeAddBigInt(startBlockNumber, incrementBlock)
	}
}

func TestRewardVerifyYearly(t *testing.T) {
	for i := 1; i <= 12; i++ {
		blockNumber := rewardStartBlock.Int64() - 1 + (blockYearly.Int64() * int64(i))
		startBlockNumber := big.NewInt(blockNumber - blockYearly.Int64())
		startReward := new(big.Int).Set(GetReward(startBlockNumber))

		r1 := params.WeiToEther(getTestReward(startBlockNumber))
		r2 := params.WeiToEther(startReward)
		assert.Equal(t, r1, r2)

		endBlockNumber := big.NewInt(blockNumber - 1)
		endReward := new(big.Int).Set(GetReward(endBlockNumber))

		r1 = params.WeiToEther(getTestReward(endBlockNumber))
		r2 = params.WeiToEther(endReward)
		assert.Equal(t, r1, r2)
	}
}

func TestRewardVerifyBlocks(t *testing.T) {
	startBlockNumber := big.NewInt(1497600 - 1000)
	endBlockNumber := big.NewInt(1497500)
	incrementBlock := big.NewInt(1)

	for startBlockNumber.Int64() <= endBlockNumber.Int64() {
		reward := new(big.Int).Set(GetReward(startBlockNumber))
		r1 := params.WeiToEther(getTestReward(startBlockNumber))
		r2 := params.WeiToEther(reward)
		assert.Equal(t, r1, r2)
		startBlockNumber = common.SafeAddBigInt(startBlockNumber, incrementBlock)
	}
}

func getTestReward(blockNumber *big.Int) *big.Int {
	var reward = big.NewInt(0)
	if blockStartRang <= blockNumber.Int64() {
		var i = 0
		for i < len(blockEndRange) {
			b := blockEndRange[i] - 1
			if blockNumber.Int64() <= b {
				reward = etherToWeiFloat(big.NewFloat(blockRewardTotal[i]))
				break
			}
			i = i + 1
		}
		return reward
	}
	return reward
}
