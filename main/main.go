package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/L7-MCPE/lav7"
	"github.com/L7-MCPE/lav7/command"
	"github.com/L7-MCPE/lav7/raknet"
	"github.com/L7-MCPE/lav7/util"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	raknet.Players = lav7.Players
	lav7.GetDefaultLevel().Init()
	util.Debug("Generating chunks")
	wg := new(sync.WaitGroup)
	wg.Add(25)
	start := time.Now()
	for x := int32(-2); x <= 2; x++ {
		for z := int32(-2); z <= 2; z++ {
			go func(x, z int32) { lav7.GetDefaultLevel().GetChunk(x, z, true); wg.Done() }(x, z)
		}
	}
	wg.Wait()
	util.Debug("Elapsed time:", time.Since(start).Seconds(), "seconds")
	var r *raknet.Router
	var err error
	if r, err = raknet.CreateRouter(lav7.AddPlayer, 19132); err != nil {
		fmt.Println(err)
		return
	}
	r.Start()
	command.HandleCommand()
}
