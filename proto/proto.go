// Package proto provides MCPE network protocol, defined by Mojang.
package proto

import (
	"log"

	"github.com/L7-MCPE/lav7/level"
	"github.com/L7-MCPE/lav7/raknet"
	"github.com/L7-MCPE/lav7/types"
	"github.com/L7-MCPE/lav7/util"
	"github.com/L7-MCPE/lav7/util/buffer"
)

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
	_ // 0xbe is skipped
	FullChunkDataHead
	SetDifficultyHead
	_ // 0xc1 is skipped
	SetPlayerGametypeHead
	PlayerListHead
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
}

type Packet interface {
	Pid() byte
	Read(*buffer.Buffer)
	Write() *buffer.Buffer
}

// GetPackets returns Packet struct with given pid.
func GetPacket(pid byte) Packet {
	pk, _ := packets[pid]
	return pk
}

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

func (i Login) Pid() byte { return LoginHead } // 0x8f

func (i *Login) Read(buf *buffer.Buffer) {
	buf.BatchRead(&i.Username, &i.Proto1)
	if i.Proto1 < raknet.MinecraftProtocol { // Old protocol
		return
	}
	buf.BatchRead(&i.Proto2, &i.ClientID)
	copy(i.RawUUID[:], buf.Read(16))
	buf.BatchRead(&i.ServerAddress, &i.ClientSecret, &i.SkinName)
	i.Skin = []byte(buf.ReadString())
}

func (i Login) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.BatchWrite(i.Username, i.Proto1, i.Proto2,
		i.ClientID, i.RawUUID[:], i.ServerAddress,
		i.ClientSecret, i.SkinName, string(i.Skin))
	return buf
}

const (
	LoginSuccess uint32 = iota
	LoginFailedClient
	LoginFailedServer
	PlayerSpawn
)

type PlayStatus struct {
	Status uint32
}

func (i *PlayStatus) Pid() byte { return PlayStatusHead }

func (i *PlayStatus) Read(buf *buffer.Buffer) {
	i.Status = buf.ReadInt()
}

func (i *PlayStatus) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteInt(i.Status)
	return buf
}

type Disconnect struct {
	Message string
}

func (i *Disconnect) Pid() byte { return DisconnectHead }

func (i *Disconnect) Read(buf *buffer.Buffer) {
	i.Message = buf.ReadString()
}

func (i *Disconnect) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteString(i.Message)
	return buf
}

type Batch struct {
	Payloads [][]byte
}

func (i Batch) Pid() byte { return BatchHead } // 0x92

func (i *Batch) Read(buf *buffer.Buffer) {
	i.Payloads = make([][]byte, 0)
	payload, err := util.DecodeDeflate(buf.Read(uint32(buf.ReadInt())))
	if err != nil {
		log.Println("Error while decompressing Batch payload:", err)
		return
	}
	b := buffer.FromBytes(payload)
	for b.Require(4) {
		size := b.ReadInt()
		pk := b.Read(size)
		if pk[0] == 0x92 {
			panic("Invalid BatchPacket inside BatchPacket")
		}
		i.Payloads = append(i.Payloads, pk)
	}
}

func (i Batch) Write() *buffer.Buffer {
	b := new(buffer.Buffer)
	for _, pk := range i.Payloads {
		b.WriteInt(uint32(len(pk)))
		b.Write(pk)
	}
	payload := util.EncodeDeflate(b.Done())
	buf := new(buffer.Buffer)
	buf.BatchWrite(uint32(len(payload)), payload)
	return buf
}

const (
	TextTypeRaw byte = iota
	TextTypeChat
	TextTypeTranslation
	TextTypePopup
	TextTypeTip
	TextTypeSystem
)

type Text struct {
	TextType byte
	Source   string
	Message  string
	Params   []string
}

func (i Text) Pid() byte { return TextHead } // 0x93

func (i *Text) Read(buf *buffer.Buffer) {
	i.TextType = buf.ReadByte()
	switch i.TextType {
	case TextTypePopup, TextTypeChat:
		buf.ReadAny(&i.Source)
		fallthrough
	case TextTypeRaw, TextTypeTip, TextTypeSystem:
		buf.ReadAny(&i.Message)
	case TextTypeTranslation:
		buf.ReadAny(&i.Message)
		cnt := buf.ReadByte()
		i.Params = make([]string, cnt)
		for k := byte(0); k < cnt; k++ {
			i.Params[k] = buf.ReadString()
		}
	}
}

func (i Text) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteByte(i.TextType)
	switch i.TextType {
	case TextTypePopup, TextTypeChat:
		buf.WriteAny(i.Source)
		fallthrough
	case TextTypeRaw, TextTypeTip, TextTypeSystem:
		buf.WriteAny(i.Message)
	case TextTypeTranslation:
		buf.WriteAny(&i.Message)
		buf.WriteByte(byte(len(i.Params)))
		for _, p := range i.Params {
			buf.WriteAny(p)
		}
	}
	return buf
}

type SetTime struct {
	Time    uint32
	Started bool
}

func (i SetTime) Pid() byte { return SetTimeHead }

func (i *SetTime) Read(buf *buffer.Buffer) {
	i.Time = uint32((buf.ReadInt() / 19200) * level.FullTime)
	i.Started = buf.ReadBool()
}

func (i SetTime) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteInt(uint32((i.Time / level.FullTime) * 19200))
	buf.WriteBool(i.Started)
	return buf
}

type StartGame struct {
	Seed                   uint32
	Dimension              byte
	Generator              uint32
	Gamemode               uint32
	EntityID               uint64
	SpawnX, SpawnY, SpawnZ uint32
	X, Y, Z                float32
}

func (i StartGame) Pid() byte { return StartGameHead } // 0x95

func (i *StartGame) Read(buf *buffer.Buffer) {
	buf.BatchRead(&i.Seed, &i.Dimension, &i.Generator,
		&i.Gamemode, &i.EntityID, &i.SpawnX,
		&i.SpawnY, &i.SpawnZ, &i.X,
		&i.Y, &i.Z)
}

func (i StartGame) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.BatchWrite(i.Seed, i.Dimension, i.Generator,
		i.Gamemode, i.EntityID, i.SpawnX,
		i.SpawnY, i.SpawnZ, i.X,
		i.Y, i.Z)
	buf.WriteByte(0)
	return buf
}

type AddPlayer struct {
	RawUUID                [16]byte
	Username               string
	EntityID               uint64
	X, Y, Z                float32
	SpeedX, SpeedY, SpeedZ float32
	BodyYaw, Yaw, Pitch    float32
	Metadata               []byte
}

func (i AddPlayer) Pid() byte { return AddPlayerHead }

func (i *AddPlayer) Read(buf *buffer.Buffer) {
	copy(i.RawUUID[:], buf.Read(16))
	buf.BatchRead(&i.Username, &i.EntityID,
		&i.X, &i.Y, &i.Z,
		&i.SpeedX, &i.SpeedY, &i.SpeedZ,
		&i.BodyYaw, &i.Yaw, &i.Pitch)
	i.Metadata = buf.Read(0)
}

func (i AddPlayer) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.BatchWrite(i.RawUUID[:], i.Username, i.EntityID,
		i.X, i.Y, i.Z,
		i.SpeedX, i.SpeedY, i.SpeedZ,
		i.BodyYaw, i.Yaw, i.Pitch, i.Metadata)
	buf.WriteByte(0x7f) // Temporal, TODO: implement metadata functions
	return buf
}

type RemovePlayer struct {
	EntityID uint64
	RawUUID  [16]byte
}

func (i RemovePlayer) Pid() byte { return RemovePlayerHead }

func (i *RemovePlayer) Read(buf *buffer.Buffer) {
	i.EntityID = buf.ReadLong()
	copy(i.RawUUID[:], buf.Read(16))
}

func (i RemovePlayer) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteLong(i.EntityID)
	buf.Write(i.RawUUID[:])
	return buf
}

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

func (i AddEntity) Pid() byte { return AddEntityHead }

func (i *AddEntity) Read(buf *buffer.Buffer) {
	buf.BatchRead(&i.EntityID, &i.Type,
		&i.X, &i.Y, &i.Z,
		&i.SpeedX, &i.SpeedY, &i.SpeedZ,
		&i.Yaw, &i.Pitch)
	i.Metadata = buf.Read(0)
	// TODO
}

func (i AddEntity) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.BatchWrite(i.EntityID, i.Type,
		i.X, i.Y, i.Z,
		i.SpeedX, i.SpeedY, i.SpeedZ,
		i.Yaw, i.Pitch)
	buf.WriteByte(0x7f)
	buf.BatchWrite(i.Link1, i.Link2, i.Link3)
	return buf
}

type RemoveEntity struct {
	EntityID uint64
}

func (i RemoveEntity) Pid() byte { return RemoveEntityHead }

func (i *RemoveEntity) Read(buf *buffer.Buffer) {
	i.EntityID = buf.ReadLong()
}

func (i RemoveEntity) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteLong(i.EntityID)
	return buf
}

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

func (i AddItemEntity) Pid() byte { return AddItemEntityHead }

func (i *AddItemEntity) Read(buf *buffer.Buffer) {
	i.EntityID = buf.ReadLong()
	i.Item = new(types.Item)
	i.Item.Read(buf)
	i.X = buf.ReadFloat()
	i.Y = buf.ReadFloat()
	i.Z = buf.ReadFloat()
	i.SpeedX = buf.ReadFloat()
	i.SpeedY = buf.ReadFloat()
	i.SpeedZ = buf.ReadFloat()
}

func (i AddItemEntity) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteLong(i.EntityID)
	buf.Write(i.Item.Write())
	buf.WriteFloat(i.X)
	buf.WriteFloat(i.Y)
	buf.WriteFloat(i.Z)
	buf.WriteFloat(i.SpeedX)
	buf.WriteFloat(i.SpeedY)
	buf.WriteFloat(i.SpeedZ)
	return buf
}

type TakeItemEntity struct {
	Target   uint64
	EntityID uint64
}

func (i TakeItemEntity) Pid() byte { return TakeItemEntityHead }

func (i *TakeItemEntity) Read(buf *buffer.Buffer) {
	i.Target = buf.ReadLong()
	i.EntityID = buf.ReadLong()
}

func (i TakeItemEntity) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteLong(i.Target)
	buf.WriteLong(i.EntityID)
	return buf
}

type MoveEntity struct {
	EntityIDs []uint64
	EntityPos [][6]float32 // X, Y, Z, Yaw, HeadYaw, Pitch
}

func (i MoveEntity) Pid() byte { return MoveEntityHead }

func (i *MoveEntity) Read(buf *buffer.Buffer) {
	entityCnt := buf.ReadInt()
	i.EntityIDs = make([]uint64, entityCnt)
	i.EntityPos = make([][6]float32, entityCnt)
	for j := uint32(0); j < entityCnt; j++ {
		i.EntityIDs[j] = buf.ReadLong()
		for k := 0; k < 6; k++ {
			i.EntityPos[j][k] = buf.ReadFloat()
		}
	}
}

func (i MoveEntity) Write() *buffer.Buffer {
	if len(i.EntityIDs) != len(i.EntityPos) {
		panic("Entity data slice length mismatch")
	}
	buf := new(buffer.Buffer)
	buf.WriteInt(uint32(len(i.EntityIDs)))
	for k, e := range i.EntityIDs {
		buf.WriteLong(e)
		for j := 0; j < 6; j++ {
			buf.WriteFloat(i.EntityPos[k][j])
		}
	}
	return buf
}

const (
	ModeNormal   byte = 0
	ModeReset    byte = 1
	ModeRotation byte = 2
)

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

func (i MovePlayer) Pid() byte { return MovePlayerHead }

func (i *MovePlayer) Read(buf *buffer.Buffer) {
	i.EntityID = buf.ReadLong()
	i.X = buf.ReadFloat()
	i.Y = buf.ReadFloat()
	i.Z = buf.ReadFloat()
	i.Yaw = buf.ReadFloat()
	i.BodyYaw = buf.ReadFloat()
	i.Pitch = buf.ReadFloat()
	i.Mode = buf.ReadByte()
	i.OnGround = buf.ReadByte()
}

func (i MovePlayer) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteLong(i.EntityID)
	buf.WriteFloat(i.X)
	buf.WriteFloat(i.Y)
	buf.WriteFloat(i.Z)
	buf.WriteFloat(i.Yaw)
	buf.WriteFloat(i.BodyYaw)
	buf.WriteFloat(i.Pitch)
	buf.WriteByte(i.Mode)
	buf.WriteByte(i.OnGround)
	return buf
}

type RemoveBlock struct {
	EntityID uint64
	X, Z     uint32
	Y        byte
}

func (i RemoveBlock) Pid() byte { return RemoveBlockHead }

func (i *RemoveBlock) Read(buf *buffer.Buffer) {
	i.EntityID = buf.ReadLong()
	i.X = buf.ReadInt()
	i.Z = buf.ReadInt()
	i.Y = buf.ReadByte()
}

func (i RemoveBlock) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteLong(i.EntityID)
	buf.WriteInt(i.X)
	buf.WriteInt(i.Z)
	buf.WriteByte(i.Y)
	return buf
}

const (
	UpdateNone byte = 1<<iota - 1
	UpdateNeighbors
	UpdateNetwork
	UpdateNographic
	UpdatePriority
	UpdateAll         = UpdateNeighbors | UpdateNetwork
	UpdateAllPriority = UpdateAll | UpdatePriority
)

type BlockRecord struct {
	X, Z  uint32
	Y     byte
	Block types.Block
	Flags byte
}

type UpdateBlock struct {
	BlockRecords []BlockRecord
}

func (i UpdateBlock) Pid() byte { return UpdateBlockHead }

func (i *UpdateBlock) Read(buf *buffer.Buffer) {
	records := buf.ReadInt()
	i.BlockRecords = make([]BlockRecord, records)
	for k := uint32(0); k < records; k++ {
		x := buf.ReadInt()
		z := buf.ReadInt()
		y := buf.ReadByte()
		id := buf.ReadByte()
		flagMeta := buf.ReadByte()
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

func (i UpdateBlock) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteInt(uint32(len(i.BlockRecords)))
	for _, record := range i.BlockRecords {
		buf.BatchWrite(record.X, record.Z, record.Y, record.Block.ID, (record.Flags<<4 | record.Block.Meta))
	}
	return buf
}

type AddPainting struct {
	EntityID  uint64
	X         uint32
	Y         uint32
	Z         uint32
	Direction uint32
	Title     string
}

func (i AddPainting) Pid() byte { return AddPaintingHead }

func (i *AddPainting) Read(buf *buffer.Buffer) {
	i.EntityID = buf.ReadLong()
	i.X = buf.ReadInt()
	i.Y = buf.ReadInt()
	i.Z = buf.ReadInt()
	i.Direction = buf.ReadInt()
	i.Title = buf.ReadString()
}

func (i AddPainting) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteLong(i.EntityID)
	buf.WriteInt(i.X)
	buf.WriteInt(i.Y)
	buf.WriteInt(i.Z)
	buf.WriteInt(i.Direction)
	buf.WriteString(i.Title)
	return buf
}

type Explode struct {
	X, Y, Z, Radius float32
	Records         [][3]byte // X, Y, Z byte
}

func (i Explode) Pid() byte { return ExplodeHead }

func (i *Explode) Read(buf *buffer.Buffer) {
	buf.BatchRead(&i.X, &i.Y, &i.Z, &i.Radius)
	cnt := buf.ReadInt()
	i.Records = make([][3]byte, cnt)
	for k := uint32(0); k < cnt; k++ {
		buf.BatchRead(&i.Records[k][0], &i.Records[k][1], &i.Records[k][2])
	}
}

func (i Explode) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.BatchWrite(i.X, i.Y, i.Z, i.Radius)
	buf.WriteInt(uint32(len(i.Records)))
	for _, r := range i.Records {
		buf.BatchWrite(r[0], r[1], r[2])
	}
	return buf
}

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

type LevelEvent struct {
	EventID uint16
	X       float32
	Y       float32
	Z       float32
	Data    uint32
}

func (i LevelEvent) Pid() byte { return LevelEventHead }

func (i *LevelEvent) Read(buf *buffer.Buffer) {
	i.EventID = buf.ReadShort()
	i.X = buf.ReadFloat()
	i.Y = buf.ReadFloat()
	i.Z = buf.ReadFloat()
	i.Data = buf.ReadInt()
}

func (i LevelEvent) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteShort(i.EventID)
	buf.WriteFloat(i.X)
	buf.WriteFloat(i.Y)
	buf.WriteFloat(i.Z)
	buf.WriteInt(i.Data)
	return buf
}

type BlockEvent struct {
	X     uint32
	Y     uint32
	Z     uint32
	Case1 uint32
	Case2 uint32
}

func (i BlockEvent) Pid() byte { return BlockEventHead }

func (i *BlockEvent) Read(buf *buffer.Buffer) {
	i.X = buf.ReadInt()
	i.Y = buf.ReadInt()
	i.Z = buf.ReadInt()
	i.Case1 = buf.ReadInt()
	i.Case2 = buf.ReadInt()
}

func (i BlockEvent) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteInt(i.X)
	buf.WriteInt(i.Y)
	buf.WriteInt(i.Z)
	buf.WriteInt(i.Case1)
	buf.WriteInt(i.Case2)
	return buf
}

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

type EntityEvent struct {
	EntityID uint64
	Event    byte
}

func (i EntityEvent) Pid() byte { return EntityEventHead }

func (i *EntityEvent) Read(buf *buffer.Buffer) {
	i.EntityID = buf.ReadLong()
	i.Event = buf.ReadByte()
}

func (i EntityEvent) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteLong(i.EntityID)
	buf.WriteByte(i.Event)
	return buf
}

const (
	EffectAdd byte = iota + 1
	EffectModify
	EffectRemove
)

type MobEffect struct {
	EntityID  uint64
	EventId   byte
	EffectId  byte
	Amplifier byte
	Particles bool
	Duration  uint32
}

func (i MobEffect) Pid() byte { return MobEffectHead }

func (i *MobEffect) Read(buf *buffer.Buffer) {
	i.EntityID = buf.ReadLong()
	i.EventId = buf.ReadByte()
	i.EffectId = buf.ReadByte()
	i.Amplifier = buf.ReadByte()
	i.Particles = buf.ReadBool()
	i.Duration = buf.ReadInt()
}

func (i MobEffect) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteLong(i.EntityID)
	buf.WriteByte(i.EventId)
	buf.WriteByte(i.EffectId)
	buf.WriteByte(i.Amplifier)
	buf.WriteBool(i.Particles)
	buf.WriteInt(i.Duration)
	return buf
}

type UpdateAttributes struct {
	// TODO: implement this after NBT is done
}

func (i UpdateAttributes) Pid() byte { return UpdateAttributesHead }

func (i *UpdateAttributes) Read(buf *buffer.Buffer) {}

func (i UpdateAttributes) Write() *buffer.Buffer { return nil }

type MobEquipment struct {
	EntityID     uint64
	Item         *types.Item
	Slot         byte
	SelectedSlot byte
}

func (i MobEquipment) Pid() byte { return MobEquipmentHead }

func (i *MobEquipment) Read(buf *buffer.Buffer) {
	i.EntityID = buf.ReadLong()
	i.Item = new(types.Item)
	i.Item.Read(buf)
	i.Slot = buf.ReadByte()
	i.SelectedSlot = buf.ReadByte()
}

func (i MobEquipment) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteLong(i.EntityID)
	buf.Write(i.Item.Write())
	buf.WriteByte(i.Slot)
	buf.WriteByte(i.SelectedSlot)
	return buf
}

type MobArmorEquipment struct {
	EntityID uint64
	Slots    [4]*types.Item
}

func (i MobArmorEquipment) Pid() byte { return MobArmorEquipmentHead }

func (i *MobArmorEquipment) Read(buf *buffer.Buffer) {
	i.EntityID = buf.ReadLong()
	for j := range i.Slots {
		i.Slots[j] = new(types.Item)
		i.Slots[j].Read(buf)
	}
}

func (i MobArmorEquipment) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteLong(i.EntityID)
	for j := range i.Slots {
		buf.Write(i.Slots[j].Write())
	}
	return buf
}

type Interact struct {
	Action byte
	Target uint64
}

func (i Interact) Pid() byte { return InteractHead }

func (i *Interact) Read(buf *buffer.Buffer) {
	i.Action = buf.ReadByte()
	i.Target = buf.ReadLong()
}

func (i Interact) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteByte(i.Action)
	buf.WriteLong(i.Target)
	return buf
}

type UseItem struct {
	X, Y, Z                uint32
	Face                   byte
	FloatX, FloatY, FloatZ float32
	PosX, PosY, PosZ       float32
	Item                   *types.Item
}

func (i UseItem) Pid() byte { return UseItemHead }

func (i *UseItem) Read(buf *buffer.Buffer) {
	buf.BatchRead(&i.X, &i.Y, &i.Z,
		&i.Face, &i.FloatX, &i.FloatY, &i.FloatZ,
		&i.PosX, &i.PosY, &i.PosZ)
	i.Item = new(types.Item)
	i.Item.Read(buf)
}

func (i UseItem) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.BatchWrite(i.X, i.Y, i.Z,
		i.Face, i.FloatX, i.FloatY, i.FloatZ,
		i.PosX, i.PosY, i.PosZ, i.Item.Write())
	return buf
}

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

type PlayerAction struct {
	EntityID uint64
	Action   uint32
	X        uint32
	Y        uint32
	Z        uint32
	Face     uint32
}

func (i PlayerAction) Pid() byte { return PlayerActionHead }

func (i *PlayerAction) Read(buf *buffer.Buffer) {
	i.EntityID = buf.ReadLong()
	i.Action = buf.ReadInt()
	i.X = buf.ReadInt()
	i.Y = buf.ReadInt()
	i.Z = buf.ReadInt()
	i.Face = buf.ReadInt()
}

func (i PlayerAction) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteLong(i.EntityID)
	buf.WriteInt(i.Action)
	buf.WriteInt(i.X)
	buf.WriteInt(i.Y)
	buf.WriteInt(i.Z)
	buf.WriteInt(i.Face)
	return buf
}

type HurtArmor struct {
	Health byte
}

func (i HurtArmor) Pid() byte { return HurtArmorHead }

func (i *HurtArmor) Read(buf *buffer.Buffer) {
	i.Health = buf.ReadByte()
}

func (i HurtArmor) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteByte(i.Health)
	return buf
}

type SetEntityData struct{} // TODO Metadata

func (i SetEntityData) Pid() byte { return SetEntityDataHead }

func (i *SetEntityData) Read(buf *buffer.Buffer) {}

func (i SetEntityData) Write() *buffer.Buffer {
	return nil
}

type SetEntityMotion struct {
	EntityIDs    []uint64
	EntityMotion [][6]float32 // X, Y, Z, Yaw, HeadYaw, Pitch
}

func (i SetEntityMotion) Pid() byte { return SetEntityMotionHead }

func (i *SetEntityMotion) Read(buf *buffer.Buffer) {
	entityCnt := buf.ReadInt()
	i.EntityIDs = make([]uint64, entityCnt)
	i.EntityMotion = make([][6]float32, entityCnt)
	for j := uint32(0); j < entityCnt; j++ {
		i.EntityIDs[j] = buf.ReadLong()
		for k := 0; k < 6; k++ {
			i.EntityMotion[j][k] = buf.ReadFloat()
		}
	}
}

func (i SetEntityMotion) Write() *buffer.Buffer {
	if len(i.EntityIDs) != len(i.EntityMotion) {
		panic("Entity data slice length mismatch")
	}
	buf := new(buffer.Buffer)
	buf.WriteInt(uint32(len(i.EntityIDs)))
	for k, e := range i.EntityIDs {
		buf.WriteLong(e)
		for j := 0; j < 6; j++ {
			buf.WriteFloat(i.EntityMotion[k][j])
		}
	}
	return buf
}

type SetEntityLink struct {
	From uint64
	To   uint64
	Type byte
}

func (i SetEntityLink) Pid() byte { return SetEntityLinkHead }

func (i *SetEntityLink) Read(buf *buffer.Buffer) {
	i.From = buf.ReadLong()
	i.To = buf.ReadLong()
	i.Type = buf.ReadByte()
}

func (i SetEntityLink) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteLong(i.From)
	buf.WriteLong(i.To)
	buf.WriteByte(i.Type)
	return buf
}

type SetHealth struct {
	Health uint32
}

func (i SetHealth) Pid() byte { return SetHealthHead }

func (i *SetHealth) Read(buf *buffer.Buffer) {
	i.Health = buf.ReadInt()
}

func (i SetHealth) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteInt(i.Health)
	return buf
}

type SetSpawnPosition struct {
	X uint32
	Y uint32
	Z uint32
}

func (i SetSpawnPosition) Pid() byte { return SetSpawnPositionHead }

func (i *SetSpawnPosition) Read(buf *buffer.Buffer) {
	i.X = buf.ReadInt()
	i.Y = buf.ReadInt()
	i.Z = buf.ReadInt()
}

func (i SetSpawnPosition) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteInt(i.X)
	buf.WriteInt(i.Y)
	buf.WriteInt(i.Z)
	return buf
}

type Animate struct {
	Action   byte
	EntityID uint64
}

func (i Animate) Pid() byte { return AnimateHead }

func (i *Animate) Read(buf *buffer.Buffer) {
	i.Action = buf.ReadByte()
	i.EntityID = buf.ReadLong()
}

func (i Animate) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteByte(i.Action)
	buf.WriteLong(i.EntityID)
	return buf
}

type Respawn struct {
	X float32
	Y float32
	Z float32
}

func (i Respawn) Pid() byte { return RespawnHead }

func (i *Respawn) Read(buf *buffer.Buffer) {
	i.X = buf.ReadFloat()
	i.Y = buf.ReadFloat()
	i.Z = buf.ReadFloat()
}

func (i Respawn) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteFloat(i.X)
	buf.WriteFloat(i.Y)
	buf.WriteFloat(i.Z)
	return buf
}

type DropItem struct {
	Type byte
	Item *types.Item
}

func (i DropItem) Pid() byte { return DropItemHead }

func (i *DropItem) Read(buf *buffer.Buffer) {
	i.Type = buf.ReadByte()
	i.Item = new(types.Item)
	i.Item.Read(buf)
}

func (i DropItem) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.BatchWrite(i.Type, i.Item.Write())
	return buf
}

type ContainerOpen struct {
	WindowID byte
	Type     byte
	Slots    uint16
	X        uint32
	Y        uint32
	Z        uint32
}

func (i ContainerOpen) Pid() byte { return ContainerOpenHead }

func (i *ContainerOpen) Read(buf *buffer.Buffer) {
	i.WindowID = buf.ReadByte()
	i.Type = buf.ReadByte()
	i.Slots = buf.ReadShort()
	i.X = buf.ReadInt()
	i.Y = buf.ReadInt()
	i.Z = buf.ReadInt()
}

func (i ContainerOpen) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteByte(i.WindowID)
	buf.WriteByte(i.Type)
	buf.WriteShort(i.Slots)
	buf.WriteInt(i.X)
	buf.WriteInt(i.Y)
	buf.WriteInt(i.Z)
	return buf
}

type ContainerClose struct {
	WindowID byte
}

func (i ContainerClose) Pid() byte { return ContainerCloseHead }

func (i *ContainerClose) Read(buf *buffer.Buffer) {
	i.WindowID = buf.ReadByte()
}

func (i ContainerClose) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteByte(i.WindowID)
	return buf
}

type ContainerSetSlot struct { // TODO: implement this after slots
	Windowid   byte
	Slot       uint16
	HotbarSlot uint16
	Item       *types.Item
}

func (i ContainerSetSlot) Pid() byte { return ContainerSetSlotHead }

func (i *ContainerSetSlot) Read(buf *buffer.Buffer) {
	i.Windowid = buf.ReadByte()
	i.Slot = buf.ReadShort()
	i.HotbarSlot = buf.ReadShort()
	i.Item = new(types.Item)
	i.Item.Read(buf)
}

func (i ContainerSetSlot) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteByte(i.Windowid)
	buf.WriteShort(i.Slot)
	buf.WriteShort(i.HotbarSlot)
	buf.Write(i.Item.Write())
	return buf
}

type ContainerSetData struct {
	WindowID byte
	Property uint16
	Value    uint16
}

func (i ContainerSetData) Pid() byte { return ContainerSetDataHead }

func (i *ContainerSetData) Read(buf *buffer.Buffer) {
	i.WindowID = buf.ReadByte()
	i.Property = buf.ReadShort()
	i.Value = buf.ReadShort()
}

func (i ContainerSetData) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteByte(i.WindowID)
	buf.WriteShort(i.Property)
	buf.WriteShort(i.Value)
	return buf
}

const (
	InventoryWindow byte = 0
	ArmorWindow     byte = 0x78
	CreativeWindow  byte = 0x79
)

type ContainerSetContent struct {
	WindowID byte
	Slots    []*types.Item
	Hotbar   []uint32
}

func (i ContainerSetContent) Pid() byte { return ContainerSetContentHead }

func (i *ContainerSetContent) Read(buf *buffer.Buffer) {
	i.WindowID = buf.ReadByte()
	count := buf.ReadShort()
	i.Slots = make([]*types.Item, count)
	for j := range i.Slots {
		if !buf.Require(0) {
			break
		}
		i.Slots[j] = new(types.Item)
		i.Slots[j].Read(buf)
	}
	if i.WindowID == InventoryWindow {
		count := buf.ReadShort()
		i.Hotbar = make([]uint32, count)
		for j := range i.Hotbar {
			i.Hotbar[j] = buf.ReadInt()
		}
	}
}

func (i ContainerSetContent) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteByte(i.WindowID)
	buf.WriteShort(uint16(len(i.Slots)))
	for _, slot := range i.Slots {
		buf.Write(slot.Write())
	}
	if i.WindowID == InventoryWindow {
		for _, h := range i.Hotbar {
			buf.WriteInt(h)
		}
	} else {
		buf.WriteShort(0)
	}
	return buf
}

type CraftingData struct{} // TODO

func (i CraftingData) Pid() byte { return CraftingDataHead }

func (i *CraftingData) Read(buf *buffer.Buffer) {}

func (i CraftingData) Write() *buffer.Buffer { return nil }

type CraftingEvent struct{} // TODO

func (i CraftingEvent) Pid() byte { return CraftingEventHead }

func (i *CraftingEvent) Read(buf *buffer.Buffer) {}

func (i CraftingEvent) Write() *buffer.Buffer { return nil }

type AdventureSettings struct {
	Flags uint32
}

func (i AdventureSettings) Pid() byte { return AdventureSettingsHead }

func (i *AdventureSettings) Read(buf *buffer.Buffer) {
	i.Flags = buf.ReadInt()
}

func (i AdventureSettings) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteInt(i.Flags)
	return buf
}

type BlockEntityData struct {
	X        uint32
	Y        uint32
	Z        uint32
	NamedTag []byte
}

func (i BlockEntityData) Pid() byte { return BlockEntityDataHead }

func (i *BlockEntityData) Read(buf *buffer.Buffer) {
	i.X = buf.ReadInt()
	i.Y = buf.ReadInt()
	i.Z = buf.ReadInt()
	i.NamedTag = buf.Read(0)
}

func (i BlockEntityData) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteInt(i.X)
	buf.WriteInt(i.Y)
	buf.WriteInt(i.Z)
	buf.Write(i.NamedTag)
	return buf
}

const (
	OrderColumns byte = 0
	OrderLayered byte = 1
)

type FullChunkData struct {
	ChunkX, ChunkZ uint32
	Order          byte
	Payload        []byte
}

func (i FullChunkData) Pid() byte { return FullChunkDataHead }

func (i *FullChunkData) Read(buf *buffer.Buffer) {
	buf.BatchRead(&i.ChunkX, &i.ChunkZ, &i.Order)
	i.Payload = buf.Read(uint32(buf.ReadInt()))
}

func (i FullChunkData) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.BatchWrite(i.ChunkX, i.ChunkZ, i.Order,
		uint32(len(i.Payload)), i.Payload)
	return buf
}

type SetDifficulty struct {
	Difficulty uint32
}

func (i SetDifficulty) Pid() byte { return SetDifficultyHead }

func (i *SetDifficulty) Read(buf *buffer.Buffer) {
	i.Difficulty = buf.ReadInt()
}

func (i SetDifficulty) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteInt(i.Difficulty)
	return buf
}

type SetPlayerGametype struct {
	Gamemode uint32
}

func (i SetPlayerGametype) Pid() byte { return SetPlayerGametypeHead }

func (i *SetPlayerGametype) Read(buf *buffer.Buffer) {
	i.Gamemode = buf.ReadInt()
}

func (i SetPlayerGametype) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteInt(i.Gamemode)
	return buf
}

type PlayerListEntry struct {
	RawUUID            [16]byte
	EntityID           uint64
	Username, Slimness string
	Skin               []byte
}

const (
	PlayerListRemove byte = 0
	PlayerListAdd    byte = 1
)

type PlayerList struct {
	Type          byte
	PlayerEntries []PlayerListEntry
}

func (i PlayerList) Pid() byte { return PlayerListHead }

func (i *PlayerList) Read(buf *buffer.Buffer) {
	i.Type = buf.ReadByte()
	entryCnt := buf.ReadInt()
	i.PlayerEntries = make([]PlayerListEntry, entryCnt)
	for k := uint32(0); k < entryCnt; k++ {
		entry := PlayerListEntry{}
		copy(entry.RawUUID[:], buf.Read(16))
		if i.Type == PlayerListRemove {
			i.PlayerEntries[k] = entry
			continue
		}
		entry.EntityID = buf.ReadLong()
		entry.Username = buf.ReadString()
		entry.Slimness = buf.ReadString()
		entry.Skin = []byte(buf.ReadString())
		i.PlayerEntries[k] = entry
	}
}

func (i PlayerList) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteByte(i.Type)
	buf.WriteInt(uint32(len(i.PlayerEntries)))
	for _, entry := range i.PlayerEntries {
		buf.Write(entry.RawUUID[:])
		if i.Type == PlayerListRemove {
			continue
		}
		buf.WriteLong(entry.EntityID)
		buf.WriteString(entry.Username)
		buf.WriteString(entry.Slimness)
		buf.WriteString(string(entry.Skin))
	}
	return buf
}
