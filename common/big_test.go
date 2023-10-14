package common

import (
	"fmt"
	"math/big"
	"testing"
)

func testAddBigInt(a int64, b int64, sum int64) bool {
	var a1 *big.Int
	var b1 *big.Int
	a1 = big.NewInt(a)
	b1 = big.NewInt(b)
	a1 = a1.Add(a1, b1)
	fmt.Println("a1", a1)
	if a1.Int64() != sum {
		return false
	}
	if b1.Int64() != b {
		return false
	}
	return true
}

func testSafeAddBigInt(a int64, b int64, result int64) bool {
	var a1 *big.Int
	var b1 *big.Int
	var result1 *big.Int
	a1 = big.NewInt(a)
	b1 = big.NewInt(b)
	result1 = SafeAddBigInt(a1, b1)
	fmt.Println("result1", result1)
	if result1.Int64() != result {
		return false
	}
	if a1.Int64() != a {
		return false
	}
	if b1.Int64() != b {
		return false
	}
	return true
}

func testSafeSubBigInt(a int64, b int64, result int64) bool {
	var a1 *big.Int
	var b1 *big.Int
	var result1 *big.Int
	a1 = big.NewInt(a)
	b1 = big.NewInt(b)
	result1 = SafeSubBigInt(a1, b1)
	fmt.Println("result1", result1)
	if result1.Int64() != result {
		return false
	}
	if a1.Int64() != a {
		return false
	}
	if b1.Int64() != b {
		return false
	}
	return true
}

func Test_SafeAddBigInt(t *testing.T) {
	if testSafeAddBigInt(123456, 101010, 224466) == false {
		t.Fatalf("failed")
	}

	if testSafeAddBigInt(-123456, -101010, -224466) == false {
		t.Fatalf("failed")
	}

	if testSafeAddBigInt(-123456, 101010, -22446) == false {
		t.Fatalf("failed")
	}

	if testSafeAddBigInt(123456, -101010, 22446) == false {
		t.Fatalf("failed")
	}

	if testAddBigInt(123456, 101010, 224466) == false {
		t.Fatalf("failed")
	}

	if testAddBigInt(-123456, -101010, -224466) == false {
		t.Fatalf("failed")
	}

	if testAddBigInt(-123456, 101010, -22446) == false {
		t.Fatalf("failed")
	}

	if testAddBigInt(123456, -101010, 22446) == false {
		t.Fatalf("failed")
	}
}

func Test_SafeSubBigInt(t *testing.T) {
	if testSafeSubBigInt(123456, 101010, 22446) == false {
		t.Fatalf("failed")
	}

	if testSafeSubBigInt(-123456, -101010, -22446) == false {
		t.Fatalf("failed")
	}

	if testSafeSubBigInt(-123456, 101010, -224466) == false {
		t.Fatalf("failed")
	}

	if testSafeSubBigInt(123456, -101010, 224466) == false {
		t.Fatalf("failed")
	}

	if testSafeSubBigInt(123456, 101010, 22446) == false {
		t.Fatalf("failed")
	}

	if testSafeSubBigInt(-123456, -101010, -22446) == false {
		t.Fatalf("failed")
	}

	if testSafeSubBigInt(-123456, 101010, -224466) == false {
		t.Fatalf("failed")
	}

	if testSafeSubBigInt(123456, -101010, 224466) == false {
		t.Fatalf("failed")
	}
}

func testSafeMulBigInt(a int64, b int64, result int64) bool {
	var a1 *big.Int
	var b1 *big.Int
	var result1 *big.Int
	a1 = big.NewInt(a)
	b1 = big.NewInt(b)
	result1 = SafeMulBigInt(a1, b1)
	fmt.Println("result1", result1)
	if result1.Int64() != result {
		return false
	}
	if a1.Int64() != a {
		return false
	}
	if b1.Int64() != b {
		return false
	}
	return true
}

func Test_SafeMulBigInt(t *testing.T) {
	if testSafeMulBigInt(123456, 101010, 12470290560) == false {
		t.Fatalf("failed")
	}

	if testSafeMulBigInt(-123456, -101010, 12470290560) == false {
		t.Fatalf("failed")
	}

	if testSafeMulBigInt(-123456, 101010, -12470290560) == false {
		t.Fatalf("failed")
	}

	if testSafeMulBigInt(123456, -101010, -12470290560) == false {
		t.Fatalf("failed")
	}
}

func testSafeDivBigInt(a int64, b int64, result int64) bool {
	var a1 *big.Int
	var b1 *big.Int
	var result1 *big.Int
	a1 = big.NewInt(a)
	b1 = big.NewInt(b)
	result1 = SafeDivBigInt(a1, b1)
	fmt.Println("result1", result1)
	if result1.Int64() != result {
		return false
	}
	if a1.Int64() != a {
		return false
	}
	if b1.Int64() != b {
		return false
	}
	return true
}

func Test_SafeDivBigInt(t *testing.T) {
	if testSafeDivBigInt(123456, 101010, 1) == false {
		t.Fatalf("failed")
	}

	if testSafeDivBigInt(-123456, -101010, 2) == false {
		t.Fatalf("failed")
	}

	if testSafeDivBigInt(-123456, 101010, -2) == false {
		t.Fatalf("failed")
	}

	if testSafeDivBigInt(123456, -101010, -1) == false {
		t.Fatalf("failed")
	}

	if testSafeDivBigInt(101010, 123456, 0) == false {
		t.Fatalf("failed")
	}

	if testSafeDivBigInt(-101010, -123456, 1) == false {
		t.Fatalf("failed")
	}

	if testSafeDivBigInt(101010, -123456, 0) == false {
		t.Fatalf("failed")
	}

	if testSafeDivBigInt(-101010, 123456, -1) == false {
		t.Fatalf("failed")
	}
}
func testSafeDivBigFloat(a float64, b float64, result float64) bool {
	var a1 *big.Float
	var b1 *big.Float
	var result1 *big.Float
	a1 = big.NewFloat(a)
	b1 = big.NewFloat(b)
	result1 = SafeDivBigFloat(a1, b1)
	fmt.Println("result1", result1)

	r, _ := result1.Float64()
	if r != result {
		return false
	}

	a2, _ := a1.Float64()
	if a2 != a {
		return false
	}
	b2, _ := b1.Float64()
	if b2 != b {
		return false
	}
	return true
}

func Test_SafeDivBigFloat(t *testing.T) {
	if testSafeDivBigFloat(123456, 101010, 1.2222156222156222) == false {
		t.Fatalf("failed")
	}

	if testSafeDivBigFloat(-123456, -101010, 1.2222156222156222) == false {
		t.Fatalf("failed")
	}

	if testSafeDivBigFloat(-123456, 101010, -1.2222156222156222) == false {
		t.Fatalf("failed")
	}

	if testSafeDivBigFloat(123456, -101010, -1.2222156222156222) == false {
		t.Fatalf("failed")
	}

	if testSafeDivBigFloat(101010, 123456, 0.8181862363919129) == false {
		t.Fatalf("failed")
	}

	if testSafeDivBigFloat(-101010, -123456, 0.8181862363919129) == false {
		t.Fatalf("failed")
	}

	if testSafeDivBigFloat(101010, -123456, -0.8181862363919129) == false {
		t.Fatalf("failed")
	}

	if testSafeDivBigFloat(22, 7, 3.142857142857143) == false {
		t.Fatalf("failed")
	}
}

func testSafePercentageOfBigInt(a int64, b int64, result int64) bool {
	var a1 *big.Int
	var b1 *big.Int
	var result1 *big.Int
	a1 = big.NewInt(a)
	b1 = big.NewInt(b)
	result1 = SafePercentageOfBigInt(a1, b1)
	fmt.Println("result1", result1)
	if result1.Int64() != result {
		return false
	}
	if a1.Int64() != a {
		return false
	}
	if b1.Int64() != b {
		return false
	}
	return true
}

func Test_SafePercentageOfBigInt(t *testing.T) {
	if testSafePercentageOfBigInt(100, 100, 100) == false {
		t.Fatalf("failed")
	}

	if testSafePercentageOfBigInt(1, 100, 1) == false {
		t.Fatalf("failed")
	}

	if testSafePercentageOfBigInt(200000, 400000, 50) == false {
		t.Fatalf("failed")
	}

	if testSafePercentageOfBigInt(70000, 100000, 70) == false {
		t.Fatalf("failed")
	}

	if testSafePercentageOfBigInt(70500, 100000, 70) == false {
		t.Fatalf("failed")
	}
}

func testSafeRelativePercentageBigInt(a int64, b int64, result int64) bool {
	var a1 *big.Int
	var b1 *big.Int
	var result1 *big.Int
	a1 = big.NewInt(a)
	b1 = big.NewInt(b)
	result1 = SafeRelativePercentageBigInt(a1, b1)
	fmt.Println("result1", result1)
	if result1.Int64() != result {
		return false
	}
	if a1.Int64() != a {
		return false
	}
	if b1.Int64() != b {
		return false
	}
	return true
}

func Test_SafeRelativePercentageBigInt(t *testing.T) {
	if testSafeRelativePercentageBigInt(100, 70, 70) == false {
		t.Fatalf("failed")
	}

	if testSafeRelativePercentageBigInt(100000, 70, 70000) == false {
		t.Fatalf("failed")
	}

	if testSafeRelativePercentageBigInt(100000, 200, 200000) == false {
		t.Fatalf("failed")
	}

	if testSafeRelativePercentageBigInt(400000, 75, 300000) == false {
		t.Fatalf("failed")
	}

	if testSafeRelativePercentageBigInt(70000, 70, 49000) == false {
		t.Fatalf("failed")
	}
}
