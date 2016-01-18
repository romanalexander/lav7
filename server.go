package lav7

import (
	"fmt"
	"net"
	"sync/atomic"

	"github.com/L7-MCPE/lav7/level"
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
	p.Level = GetDefaultLevel()
	p.EntityID = atomic.AddUint64(&lastEntityID, 1)
	Players[identifier] = p
	return p.HandlePacket
}

// GetLevel returns level reference with given name if exists, or nil.
func GetLevel(name string) level.Level {
	if l, ok := levels[name]; ok {
		return l
	}
	return nil
}

// GetDefaultLevel returns default level reference.
func GetDefaultLevel() level.Level {
	return levels[defaultLvl]
}
