package util

import (
	"fmt"
	"github.com/DogeProtocol/dp/common"
	"testing"
)

func Test_Temp(t *testing.T) {
	tempData, err := common.Hex2BytesWithErrorCheck("00000000000000000000000000000000000000000005ca4ec2a79a7f6700000000000000000000000000000000000000000000000000000000000000000676c0000000000000000000000000000000000000000000000000000000006650c7a8")
	if err != nil {
		fmt.Println("err", err)
		return
	}
	fmt.Println(tempData)

	fmt.Println(common.Hex2Bytes("731f750d43dc0fb62f1251286479ed4f420f30d4ec593422dff43936b7df49a8036c6864"))
}
