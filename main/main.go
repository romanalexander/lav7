package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/L7-MCPE/lav7"
	"github.com/L7-MCPE/lav7/command"
	"github.com/L7-MCPE/raknet"
	"github.com/L7-MCPE/util"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	raknet.Players = lav7.Players
	lav7.GetDefaultLevel().Init()
	util.Debug("Generating chunks")
	start := time.Now()
	for x := int32(-7); x <= 7; x++ {
		for z := int32(-7); z <= 7; z++ {
			go func() { lav7.GetDefaultLevel().GetChunk(x, z, true) }()
		}
	}
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
