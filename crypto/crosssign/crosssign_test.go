package crosssign

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	error_1   = "[{\"This is my wallet\",\r\n  \"sig\": \"0x6c679b0b1f47bfc668bd67a217ffef3a0d645653b8ae7f25242bc1e7fac1511d7fc514e8a6c6199fc2c9a50648022cd549ba9a6a9def315a106fcc924a1903941c\",\r\n  \"versionTest\": \"2\"\r\n}]"
	error_1_1 = "{\r\n  \"addressTest\": \"0x4667Fc33E26911Ba828cc9F1bB4A0F47B1e383B1\",\r\n  \"msgTest\": \"This is my wallet\",\r\n  \"sig\": \"0x6c679b0b1f47bfc668bd67a217ffef3a0d645653b8ae7f25242bc1e7fac1511d7fc514e8a6c6199fc2c9a50648022cd549ba9a6a9def315a106fcc924a1903941c\",\r\n  \"versionTest\": \"2\"\r\n}"
	error_2   = "{\r\n  \"address\": \"0x4667Fc33E26911Ba828cc9F1bB4A0F47B1e383B1\",\r\n  \"msg\": \"This is my wallet\",\r\n  \"sig\": \"0x6c679b0b1f47bfc668bd67a217ffef3a0d645653b8ae7f25242bc1e7fac1511d7fc514e8a6c6199fc2c9a50648022cd549ba9a6a9def315a106fcc924a190394\",\r\n  \"version\": \"2\"\r\n}"
	error_3   = "{\r\n  \"address\": \"0x4667Fc33E26911Ba828cc9F1bB4A0F47B1e383B1\",\r\n  \"msg\": \"This is my wallet\",\r\n  \"sig\": \"0x6c679b0b1f47bfc668bd67a217ffef3a0d645653b8ae7f25242bc1e7fac1511d7fc514e8a6c6199fc2c9a50648022cd549ba9a6a9def315a106fcc924a1903942b\",\r\n  \"version\": \"2\"\r\n}"

	error_4   = "{\r\n  \"address\": \"0x5667Fc33E26911Ba828cc9F1bB4A0F47B1e383B1\",\r\n  \"msg\": \"This is my wallet\",\r\n  \"sig\": \"0x6c679b0b1f47bfc668bd67a217ffef3a0d645653b8ae7f25242bc1e7fac1511d7fc514e8a6c6199fc2c9a50648022cd549ba9a6a9def315a106fcc924a1903941c\",\r\n  \"version\": \"2\"\r\n}"
	error_4_1 = "{\r\n  \"address\": \"0x4667Fc33E26911Ba828cc9F1bB4A0F47B1e383B1\",\r\n  \"msg\": \"This is my wallet\",\r\n  \"sig\": \"0x8c679b0b1f47bfc668bd67a217ffef3a0d645653b8ae7f25242bc1e7fac1511d7fc514e8a6c6199fc2c9a50648022cd549ba9a6a9def315a106fcc924a1903941c\",\r\n  \"version\": \"2\"\r\n}"
	error_4_2 = "{\r\n  \"address\": \"0x4667Fc33E26911Ba828cc9F1bB4A0F47B1e383B1\",\r\n  \"msg\": \"This is my wallet1\",\r\n  \"sig\": \"0x8c679b0b1f47bfc668bd67a217ffef3a0d645653b8ae7f25242bc1e7fac1511d7fc514e8a6c6199fc2c9a50648022cd549ba9a6a9def315a106fcc924a1903941c\",\r\n  \"version\": \"2\"\r\n}"

	error_5 = "{\r\n  \"address\": \"0x466c33E26911Ba828cc9F1bB4A0F47B1e383B1\",\r\n  \"msg\": \"This is my wallet\",\r\n  \"sig\": \"0x6c679b0b1f47bfc668bd67a217ffef3a0d645653b8ae7f25242bc1e7fac1511d7fc514e8a6c6199fc2c9a50648022cd549ba9a6a9def315a106fcc924a1903941c\",\r\n  \"version\": \"2\"\r\n}"

	error_7 = "{\r\n  \"address\": \"0x5667Fc33E26911Ba828cc9F1bB4A0F47B1e383B1\",\r\n  \"msg\": \"This is my wallet\",\r\n  \"sig\": \"0x6c679b0b1f47bfc668bd67a217ffef3a0d645653b8ae7f25242bc1e7fac1511d7fc514e8a6c6199fc2c9a50648022cd549ba9a6a9def315a106fcc924a1903941c\",\r\n  \"version\": \"2\"\r\n}"

	success_1 = "{\r\n  \"address\": \"0xF422Ec881E87B934A165DB64132a87fbd1753daD\",\r\n  \"msg\": \"Test message waller\",\r\n  \"sig\": \"0x5c73e35d19d6656f826c82513a4523a8c789762bacfd1ce5127f24c1e61cd59f7779132c3a390294db158735e398c4e87b726b87bef44ad840a47ac6ca06ef8d1b\",\r\n  \"version\": \"2\"\r\n}"
	success_2 = "{\r\n  \"address\": \"0x4667Fc33E26911Ba828cc9F1bB4A0F47B1e383B1\",\r\n  \"msg\": \"This is my wallet\",\r\n  \"sig\": \"0x6c679b0b1f47bfc668bd67a217ffef3a0d645653b8ae7f25242bc1e7fac1511d7fc514e8a6c6199fc2c9a50648022cd549ba9a6a9def315a106fcc924a1903941c\",\r\n  \"version\": \"2\"\r\n}"
)

func TestCrossSign_Verification(t *testing.T) {
	err1 := CrossSignVerification(error_1)
	err11 := CrossSignVerification(error_1_1)

	err2 := CrossSignVerification(error_2)
	err3 := CrossSignVerification(error_3)

	err4 := CrossSignVerification(error_4)
	err41 := CrossSignVerification(error_4_1)
	err42 := CrossSignVerification(error_4_2)

	err5 := CrossSignVerification(error_5)
	err7 := CrossSignVerification(error_7)

	err8 := CrossSignVerification(success_1)
	err9 := CrossSignVerification(success_2)

	assert.Errorf(t, err1, err1.Error())
	assert.Errorf(t, err11, err11.Error())

	assert.Errorf(t, err2, err2.Error())
	assert.Errorf(t, err3, err3.Error())

	assert.Errorf(t, err4, err4.Error())
	assert.Errorf(t, err41, err41.Error())
	assert.Errorf(t, err42, err42.Error())

	assert.Errorf(t, err5, err5.Error())

	assert.Errorf(t, err7, err7.Error())

	assert.Equal(t, nil, err8)
	assert.Equal(t, nil, err9)
}
