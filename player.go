package lav7

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/L7-MCPE/lav7/level"
	. "github.com/L7-MCPE/lav7/proto"
	"github.com/L7-MCPE/lav7/raknet"
	"github.com/L7-MCPE/lav7/types"
	"github.com/L7-MCPE/lav7/util"
	"github.com/L7-MCPE/lav7/util/buffer"
)

type PlayerCallback struct {
	Call func(*Player, interface{})
	Arg  interface{}
}

// Player is a struct for handling/containing MCPE client specific things.
type Player struct {
	Address  *net.UDPAddr
	Username string
	ID       uint64
	UUID     [16]byte
	Secret   string
	EntityID uint64
	Skin     []byte
	SkinName string

	Position            util.Vector3
	Level               *level.Level
	Yaw, BodyYaw, Pitch float32

	playerShown map[uint64]struct{}

	sentChunks   map[[2]int32]struct{}
	chunkRequest chan [2]int32
	chunkBusy    bool
	chunkStop    chan struct{}

	recvChan     chan *buffer.Buffer
	raknetChan   chan<- *raknet.EncapsulatedPacket
	callbackChan chan PlayerCallback
	updateTicker *time.Ticker

	loggedIn bool
	spawned  bool
	closed   bool
}

func (p *Player) process() {
	radius := int32(3)
	for {
		select {
		case buf, ok := <-p.recvChan:
			if !ok {
				return
			}
			p.HandlePacket(buf)
		case callback := <-p.callbackChan:
			callback.Call(p, callback.Arg)
		case <-p.updateTicker.C:
			cx, cz := int32(p.Position.X)>>4, int32(p.Position.Z)>>4
			for sx := cx - radius; sx <= cx+radius; sx++ {
				for sz := cz - radius; sz <= cz+radius; sz++ {
					p.chunkRequest <- [2]int32{sx, sz}
				}
			}
		}
	}
}

// NOTE: Do NOT execute. This is an internal function.
func (p *Player) updateChunk() {
	p.sentChunks = make(map[[2]int32]struct{})
	for {
		select {
		case <-p.chunkStop:
			return
		case req := <-p.chunkRequest:
			if _, ok := p.sentChunks[[2]int32{req[0], req[1]}]; !ok {
				go p.sendChunk(req[0], req[1], p.Level.GetChunk(req[0], req[1], true))
				p.sentChunks[req] = struct{}{}
			}
		}
	}
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

		p.ID = pk.ClientID
		p.UUID = pk.RawUUID
		p.Secret = pk.ClientSecret
		p.SkinName = pk.SkinName
		p.Skin = pk.Skin

		p.SendPacket(&StartGame{
			Seed:      0xffffffff, // -1
			Dimension: 0,
			Generator: 1, // 0: old, 1: infinite, 2: flat
			Gamemode:  1, // 0: Survival, 1: Creative
			EntityID:  0, // Player eid set to 0
			SpawnX:    0,
			SpawnY:    65,
			SpawnZ:    0,
			X:         0,
			Y:         65,
			Z:         0,
		})
		p.loggedIn = true

		// TODO: Send SetTime/SpawnPosition/Health/Difficulty packets
		p.chunkBusy = true
		p.firstSpawn()
		log.Println(p.Username + " joined the game")
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
		Message(fmt.Sprintf("<%s> %s", p.Username, pk.Message))

	case *MovePlayer:
		pk := pk.(*MovePlayer)
		// log.Println("Player move:", pk.X, pk.Y, pk.Z, pk.Yaw, pk.BodyYaw, pk.Pitch)
		p.updateMove(pk)

	case *RemoveBlock:
		pk := pk.(*RemoveBlock)
		p.Level.SetBlock(int32(pk.X), int32(pk.Y), int32(pk.Z), 0) // Air
		p.BroadcastOthers(&UpdateBlock{
			BlockRecords: []BlockRecord{
				BlockRecord{
					X:     uint32(pk.X),
					Y:     byte(pk.Y),
					Z:     uint32(pk.Z),
					Block: types.Block{},
				},
			},
		})

	case *UseItem:
		pk := pk.(*UseItem)
		px, py, pz := int32(pk.X), int32(pk.Y), int32(pk.Z)
		if !p.Level.OnUseItem(&px, &py, &pz, pk.Face, pk.Item) {
			p.BroadcastOthers(&UpdateBlock{
				BlockRecords: []BlockRecord{
					BlockRecord{
						X: uint32(px),
						Y: byte(py),
						Z: uint32(pz),
						Block: types.Block{
							ID:   byte(pk.Item.ID),
							Meta: byte(pk.Item.Meta),
						},
					},
				},
			})
		} else {
			BroadcastPacket(&UpdateBlock{
				BlockRecords: []BlockRecord{
					BlockRecord{
						X:     uint32(pk.X),
						Y:     byte(pk.Y),
						Z:     uint32(pk.Z),
						Block: types.Block{},
					},
				},
			})
		}
		//spew.Dump(pk)

	case *Animate:
		pk := pk.(*Animate)
		pk.EntityID = p.EntityID
		p.BroadcastOthers(pk)

	default:
		// log.Println("0x" + hex.EncodeToString([]byte{pk.Pid()}) + "is unimplemented:")
		// spew.Dump(pk)
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

//NOTE: This function is NOT goroutine-safe. Only for internal use.
func (p *Player) sendChunk(chunkX, chunkZ int32, c *types.Chunk) {
	c.Mutex().RLock()
	i := &FullChunkData{
		ChunkX:  uint32(chunkX),
		ChunkZ:  uint32(chunkZ),
		Order:   OrderLayered,
		Payload: c.FullChunkData(),
	}
	c.Mutex().RUnlock()
	p.SendCompressed(i)
}

// ShowPlayer shows given player struct to player.
func (p *Player) ShowPlayer(player *Player) {
	if p.IsVisible(player) || p.IsSelf(player) {
		return
	}
	p.SendPacket(&AddPlayer{
		RawUUID:  player.UUID,
		Username: player.Username,
		EntityID: player.EntityID,
		X:        player.Position.X,
		Y:        player.Position.Y,
		Z:        player.Position.Z,
		SpeedX:   0,
		SpeedY:   0,
		SpeedZ:   0,
		BodyYaw:  player.BodyYaw,
		Yaw:      player.Yaw,
		Pitch:    player.Pitch,
	})
	p.playerShown[player.EntityID] = struct{}{}
}

// HidePlayer hides given player struct from player.
func (p *Player) HidePlayer(player *Player) {
	if !p.IsVisible(player) || p.IsSelf(player) {
		return
	}
	p.SendPacket(&RemovePlayer{
		EntityID: player.EntityID,
		RawUUID:  player.UUID,
	})
	delete(p.playerShown, player.EntityID)
}

// IsVisible determines if the player can see given player struct
func (p *Player) IsVisible(player *Player) bool {
	_, ok := p.playerShown[player.EntityID]
	return ok
}

// IsSelf determines if given player is the player self
func (p *Player) IsSelf(player *Player) bool {
	return p.EntityID == player.EntityID
}

func (p *Player) updateMove(pk *MovePlayer) {
	p.Position.X, p.Position.Y, p.Position.Z = pk.X, pk.Y, pk.Z
	p.Yaw, p.BodyYaw, p.Pitch = pk.Yaw, pk.BodyYaw, pk.Pitch

	go BroadcastCallback(PlayerCallback{
		Call: func(pl *Player, arg interface{}) {
			if pl.IsVisible(p) {
				pl.SendPacket(&MoveEntity{
					EntityIDs: []uint64{p.EntityID},
					EntityPos: [][6]float32{[6]float32{
						pk.X,
						pk.Y,
						pk.Z,
						pk.BodyYaw,
						pk.Yaw,
						pk.Pitch,
					}},
				})
			}
		},
	})
}

func (p *Player) firstSpawn() {
	if p.spawned {
		return
	}

	BroadcastCallback(PlayerCallback{
		Call: func(player *Player, arg interface{}) {
			player.ShowPlayer(p)
			player.SendPacket(&PlayerList{
				Type: PlayerListAdd,
				PlayerEntries: []PlayerListEntry{PlayerListEntry{
					RawUUID:  p.UUID,
					EntityID: p.EntityID,
					Username: p.Username,
					SkinName: p.SkinName,
					Skin:     p.Skin,
				}},
			})
		},
	})

	entries := make([]PlayerListEntry, 0)
	AsPlayers(func(pl *Player) {
		p.ShowPlayer(pl)
		entries = append(entries, PlayerListEntry{
			RawUUID:  pl.UUID,
			EntityID: pl.EntityID,
			Username: pl.Username,
			SkinName: pl.SkinName,
			Skin:     pl.Skin,
		})
	})

	p.SendPacket(&PlayerList{
		Type:          PlayerListAdd,
		PlayerEntries: entries,
	})
	p.SendPacket(&PlayStatus{
		Status: PlayerSpawn,
	})

	SpawnPlayer(p)
	p.spawned = true
	Message(p.Username + " joined")
}

// Kick kicks player from server.
func (p *Player) Kick(reason string) {
	p.disconnect("Kicked: " + reason)
}

func (p *Player) disconnect(msg string) {
	p.SendDirect(&Disconnect{
		Message: msg,
	})

	raknet.SessionLock.Lock()
	s, ok := raknet.Sessions[p.Address.String()]
	raknet.SessionLock.Unlock()
	if ok {
		s.Close(msg)
	}
}

// BroadcastOthers broadcasts packet except player self.
func (p *Player) BroadcastOthers(pk Packet) {
	AsPlayers(func(pl *Player) {
		if !pl.IsSelf(p) {
			pl.SendPacket(pk)
		}
	})
}

// SendPacket sends given packet to client.
func (p *Player) SendPacket(pk Packet) {
	buf := buffer.FromBytes([]byte{0x8e, pk.Pid()})
	buf.Write(pk.Write().Done())
	p.send(buf)
}

// SendDirect sends given packet without passing to raknetChan channel.
func (p *Player) SendDirect(pk Packet) {
	buf := buffer.FromBytes([]byte{0x8e, pk.Pid()})
	buf.Write(pk.Write().Done())

	ep := new(raknet.EncapsulatedPacket)
	ep.Reliability = 2
	ep.Buffer = buf
	ep.Buffer.Offset = 0

	raknet.SessionLock.Lock()
	s, ok := raknet.Sessions[p.Address.String()]
	raknet.SessionLock.Unlock()

	if ok {
		s.SendEncapsulated(ep)
	}
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

// RunAs runs given callback on the player's goroutine.
// You should use this if you need access to other players' fields.
func (p *Player) RunAs(callback PlayerCallback) {
	p.callbackChan <- callback
}

// Do not use this method for sending packet to client, this is an internal function.
func (p *Player) send(buf *buffer.Buffer) {
	ep := new(raknet.EncapsulatedPacket)
	ep.Reliability = 2
	ep.Buffer = buf
	ep.Buffer.Offset = 0
	p.raknetChan <- ep
}
