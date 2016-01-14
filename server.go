package lav7

import (
	"fmt"
	"net"

	"github.com/L7-MCPE/util/buffer"
)

// AddPlayer registers player to the server and returns packet handler function for it.
func AddPlayer(addr *net.UDPAddr) (handlerFunc func(*buffer.Buffer) error) {
	identifier := addr.String()
	if _, ok := Players[identifier]; ok {
		fmt.Println("Duplicate authentication from", addr)
		Players[identifier].(*Player).disconnect("Logged in from another location")
	}
	p := new(Player)
	p.Address = addr
	Players[identifier] = p
	fmt.Println(Players)
	return p.HandlePacket
}
