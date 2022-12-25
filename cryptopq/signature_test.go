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

//This file was added for go-dogep project (Doge Protocol Platform)

package cryptopq

import (
	"bytes"
	"github.com/ethereum/go-ethereum/cryptopq/oqs"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var (
	testmsg     = hexutil.MustDecode("0x68692074686572656f636b636861696e")
	testsig     = hexutil.MustDecode("0x39c373f34925c9254b05b533ebb167d31fa28a4f88f155df8000ac54ea452474ccfe43a9dc85ed1585c3f2bbfe0c63377bcfa03a042832da57a76f2f5cb4ef5ec23f9077ab588e444c4d169b663613c875891d731d012fb3420b55afc5f4baf47217b7afc2dc4f1dc7a36ba4511dbe5abd9a41551d42e3d3c416891b26bccb2d1d244d28fd2b1cfc29902996562c2b08b5628773709aff66a94662bc12bacafda5452df66f0891cd3198c7edc26f49537fa86b6c9f3a4a3723c43866192aaf6dd5c5e9e063d08407ece3362ba4f9dd5f75cda2dd9f50792a527f1361fa58864ef19d31732d45df88a3f0db95c7276795bf962e3611bcabda12f24d91a77db3e75fc89cc1d504d28cba2b2b8a5451cf0eeede93e9fc87aba79c47cdb35ae2b1dc5873472ee923b01394596a7c9297cbf5a6134ef7027110bc62d1a498844cd9bc6f554154e7f3a3fbcddf55e3aae462714ea823e8cb84e753f2ece7137341ec42c6dd6be5e72c27d9fab4b818f914d791f1858d2504eaba365dda0f76efa891cdc19665b38d8c15f87211e284eb5d67c278bffbeaf1cafaeb4289af5bf27e90160ef261177d3ac9bb0be658ac867e9269aa9724c51e783f79f1c2a8b9c4ea9acac3b35861435874c65b61e24c2aac7304a75e50a8deb94de4c6e3ba45ce411b2f9e8ce472248d10f7a3b72dcbbb947c0cfa419e758805276879671edc6b3cb025efaa6f5aed8c8a65114329ee0b49d765940322b4a975f5bd7c64ddf63683e987948fd51353486dda52a489e8a99e5cfdcfe4b3ca0c42b88dd7f4f5845a5adad11b04571225d9f85fcf6894c26c4e8f9b228034c8f98d1804eb771a651025e2e9a02a5d0ea57d517478297a67df6b9af48ae1dc70d72555699911bc81d8bb510d014bd645efee8b1a8a2e146455ce5a9009aba20b5b6645d1090d60786e612d3f4ad126cb1c9584386a0499297c5c4d7db100a5d27e592953c41dd8b53553a3e42c969ace8bb809ed8e4f44393a9b2a24704edb94f5d58422edd656a191228cede858b6e9c9609821127e9ac17e4ae6282d6e16378cf346ab94570b4ae4582a249aa0f5d1585f2da471ae8c09000c0335183cca0012c7b5f4f253dcda214afa8d107e04b402e56947a1021f366ee6016c0a6342c6d997ec03f919aebc69f752ed816c225623292c00589db523256e5dd478f4f1c1007cc48608e758a527ad1905470f4d4e0a690b618c4f348c94bb984b040c4245d6bcc6f802516cb828c6199766f3d213a1494a7c4971f26f8384350e8aaae0372a22149ad62a9092a0d35b7b5cc834d92642261fd43da1f3919d135a666ccb601f501c1181d67451e4f4255db6aa541ec080426279275b35a9629961b6247350851b12411b242378c3d7eb96e108246ec65cc9af37f61bd2814e54853e8e25015e9e780ed644c58a0574bdc0d2bc82945ce4bf82564159b6028d751fa47f7664ff77d29eab337d1504ae9edd2cfeab45c44975277abc77dd619ba93028e31952047c95c475ebaa5c60d53754711a9d29842c0603eec37481c4c31e4d0e04cc6b73d25cf73f06d385846adeb092b13967dfb924cac8e8801d883a206522796909068d8c622490112015b9ba9d52924dd830b168de0b4b20b2915257f477ad1215e162e18cd6b93958e16e76e4490078e8907709b347729461a8a92ce169aea58834f70830ec9a905d893097887f79bcbfd999a4774f898ee1b1a713ad8ca669f865713d01c3a5785040831e7cb840a6a6b975d915abce0d97e18ad03be3eca3ae180895236d9de59044968b9aad0566531e4bb385253e27c6ae3a13b0e7a44898d9b6ad3d30f75e7582da70faaae06ba4d47ed07b5a5ca73022991f0b631be58348c423513d2c9c835440e734a13c36fd4e1c7546034428551c7b883360569d191c5247b48817919070a2313066c73a79c8e5afde80b4510fe99d56ccac60f5bb1427553a209784402c3493784b87f3a04351164609e3450156498cb8a496434bbe086aadd96a686c0d2aa147298ce6f4884a0cc59c1a11cba96e2775ac5b2149e91a8518315af2cbd7a270034301c98f402c607891a08e2b181d66fa8dc48ef9aec69683816472f1cacc83dc2745c7ca88e3b62fc9da518a0f5e65f204b786a103a7d9322641e3227d19b9ca00979294a513d1b44db22996a93119aa9b")
	testpubkey  = hexutil.MustDecode("0x09aba20b5b6645d1090d60786e612d3f4ad126cb1c9584386a0499297c5c4d7db100a5d27e592953c41dd8b53553a3e42c969ace8bb809ed8e4f44393a9b2a24704edb94f5d58422edd656a191228cede858b6e9c9609821127e9ac17e4ae6282d6e16378cf346ab94570b4ae4582a249aa0f5d1585f2da471ae8c09000c0335183cca0012c7b5f4f253dcda214afa8d107e04b402e56947a1021f366ee6016c0a6342c6d997ec03f919aebc69f752ed816c225623292c00589db523256e5dd478f4f1c1007cc48608e758a527ad1905470f4d4e0a690b618c4f348c94bb984b040c4245d6bcc6f802516cb828c6199766f3d213a1494a7c4971f26f8384350e8aaae0372a22149ad62a9092a0d35b7b5cc834d92642261fd43da1f3919d135a666ccb601f501c1181d67451e4f4255db6aa541ec080426279275b35a9629961b6247350851b12411b242378c3d7eb96e108246ec65cc9af37f61bd2814e54853e8e25015e9e780ed644c58a0574bdc0d2bc82945ce4bf82564159b6028d751fa47f7664ff77d29eab337d1504ae9edd2cfeab45c44975277abc77dd619ba93028e31952047c95c475ebaa5c60d53754711a9d29842c0603eec37481c4c31e4d0e04cc6b73d25cf73f06d385846adeb092b13967dfb924cac8e8801d883a206522796909068d8c622490112015b9ba9d52924dd830b168de0b4b20b2915257f477ad1215e162e18cd6b93958e16e76e4490078e8907709b347729461a8a92ce169aea58834f70830ec9a905d893097887f79bcbfd999a4774f898ee1b1a713ad8ca669f865713d01c3a5785040831e7cb840a6a6b975d915abce0d97e18ad03be3eca3ae180895236d9de59044968b9aad0566531e4bb385253e27c6ae3a13b0e7a44898d9b6ad3d30f75e7582da70faaae06ba4d47ed07b5a5ca73022991f0b631be58348c423513d2c9c835440e734a13c36fd4e1c7546034428551c7b883360569d191c5247b48817919070a2313066c73a79c8e5afde80b4510fe99d56ccac60f5bb1427553a209784402c3493784b87f3a04351164609e3450156498cb8a496434bbe086aadd96a686c0d2aa147298ce6f4884a0cc59c1a11cba96e2775ac5b2149e91a8518315af2cbd7a270034301c98f402c607891a08e2b181d66fa8dc48ef9aec69683816472f1cacc83dc2745c7ca88e3b62fc9da518a0f5e65f204b786a103a7d9322641e3227d19b9ca00979294a513d1b44db22996a93119aa9b")
	testpubkeyc = hexutil.MustDecode("0x09aba20b5b6645d1090d60786e612d3f4ad126cb1c9584386a0499297c5c4d7db100a5d27e592953c41dd8b53553a3e42c969ace8bb809ed8e4f44393a9b2a24704edb94f5d58422edd656a191228cede858b6e9c9609821127e9ac17e4ae6282d6e16378cf346ab94570b4ae4582a249aa0f5d1585f2da471ae8c09000c0335183cca0012c7b5f4f253dcda214afa8d107e04b402e56947a1021f366ee6016c0a6342c6d997ec03f919aebc69f752ed816c225623292c00589db523256e5dd478f4f1c1007cc48608e758a527ad1905470f4d4e0a690b618c4f348c94bb984b040c4245d6bcc6f802516cb828c6199766f3d213a1494a7c4971f26f8384350e8aaae0372a22149ad62a9092a0d35b7b5cc834d92642261fd43da1f3919d135a666ccb601f501c1181d67451e4f4255db6aa541ec080426279275b35a9629961b6247350851b12411b242378c3d7eb96e108246ec65cc9af37f61bd2814e54853e8e25015e9e780ed644c58a0574bdc0d2bc82945ce4bf82564159b6028d751fa47f7664ff77d29eab337d1504ae9edd2cfeab45c44975277abc77dd619ba93028e31952047c95c475ebaa5c60d53754711a9d29842c0603eec37481c4c31e4d0e04cc6b73d25cf73f06d385846adeb092b13967dfb924cac8e8801d883a206522796909068d8c622490112015b9ba9d52924dd830b168de0b4b20b2915257f477ad1215e162e18cd6b93958e16e76e4490078e8907709b347729461a8a92ce169aea58834f70830ec9a905d893097887f79bcbfd999a4774f898ee1b1a713ad8ca669f865713d01c3a5785040831e7cb840a6a6b975d915abce0d97e18ad03be3eca3ae180895236d9de59044968b9aad0566531e4bb385253e27c6ae3a13b0e7a44898d9b6ad3d30f75e7582da70faaae06ba4d47ed07b5a5ca73022991f0b631be58348c423513d2c9c835440e734a13c36fd4e1c7546034428551c7b883360569d191c5247b48817919070a2313066c73a79c8e5afde80b4510fe99d56ccac60f5bb1427553a209784402c3493784b87f3a04351164609e3450156498cb8a496434bbe086aadd96a686c0d2aa147298ce6f4884a0cc59c1a11cba96e2775ac5b2149e91a8518315af2cbd7a270034301c98f402c607891a08e2b181d66fa8dc48ef9aec69683816472f1cacc83dc2745c7ca88e3b62fc9da518a0f5e65f204b786a103a7d9322641e3227d19b9ca00979294a513d1b44db22996a93119aa9b")
)

func TestVerifySignature(t *testing.T) {
	sig := testsig[:len(testsig)]

	if VerifySignature(nil, testmsg, sig) {
		t.Errorf("signature valid with no key")
	}
	if VerifySignature(testpubkey, nil, sig) {
		t.Errorf("signature valid with no message")
	}
	if VerifySignature(testpubkey, testmsg, nil) {
		t.Errorf("nil signature valid")
	}
	if VerifySignature(testpubkey, testmsg, append(common.CopyBytes(sig), 1, 2, 3)) {
		t.Errorf("signature valid with extra bytes at the end")
	}
	if VerifySignature(testpubkey, testmsg, sig[:len(sig)-2]) {
		t.Errorf("signature valid even though it's incomplete")
	}
	wrongkey := common.CopyBytes(testpubkey)
	wrongkey[10]++
	if VerifySignature(wrongkey, testmsg, sig) {
		t.Errorf("signature valid with with wrong public key")
	}
}

func GenerateKeyPairTest() (pubkey, privkey []byte) {
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

func TestDecompressPubkey(t *testing.T) {
	key, err := DecompressPubkey(testpubkeyc)
	if err != nil {
		t.Fatal(err)
	}
	uncompressed, err := FromOQSPub(key)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(uncompressed, testpubkey) {
		t.Errorf("wrong public key result: got %x, want %x", uncompressed, testpubkey)
	}
	if _, err := DecompressPubkey(nil); err == nil {
		t.Errorf("no error for nil pubkey")
	}
	if _, err := DecompressPubkey(testpubkeyc[:5]); err == nil {
		t.Errorf("no error for incomplete pubkey")
	}
	if _, err := DecompressPubkey(append(common.CopyBytes(testpubkeyc), 1, 2, 3)); err == nil {
		t.Errorf("no error for pubkey with extra bytes at the end")
	}
}

func TestCompressPubkey(t *testing.T) {
	key := &oqs.PublicKey{
		N: new(big.Int).SetBytes(testpubkey),
	}
	compressed := CompressPubkey(key)
	if !bytes.Equal(compressed, testpubkeyc) {
		t.Errorf("wrong public key result: got %x, want %x", compressed, testpubkeyc)
	}
}

func TestPubkeyRandom(t *testing.T) {
	const runs = 200

	for i := 0; i < runs; i++ {
		key, err := GenerateKey()
		if err != nil {
			t.Fatalf("iteration %d: %v", i, err)
		}
		pubkey2, err := DecompressPubkey(CompressPubkey(&key.PublicKey))
		if err != nil {
			t.Fatalf("iteration %d: %v", i, err)
		}
		if !reflect.DeepEqual(key.PublicKey, *pubkey2) {
			t.Fatalf("iteration %d: keys not equal", i)
		}
	}
}

func BenchmarkRecover(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := RecoverPublicKey(testmsg, testsig); err != nil {
			b.Fatal("RecoverPublicKey error", err)
		}
	}
}

func BenchmarkVerifySignature(b *testing.B) {
	sig := testsig[:len(testsig)-1] // remove recovery id
	for i := 0; i < b.N; i++ {
		if !VerifySignature(testpubkey, testmsg, sig) {
			b.Fatal("verify error")
		}
	}
}

func BenchmarkDecompressPubkey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := DecompressPubkey(testpubkeyc); err != nil {
			b.Fatal(err)
		}
	}
}
