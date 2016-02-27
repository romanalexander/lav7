// Package proto provides MCPE network protocol, defined by Mojang.
package proto

import (
	"bytes"
	"log"

	"github.com/L7-MCPE/lav7/raknet"
	"github.com/L7-MCPE/lav7/types"
	"github.com/L7-MCPE/lav7/util"
	"github.com/L7-MCPE/lav7/util/buffer"
)

// Packet IDs
const (
	LoginHead byte = 0x8f + iota
	PlayStatusHead
	DisconnectHead
	BatchHead
	TextHead
	SetTimeHead
	StartGameHead
	AddPlayerHead
	RemovePlayerHead
	AddEntityHead
	RemoveEntityHead
	AddItemEntityHead
	TakeItemEntityHead
	MoveEntityHead
	MovePlayerHead
	RemoveBlockHead
	UpdateBlockHead
	AddPaintingHead
	ExplodeHead
	LevelEventHead
	BlockEventHead
	EntityEventHead
	MobEffectHead
	UpdateAttributesHead
	MobEquipmentHead
	MobArmorEquipmentHead
	InteractHead
	UseItemHead
	PlayerActionHead
	HurtArmorHead
	SetEntityDataHead
	SetEntityMotionHead
	SetEntityLinkHead
	SetHealthHead
	SetSpawnPositionHead
	AnimateHead
	RespawnHead
	DropItemHead
	ContainerOpenHead
	ContainerCloseHead
	ContainerSetSlotHead
	ContainerSetDataHead
	ContainerSetContentHead
	CraftingDataHead
	CraftingEventHead
	AdventureSettingsHead
	BlockEntityDataHead
	_ // 0xbe is skipped: PlayerInput
	FullChunkDataHead
	SetDifficultyHead
	_ // 0xc1 is skipped: ChangeDimension
	SetPlayerGametypeHead
	PlayerListHead
	_ // TelemetryEvent
	_ // SpawnExperienceOrb
	_ // ClientboundMapItemData
	_ // MapInfoRequest
	RequestChunkRadiusHead
	ChunkRadiusUpdateHead
	_ // ItemFrameDrop
	_ // ReplaceSelectedItem
)

var packets = map[byte]Packet{
	LoginHead:               new(Login),
	PlayStatusHead:          new(PlayStatus),
	DisconnectHead:          new(Disconnect),
	BatchHead:               new(Batch),
	TextHead:                new(Text),
	SetTimeHead:             new(SetTime),
	StartGameHead:           new(StartGame),
	AddPlayerHead:           new(AddPlayer),
	RemovePlayerHead:        new(RemovePlayer),
	AddEntityHead:           new(AddEntity),
	RemoveEntityHead:        new(RemoveEntity),
	AddItemEntityHead:       new(AddItemEntity),
	TakeItemEntityHead:      new(TakeItemEntity),
	MoveEntityHead:          new(MoveEntity),
	MovePlayerHead:          new(MovePlayer),
	RemoveBlockHead:         new(RemoveBlock),
	UpdateBlockHead:         new(UpdateBlock),
	AddPaintingHead:         new(AddPainting),
	ExplodeHead:             new(Explode),
	LevelEventHead:          new(LevelEvent),
	BlockEventHead:          new(BlockEvent),
	EntityEventHead:         new(EntityEvent),
	MobEffectHead:           new(MobEffect),
	UpdateAttributesHead:    new(UpdateAttributes),
	MobEquipmentHead:        new(MobEquipment),
	MobArmorEquipmentHead:   new(MobArmorEquipment),
	InteractHead:            new(Interact),
	UseItemHead:             new(UseItem),
	PlayerActionHead:        new(PlayerAction),
	HurtArmorHead:           new(HurtArmor),
	SetEntityDataHead:       new(SetEntityData),
	SetEntityMotionHead:     new(SetEntityMotion),
	SetEntityLinkHead:       new(SetEntityLink),
	SetHealthHead:           new(SetHealth),
	SetSpawnPositionHead:    new(SetSpawnPosition),
	AnimateHead:             new(Animate),
	RespawnHead:             new(Respawn),
	DropItemHead:            new(DropItem),
	ContainerOpenHead:       new(ContainerOpen),
	ContainerCloseHead:      new(ContainerClose),
	ContainerSetSlotHead:    new(ContainerSetSlot),
	ContainerSetDataHead:    new(ContainerSetData),
	ContainerSetContentHead: new(ContainerSetContent),
	CraftingDataHead:        new(CraftingData),
	CraftingEventHead:       new(CraftingEvent),
	AdventureSettingsHead:   new(AdventureSettings),
	BlockEntityDataHead:     new(BlockEntityData),
	FullChunkDataHead:       new(FullChunkData),
	SetDifficultyHead:       new(SetDifficulty),
	SetPlayerGametypeHead:   new(SetPlayerGametype),
	PlayerListHead:          new(PlayerList),
	RequestChunkRadiusHead:  new(RequestChunkRadius),
	ChunkRadiusUpdateHead:   new(ChunkRadiusUpdate),
}

// Packet is an interface for decoding/encoding MCPE packets.
type Packet interface {
	Pid() byte
	Read(*bytes.Buffer)
	Write() *bytes.Buffer
}

// GetPacket returns Packet struct with given pid.
func GetPacket(pid byte) Packet {
	pk, _ := packets[pid]
	return pk
}

// Login needs to be documented.
type Login struct {
	Username       string
	Proto1, Proto2 uint32
	ClientID       uint64
	RawUUID        [16]byte
	ServerAddress  string
	ClientSecret   string
	SkinName       string
	Skin           []byte
}

// Pid implements proto.Packet interface.
func (i Login) Pid() byte { return LoginHead } // 0x8f

// Read implements proto.Packet interface.
func (i *Login) Read(buf *bytes.Buffer) {
	buffer.BatchRead(buf, &i.Username, &i.Proto1)
	if i.Proto1 < raknet.MinecraftProtocol { // Old protocol
		return
	}
	buffer.BatchRead(buf, &i.Proto2, &i.ClientID)
	copy(i.RawUUID[:], buf.Next(16))
	buffer.BatchRead(buf, &i.ServerAddress, &i.ClientSecret, &i.SkinName)
	i.Skin = []byte(buffer.ReadString(buf))
}

// Write implements proto.Packet interface.
func (i Login) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.BatchWrite(buf, i.Username, i.Proto1, i.Proto2,
		i.ClientID, i.RawUUID[:], i.ServerAddress,
		i.ClientSecret, i.SkinName, string(i.Skin))
	return buf
}

// Packet-specific constants
const (
	LoginSuccess uint32 = iota
	LoginFailedClient
	LoginFailedServer
	PlayerSpawn
)

// PlayStatus needs to be documented.
type PlayStatus struct {
	Status uint32
}

// Pid implements proto.Packet interface.
func (i *PlayStatus) Pid() byte { return PlayStatusHead }

// Read implements proto.Packet interface.
func (i *PlayStatus) Read(buf *bytes.Buffer) {
	i.Status = buffer.ReadInt(buf)
}

// Write implements proto.Packet interface.
func (i *PlayStatus) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.WriteInt(buf, i.Status)
	return buf
}

// Disconnect needs to be documented.
type Disconnect struct {
	Message string
}

// Pid implements proto.Packet interface.
func (i *Disconnect) Pid() byte { return DisconnectHead }

// Read implements proto.Packet interface.
func (i *Disconnect) Read(buf *bytes.Buffer) {
	i.Message = buffer.ReadString(buf)
}

// Write implements proto.Packet interface.
func (i *Disconnect) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.WriteString(buf, i.Message)
	return buf
}

// Batch needs to be documented.
type Batch struct {
	Payloads [][]byte
}

// Pid implements proto.Packet interface.
func (i Batch) Pid() byte { return BatchHead } // 0x92

// Read implements proto.Packet interface.
func (i *Batch) Read(buf *bytes.Buffer) {
	i.Payloads = make([][]byte, 0)
	payload, err := util.DecodeDeflate(buf.Next(int(buffer.ReadInt(buf))))
	if err != nil {
		log.Println("Error while decompressing Batch payload:", err)
		return
	}
	b := bytes.NewBuffer(payload)
	for b.Len() > 4 {
		size := buffer.ReadInt(b)
		pk := b.Next(int(size))
		if pk[0] == 0x92 {
			panic("Invalid BatchPacket inside BatchPacket")
		}
		i.Payloads = append(i.Payloads, pk)
	}
}

// Write implements proto.Packet interface.
func (i Batch) Write() *bytes.Buffer {
	b := new(bytes.Buffer)
	for _, pk := range i.Payloads {
		buffer.WriteInt(b, uint32(len(pk)))
		buffer.Write(b, pk)
	}
	payload := util.EncodeDeflate(b.Bytes())
	buf := new(bytes.Buffer)
	buffer.BatchWrite(buf, uint32(len(payload)), payload)
	return buf
}

// Packet-specific constants
const (
	TextTypeRaw byte = iota
	TextTypeChat
	TextTypeTranslation
	TextTypePopup
	TextTypeTip
	TextTypeSystem
)

// Text needs to be documented.
type Text struct {
	TextType byte
	Source   string
	Message  string
	Params   []string
}

// Pid implements proto.Packet interface.
func (i Text) Pid() byte { return TextHead } // 0x93

// Read implements proto.Packet interface.
func (i *Text) Read(buf *bytes.Buffer) {
	i.TextType = buffer.ReadByte(buf)
	switch i.TextType {
	case TextTypePopup, TextTypeChat:
		buffer.ReadAny(buf, &i.Source)
		fallthrough
	case TextTypeRaw, TextTypeTip, TextTypeSystem:
		buffer.ReadAny(buf, &i.Message)
	case TextTypeTranslation:
		buffer.ReadAny(buf, &i.Message)
		cnt := buffer.ReadByte(buf)
		i.Params = make([]string, cnt)
		for k := byte(0); k < cnt; k++ {
			i.Params[k] = buffer.ReadString(buf)
		}
	}
}

// Write implements proto.Packet interface.
func (i Text) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.WriteByte(buf, i.TextType)
	switch i.TextType {
	case TextTypePopup, TextTypeChat:
		buffer.WriteAny(buf, i.Source)
		fallthrough
	case TextTypeRaw, TextTypeTip, TextTypeSystem:
		buffer.WriteAny(buf, i.Message)
	case TextTypeTranslation:
		buffer.WriteAny(buf, &i.Message)
		buffer.WriteByte(buf, byte(len(i.Params)))
		for _, p := range i.Params {
			buffer.WriteAny(buf, p)
		}
	}
	return buf
}

// Packet-specific constants
const (
	DayTime     = 0
	SunsetTime  = 12000
	NightTime   = 14000
	SunriseTime = 23000
	FullTime    = 24000
)

// SetTime needs to be documented.
type SetTime struct {
	Time    uint32
	Started bool
}

// Pid implements proto.Packet interface.
func (i SetTime) Pid() byte { return SetTimeHead }

// Read implements proto.Packet interface.
func (i *SetTime) Read(buf *bytes.Buffer) {
	i.Time = uint32((buffer.ReadInt(buf) / 19200) * FullTime)
	i.Started = buffer.ReadBool(buf)
}

// Write implements proto.Packet interface.
func (i SetTime) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.WriteInt(buf, uint32((i.Time*19200)/FullTime))
	buffer.WriteBool(buf, i.Started)
	return buf
}

// StartGame needs to be documented.
type StartGame struct {
	Seed                   uint32
	Dimension              byte
	Generator              uint32
	Gamemode               uint32
	EntityID               uint64
	SpawnX, SpawnY, SpawnZ uint32
	X, Y, Z                float32
}

// Pid implements proto.Packet interface.
func (i StartGame) Pid() byte { return StartGameHead } // 0x95

// Read implements proto.Packet interface.
func (i *StartGame) Read(buf *bytes.Buffer) {
	buffer.BatchRead(buf, &i.Seed, &i.Dimension, &i.Generator,
		&i.Gamemode, &i.EntityID, &i.SpawnX,
		&i.SpawnY, &i.SpawnZ, &i.X,
		&i.Y, &i.Z)
}

// Write implements proto.Packet interface.
func (i StartGame) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.BatchWrite(buf, i.Seed, i.Dimension, i.Generator,
		i.Gamemode, i.EntityID, i.SpawnX,
		i.SpawnY, i.SpawnZ, i.X,
		i.Y, i.Z)
	buffer.WriteByte(buf, 0)
	return buf
}

// AddPlayer needs to be documented.
type AddPlayer struct {
	RawUUID                [16]byte
	Username               string
	EntityID               uint64
	X, Y, Z                float32
	SpeedX, SpeedY, SpeedZ float32
	BodyYaw, Yaw, Pitch    float32
	Metadata               []byte
}

// Pid implements proto.Packet interface.
func (i AddPlayer) Pid() byte { return AddPlayerHead }

// Read implements proto.Packet interface.
func (i *AddPlayer) Read(buf *bytes.Buffer) {
	copy(i.RawUUID[:], buf.Next(16))
	buffer.BatchRead(buf, &i.Username, &i.EntityID,
		&i.X, &i.Y, &i.Z,
		&i.SpeedX, &i.SpeedY, &i.SpeedZ,
		&i.BodyYaw, &i.Yaw, &i.Pitch)
	i.Metadata = buf.Bytes()
}

// Write implements proto.Packet interface.
func (i AddPlayer) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.BatchWrite(buf, i.RawUUID[:], i.Username, i.EntityID,
		i.X, i.Y, i.Z,
		i.SpeedX, i.SpeedY, i.SpeedZ,
		i.BodyYaw, i.Yaw, i.Pitch, i.Metadata)
	buffer.WriteByte(buf, 0x7f) // Temporal, TODO: implement metadata functions
	return buf
}

// RemovePlayer needs to be documented.
type RemovePlayer struct {
	EntityID uint64
	RawUUID  [16]byte
}

// Pid implements proto.Packet interface.
func (i RemovePlayer) Pid() byte { return RemovePlayerHead }

// Read implements proto.Packet interface.
func (i *RemovePlayer) Read(buf *bytes.Buffer) {
	i.EntityID = buffer.ReadLong(buf)
	copy(i.RawUUID[:], buf.Next(16))
}

// Write implements proto.Packet interface.
func (i RemovePlayer) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.WriteLong(buf, i.EntityID)
	buf.Write(i.RawUUID[:])
	return buf
}

// AddEntity needs to be documented.
type AddEntity struct {
	EntityID               uint64
	Type                   uint32
	X, Y, Z                float32
	SpeedX, SpeedY, SpeedZ float32
	Yaw, Pitch             float32
	Metadata               []byte
	Link1, Link2           uint64
	Link3                  byte
}

// Pid implements proto.Packet interface.
func (i AddEntity) Pid() byte { return AddEntityHead }

// Read implements proto.Packet interface.
func (i *AddEntity) Read(buf *bytes.Buffer) {
	buffer.BatchRead(buf, &i.EntityID, &i.Type,
		&i.X, &i.Y, &i.Z,
		&i.SpeedX, &i.SpeedY, &i.SpeedZ,
		&i.Yaw, &i.Pitch)
	i.Metadata = buf.Bytes()
	// TODO
}

// Write implements proto.Packet interface.
func (i AddEntity) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.BatchWrite(buf, i.EntityID, i.Type,
		i.X, i.Y, i.Z,
		i.SpeedX, i.SpeedY, i.SpeedZ,
		i.Yaw, i.Pitch)
	buffer.WriteByte(buf, 0x7f)
	buffer.BatchWrite(buf, i.Link1, i.Link2, i.Link3)
	return buf
}

// RemoveEntity needs to be documented.
type RemoveEntity struct {
	EntityID uint64
}

// Pid implements proto.Packet interface.
func (i RemoveEntity) Pid() byte { return RemoveEntityHead }

// Read implements proto.Packet interface.
func (i *RemoveEntity) Read(buf *bytes.Buffer) {
	i.EntityID = buffer.ReadLong(buf)
}

// Write implements proto.Packet interface.
func (i RemoveEntity) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.WriteLong(buf, i.EntityID)
	return buf
}

// AddItemEntity needs to be documented.
type AddItemEntity struct {
	EntityID uint64
	Item     *types.Item
	X        float32
	Y        float32
	Z        float32
	SpeedX   float32
	SpeedY   float32
	SpeedZ   float32
}

// Pid implements proto.Packet interface.
func (i AddItemEntity) Pid() byte { return AddItemEntityHead }

// Read implements proto.Packet interface.
func (i *AddItemEntity) Read(buf *bytes.Buffer) {
	i.EntityID = buffer.ReadLong(buf)
	i.Item = new(types.Item)
	i.Item.Read(buf)
	i.X = buffer.ReadFloat(buf)
	i.Y = buffer.ReadFloat(buf)
	i.Z = buffer.ReadFloat(buf)
	i.SpeedX = buffer.ReadFloat(buf)
	i.SpeedY = buffer.ReadFloat(buf)
	i.SpeedZ = buffer.ReadFloat(buf)
}

// Write implements proto.Packet interface.
func (i AddItemEntity) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.WriteLong(buf, i.EntityID)
	buf.Write(i.Item.Write())
	buffer.WriteFloat(buf, i.X)
	buffer.WriteFloat(buf, i.Y)
	buffer.WriteFloat(buf, i.Z)
	buffer.WriteFloat(buf, i.SpeedX)
	buffer.WriteFloat(buf, i.SpeedY)
	buffer.WriteFloat(buf, i.SpeedZ)
	return buf
}

// TakeItemEntity needs to be documented.
type TakeItemEntity struct {
	Target   uint64
	EntityID uint64
}

// Pid implements proto.Packet interface.
func (i TakeItemEntity) Pid() byte { return TakeItemEntityHead }

// Read implements proto.Packet interface.
func (i *TakeItemEntity) Read(buf *bytes.Buffer) {
	i.Target = buffer.ReadLong(buf)
	i.EntityID = buffer.ReadLong(buf)
}

// Write implements proto.Packet interface.
func (i TakeItemEntity) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.WriteLong(buf, i.Target)
	buffer.WriteLong(buf, i.EntityID)
	return buf
}

// MoveEntity needs to be documented.
type MoveEntity struct {
	EntityIDs []uint64
	EntityPos [][6]float32 // X, Y, Z, Yaw, HeadYaw, Pitch
}

// Pid implements proto.Packet interface.
func (i MoveEntity) Pid() byte { return MoveEntityHead }

// Read implements proto.Packet interface.
func (i *MoveEntity) Read(buf *bytes.Buffer) {
	entityCnt := buffer.ReadInt(buf)
	i.EntityIDs = make([]uint64, entityCnt)
	i.EntityPos = make([][6]float32, entityCnt)
	for j := uint32(0); j < entityCnt; j++ {
		i.EntityIDs[j] = buffer.ReadLong(buf)
		for k := 0; k < 6; k++ {
			i.EntityPos[j][k] = buffer.ReadFloat(buf)
		}
	}
}

// Write implements proto.Packet interface.
func (i MoveEntity) Write() *bytes.Buffer {
	if len(i.EntityIDs) != len(i.EntityPos) {
		panic("Entity data slice length mismatch")
	}
	buf := new(bytes.Buffer)
	buffer.WriteInt(buf, uint32(len(i.EntityIDs)))
	for k, e := range i.EntityIDs {
		buffer.WriteLong(buf, e)
		for j := 0; j < 6; j++ {
			buffer.WriteFloat(buf, i.EntityPos[k][j])
		}
	}
	return buf
}

// Packet-specific constants
const (
	ModeNormal   byte = 0
	ModeReset    byte = 1
	ModeRotation byte = 2
)

// MovePlayer needs to be documented.
type MovePlayer struct {
	EntityID uint64
	X        float32
	Y        float32
	Z        float32
	Yaw      float32
	BodyYaw  float32
	Pitch    float32
	Mode     byte
	OnGround byte
}

// Pid implements proto.Packet interface.
func (i MovePlayer) Pid() byte { return MovePlayerHead }

// Read implements proto.Packet interface.
func (i *MovePlayer) Read(buf *bytes.Buffer) {
	i.EntityID = buffer.ReadLong(buf)
	i.X = buffer.ReadFloat(buf)
	i.Y = buffer.ReadFloat(buf)
	i.Z = buffer.ReadFloat(buf)
	i.Yaw = buffer.ReadFloat(buf)
	i.BodyYaw = buffer.ReadFloat(buf)
	i.Pitch = buffer.ReadFloat(buf)
	i.Mode = buffer.ReadByte(buf)
	i.OnGround = buffer.ReadByte(buf)
}

// Write implements proto.Packet interface.
func (i MovePlayer) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.WriteLong(buf, i.EntityID)
	buffer.WriteFloat(buf, i.X)
	buffer.WriteFloat(buf, i.Y)
	buffer.WriteFloat(buf, i.Z)
	buffer.WriteFloat(buf, i.Yaw)
	buffer.WriteFloat(buf, i.BodyYaw)
	buffer.WriteFloat(buf, i.Pitch)
	buffer.WriteByte(buf, i.Mode)
	buffer.WriteByte(buf, i.OnGround)
	return buf
}

// RemoveBlock needs to be documented.
type RemoveBlock struct {
	EntityID uint64
	X, Z     uint32
	Y        byte
}

// Pid implements proto.Packet interface.
func (i RemoveBlock) Pid() byte { return RemoveBlockHead }

// Read implements proto.Packet interface.
func (i *RemoveBlock) Read(buf *bytes.Buffer) {
	i.EntityID = buffer.ReadLong(buf)
	i.X = buffer.ReadInt(buf)
	i.Z = buffer.ReadInt(buf)
	i.Y = buffer.ReadByte(buf)
}

// Write implements proto.Packet interface.
func (i RemoveBlock) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.WriteLong(buf, i.EntityID)
	buffer.WriteInt(buf, i.X)
	buffer.WriteInt(buf, i.Z)
	buffer.WriteByte(buf, i.Y)
	return buf
}

// Packet-specific constants
const (
	UpdateNone byte = (1 << iota) >> 1
	UpdateNeighbors
	UpdateNetwork
	UpdateNographic
	UpdatePriority
	UpdateAll         = UpdateNeighbors | UpdateNetwork
	UpdateAllPriority = UpdateAll | UpdatePriority
)

// BlockRecord needs to be documented.
type BlockRecord struct {
	X, Z  uint32
	Y     byte
	Block types.Block
	Flags byte
}

// UpdateBlock needs to be documented.
type UpdateBlock struct {
	BlockRecords []BlockRecord
}

// Pid implements proto.Packet interface.
func (i UpdateBlock) Pid() byte { return UpdateBlockHead }

// Read implements proto.Packet interface.
func (i *UpdateBlock) Read(buf *bytes.Buffer) {
	records := buffer.ReadInt(buf)
	i.BlockRecords = make([]BlockRecord, records)
	for k := uint32(0); k < records; k++ {
		x := buffer.ReadInt(buf)
		z := buffer.ReadInt(buf)
		y := buffer.ReadByte(buf)
		id := buffer.ReadByte(buf)
		flagMeta := buffer.ReadByte(buf)
		i.BlockRecords[k] = BlockRecord{X: x,
			Y: y,
			Z: z,
			Block: types.Block{
				ID:   id,
				Meta: flagMeta & 0x0f,
			},
			Flags: (flagMeta >> 4) & 0x0f,
		}
	}
}

// Write implements proto.Packet interface.
func (i UpdateBlock) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.WriteInt(buf, uint32(len(i.BlockRecords)))
	for _, record := range i.BlockRecords {
		buffer.BatchWrite(buf, record.X, record.Z, record.Y, record.Block.ID, (record.Flags<<4 | record.Block.Meta))
	}
	return buf
}

// AddPainting needs to be documented.
type AddPainting struct {
	EntityID  uint64
	X         uint32
	Y         uint32
	Z         uint32
	Direction uint32
	Title     string
}

// Pid implements proto.Packet interface.
func (i AddPainting) Pid() byte { return AddPaintingHead }

// Read implements proto.Packet interface.
func (i *AddPainting) Read(buf *bytes.Buffer) {
	i.EntityID = buffer.ReadLong(buf)
	i.X = buffer.ReadInt(buf)
	i.Y = buffer.ReadInt(buf)
	i.Z = buffer.ReadInt(buf)
	i.Direction = buffer.ReadInt(buf)
	i.Title = buffer.ReadString(buf)
}

// Write implements proto.Packet interface.
func (i AddPainting) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.WriteLong(buf, i.EntityID)
	buffer.WriteInt(buf, i.X)
	buffer.WriteInt(buf, i.Y)
	buffer.WriteInt(buf, i.Z)
	buffer.WriteInt(buf, i.Direction)
	buffer.WriteString(buf, i.Title)
	return buf
}

// Explode needs to be documented.
type Explode struct {
	X, Y, Z, Radius float32
	Records         [][3]byte // X, Y, Z byte
}

// Pid implements proto.Packet interface.
func (i Explode) Pid() byte { return ExplodeHead }

// Read implements proto.Packet interface.
func (i *Explode) Read(buf *bytes.Buffer) {
	buffer.BatchRead(buf, &i.X, &i.Y, &i.Z, &i.Radius)
	cnt := buffer.ReadInt(buf)
	i.Records = make([][3]byte, cnt)
	for k := uint32(0); k < cnt; k++ {
		buffer.BatchRead(buf, &i.Records[k][0], &i.Records[k][1], &i.Records[k][2])
	}
}

// Write implements proto.Packet interface.
func (i Explode) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.BatchWrite(buf, i.X, i.Y, i.Z, i.Radius)
	buffer.WriteInt(buf, uint32(len(i.Records)))
	for _, r := range i.Records {
		buffer.BatchWrite(buf, r[0], r[1], r[2])
	}
	return buf
}

// Packet-specific constants
const (
	EventSoundClick            = 1000
	EventSoundClickFail        = 1001
	EventSoundShoot            = 1002
	EventSoundDoor             = 1003
	EventSoundFizz             = 1004
	EventSoundGhast            = 1007
	EventSoundGhastShoot       = 1008
	EventSoundBlazeShoot       = 1009
	EventSoundDoorBump         = 1010
	EventSoundDoorCrash        = 1012
	EventSoundBatFly           = 1015
	EventSoundZombieInfect     = 1016
	EventSoundZombieHeal       = 1017
	EventSoundEndermanTeleport = 1018
	EventSoundAnvilBreak       = 1020
	EventSoundAnvilUse         = 1021
	EventSoundAnvilFall        = 1022
	EventParticleShoot         = 2000
	EventParticleDestroy       = 2001
	EventParticleSplash        = 2002
	EventParticleEyeDespawn    = 2003
	EventParticleSpawn         = 2004
	EventStartRain             = 3001
	EventStartThunder          = 3002
	EventStopRain              = 3003
	EventStopThunder           = 3004
	EventSetData               = 4000
	EventPlayersSleeping       = 9800
)

// LevelEvent needs to be documented.
type LevelEvent struct {
	EventID uint16
	X       float32
	Y       float32
	Z       float32
	Data    uint32
}

// Pid implements proto.Packet interface.
func (i LevelEvent) Pid() byte { return LevelEventHead }

// Read implements proto.Packet interface.
func (i *LevelEvent) Read(buf *bytes.Buffer) {
	i.EventID = buffer.ReadShort(buf)
	i.X = buffer.ReadFloat(buf)
	i.Y = buffer.ReadFloat(buf)
	i.Z = buffer.ReadFloat(buf)
	i.Data = buffer.ReadInt(buf)
}

// Write implements proto.Packet interface.
func (i LevelEvent) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.WriteShort(buf, i.EventID)
	buffer.WriteFloat(buf, i.X)
	buffer.WriteFloat(buf, i.Y)
	buffer.WriteFloat(buf, i.Z)
	buffer.WriteInt(buf, i.Data)
	return buf
}

// BlockEvent needs to be documented.
type BlockEvent struct {
	X     uint32
	Y     uint32
	Z     uint32
	Case1 uint32
	Case2 uint32
}

// Pid implements proto.Packet interface.
func (i BlockEvent) Pid() byte { return BlockEventHead }

// Read implements proto.Packet interface.
func (i *BlockEvent) Read(buf *bytes.Buffer) {
	i.X = buffer.ReadInt(buf)
	i.Y = buffer.ReadInt(buf)
	i.Z = buffer.ReadInt(buf)
	i.Case1 = buffer.ReadInt(buf)
	i.Case2 = buffer.ReadInt(buf)
}

// Write implements proto.Packet interface.
func (i BlockEvent) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.WriteInt(buf, i.X)
	buffer.WriteInt(buf, i.Y)
	buffer.WriteInt(buf, i.Z)
	buffer.WriteInt(buf, i.Case1)
	buffer.WriteInt(buf, i.Case2)
	return buf
}

// Packet-specific constants
const (
	EventHurtAnimation byte = iota + 2
	EventDeathAnimation
	_
	_
	EventTameFail
	EventTameSuccess
	EventShakeWet
	EventUseItem
	EventEatGrassAnimation
	EventFishHookBubble
	EventFishHookPosition
	EventFishHookHook
	EventFishHookTease
	EventSquidInkCloud
	EventAmbientSound
	EventRespawn
)

// EntityEvent needs to be documented.
type EntityEvent struct {
	EntityID uint64
	Event    byte
}

// Pid implements proto.Packet interface.
func (i EntityEvent) Pid() byte { return EntityEventHead }

// Read implements proto.Packet interface.
func (i *EntityEvent) Read(buf *bytes.Buffer) {
	i.EntityID = buffer.ReadLong(buf)
	i.Event = buffer.ReadByte(buf)
}

// Write implements proto.Packet interface.
func (i EntityEvent) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.WriteLong(buf, i.EntityID)
	buffer.WriteByte(buf, i.Event)
	return buf
}

// Packet-specific constants
const (
	EffectAdd byte = iota + 1
	EffectModify
	EffectRemove
)

// MobEffect needs to be documented.
type MobEffect struct {
	EntityID  uint64
	EventID   byte
	EffectID  byte
	Amplifier byte
	Particles bool
	Duration  uint32
}

// Pid implements proto.Packet interface.
func (i MobEffect) Pid() byte { return MobEffectHead }

// Read implements proto.Packet interface.
func (i *MobEffect) Read(buf *bytes.Buffer) {
	i.EntityID = buffer.ReadLong(buf)
	i.EventID = buffer.ReadByte(buf)
	i.EffectID = buffer.ReadByte(buf)
	i.Amplifier = buffer.ReadByte(buf)
	i.Particles = buffer.ReadBool(buf)
	i.Duration = buffer.ReadInt(buf)
}

// Write implements proto.Packet interface.
func (i MobEffect) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.WriteLong(buf, i.EntityID)
	buffer.WriteByte(buf, i.EventID)
	buffer.WriteByte(buf, i.EffectID)
	buffer.WriteByte(buf, i.Amplifier)
	buffer.WriteBool(buf, i.Particles)
	buffer.WriteInt(buf, i.Duration)
	return buf
}

// UpdateAttributes needs to be documented.
type UpdateAttributes struct {
	// TODO: implement this after NBT is done
}

// Pid implements proto.Packet interface.
func (i UpdateAttributes) Pid() byte { return UpdateAttributesHead }

// Read implements proto.Packet interface.
func (i *UpdateAttributes) Read(buf *bytes.Buffer) {}

// Write implements proto.Packet interface.
func (i UpdateAttributes) Write() *bytes.Buffer { return nil }

// MobEquipment needs to be documented.
type MobEquipment struct {
	EntityID     uint64
	Item         *types.Item
	Slot         byte
	SelectedSlot byte
}

// Pid implements proto.Packet interface.
func (i MobEquipment) Pid() byte { return MobEquipmentHead }

// Read implements proto.Packet interface.
func (i *MobEquipment) Read(buf *bytes.Buffer) {
	i.EntityID = buffer.ReadLong(buf)
	i.Item = new(types.Item)
	i.Item.Read(buf)
	i.Slot = buffer.ReadByte(buf)
	i.SelectedSlot = buffer.ReadByte(buf)
}

// Write implements proto.Packet interface.
func (i MobEquipment) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.WriteLong(buf, i.EntityID)
	buf.Write(i.Item.Write())
	buffer.WriteByte(buf, i.Slot)
	buffer.WriteByte(buf, i.SelectedSlot)
	return buf
}

// MobArmorEquipment needs to be documented.
type MobArmorEquipment struct {
	EntityID uint64
	Slots    [4]*types.Item
}

// Pid implements proto.Packet interface.
func (i MobArmorEquipment) Pid() byte { return MobArmorEquipmentHead }

// Read implements proto.Packet interface.
func (i *MobArmorEquipment) Read(buf *bytes.Buffer) {
	i.EntityID = buffer.ReadLong(buf)
	for j := range i.Slots {
		i.Slots[j] = new(types.Item)
		i.Slots[j].Read(buf)
	}
}

// Write implements proto.Packet interface.
func (i MobArmorEquipment) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.WriteLong(buf, i.EntityID)
	for j := range i.Slots {
		buf.Write(i.Slots[j].Write())
	}
	return buf
}

// Interact needs to be documented.
type Interact struct {
	Action byte
	Target uint64
}

// Pid implements proto.Packet interface.
func (i Interact) Pid() byte { return InteractHead }

// Read implements proto.Packet interface.
func (i *Interact) Read(buf *bytes.Buffer) {
	i.Action = buffer.ReadByte(buf)
	i.Target = buffer.ReadLong(buf)
}

// Write implements proto.Packet interface.
func (i Interact) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.WriteByte(buf, i.Action)
	buffer.WriteLong(buf, i.Target)
	return buf
}

// UseItem needs to be documented.
type UseItem struct {
	X, Y, Z                uint32
	Face                   byte
	FloatX, FloatY, FloatZ float32
	PosX, PosY, PosZ       float32
	Item                   *types.Item
}

// Pid implements proto.Packet interface.
func (i UseItem) Pid() byte { return UseItemHead }

// Read implements proto.Packet interface.
func (i *UseItem) Read(buf *bytes.Buffer) {
	buffer.BatchRead(buf, &i.X, &i.Y, &i.Z,
		&i.Face, &i.FloatX, &i.FloatY, &i.FloatZ,
		&i.PosX, &i.PosY, &i.PosZ)
	i.Item = new(types.Item)
	i.Item.Read(buf)
}

// Write implements proto.Packet interface.
func (i UseItem) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.BatchWrite(buf, i.X, i.Y, i.Z,
		i.Face, i.FloatX, i.FloatY, i.FloatZ,
		i.PosX, i.PosY, i.PosZ, i.Item.Write())
	return buf
}

// Packet-specific constants
const (
	ActionStartBreak uint32 = iota
	ActionAbortBreak
	ActionStopBreak
	_
	_
	ActionReleaseItem
	ActionStopSleeping
	ActionRespawn
	ActionJump
	ActionStartSprint
	ActionStopSprint
	ActionStartSneak
	ActionStopSneak
	ActionDimensionChange
)

// PlayerAction needs to be documented.
type PlayerAction struct {
	EntityID uint64
	Action   uint32
	X        uint32
	Y        uint32
	Z        uint32
	Face     uint32
}

// Pid implements proto.Packet interface.
func (i PlayerAction) Pid() byte { return PlayerActionHead }

// Read implements proto.Packet interface.
func (i *PlayerAction) Read(buf *bytes.Buffer) {
	i.EntityID = buffer.ReadLong(buf)
	i.Action = buffer.ReadInt(buf)
	i.X = buffer.ReadInt(buf)
	i.Y = buffer.ReadInt(buf)
	i.Z = buffer.ReadInt(buf)
	i.Face = buffer.ReadInt(buf)
}

// Write implements proto.Packet interface.
func (i PlayerAction) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.WriteLong(buf, i.EntityID)
	buffer.WriteInt(buf, i.Action)
	buffer.WriteInt(buf, i.X)
	buffer.WriteInt(buf, i.Y)
	buffer.WriteInt(buf, i.Z)
	buffer.WriteInt(buf, i.Face)
	return buf
}

// HurtArmor needs to be documented.
type HurtArmor struct {
	Health byte
}

// Pid implements proto.Packet interface.
func (i HurtArmor) Pid() byte { return HurtArmorHead }

// Read implements proto.Packet interface.
func (i *HurtArmor) Read(buf *bytes.Buffer) {
	i.Health = buffer.ReadByte(buf)
}

// Write implements proto.Packet interface.
func (i HurtArmor) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.WriteByte(buf, i.Health)
	return buf
}

// SetEntityData needs to be documented.
type SetEntityData struct{} // TODO Metadata

// Pid implements proto.Packet interface.
func (i SetEntityData) Pid() byte { return SetEntityDataHead }

// Read implements proto.Packet interface.
func (i *SetEntityData) Read(buf *bytes.Buffer) {}

// Write implements proto.Packet interface.
func (i SetEntityData) Write() *bytes.Buffer {
	return nil
}

// SetEntityMotion needs to be documented.
type SetEntityMotion struct {
	EntityIDs    []uint64
	EntityMotion [][6]float32 // X, Y, Z, Yaw, HeadYaw, Pitch
}

// Pid implements proto.Packet interface.
func (i SetEntityMotion) Pid() byte { return SetEntityMotionHead }

// Read implements proto.Packet interface.
func (i *SetEntityMotion) Read(buf *bytes.Buffer) {
	entityCnt := buffer.ReadInt(buf)
	i.EntityIDs = make([]uint64, entityCnt)
	i.EntityMotion = make([][6]float32, entityCnt)
	for j := uint32(0); j < entityCnt; j++ {
		i.EntityIDs[j] = buffer.ReadLong(buf)
		for k := 0; k < 6; k++ {
			i.EntityMotion[j][k] = buffer.ReadFloat(buf)
		}
	}
}

// Write implements proto.Packet interface.
func (i SetEntityMotion) Write() *bytes.Buffer {
	if len(i.EntityIDs) != len(i.EntityMotion) {
		panic("Entity data slice length mismatch")
	}
	buf := new(bytes.Buffer)
	buffer.WriteInt(buf, uint32(len(i.EntityIDs)))
	for k, e := range i.EntityIDs {
		buffer.WriteLong(buf, e)
		for j := 0; j < 6; j++ {
			buffer.WriteFloat(buf, i.EntityMotion[k][j])
		}
	}
	return buf
}

// SetEntityLink needs to be documented.
type SetEntityLink struct {
	From uint64
	To   uint64
	Type byte
}

// Pid implements proto.Packet interface.
func (i SetEntityLink) Pid() byte { return SetEntityLinkHead }

// Read implements proto.Packet interface.
func (i *SetEntityLink) Read(buf *bytes.Buffer) {
	i.From = buffer.ReadLong(buf)
	i.To = buffer.ReadLong(buf)
	i.Type = buffer.ReadByte(buf)
}

// Write implements proto.Packet interface.
func (i SetEntityLink) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.WriteLong(buf, i.From)
	buffer.WriteLong(buf, i.To)
	buffer.WriteByte(buf, i.Type)
	return buf
}

// SetHealth needs to be documented.
type SetHealth struct {
	Health uint32
}

// Pid implements proto.Packet interface.
func (i SetHealth) Pid() byte { return SetHealthHead }

// Read implements proto.Packet interface.
func (i *SetHealth) Read(buf *bytes.Buffer) {
	i.Health = buffer.ReadInt(buf)
}

// Write implements proto.Packet interface.
func (i SetHealth) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.WriteInt(buf, i.Health)
	return buf
}

// SetSpawnPosition needs to be documented.
type SetSpawnPosition struct {
	X uint32
	Y uint32
	Z uint32
}

// Pid implements proto.Packet interface.
func (i SetSpawnPosition) Pid() byte { return SetSpawnPositionHead }

// Read implements proto.Packet interface.
func (i *SetSpawnPosition) Read(buf *bytes.Buffer) {
	i.X = buffer.ReadInt(buf)
	i.Y = buffer.ReadInt(buf)
	i.Z = buffer.ReadInt(buf)
}

// Write implements proto.Packet interface.
func (i SetSpawnPosition) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.WriteInt(buf, i.X)
	buffer.WriteInt(buf, i.Y)
	buffer.WriteInt(buf, i.Z)
	return buf
}

// Animate needs to be documented.
type Animate struct {
	Action   byte
	EntityID uint64
}

// Pid implements proto.Packet interface.
func (i Animate) Pid() byte { return AnimateHead }

// Read implements proto.Packet interface.
func (i *Animate) Read(buf *bytes.Buffer) {
	i.Action = buffer.ReadByte(buf)
	i.EntityID = buffer.ReadLong(buf)
}

// Write implements proto.Packet interface.
func (i Animate) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.WriteByte(buf, i.Action)
	buffer.WriteLong(buf, i.EntityID)
	return buf
}

// Respawn needs to be documented.
type Respawn struct {
	X float32
	Y float32
	Z float32
}

// Pid implements proto.Packet interface.
func (i Respawn) Pid() byte { return RespawnHead }

// Read implements proto.Packet interface.
func (i *Respawn) Read(buf *bytes.Buffer) {
	i.X = buffer.ReadFloat(buf)
	i.Y = buffer.ReadFloat(buf)
	i.Z = buffer.ReadFloat(buf)
}

// Write implements proto.Packet interface.
func (i Respawn) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.WriteFloat(buf, i.X)
	buffer.WriteFloat(buf, i.Y)
	buffer.WriteFloat(buf, i.Z)
	return buf
}

// DropItem needs to be documented.
type DropItem struct {
	Type byte
	Item *types.Item
}

// Pid implements proto.Packet interface.
func (i DropItem) Pid() byte { return DropItemHead }

// Read implements proto.Packet interface.
func (i *DropItem) Read(buf *bytes.Buffer) {
	i.Type = buffer.ReadByte(buf)
	i.Item = new(types.Item)
	i.Item.Read(buf)
}

// Write implements proto.Packet interface.
func (i DropItem) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.BatchWrite(buf, i.Type, i.Item.Write())
	return buf
}

// ContainerOpen needs to be documented.
type ContainerOpen struct {
	WindowID byte
	Type     byte
	Slots    uint16
	X        uint32
	Y        uint32
	Z        uint32
}

// Pid implements proto.Packet interface.
func (i ContainerOpen) Pid() byte { return ContainerOpenHead }

// Read implements proto.Packet interface.
func (i *ContainerOpen) Read(buf *bytes.Buffer) {
	i.WindowID = buffer.ReadByte(buf)
	i.Type = buffer.ReadByte(buf)
	i.Slots = buffer.ReadShort(buf)
	i.X = buffer.ReadInt(buf)
	i.Y = buffer.ReadInt(buf)
	i.Z = buffer.ReadInt(buf)
}

// Write implements proto.Packet interface.
func (i ContainerOpen) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.WriteByte(buf, i.WindowID)
	buffer.WriteByte(buf, i.Type)
	buffer.WriteShort(buf, i.Slots)
	buffer.WriteInt(buf, i.X)
	buffer.WriteInt(buf, i.Y)
	buffer.WriteInt(buf, i.Z)
	return buf
}

// ContainerClose needs to be documented.
type ContainerClose struct {
	WindowID byte
}

// Pid implements proto.Packet interface.
func (i ContainerClose) Pid() byte { return ContainerCloseHead }

// Read implements proto.Packet interface.
func (i *ContainerClose) Read(buf *bytes.Buffer) {
	i.WindowID = buffer.ReadByte(buf)
}

// Write implements proto.Packet interface.
func (i ContainerClose) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.WriteByte(buf, i.WindowID)
	return buf
}

// ContainerSetSlot needs to be documented.
type ContainerSetSlot struct { // TODO: implement this after slots
	Windowid   byte
	Slot       uint16
	HotbarSlot uint16
	Item       *types.Item
}

// Pid implements proto.Packet interface.
func (i ContainerSetSlot) Pid() byte { return ContainerSetSlotHead }

// Read implements proto.Packet interface.
func (i *ContainerSetSlot) Read(buf *bytes.Buffer) {
	i.Windowid = buffer.ReadByte(buf)
	i.Slot = buffer.ReadShort(buf)
	i.HotbarSlot = buffer.ReadShort(buf)
	i.Item = new(types.Item)
	i.Item.Read(buf)
}

// Write implements proto.Packet interface.
func (i ContainerSetSlot) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.WriteByte(buf, i.Windowid)
	buffer.WriteShort(buf, i.Slot)
	buffer.WriteShort(buf, i.HotbarSlot)
	buf.Write(i.Item.Write())
	return buf
}

// ContainerSetData needs to be documented.
type ContainerSetData struct {
	WindowID byte
	Property uint16
	Value    uint16
}

// Pid implements proto.Packet interface.
func (i ContainerSetData) Pid() byte { return ContainerSetDataHead }

// Read implements proto.Packet interface.
func (i *ContainerSetData) Read(buf *bytes.Buffer) {
	i.WindowID = buffer.ReadByte(buf)
	i.Property = buffer.ReadShort(buf)
	i.Value = buffer.ReadShort(buf)
}

// Write implements proto.Packet interface.
func (i ContainerSetData) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.WriteByte(buf, i.WindowID)
	buffer.WriteShort(buf, i.Property)
	buffer.WriteShort(buf, i.Value)
	return buf
}

// Packet-specific constants
const (
	InventoryWindow byte = 0
	ArmorWindow     byte = 0x78
	CreativeWindow  byte = 0x79
)

// ContainerSetContent needs to be documented.
type ContainerSetContent struct {
	WindowID byte
	Slots    []types.Item
	Hotbar   []uint32
}

// Pid implements proto.Packet interface.
func (i ContainerSetContent) Pid() byte { return ContainerSetContentHead }

// Read implements proto.Packet interface.
func (i *ContainerSetContent) Read(buf *bytes.Buffer) {
	i.WindowID = buffer.ReadByte(buf)
	count := buffer.ReadShort(buf)
	i.Slots = make([]types.Item, count)
	for j := range i.Slots {
		if buf.Len() < 7 {
			break
		}
		i.Slots[j] = *new(types.Item)
		(&i.Slots[j]).Read(buf)
	}
	if i.WindowID == InventoryWindow {
		count := buffer.ReadShort(buf)
		i.Hotbar = make([]uint32, count)
		for j := range i.Hotbar {
			i.Hotbar[j] = buffer.ReadInt(buf)
		}
	}
}

// Write implements proto.Packet interface.
func (i ContainerSetContent) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.WriteByte(buf, i.WindowID)
	buffer.WriteShort(buf, uint16(len(i.Slots)))
	for _, slot := range i.Slots {
		buffer.Write(buf, slot.Write())
	}
	if i.WindowID == InventoryWindow {
		for _, h := range i.Hotbar {
			buffer.WriteInt(buf, h)
		}
	} else {
		buffer.WriteShort(buf, 0)
	}
	return buf
}

// CraftingData needs to be documented.
type CraftingData struct{} // TODO

// Pid implements proto.Packet interface.
func (i CraftingData) Pid() byte { return CraftingDataHead }

// Read implements proto.Packet interface.
func (i *CraftingData) Read(buf *bytes.Buffer) {}

// Write implements proto.Packet interface.
func (i CraftingData) Write() *bytes.Buffer { return nil }

// CraftingEvent needs to be documented.
type CraftingEvent struct{} // TODO

// Pid implements proto.Packet interface.
func (i CraftingEvent) Pid() byte { return CraftingEventHead }

// Read implements proto.Packet interface.
func (i *CraftingEvent) Read(buf *bytes.Buffer) {}

// Write implements proto.Packet interface.
func (i CraftingEvent) Write() *bytes.Buffer { return nil }

// AdventureSettings needs to be documented.
type AdventureSettings struct {
	Flags uint32
}

// Pid implements proto.Packet interface.
func (i AdventureSettings) Pid() byte { return AdventureSettingsHead }

// Read implements proto.Packet interface.
func (i *AdventureSettings) Read(buf *bytes.Buffer) {
	i.Flags = buffer.ReadInt(buf)
}

// Write implements proto.Packet interface.
func (i AdventureSettings) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.WriteInt(buf, i.Flags)
	return buf
}

// BlockEntityData needs to be documented.
type BlockEntityData struct {
	X        uint32
	Y        uint32
	Z        uint32
	NamedTag []byte
}

// Pid implements proto.Packet interface.
func (i BlockEntityData) Pid() byte { return BlockEntityDataHead }

// Read implements proto.Packet interface.
func (i *BlockEntityData) Read(buf *bytes.Buffer) {
	i.X = buffer.ReadInt(buf)
	i.Y = buffer.ReadInt(buf)
	i.Z = buffer.ReadInt(buf)
	i.NamedTag = buf.Bytes()
}

// Write implements proto.Packet interface.
func (i BlockEntityData) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.WriteInt(buf, i.X)
	buffer.WriteInt(buf, i.Y)
	buffer.WriteInt(buf, i.Z)
	buf.Write(i.NamedTag)
	return buf
}

// Packet-specific constants
const (
	OrderColumns byte = 0
	OrderLayered byte = 1
)

// FullChunkData needs to be documented.
type FullChunkData struct {
	ChunkX, ChunkZ uint32
	Order          byte
	Payload        []byte
}

// Pid implements proto.Packet interface.
func (i FullChunkData) Pid() byte { return FullChunkDataHead }

// Read implements proto.Packet interface.
func (i *FullChunkData) Read(buf *bytes.Buffer) {
	buffer.BatchRead(buf, &i.ChunkX, &i.ChunkZ, &i.Order)
	i.Payload = buf.Next(int(buffer.ReadInt(buf)))
}

// Write implements proto.Packet interface.
func (i FullChunkData) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.BatchWrite(buf, i.ChunkX, i.ChunkZ, i.Order,
		uint32(len(i.Payload)), i.Payload)
	return buf
}

// SetDifficulty needs to be documented.
type SetDifficulty struct {
	Difficulty uint32
}

// Pid implements proto.Packet interface.
func (i SetDifficulty) Pid() byte { return SetDifficultyHead }

// Read implements proto.Packet interface.
func (i *SetDifficulty) Read(buf *bytes.Buffer) {
	i.Difficulty = buffer.ReadInt(buf)
}

// Write implements proto.Packet interface.
func (i SetDifficulty) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.WriteInt(buf, i.Difficulty)
	return buf
}

// SetPlayerGametype needs to be documented.
type SetPlayerGametype struct {
	Gamemode uint32
}

// Pid implements proto.Packet interface.
func (i SetPlayerGametype) Pid() byte { return SetPlayerGametypeHead }

// Read implements proto.Packet interface.
func (i *SetPlayerGametype) Read(buf *bytes.Buffer) {
	i.Gamemode = buffer.ReadInt(buf)
}

// Write implements proto.Packet interface.
func (i SetPlayerGametype) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.WriteInt(buf, i.Gamemode)
	return buf
}

// PlayerListEntry needs to be documented.
type PlayerListEntry struct {
	RawUUID            [16]byte
	EntityID           uint64
	Username, Skinname string
	Skin               []byte
}

// Packet-specific constants
const (
	PlayerListRemove byte = 0 // UUID only
	PlayerListAdd    byte = 1 // Everything!
)

// PlayerList needs to be documented.
type PlayerList struct {
	Type          byte
	PlayerEntries []PlayerListEntry
}

// Pid implements proto.Packet interface.
func (i PlayerList) Pid() byte { return PlayerListHead }

// Read implements proto.Packet interface.
func (i *PlayerList) Read(buf *bytes.Buffer) {
	i.Type = buffer.ReadByte(buf)
	entryCnt := buffer.ReadInt(buf)
	i.PlayerEntries = make([]PlayerListEntry, entryCnt)
	for k := uint32(0); k < entryCnt; k++ {
		entry := PlayerListEntry{}
		copy(entry.RawUUID[:], buf.Next(16))
		if i.Type == PlayerListRemove {
			i.PlayerEntries[k] = entry
			continue
		}
		entry.EntityID = buffer.ReadLong(buf)
		entry.Username = buffer.ReadString(buf)
		entry.Skinname = buffer.ReadString(buf)
		entry.Skin = []byte(buffer.ReadString(buf))
		i.PlayerEntries[k] = entry
	}
}

// Write implements proto.Packet interface.
func (i PlayerList) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.WriteByte(buf, i.Type)
	buffer.WriteInt(buf, uint32(len(i.PlayerEntries)))
	for _, entry := range i.PlayerEntries {
		buf.Write(entry.RawUUID[:])
		if i.Type == PlayerListRemove {
			continue
		}
		buffer.WriteLong(buf, entry.EntityID)
		buffer.WriteString(buf, entry.Username)
		buffer.WriteString(buf, entry.Skinname)
		buffer.WriteShort(buf, uint16(len(entry.Skin)))
		buffer.Write(buf, entry.Skin)
	}
	return buf
}

// RequestChunkRadius needs to be documented.
type RequestChunkRadius struct {
	Radius uint32
}

// Pid implements proto.Packet interface.
func (i RequestChunkRadius) Pid() byte { return RequestChunkRadiusHead }

// Read implements proto.Packet interface.
func (i *RequestChunkRadius) Read(buf *bytes.Buffer) {
	i.Radius = buffer.ReadInt(buf)
}

// Write implements proto.Packet interface.
func (i RequestChunkRadius) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.WriteInt(buf, i.Radius)
	return buf
}

// ChunkRadiusUpdate needs to be documented.
type ChunkRadiusUpdate struct {
	Radius uint32
}

// Pid implements proto.Packet interface.
func (i ChunkRadiusUpdate) Pid() byte { return ChunkRadiusUpdateHead }

// Read implements proto.Packet interface.
func (i *ChunkRadiusUpdate) Read(buf *bytes.Buffer) {
	i.Radius = buffer.ReadInt(buf)
}

// Write implements proto.Packet interface.
func (i ChunkRadiusUpdate) Write() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buffer.WriteInt(buf, i.Radius)
	return buf
}
