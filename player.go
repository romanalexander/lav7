package lav7

import (
	"encoding/hex"
	"fmt"

	"github.com/L7-MCPE/util"
	"github.com/L7-MCPE/util/buffer"
)

// Player is a struct for handling/containing MCPE client specific things.
type Player struct {
	x, y, z  float64
	SendChan chan *buffer.Buffer
}

// HandlePacket handles received MCPE packet after raknet connection is established.
func (p *Player) HandlePacket(b *buffer.Buffer) error {
	pid, err := b.ReadByte()
	if err != nil {
		return err
	}
	switch pid {
	case 0x92: // BatchPacket
		size, err := b.ReadInt()
		if err != nil {
			return err
		}
		payload, err := b.Read(uint64(size))
		if err != nil {
			return err
		}
		b, err := util.DecodeDeflate(payload)
		if err != nil {
			return err
		}
		buf := buffer.FromBytes(b)
		for buf.Require(4) {
			size, err := buf.ReadInt()
			if err != nil {
				return err
			}
			b, err := buf.Read(uint64(size))
			if err != nil {
				return err
			}
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

func (p *Player) sendPacket(buf *buffer.Buffer) {
	if p.SendChan == nil {
		return
	}
	p.SendChan <- buf
}
