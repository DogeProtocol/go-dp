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

//This file was added for go-dogep project (Doge Protocol Platform)

package cryptopq

import (
	"bytes"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/cryptopq/oqs"
	"io/ioutil"
	"math/big"
	"os"
	"reflect"
	"testing"
)

var testAddrHex = "09acbabfc8c921d09be6ecf48b6f6c4cb5e6b6e18e474696da80dfe2b452464d57a46941b10031a1d7ea426f86fb9a88202b21e100dbf022d9877e50a3a9e2599be23541821c56a1fc3c057447520c62fbdb22abbe0bfcfc79c204b25c54d714474f7954895d95211fff623018163769ac809ed8096e5d1e7f86f978b8e12a369a75b7de26273e5565870a1617b0cdcf4a16945b96490fc6b1aa53bffee0d04b50cea0f16b90e60fcd8abe910b8c9d471d5db23b4d1f781e88142a184be499e9a29a6fe27056da5270dae130a1dd6264371c9540bc40e2d468c725e0c4d5faee190e886ce3c22ecc4f55c2d15ec2c9422468781397091b78389a1e44126fda8b9110a252affce49296e16f31d9c858d4ef0156b50f878813152ee48d95095d845f6d2202cf6298043da7faf6c62d5708c7c662dd284c1cfa55bacb0b9a4866db730ff541a14a8a51982655a5be11eb45b39a01741634d8538b3453159dd11f7216fea94685a7bb5a32706326fdb2c950d1120a2eb45125039a5cc137ac304079068762c90866085bd210c564a5bfa047a54ab5f4ac5eac3870dc3ad7d7e8ba2905e0f9d4715c6a43c9ab6283a3e667d7124ab51b252398d7cdd83bd93afb1b23724e2debd88913919a4a7bdb75dc2636f05c175f86fab74d9bea5cca2da3813947308f88a5f844078b2fc9ad910f95c8490a302494225aed97905d49465980b6be1464e3c0af869d042bd046424c328e880aa8623714328139c49a6c777eb45090915e628282b673d90771a5b0604b882d8290b0a80fae2189c5a8e9f3990b6e614f713f86a1b8ca7837740ff83a56b9aac2a6cb56a0551198ebd3842917b324a4d4309483bb61515d7384008a1286bb4df6c5706b80d6a9cc82d73a993e7ccd2ab6f1805803ec64b6939d47bad8e37a15c5a88bd29e020fac3438bdc9d41dc08654440da91cc2812a44cf46f5096a4a80d3ea4d49c2432edc975a4e5097a176820e0aa90a760d39a3248595057ed5add24c3796cffb3893f028565a65f38b58a9faca8a0954f5a6057f5a7e288f58d8b4657c68ce5bbe7130552b1051a54a30bad35a855a29f6b30af9358ea85a7a644817fd33a0984f153d49df034098a5728d722bdbef3e914416967ec56556d667d8161a412f19d71a0c4df1bf6cd3eb4b18167798de4abb6ab5b9d699c2582b12862a6cca4667e8b00da2503252231b32e33831a77e9906af13db9d9713b93fc0cf933c3ffab95265ebd70d6c265b7fcb4b"
var testPrivHex = "59f820410bfff8e7fdc3efdfbcf3c07ff4223d040f060030441800b9ebff3af81205f7df4303f0bf07d07f048f05185e7d086f7df81fc617cf03e8117e0fa184f46f4518027dfb9f7cf3f0c407ef3908213e13a080e8300307f0c50b9ec6f790400c1f0208407e0420bbfc5f000fefc700508aec4ffbebe08613f0c007e03b1c003f10100007d17f073fbd00300203f1c3083f4708003b1840ff143fbd2020b703df43ffe23dd81fbd1ff03c081f3b0050bff42f440b318007f0400c11030c00bde410c4f8117efc3efe142182dfae7a140ffedc3e42002082f771060befbfec608317f03b07d0befbf088e420bdfc3081179046f070411400bb184fc20040bb080084001f41ebc07833e10503fdbbf83f83f431c8003007f7dfc507e03c04503d0bfef90c51c5ebf1f9ec0fbd101ebe0bb0420bd041d40fb7ec1f01081ffc07d0fdebffbbfbaf7df06f42f7a2c1ffffffd83ebdf77001fc5045248fff0feefff85fc0188f42f83200f060410bd07e080fc2001fc3e8007f0ff0c3081f02d800c0102005fc7ffff84fbd1020be042ebc0fbf46042f40045df61bff02004f7bfbe1420fbe7defa143f04143041f82f43f83fbf07f040ffdf40141eb6081fc3fc10fe03ef030ff07a0fbec403d0c0fc4083087e44ffc03ffff17f040f7d202f7c0c0002e3bf7eec10fdf810c7f00101f41181efc07fff9f04ffe084f8100007c17c0bdebcf02003fbf0cad3e07e1c21fd002103139f7d103f40ffb0391001c1fffebafbf07bf8204207a13f0bd002040142f44105f7b181e80ec6081fc5f4103a14108003c102006ebbfbb0420fff861c6f80f85f00fbce07fc10020fa0bc0810fd084e860810f8142fbf140007004f410ba0c100013cf06006fbe14608700913f1430791bcf02f031430c0049fbc0401450451bdffa0fdffc101ffa146e0afc0fb4f4510507b07c03f1030fe039ff6076f3bf451c4001fc004af85038dc203ef80e3b085082000f79003efb184081003e7ef43ffbfbaff9fff04107cf000800450c30c50462401c4f431bbf00ef8e4103e07f0840c6f7e1c2042e7a006e7e17d003f3fe1300f1c4e50a0508121d0403e2dbe3fdde11131a08010202eb130a22ee0a1202f00d1102ed2dd407ed1e061fdc01eb03261822d43325fd2602f6fadd22ee11000112e7fc450134060122faf1dbe417df1ff93721efe90de0e6e90030eb2df90f24f3f2fbf41618fbe5f4dff312fafe0cf8d9e9da31d5ddfd00fbd2f8e2f704d40a200c1317fad42c1ae31e09f4e2efe0c6d421bcdf34fc0cf4dc01e90f0123df09f63f0fc4032801fa20d5df0cf91a0113f8c62af61bda050900100612f22a0e01f8d70ee7f509ec2c38f1e816cb0de916561604f0f10a2a22d91ffb0afc1ffb1800fcf5202ffdfb14d923d7f5410f0feedf25d320f6ee3c2ce6fd3b05d7e8171106e409d3f922cef1131437e91e0fd1f2110d211af8e7b60bdee42a1420200c030ad1dcf5f10df402e6fdcf0712f8fe130805ce09c71f14d2f5f719e8dc070514f52af611ed07f4fa0a1e110e280b07dd14f3d2e829f0dc04f4c51bfafbf41dfb041301142008f3f600f43808fd0ef0dd04050a1506dff72afbc74c010a1efdf5d4ff0c0b0af4f506011af30cfaef07d006fcecec0b1cfc080cd2de0ee9edfcfff10ffcfb191c2ffcfee8f5fcedd7f0f20afed8e5fceff8f20b2803ebf4fbfaf0dafcf717f0d5e1e90028ec0917f8f4100d09f21ff21ee1f5f209e6f40528cfec10f007cb072520fe0edf1b23ede92907ea280ddd0eebf4f910e712f9fa"

// These tests are sanity checks.
// They should ensure that we don't e.g. use Sha3-224 instead of Sha3-256
// and that the sha3 library uses keccak-f permutation.
func TestKeccak256Hash(t *testing.T) {
	msg := []byte("abc")
	exp, _ := hex.DecodeString("4e03657aea45a94fc7d47ba826c8d667c0d1e6e33a64a036ec44f58fa12d6c45")
	checkhash(t, "Sha3-256-array", func(in []byte) []byte { h := crypto.Keccak256Hash(in); return h[:] }, msg, exp)
}

func TestKeccak256Hasher(t *testing.T) {
	msg := []byte("abc")
	exp, _ := hex.DecodeString("4e03657aea45a94fc7d47ba826c8d667c0d1e6e33a64a036ec44f58fa12d6c45")
	hasher := crypto.NewKeccakState()
	checkhash(t, "Sha3-256-array", func(in []byte) []byte { h := crypto.HashData(hasher, in); return h[:] }, msg, exp)
}

func TestToOQSErrors(t *testing.T) {
	if _, err := HexToOQS("0000000000000000000000000000000000000000000000000000000000000000"); err == nil {
		t.Fatal("HexToOQS should've returned error")
	}
	if _, err := HexToOQS("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"); err == nil {
		t.Fatal("HexToOQS should've returned error")
	}
}

func BenchmarkSha3(b *testing.B) {
	a := []byte("hello world")
	for i := 0; i < b.N; i++ {
		crypto.Keccak256(a)
	}
}

func TestUnmarshalPubkey(t *testing.T) {
	key, err := UnmarshalPubkey(nil)
	if err != oqs.ErrInvalidPublicKeyLen || key != nil {
		t.Fatalf("expected error1, got %v, %v", err, key)
	}
	key, err = UnmarshalPubkey([]byte{1, 2, 3})
	if err != oqs.ErrInvalidPublicKeyLen || key != nil {
		t.Fatalf("expected error, got %v, %v", err, key)
	}
	var (
		enc, _ = hexutil.Decode("0x09aba20b5b6645d1090d60786e612d3f4ad126cb1c9584386a0499297c5c4d7db100a5d27e592953c41dd8b53553a3e42c969ace8bb809ed8e4f44393a9b2a24704edb94f5d58422edd656a191228cede858b6e9c9609821127e9ac17e4ae6282d6e16378cf346ab94570b4ae4582a249aa0f5d1585f2da471ae8c09000c0335183cca0012c7b5f4f253dcda214afa8d107e04b402e56947a1021f366ee6016c0a6342c6d997ec03f919aebc69f752ed816c225623292c00589db523256e5dd478f4f1c1007cc48608e758a527ad1905470f4d4e0a690b618c4f348c94bb984b040c4245d6bcc6f802516cb828c6199766f3d213a1494a7c4971f26f8384350e8aaae0372a22149ad62a9092a0d35b7b5cc834d92642261fd43da1f3919d135a666ccb601f501c1181d67451e4f4255db6aa541ec080426279275b35a9629961b6247350851b12411b242378c3d7eb96e108246ec65cc9af37f61bd2814e54853e8e25015e9e780ed644c58a0574bdc0d2bc82945ce4bf82564159b6028d751fa47f7664ff77d29eab337d1504ae9edd2cfeab45c44975277abc77dd619ba93028e31952047c95c475ebaa5c60d53754711a9d29842c0603eec37481c4c31e4d0e04cc6b73d25cf73f06d385846adeb092b13967dfb924cac8e8801d883a206522796909068d8c622490112015b9ba9d52924dd830b168de0b4b20b2915257f477ad1215e162e18cd6b93958e16e76e4490078e8907709b347729461a8a92ce169aea58834f70830ec9a905d893097887f79bcbfd999a4774f898ee1b1a713ad8ca669f865713d01c3a5785040831e7cb840a6a6b975d915abce0d97e18ad03be3eca3ae180895236d9de59044968b9aad0566531e4bb385253e27c6ae3a13b0e7a44898d9b6ad3d30f75e7582da70faaae06ba4d47ed07b5a5ca73022991f0b631be58348c423513d2c9c835440e734a13c36fd4e1c7546034428551c7b883360569d191c5247b48817919070a2313066c73a79c8e5afde80b4510fe99d56ccac60f5bb1427553a209784402c3493784b87f3a04351164609e3450156498cb8a496434bbe086aadd96a686c0d2aa147298ce6f4884a0cc59c1a11cba96e2775ac5b2149e91a8518315af2cbd7a270034301c98f402c607891a08e2b181d66fa8dc48ef9aec69683816472f1cacc83dc2745c7ca88e3b62fc9da518a0f5e65f204b786a103a7d9322641e3227d19b9ca00979294a513d1b44db22996a93119aa9b")
		dec    = &oqs.PublicKey{
			N: new(big.Int).SetBytes(hexutil.MustDecode("0x09aba20b5b6645d1090d60786e612d3f4ad126cb1c9584386a0499297c5c4d7db100a5d27e592953c41dd8b53553a3e42c969ace8bb809ed8e4f44393a9b2a24704edb94f5d58422edd656a191228cede858b6e9c9609821127e9ac17e4ae6282d6e16378cf346ab94570b4ae4582a249aa0f5d1585f2da471ae8c09000c0335183cca0012c7b5f4f253dcda214afa8d107e04b402e56947a1021f366ee6016c0a6342c6d997ec03f919aebc69f752ed816c225623292c00589db523256e5dd478f4f1c1007cc48608e758a527ad1905470f4d4e0a690b618c4f348c94bb984b040c4245d6bcc6f802516cb828c6199766f3d213a1494a7c4971f26f8384350e8aaae0372a22149ad62a9092a0d35b7b5cc834d92642261fd43da1f3919d135a666ccb601f501c1181d67451e4f4255db6aa541ec080426279275b35a9629961b6247350851b12411b242378c3d7eb96e108246ec65cc9af37f61bd2814e54853e8e25015e9e780ed644c58a0574bdc0d2bc82945ce4bf82564159b6028d751fa47f7664ff77d29eab337d1504ae9edd2cfeab45c44975277abc77dd619ba93028e31952047c95c475ebaa5c60d53754711a9d29842c0603eec37481c4c31e4d0e04cc6b73d25cf73f06d385846adeb092b13967dfb924cac8e8801d883a206522796909068d8c622490112015b9ba9d52924dd830b168de0b4b20b2915257f477ad1215e162e18cd6b93958e16e76e4490078e8907709b347729461a8a92ce169aea58834f70830ec9a905d893097887f79bcbfd999a4774f898ee1b1a713ad8ca669f865713d01c3a5785040831e7cb840a6a6b975d915abce0d97e18ad03be3eca3ae180895236d9de59044968b9aad0566531e4bb385253e27c6ae3a13b0e7a44898d9b6ad3d30f75e7582da70faaae06ba4d47ed07b5a5ca73022991f0b631be58348c423513d2c9c835440e734a13c36fd4e1c7546034428551c7b883360569d191c5247b48817919070a2313066c73a79c8e5afde80b4510fe99d56ccac60f5bb1427553a209784402c3493784b87f3a04351164609e3450156498cb8a496434bbe086aadd96a686c0d2aa147298ce6f4884a0cc59c1a11cba96e2775ac5b2149e91a8518315af2cbd7a270034301c98f402c607891a08e2b181d66fa8dc48ef9aec69683816472f1cacc83dc2745c7ca88e3b62fc9da518a0f5e65f204b786a103a7d9322641e3227d19b9ca00979294a513d1b44db22996a93119aa9b")),
		}
	)
	key, err = UnmarshalPubkey(enc)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !reflect.DeepEqual(key, dec) {
		t.Fatal("wrong result")
	}
}

func GenKeyPairTest() (pubkey, privkey []byte) {
	key, err := GenerateKey()

	if err != nil {
		panic(err)
	}
	pubkey = make([]byte, oqs.PublicKeyLen)
	pub := key.N.Bytes()
	copy(pubkey[oqs.PublicKeyLen-len(pub):], pub)

	privkey = make([]byte, oqs.PrivateKeyLen)
	priv := key.D.Bytes()
	copy(privkey[oqs.PrivateKeyLen-len(priv):], priv)

	return pubkey, privkey
}

func TestSign(t *testing.T) {
	key, err := HexToOQS(testPrivHex)
	if err != nil {
		t.Errorf("HexToOQS")
	}

	pubBytes, err := hex.DecodeString(testAddrHex)
	if err != nil {
		t.Errorf("DecodeString")
	}

	addr := common.BytesToAddress(crypto.Keccak256(pubBytes[1:])[12:])

	msg := crypto.Keccak256([]byte("foo"))
	data, err := oqs.ExportPrivateKey(key)
	if err != nil {
		t.Errorf("ExportPrivateKey")
	}

	sig, err := oqs.Sign(msg, data)
	if err != nil {
		t.Errorf("Sign error: %s", err)
	}

	recoveredPub, err := RecoverPublicKey(msg, sig)
	if err != nil {
		t.Errorf("ECRecover error: %s", err)
	}
	pubKey, _ := UnmarshalPubkey(recoveredPub)
	recoveredAddr, err := PubkeyToAddress(*pubKey)
	if err != nil {
		t.Errorf("PubkeyToAddress error: %s", err)
	}

	if addr != recoveredAddr {
		t.Errorf("Address mismatch: want: %x have: %x", addr, recoveredAddr)
	}

	// should be equal to SigToPub
	recoveredPub2, err := SigToPub(msg, sig)
	if err != nil {
		t.Errorf("SigToPub error: %s", err)
	}
	recoveredAddr2, err := PubkeyToAddress(*recoveredPub2)
	if err != nil {
		t.Errorf("PubkeyToAddress error: %s", err)
	}

	if addr != recoveredAddr2 {
		t.Errorf("Address mismatch: want: %x have: %x", addr, recoveredAddr2)
	}
}

func TestInvalidSign(t *testing.T) {
	if _, err := oqs.Sign(make([]byte, 1), nil); err == nil {
		t.Errorf("expected sign with hash 1 byte to error")
	}
	if _, err := oqs.Sign(make([]byte, 33), nil); err == nil {
		t.Errorf("expected sign with hash 33 byte to error")
	}
}

func TestNewContractAddress(t *testing.T) {
	pubBytes, err := hex.DecodeString(testAddrHex)
	if err != nil {
		t.Errorf("DecodeString")
	}
	addr := common.BytesToAddress(crypto.Keccak256(pubBytes[1:]))
	//addr := common.HexToAddress(testAddrHex)
	pubKey, _ := UnmarshalPubkey(pubBytes)
	genAddr, err := PubkeyToAddress(*pubKey)
	if err != nil {
		t.Errorf("PubkeyToAddress")
	}
	// sanity check before using addr to create contract address
	checkAddr(t, genAddr, addr)

	caddr0 := crypto.CreateAddress(addr, 0)
	caddr1 := crypto.CreateAddress(addr, 1)
	caddr2 := crypto.CreateAddress(addr, 2)
	t.Logf("caddr0: %s", caddr0)
	t.Logf("caddr1: %s", caddr1)
	t.Logf("caddr2: %s", caddr2)
}

func TestSaveOQS(t *testing.T) {
	f, err := ioutil.TempFile("", "saveOQS_test.*.txt")
	if err != nil {
		t.Fatal(err)
	}
	file := f.Name()
	f.Close()
	defer os.Remove(file)

	key, _ := HexToOQS(testPrivHex)
	if err := SaveOQS(file, key); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadOQS(file)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(key, loaded) {
		t.Fatal("loaded key not equal to saved key")
	}
}

func checkhash(t *testing.T, name string, f func([]byte) []byte, msg, exp []byte) {
	sum := f(msg)
	if !bytes.Equal(exp, sum) {
		t.Fatalf("hash %s mismatch: want: %x have: %x", name, exp, sum)
	}
}

func checkAddr(t *testing.T, addr0, addr1 common.Address) {
	if addr0 != addr1 {
		t.Fatalf("address mismatch: want: %x have: %x", addr0, addr1)
	}
}
