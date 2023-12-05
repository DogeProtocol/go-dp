package main

import (
	"fmt"
	"github.com/DogeProtocol/dp/common/hexutil"
	"github.com/DogeProtocol/dp/crypto/hybrideds"
	"os"
	"runtime/pprof"
	"strconv"
	"time"
)

func main() {
	fmt.Println("Perf starting 1")
	f, err := os.Create("hybrid.prof")
	if err != nil {
		fmt.Println("profiling failed", err)
		return
	}
	fmt.Println("Perf starting 2")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	sig := hybrideds.CreateHybridedsSig(true)
	keypair, err := sig.GenerateKey()
	if err != nil {
		fmt.Println("GenerateKey failed", err)
		return
	}
	fmt.Println("Perf starting 3")
	testmsg := hexutil.MustDecode("0x68692074686572656f636b636861696e62626262626262626262626262626262")
	digestHash := []byte(testmsg)

	signature, err := sig.Sign(digestHash, keypair)
	if err != nil {
		fmt.Println("Sign failed", err)
		return
	}
	fmt.Println("Perf starting 4")
	pubBytes, err := sig.SerializePublicKey(&keypair.PublicKey)
	if err != nil {
		fmt.Println("SerializePublicKey failed", err)
		return
	}

	if len(os.Args) < 2 {
		fmt.Println("enter count")
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
