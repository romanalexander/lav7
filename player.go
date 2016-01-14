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
	case 0x8f:
		if p.loggedIn {
			return
		}
		if len(Players) >= raknet.MaxPlayers {
			p.disconnect("Server is full!")
		}
		p.Username = b.ReadString()
		proto := b.ReadInt()
		if proto > raknet.MinecraftProtocol {
			buf := buffer.FromBytes([]byte{0x90}) // LoginStatusPacket
			buf.WriteInt(2)                       // Failed by server
			p.sendPacket(buf)
			p.disconnect("Outdated server")
			return nil
		} else if proto < raknet.MinecraftProtocol {
			buf := buffer.FromBytes([]byte{0x90}) // LoginStatusPacket
			buf.WriteInt(1)                       // Failed by client
			p.sendPacket(buf)
			p.disconnect("Outdated client")
			return nil
		}
		p.disconnect("Success")
		b.Read(4)
		p.ClientID = b.ReadLong()
		var uuid []byte
		uuid = b.Read(16)
		copy(p.ClientUUIDRaw[:], uuid)
		b.ReadString() // Skip sever address field
		p.ClientSecret = b.ReadString()
		p.SkinName = b.ReadString()

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

// Kick kicks player from server.
func (p *Player) Kick(reason string) {

}

func (p *Player) disconnect(msg string) {
	buf := buffer.FromBytes([]byte{0x91})
	buf.WriteString(msg)
	p.sendPacket(buf)
	raknet.Sessions[p.Address.String()].Close("disconnected from server: " + msg)
}

func (p *Player) sendPacket(buf *buffer.Buffer) {
	if session, ok := raknet.Sessions[p.Address.String()]; ok {
		ep := new(raknet.EncapsulatedPacket)
		ep.Reliability = 2
		ep.Buffer = buf
		ep.Buffer.Offset = 0
		session.SendEncapsulated(ep)
	} else {
		fmt.Println("Oops?", p.Address.String())
		fmt.Println(raknet.Sessions)
	}
}
