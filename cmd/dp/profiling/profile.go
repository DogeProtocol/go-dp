package profiling

import (
	"fmt"
	"github.com/DogeProtocol/dp/log"
	"net/http"
	"os"
	"runtime/pprof"
	"strconv"
)

func StartProfiling(port int) {
	cpuProf := os.Getenv("CPU_PROF")
	if len(cpuProf) > 0 {
		fmt.Println("CPU_PROF enabled. Starting CPU profiling.")
		f, err := os.Create("cpu.prof")
		if err != nil {
			log.Error("profiling failed")
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	listenStack := "localhost:" + strconv.Itoa(port)
	log.Info("Profile Flag is set", "url", "http://"+listenStack)
	http.ListenAndServe(listenStack, nil)
}
