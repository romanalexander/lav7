package main

import (
	"fmt"

	"github.com/L7-MCPE/lav7"
	"github.com/L7-MCPE/lav7/level/format/dummy"
	"github.com/L7-MCPE/raknet"
	"github.com/L7-MCPE/util"
)

func main() {
	dummy.InitDummyChunk()
	raknet.Players = lav7.Players
	var r *raknet.Router
	var err error
	if r, err = raknet.CreateRouter(lav7.AddPlayer, 19132); err != nil {
		fmt.Println(err)
		return
	}
	r.Start()
	util.Suspend()
}
