package main

import (
	"fmt"
	"runtime"

	"github.com/L7-MCPE/lav7"
	"github.com/L7-MCPE/lav7/command"
	"github.com/L7-MCPE/raknet"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	raknet.Players = lav7.Players
	lav7.GetDefaultLevel().Init()
	var r *raknet.Router
	var err error
	if r, err = raknet.CreateRouter(lav7.AddPlayer, 19132); err != nil {
		fmt.Println(err)
		return
	}
	r.Start()
	command.HandleCommand()
}
