package lav7

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"

	"github.com/L7-MCPE/lav7/level"
	"github.com/L7-MCPE/lav7/raknet"
	"github.com/L7-MCPE/lav7/util"
	"github.com/L7-MCPE/lav7/util/buffer"
	"github.com/davecgh/go-spew/spew"
)

// Player is a struct for handling/containing MCPE client specific things.
type Player struct {
	Address       *net.UDPAddr
	Username      string
	ClientID      uint64
	ClientUUIDRaw [16]byte
	ClientSecret  string
	EntityID      uint64
	Skin          []byte
	SkinName      string

	Position            util.Vector3
	Level               *level.Level
	Yaw, BodyYaw, Pitch float32

	sentChunks map[[2]int32]bool
	loggedIn   bool
	closed     bool
}

// HandlePacket handles received MCPE packet after raknet connection is established.
func (p *Player) HandlePacket(b *buffer.Buffer) (err error) {
	pid := b.ReadByte()
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

		ret := &PlayStatus{}
		if pk.Proto1 > raknet.MinecraftProtocol {
			ret.Status = LoginFailedServer
			p.SendPacket(ret)
			p.disconnect("Outdated server")
			return
		} else if pk.Proto1 < raknet.MinecraftProtocol {
			ret.Status = LoginFailedClient
			p.SendPacket(ret)
			p.disconnect("Outdated client")
			return
		}
		ret.Status = LoginSuccess
		p.SendPacket(ret)

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
		xRadius := int32(2)
		zRadius := int32(2)
		chunkChan := make(chan struct {
			x, z int32
			c    level.Chunk
		}, (xRadius*2+1)*(zRadius*2+1))
		go func() {
			for x := -xRadius; x <= xRadius; x++ {
				for z := -zRadius; z <= zRadius; z++ {
					chunkChan <- struct {
						x, z int32
						c    level.Chunk
					}{x, z, p.Level.GetChunk(x, z, true)}
				}
			}
		}()
		go func() {
			for chunks := (xRadius*2 + 1) * (zRadius*2 + 1); chunks > 0; chunks-- {
				s := <-chunkChan
				p.SendChunk(s.x, s.z, s.c)
			}
			p.firstSpawn()
			log.Println(p.Username + " joined the game")
			p.SendMessage("Hello, this is lav7 test server!")
		}()
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
		if pk.Message[:1] == "/" {
			if pk.Message[1:] == "stop" {
				Stop("issued by player " + p.Username)
			}
		}
		log.Println(fmt.Sprintf("<%s> %s", p.Username, pk.Message))
		AsPlayers(func(pp *Player) error { pp.SendMessage(fmt.Sprintf("<%s> %s", p.Username, pk.Message)); return nil })
	case *MovePlayer:
		pk := pk.(*MovePlayer)
		log.Println("Player move:", pk.X, pk.Y, pk.Z, pk.Yaw, pk.BodyYaw, pk.Pitch)
	case *RemoveBlock:
		pk := pk.(*RemoveBlock)
		p.Level.SetBlock(int32(pk.X), int32(pk.Y), int32(pk.Z), 0) // Air
	case *UseItem:
		pk := pk.(*UseItem)
		p.Level.OnUseItem(int32(pk.X), int32(pk.Y), int32(pk.Z), pk.Face, pk.Item)
		spew.Dump(pk)
	default:
		log.Println("0x" + hex.EncodeToString([]byte{pk.Pid()}) + "is unimplemented:")
		spew.Dump(pk)
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
	if _, exists := p.sentChunks[[2]int32{chunkX, chunkZ}]; exists {
		return
	}
	i := &FullChunkData{
		ChunkX:  uint32(chunkX),
		ChunkZ:  uint32(chunkZ),
		Order:   1,
		Payload: c.FullChunkData(),
	}
	p.SendCompressed(i)
	p.sentChunks[[2]int32{chunkX, chunkZ}] = true
}

func (p *Player) updateMove(pk *MovePlayer) {
	p.Position.X, p.Position.Y, p.Position.Z = pk.X, pk.Y, pk.Z
	p.Yaw, p.BodyYaw, p.Pitch = pk.Yaw, pk.BodyYaw, pk.Pitch
}

func (p *Player) firstSpawn() {
	pk := &PlayStatus{
		Status: PlayerSpawn,
	}
	p.SendPacket(pk)
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
	buf := buffer.FromBytes([]byte{0x8e, pk.Pid()})
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

// Do not use this method for sending packet to client, this is an internal function.
func (p *Player) send(buf *buffer.Buffer) {
	if session, ok := raknet.Sessions[p.Address.String()]; ok {
		ep := new(raknet.EncapsulatedPacket)
		ep.Reliability = 2
		ep.Buffer = buf
		ep.Buffer.Offset = 0
		session.SendEncapsulated(ep)
	}
}
