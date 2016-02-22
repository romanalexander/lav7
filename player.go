package lav7

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/L7-MCPE/lav7/proto"
	"github.com/L7-MCPE/lav7/raknet"
	"github.com/L7-MCPE/lav7/types"
	"github.com/L7-MCPE/lav7/util"
	"github.com/L7-MCPE/lav7/util/buffer"
)

// ChunkRadius is a chunk radius that the player can handle.
const ChunkRadius int32 = 5

// PlayerCallback is a struct for delivering callbacks to other player goroutines;
// It is usually used to bypass race issues.
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
	Level               *Level
	Yaw, BodyYaw, Pitch float32

	playerShown map[uint64]struct{}

	fastChunks     map[[2]int32]*types.Chunk
	fastChunkMutex util.Locker
	chunkRequest   chan [2]int32
	chunkStop      chan struct{}
	chunkNotify    chan types.ChunkDelivery
	pending        map[[2]int32]time.Time

	inventory *PlayerInventory

	recvChan     chan *bytes.Buffer
	raknetChan   chan<- *raknet.EncapsulatedPacket
	callbackChan chan PlayerCallback
	updateTicker *time.Ticker

	loggedIn bool
	spawned  bool
	closed   bool
}

func (p *Player) process() {
	p.pending = make(map[[2]int32]time.Time)
	p.chunkRequest = make(chan [2]int32, (ChunkRadius*2+1)*(ChunkRadius*2+1))
	// resendTicker := time.NewTicker(time.Second * 3)
	// defer resendTicker.Stop()
	go p.updateChunk()
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
			chunkHold := make(map[[2]int32]struct{})
			for ccx := cx - ChunkRadius; ccx <= cx+ChunkRadius; ccx++ {
				for ccz := cz - ChunkRadius; ccz <= cz+ChunkRadius; ccz++ {
					chunkHold[[2]int32{ccx, ccz}] = struct{}{}
				}
			}
			p.fastChunkMutex.Lock()
			for cc := range p.fastChunks {
				if _, ok := chunkHold[cc]; ok {
					delete(chunkHold, cc)
				} else {
					delete(p.fastChunks, cc)
				}
			}
			p.fastChunkMutex.Unlock()
			for cc := range chunkHold {
				if timeout, ok := p.pending[cc]; !ok || timeout.Before(time.Now()) {
					p.requestChunk(cc)
				}
			}
		case c := <-p.chunkNotify:
			p.fastChunkMutex.Lock()
			if _, ok := p.fastChunks[[2]int32{c.X, c.Z}]; ok {
				p.fastChunkMutex.Unlock()
				break
			}
			p.fastChunks[[2]int32{c.X, c.Z}] = c.Chunk
			p.fastChunkMutex.Unlock()
			delete(p.pending, [2]int32{c.X, c.Z})
			p.sendChunk(c)
			/*
				case <-resendTicker.C:
					for cx := int32(p.Position.X) - ChunkRadius; cx <= int32(p.Position.X)+ChunkRadius; cx++ {
						for cz := int32(p.Position.Z) - ChunkRadius; cz <= int32(p.Position.Z)+ChunkRadius; cz++ {
							p.requestChunk([2]int32{cx, cz})
						}
					}
			*/
		}
	}
}

// SendNearChunk sends chunks near the player in radius.
// This function should be run on player process goroutine, or RunAs().
func (p *Player) SendNearChunk() {
	cx, cz := int32(p.Position.X)>>4, int32(p.Position.Z)>>4
	for ccx := cx - ChunkRadius; ccx <= cx+ChunkRadius; ccx++ {
		for ccz := cz - ChunkRadius; ccz <= cz+ChunkRadius; ccz++ {
			p.requestChunk([2]int32{ccx, ccz})
		}
	}
}

// NOTE: Do NOT execute outside player process goroutine.
func (p *Player) requestChunk(cc [2]int32) {
	go func(cc [2]int32) {
		p.chunkRequest <- cc
	}(cc)
	p.pending[cc] = time.Now().Add(time.Second * 5)
}

// NOTE: Do NOT execute. This is an internal function.
func (p *Player) updateChunk() {
	for {
		select {
		case <-p.chunkStop:
			return
		case req := <-p.chunkRequest:
			if c := p.getFastChunk(req[0], req[1]); c != nil {
				p.chunkNotify <- types.ChunkDelivery{
					X:     req[0],
					Z:     req[1],
					Chunk: c,
				}
				continue
			}
			if ch := p.Level.CreateChunk(req[0], req[1]); ch != nil {
				go func(cx, cz int32, done <-chan struct{}) {
					<-done
					p.chunkNotify <- types.ChunkDelivery{
						X:     cx,
						Z:     cz,
						Chunk: p.Level.GetChunk(cx, cz),
					}
				}(req[0], req[1], ch)
			}
		}
	}
}

// NOTE: Do NOT execute outside updateChunk goroutine. It could make data races.
func (p *Player) getFastChunk(cx, cz int32) *types.Chunk {
	p.fastChunkMutex.Lock()
	defer p.fastChunkMutex.Unlock()
	if c, ok := p.fastChunks[[2]int32{cx, cz}]; ok {
		return c
	}
	return p.Level.GetChunk(cx, cz)
}

// HandlePacket handles received MCPE packet after raknet connection is established.
func (p *Player) HandlePacket(b *bytes.Buffer) (err error) {
	pid := buffer.ReadByte(b)
	var pk proto.Packet
	if pk = proto.GetPacket(pid); pk == nil {
		return
	}
	pk.Read(b)
	return p.handleDataPacket(pk)
}

func (p *Player) handleDataPacket(pk proto.Packet) (err error) {
	switch pk.(type) {
	case *proto.Login:
		pk := pk.(*proto.Login)
		if p.loggedIn {
			return
		}
		iteratorLock.Lock()
		if len(Players) >= MaxPlayers {
			iteratorLock.Unlock()
			p.disconnect("Server is full!")
			return
		}
		iteratorLock.Unlock()
		p.Username = pk.Username

		ret := &proto.PlayStatus{}
		if pk.Proto1 > raknet.MinecraftProtocol {
			ret.Status = proto.LoginFailedServer
			p.SendPacket(ret)
			p.disconnect("Outdated server")
			return
		} else if pk.Proto1 < raknet.MinecraftProtocol {
			ret.Status = proto.LoginFailedClient
			p.SendPacket(ret)
			p.disconnect("Outdated client")
			return
		}
		ret.Status = proto.LoginSuccess
		p.SendPacket(ret)

		p.ID = pk.ClientID
		p.UUID = pk.RawUUID
		p.Secret = pk.ClientSecret
		p.SkinName = pk.SkinName
		p.Skin = pk.Skin

		p.SendPacket(&proto.StartGame{
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
		p.Position = util.Vector3{X: 0, Y: 65, Z: 0}
		p.loggedIn = true

		p.inventory.Holder = p
		p.inventory.Init()
		// TODO: Send SetTime/SpawnPosition/Health/Difficulty packets
		p.firstSpawn()
		go func() {
			<-time.After(time.Second * 1)

			p.SendPacket(&proto.PlayStatus{
				Status: proto.PlayerSpawn,
			})

			SpawnPlayer(p)
			p.spawned = true

			Message(p.Username + " joined")
			log.Println(p.Username + " joined the game")
			p.SendMessage("Hello, this is lav7 test server!")
		}()

	case *proto.Batch:
		pk := pk.(*proto.Batch)
		for _, pp := range pk.Payloads {
			if err = p.HandlePacket(bytes.NewBuffer(pp)); err != nil {
				return
			}
		}

	case *proto.Text:
		pk := pk.(*proto.Text)
		if pk.TextType == proto.TextTypeTranslation {
			return
		}
		Message(fmt.Sprintf("<%s> %s", p.Username, pk.Message))

	case *proto.MovePlayer:
		pk := pk.(*proto.MovePlayer)
		// log.Println("Player move:", pk.X, pk.Y, pk.Z, pk.Yaw, pk.BodyYaw, pk.Pitch)
		p.updateMove(pk)

	case *proto.RemoveBlock:
		pk := pk.(*proto.RemoveBlock)
		p.Level.SetBlock(int32(pk.X), int32(pk.Y), int32(pk.Z), 0) // Air
		p.BroadcastOthers(&proto.UpdateBlock{
			BlockRecords: []proto.BlockRecord{
				{
					X:     uint32(pk.X),
					Y:     byte(pk.Y),
					Z:     uint32(pk.Z),
					Block: types.Block{},
				},
			},
		})

	case *proto.UseItem:
		pk := pk.(*proto.UseItem)
		px, py, pz := int32(pk.X), int32(pk.Y), int32(pk.Z)
		if !p.Level.OnUseItem(&px, &py, &pz, pk.Face, pk.Item) {
			p.BroadcastOthers(&proto.UpdateBlock{
				BlockRecords: []proto.BlockRecord{
					{
						X: uint32(px),
						Y: byte(py),
						Z: uint32(pz),
						Block: types.Block{
							ID:   byte(pk.Item.ID),
							Meta: byte(pk.Item.Meta),
						},
						Flags: proto.UpdateAllPriority,
					},
				},
			})
		} else {
			p.SendPacket(&proto.UpdateBlock{
				BlockRecords: []proto.BlockRecord{
					{
						X:     uint32(pk.X),
						Y:     byte(pk.Y),
						Z:     uint32(pk.Z),
						Block: types.Block{},
						Flags: proto.UpdateAllPriority,
					},
				},
			})
		}
		//spew.Dump(pk)
	case *proto.ContainerSetSlot:
		//pk := pk.(*proto.ContainerSetSlot)
		//spew.Dump(pk)

	case *proto.Animate:
		pk := pk.(*proto.Animate)
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
	p.SendPacket(&proto.Text{
		TextType: proto.TextTypeRaw,
		Message:  msg,
	})
}

//NOTE: This function is NOT goroutine-safe. Only for internal use.
func (p *Player) sendChunk(c types.ChunkDelivery) {
	c.Chunk.Mutex().RLock()
	i := &proto.FullChunkData{
		ChunkX:  uint32(c.X),
		ChunkZ:  uint32(c.Z),
		Order:   proto.OrderLayered,
		Payload: c.Chunk.FullChunkData(),
	}
	c.Chunk.Mutex().RUnlock()
	p.SendCompressed(i)
}

// ShowPlayer shows given player struct to player.
func (p *Player) ShowPlayer(player *Player) {
	if p.IsVisible(player) || p.IsSelf(player) {
		return
	}
	p.SendPacket(&proto.AddPlayer{
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
	p.SendPacket(&proto.RemovePlayer{
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

func (p *Player) updateMove(pk *proto.MovePlayer) {
	p.Position.X, p.Position.Y, p.Position.Z = pk.X, pk.Y, pk.Z
	p.Yaw, p.BodyYaw, p.Pitch = pk.Yaw, pk.BodyYaw, pk.Pitch

	go BroadcastCallback(PlayerCallback{
		Call: func(pl *Player, arg interface{}) {
			if pl.IsVisible(p) {
				pl.SendPacket(&proto.MoveEntity{
					EntityIDs: []uint64{p.EntityID},
					EntityPos: [][6]float32{{
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
			player.SendPacket(&proto.PlayerList{
				Type: proto.PlayerListAdd,
				PlayerEntries: []proto.PlayerListEntry{{
					RawUUID:  p.UUID,
					EntityID: p.EntityID,
					Username: p.Username,
					Skin:     p.Skin,
				}},
			})
		},
	})

	var entries []proto.PlayerListEntry
	AsPlayers(func(pl *Player) {
		p.ShowPlayer(pl)
		entries = append(entries, proto.PlayerListEntry{
			RawUUID:  pl.UUID,
			EntityID: pl.EntityID,
			Username: pl.Username,
			Skin:     pl.Skin,
		})
	})

	for cx := int32(p.Position.X) - ChunkRadius; cx <= int32(p.Position.X)+ChunkRadius; cx++ {
		for cz := int32(p.Position.Z) - ChunkRadius; cz <= int32(p.Position.Z)+ChunkRadius; cz++ {
			p.requestChunk([2]int32{cx, cz})
		}
	}

	p.SendPacket(&proto.PlayerList{
		Type:          proto.PlayerListAdd,
		PlayerEntries: entries,
	})
}

// Kick kicks player from server.
func (p *Player) Kick(reason string) {
	p.disconnect("Kicked: " + reason)
}

func (p *Player) disconnect(msg string) {
	p.SendDirect(&proto.Disconnect{
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
func (p *Player) BroadcastOthers(pk proto.Packet) {
	AsPlayers(func(pl *Player) {
		if !pl.IsSelf(p) {
			pl.SendPacket(pk)
		}
	})
}

// SendPacket sends given packet to client.
func (p *Player) SendPacket(pk proto.Packet) {
	buf := bytes.NewBuffer([]byte{0x8e, pk.Pid()})
	buffer.Write(buf, pk.Write().Bytes())
	p.Send(buf)
}

// SendDirect sends given packet without passing to raknetChan channel.
func (p *Player) SendDirect(pk proto.Packet) {
	buf := bytes.NewBuffer([]byte{0x8e, pk.Pid()})
	buffer.Write(buf, pk.Write().Bytes())

	ep := new(raknet.EncapsulatedPacket)
	ep.Reliability = 2
	ep.Buffer = buf

	raknet.SessionLock.Lock()
	s, ok := raknet.Sessions[p.Address.String()]
	raknet.SessionLock.Unlock()

	if ok {
		s.SendEncapsulated(ep)
	}
}

// SendCompressed sends packed BatchPacket with given packets.
func (p *Player) SendCompressed(pks ...proto.Packet) {
	batch := &proto.Batch{
		Payloads: make([][]byte, len(pks)),
	}
	for i, pk := range pks {
		batch.Payloads[i] = append([]byte{pk.Pid()}, pk.Write().Bytes()...)
	}
	p.SendPacket(batch)
}

// RunAs runs given callback on the player's goroutine.
// You should use this if you need access to other players' fields.
func (p *Player) RunAs(callback PlayerCallback) {
	p.callbackChan <- callback
}

// Send sends bytes buffer to client.
// Do not use this method for sending packet to client, this is an internal function.
func (p *Player) Send(buf *bytes.Buffer) {
	ep := new(raknet.EncapsulatedPacket)
	ep.Reliability = 2
	ep.Buffer = buf
	p.raknetChan <- ep
}
