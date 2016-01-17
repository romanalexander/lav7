package lav7

import (
	"encoding/hex"
	"fmt"
	"net"

	"github.com/L7-MCPE/lav7/level"
	"github.com/L7-MCPE/lav7/level/format/dummy"
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
	x, y, z       float32
	loggedIn      bool
	closed        bool
}

// HandlePacket handles received MCPE packet after raknet connection is established.
func (p *Player) HandlePacket(b *buffer.Buffer) (err error) {
	pid := b.Payload[0]
	b.Offset = 1
	var pk Packet
	if pk = GetPacket(pid); pk == nil {
		return
	}
	pk.Read(b)
	return p.handleDataPacket(pk)
}

func (p *Player) handleDataPacket(pk Packet) (err error) {
	switch pk.(type) {
	case *Login:
		pk := pk.(*Login)
		if p.loggedIn {
			return
		}
		if len(Players) >= raknet.MaxPlayers {
			p.disconnect("Server is full!")
		}
		p.Username = pk.Username

		buf := buffer.FromBytes([]byte{0x90}) // PlayStatusPacket
		if pk.Proto1 > raknet.MinecraftProtocol {
			buf.WriteInt(2) // Failed by server
			p.send(buf)
			p.disconnect("Outdated server")
			return
		} else if pk.Proto1 < raknet.MinecraftProtocol {
			buf := buffer.FromBytes([]byte{0x90}) // PlayStatusPacket
			buf.WriteInt(1)                       // Failed by client
			p.send(buf)
			p.disconnect("Outdated client")
			return
		}
		buf.WriteInt(0) // Success
		p.send(buf)

		p.ClientID = pk.ClientID
		p.ClientUUIDRaw = pk.RawUUID
		p.ClientSecret = pk.ClientSecret
		p.SkinName = pk.SkinName
		p.Skin = pk.Skin

		p.SendPacket(&StartGame{
			Seed:      0xffffffff, // -1
			Dimension: 0,
			Generator: 1, // 0: old, 1: infinite, 2: flat
			Gamemode:  1, // 0: Survival, 1: Creative
			EntityID:  0, // Player eid set to 0
			SpawnX:    0,
			SpawnY:    120,
			SpawnZ:    0,
			X:         0,
			Y:         120,
			Z:         0,
		})

		// TODO: Send SetTime/SpawnPosition/Health/Difficulty packets
		for x := int32(-5); x <= 5; x++ {
			for z := int32(-5); z <= 5; z++ {
				p.SendChunk(x, z, dummy.DummyChunk)
			}
		}
		p.firstSpawn()
		fmt.Println(p.Username + " joined the game")
		p.SendMessage("Hello, this is lav7 test server!")
	case *Batch:
		pk := pk.(*Batch)
		for _, pp := range pk.Payloads {
			if err = p.HandlePacket(buffer.FromBytes(pp)); err != nil {
				return
			}
		}
	case *Text:
		pk := pk.(*Text)
		if pk.TextType == TextTypeTranslation {
			return
		}
		p.SendMessage("Echo from: " + pk.Message)
	default:
		util.Debug("0x" + hex.EncodeToString([]byte{pk.Pid()}) + "is unimplemented: " + fmt.Sprint(pk))
	}
	return
}

// SendMessage sends text to player.
func (p *Player) SendMessage(msg string) {
	p.SendPacket(&Text{
		TextType: TextTypeRaw,
		Message:  msg,
	})
}

// SendChunk sends given Chunk struct to client.
func (p *Player) SendChunk(chunkX, chunkZ int32, c level.Chunk) {
	i := &FullChunkData{
		ChunkX:  uint32(chunkX),
		ChunkZ:  uint32(chunkZ),
		Order:   1,
		Payload: c.FullChunkData(),
	}
	p.SendCompressed(i)
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

// SendPacket sends given packet to client.
func (p *Player) SendPacket(pk Packet) {
	buf := buffer.FromBytes([]byte{pk.Pid()})
	buf.Write(pk.Write().Done())
	p.send(buf)
}

// SendCompressed sends packed BatchPacket with given packets.
func (p *Player) SendCompressed(pks ...Packet) {
	batch := &Batch{
		Payloads: make([][]byte, len(pks)),
	}
	for i, pk := range pks {
		batch.Payloads[i] = append([]byte{pk.Pid()}, pk.Write().Done()...)
	}
	p.SendPacket(batch)
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
