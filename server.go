package lav7

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync/atomic"

	"github.com/L7-MCPE/lav7/level"
	"github.com/L7-MCPE/lav7/proto"
	"github.com/L7-MCPE/lav7/raknet"
	"github.com/L7-MCPE/lav7/util/buffer"
)

// RegisterPlayer registers player to the server and returns packet handler function for it.
func RegisterPlayer(addr *net.UDPAddr) (handlerChan chan<- *buffer.Buffer) {
	identifier := addr.String()
	if _, ok := Players[identifier]; ok {
		fmt.Println("Duplicate authentication from", addr)
		Players[identifier].disconnect("Logged in from another location")
	}
	p := new(Player)
	p.Address = addr
	p.Level = GetDefaultLevel()
	p.EntityID = atomic.AddUint64(&LastEntityID, 1)
	p.playerShown = make(map[uint64]struct{})
	p.sentChunks = make(map[[2]int32]bool)
	ch := make(chan *buffer.Buffer, 64)
	p.recvChan = ch
	p.raknetChan = raknet.Sessions[identifier].PlayerChan
	iteratorLock.Lock()
	Players[identifier] = p
	iteratorLock.Unlock()
	go p.process()
	return ch
}

// UnregisterPlayer removes player from server.
func UnregisterPlayer(addr *net.UDPAddr) error {
	identifier := addr.String()
	if p, ok := Players[identifier]; ok {
		close(p.recvChan)
		AsPlayers(func(pl *Player) {
			pl.HidePlayer(p) //FIXME: semms not working
		})
		delete(Players, identifier)
		Message(p.Username + " disconnected")
		return nil
	}
	return fmt.Errorf("Tried to remove nonexistent player: %v", addr)
}

// AsPlayers executes given callback with every online players.
func AsPlayers(callback func(*Player)) {
	iteratorLock.Lock()
	defer iteratorLock.Unlock()
	for _, p := range Players {
		callback(p)
	}
}

// AsPlayersAsync is similar to AsPlayers, buf spawns new goroutine for each players.
// Warning: this could be a lot of overhead. Use with caution.
func AsPlayersAsync(callback func(*Player)) {
	for _, p := range Players {
		go func(pp *Player) {
			iteratorLock.Lock()
			defer iteratorLock.Unlock()
			callback(pp)
		}(p)
	}
}

// AsPlayersError is similar to AsPlayers, but breaks iteration if callback returns error
func AsPlayersError(callback func(*Player) error) error {
	iteratorLock.Lock()
	defer iteratorLock.Unlock()
	for _, p := range Players {
		if err := callback(p); err != nil {
			return err
		}
	}
	return nil
}

// Message broadcasts message, and logs to console.
func Message(msg string) {
	AsPlayers(func(pl *Player) {
		pl.SendMessage(msg)
	})
	log.Println(msg)
}

// SpawnPlayer shows given player to all players, except given player itself.
func SpawnPlayer(player *Player) {
	AsPlayers(func(p *Player) {
		if p.spawned && p.EntityID != player.EntityID {
			p.ShowPlayer(player)
		}
	})
}

// BroadcastPacket sends given packet to all online players.
func BroadcastPacket(pk proto.Packet) {
	for _, p := range Players {
		p.SendPacket(pk)
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
	AsPlayers(func(p *Player) { p.Kick("Server stop: " + reason) })
	for _, l := range levels {
		l.Save()
	}
	os.Exit(0)
}
