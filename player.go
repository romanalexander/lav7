package lav7

import (
	"encoding/hex"
	"fmt"
	"net"

	"github.com/L7-MCPE/raknet"
	"github.com/L7-MCPE/util"
	"github.com/L7-MCPE/util/buffer"
)

// Player is a struct for handling/containing MCPE client specific things.
type Player struct {
	Address       *net.UDPAddr
	Username      string
	ClientID      uint64
	ClientUUIDRaw [16]byte
	ClientSecret  string
	Skin          []byte
	SkinName      string
	x, y, z       float64
	loggedIn      bool
	closed        bool
}

// HandlePacket handles received MCPE packet after raknet connection is established.
func (p *Player) HandlePacket(b *buffer.Buffer) (err error) {
	pid := b.ReadByte()
	if err != nil {
		return err
	}
	switch pid {
	case 0x93: // TextPacket
		if b.ReadByte() == 1 { // Type: chat
			b.ReadString()
		}
		p.SendMessage("<" + p.Username + "> " + b.ReadString())
	case 0x8f: // LoginPacket
		if p.loggedIn {
			return
		}
		if len(Players) >= raknet.MaxPlayers {
			p.disconnect("Server is full!")
		}
		p.Username = b.ReadString()

		proto := b.ReadInt()                  // Protocol version check
		buf := buffer.FromBytes([]byte{0x90}) // PlayStatusPacket
		if proto > raknet.MinecraftProtocol {
			buf.WriteInt(2) // Failed by server
			p.send(buf)
			p.disconnect("Outdated server")
			return nil
		} else if proto < raknet.MinecraftProtocol {
			buf := buffer.FromBytes([]byte{0x90}) // PlayStatusPacket
			buf.WriteInt(1)                       // Failed by client
			p.send(buf)
			p.disconnect("Outdated client")
			return nil
		}
		buf.WriteInt(0) // Success
		p.send(buf)

		b.Read(4) // Skip proto2
		p.ClientID = b.ReadLong()
		var uuid []byte
		uuid = b.Read(16)
		copy(p.ClientUUIDRaw[:], uuid)
		b.ReadString() // Skip sever address field
		p.ClientSecret = b.ReadString()
		p.SkinName = b.ReadString()

		buf = buffer.FromBytes([]byte("\x95\xff\xff\xff\xff\x00")) // StartGamePacket (seed -1, dimension 0)
		buf.WriteInt(1)                                            // Generator - 0: old, 1: infinite, 2: flat
		buf.WriteInt(1)                                            // 0: Survival, 1: Creative
		buf.WriteLong(0)                                           // Player eid is forced to be zero
		buf.WriteInt(0)                                            // Spawnpoint X
		buf.WriteInt(64)                                           // Spawnpoint Y
		buf.WriteInt(0)                                            // Spawnpoint Z
		buf.WriteFloat(0)                                          // X
		buf.WriteFloat(64)                                         // Y
		buf.WriteFloat(0)                                          // Z
		buf.WriteByte(0)                                           // Unknown
		p.send(buf)

		// TODO: Send SetTime/SpawnPosition/Health/Difficulty packets

		p.firstSpawn()
		fmt.Println(p.Username + " joined the game")
		p.SendMessage("Hell-O from the server!")
	case 0x92: // BatchPacket
		size := b.ReadInt()
		payload := b.Read(uint64(size))
		b, err := util.DecodeDeflate(payload)
		if err != nil {
			return err
		}
		buf := buffer.FromBytes(b)
		for buf.Require(4) {
			size := buf.ReadInt()
			b := buf.Read(uint64(size))
			if b[0] == 0x92 {
				return fmt.Errorf("Invalid BatchPacket inside BatchPacket")
			}
			if err := p.HandlePacket(buffer.FromBytes(b)); err != nil {
				return err
			}
		}
	default:
		util.Debug("Unimplemented!")
		fmt.Print(hex.Dump(b.Payload))
	}
	return nil
}

// SendMessage sends text to player.
func (p *Player) SendMessage(msg string) {
	buf := buffer.FromBytes([]byte{0x93}) // TextPacket
	buf.WriteByte(0)                      // Type: Raw
	buf.WriteString(msg)
	p.send(buf)
}

func (p *Player) firstSpawn() {
	buf := buffer.FromBytes([]byte{0x90}) // PlayStatusPacket
	buf.WriteInt(3)                       // Player spawn
	p.send(buf)
}

// Kick kicks player from server.
func (p *Player) Kick(reason string) {
	p.disconnect("Kicked: " + reason)
}

func (p *Player) disconnect(msg string) {
	buf := buffer.FromBytes([]byte{0x91})
	buf.WriteString(msg)
	p.send(buf)
	raknet.Sessions[p.Address.String()].Close("disconnected from server: " + msg)
}

func (p *Player) send(buf *buffer.Buffer) {
	if session, ok := raknet.Sessions[p.Address.String()]; ok {
		ep := new(raknet.EncapsulatedPacket)
		ep.Reliability = 2
		ep.Buffer = buf
		ep.Buffer.Offset = 0
		session.SendEncapsulated(ep)
	}
}
