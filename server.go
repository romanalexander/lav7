package lav7

import (
	"fmt"
	"net"
	"os"
	"sync/atomic"

	"github.com/L7-MCPE/lav7/level"
	"github.com/L7-MCPE/lav7/util/buffer"
)

// RegisterPlayer registers player to the server and returns packet handler function for it.
func RegisterPlayer(addr *net.UDPAddr) (handlerFunc func(*buffer.Buffer) error) {
	identifier := addr.String()
	if _, ok := Players[identifier]; ok {
		fmt.Println("Duplicate authentication from", addr)
		Players[identifier].(*Player).disconnect("Logged in from another location")
	}
	p := new(Player)
	p.Address = addr
	p.Level = GetDefaultLevel()
	p.EntityID = atomic.AddUint64(&lastEntityID, 1)
	p.sentChunks = make(map[[2]int32]bool)
	Players[identifier] = p
	return p.HandlePacket
}

// AsPlayers executes given callback with every online players.
func AsPlayers(callback func(*Player) error) error {
	for _, p := range Players {
		if err := callback(p.(*Player)); err != nil {
			return err
		}
	}
	return nil
}

// BroadcastPacket sends given packet to all online players.
func BroadcastPacket(pk Packet) {
	for _, p := range Players {
		p.(*Player).SendPacket(pk)
	}
}

// GetLevel returns level reference with given name if exists, or nil.
func GetLevel(name string) *level.Level {
	if l, ok := levels[name]; ok {
		return l
	}
	return nil
}

// GetDefaultLevel returns default level reference.
func GetDefaultLevel() *level.Level {
	return levels[defaultLvl]
}

// Stop stops entire server.
func Stop(reason string) {
	if reason == "" {
		reason = "no reason"
	}
	fmt.Println("Stopping server: " + reason)
	AsPlayers(func(p *Player) error {
		p.Kick("Server stop: " + reason)
		return nil
	})
	for _, l := range levels {
		l.Save()
	}
	os.Exit(0)
}
