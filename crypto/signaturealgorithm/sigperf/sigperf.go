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
	"sync"
	"time"
)

func main() {
	sig1 := hybrid.CreateHybridSig(false)
	sig2 := oqs.InitDilithium()
	sig3 := falcon.CreateFalconSig()
	sig4 := oqs.InitFalcon()
	sig5 := hybrid.CreateHybridSig(true)

	if len(os.Args) > 2 {
		var wg sync.WaitGroup

		fmt.Println("Multi routine test start")
		for i := 0; i <= 32; i++ {
			wg.Add(1)
			go SigPerf("hybrid", sig1, &wg)
		}
		wg.Wait()

		for i := 0; i <= 32; i++ {
			wg.Add(1)
			go SigPerf("dilithiumoqs", sig2, &wg)
		}
		wg.Wait()

		for i := 0; i <= 32; i++ {
			wg.Add(1)
			go SigPerf("falcon", sig3, &wg)
		}
		wg.Wait()

		for i := 0; i <= 32; i++ {
			wg.Add(1)
			go SigPerf("falconoqs", sig4, &wg)
		}
		wg.Wait()

		for i := 0; i <= 32; i++ {
			wg.Add(1)
			go SigPerf("hybrid native", sig5, &wg)
		}
		wg.Wait()
	}

	fmt.Println("Multi routine test done")

	SigPerf("hybrid", sig1, nil)
	SigPerf("dilithiumoqs", sig2, nil)
	SigPerf("falcon", sig3, nil)
	SigPerf("falconoqs", sig4, nil)
	SigPerf("hybrid native", sig5, nil)
}

func SigPerf(name string, sig signaturealgorithm.SignatureAlgorithm, wg *sync.WaitGroup) {
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

	fmt.Println("Running verify...", name, count)
	start := time.Now()
	for i := 0; i < count; i++ {
		if sig.Verify(pubBytes, digestHash, signature) != true {
			fmt.Println("Verify failed", err)
			return
		}
	}
	duration := time.Since(start)
	fmt.Println("Verify Time Taken", "sigalg", name, "iterations", count, "totaltime ms", duration.Milliseconds(), "avg time ms", float64(duration.Milliseconds())/float64(count))
	if wg != nil {
		wg.Done()
	}
}
