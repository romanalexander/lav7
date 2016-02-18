package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/L7-MCPE/lav7"
	"github.com/L7-MCPE/lav7/command"
	"github.com/L7-MCPE/lav7/level/format/dummy"
	"github.com/L7-MCPE/lav7/level/gen"
	"github.com/L7-MCPE/lav7/raknet"
)

func main() {
	ppr := flag.Bool("pprof", false, "starts pprof debug server on :8080")
	flag.Parse()
	if *ppr {
		go http.ListenAndServe(":8080", nil)
		log.Println("debug: pprof server is running on :8080")
	}
	log.Printf("Starting lav7 version %s(git commit %s)", lav7.Version, lav7.GitCommit)
	log.Println("lav7 is licensed under the GPLv3 License; see http://rawgit.com/L7-MCPE/lav7/master/LICENSE.txt")
	log.Printf("Build timestamp: %s", lav7.BuildTime)
	start := time.Now()
	runtime.GOMAXPROCS(runtime.NumCPU())
	initLevel(5)
	initRaknet()
	startRouter()
	log.Println("All done! Elapsed time:", time.Since(start).Seconds(), "seconds")
	log.Println("Server is ready.")
	command.HandleCommand()
}

func initLevel(genRadius int32) {
	g := new(gen.SampleGenerator)
	log.Println("Generator type:", reflect.TypeOf(g))
	g.Init()
	log.Println("Generator init done. Initializing level...")
	p := new(dummy.Provider)
	p.Init(lav7.GetDefaultLevel().Name)
	lav7.GetDefaultLevel().Init(p)
	lav7.GetDefaultLevel().Gen = g.Gen
	// log.Printf("Level init done")

	log.Printf("Level init done. Preparing chunks(initial radius: %d)", genRadius)
	chunks := int((genRadius*2 + 1) * (genRadius*2 + 1))
	wg := new(sync.WaitGroup)
	wg.Add(chunks)
	for x := -genRadius; x <= genRadius; x++ {
		for z := -genRadius; z <= genRadius; z++ {
			go func(x, z int32) { lav7.GetDefaultLevel().GetChunk(x, z, true); wg.Done() }(x, z)
		}
	}
	wg.Wait()
	log.Printf("%d chunks cached in memory.", chunks)
}

func initRaknet() {
	raknet.ServerName = lav7.ServerName
	raknet.MaxPlayers = lav7.MaxPlayers
}

func startRouter() {
	log.Println("Starting raknet router, version", raknet.Version)
	var r *raknet.Router
	var err error
	if r, err = raknet.CreateRouter(lav7.RegisterPlayer, lav7.UnregisterPlayer, 19132); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	r.Start()
}
