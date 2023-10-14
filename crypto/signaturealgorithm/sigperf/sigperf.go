package main

import (
	"fmt"
	"github.com/DogeProtocol/dp/common/hexutil"
	"github.com/DogeProtocol/dp/crypto/falcon"
	"github.com/DogeProtocol/dp/crypto/hybrid"
	"github.com/DogeProtocol/dp/crypto/oqs"
	"github.com/DogeProtocol/dp/crypto/signaturealgorithm"
	"os"
	"runtime/pprof"
	"strconv"
	"time"
)

func main() {
	sig1 := hybrid.CreateHybridSig()
	SigPerf("hybrid", sig1)

	sig2 := oqs.InitDilithium()
	SigPerf("dilithiumoqs", sig2)

	sig3 := falcon.CreateFalconSig()
	SigPerf("falcon", sig3)

	sig4 := oqs.InitFalcon()
	SigPerf("falconoqs", sig4)
}

func SigPerf(name string, sig signaturealgorithm.SignatureAlgorithm) {
	fmt.Println("SigPerf", name)
	f, err := os.Create(name + ".prof")
	if err != nil {
		fmt.Println("profiling failed", err)
		return
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	keypair, err := sig.GenerateKey()
	if err != nil {
		fmt.Println("GenerateKey failed", err)
		return
	}
	testmsg := hexutil.MustDecode("0x68692074686572656f636b636861696e62626262626262626262626262626262")
	digestHash := []byte(testmsg)

	signature, err := sig.Sign(digestHash, keypair)
	if err != nil {
		fmt.Println("Sign failed", err)
		return
	}
	pubBytes, err := sig.SerializePublicKey(&keypair.PublicKey)
	if err != nil {
		fmt.Println("SerializePublicKey failed", err)
		return
	}
	if len(os.Args) < 2 {
		fmt.Println("enter test iteration count")
		return
	}
	count, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Atoi failed", err)
		return
	}

	fmt.Println("Running verify...", count)
	start := time.Now()
	for i := 0; i < count; i++ {
		if sig.Verify(pubBytes, digestHash, signature) != true {
			fmt.Println("Verify failed", err)
			return
		}
	}
	duration := time.Since(start)
	fmt.Println("Verify Time Taken", duration)
}
