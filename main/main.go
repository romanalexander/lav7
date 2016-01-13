package main

import (
	"fmt"

	"github.com/L7-MCPE/raknet"
	"github.com/L7-MCPE/util"
)

func main() {
	var r *raknet.Router
	var err error
	if r, err = raknet.CreateRouter(19132); err != nil {
		fmt.Println(err)
		return
	}
	r.Start()
	util.Suspend()
}
