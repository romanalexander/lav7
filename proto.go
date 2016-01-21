package lav7

import (
	"encoding/hex"

	"github.com/L7-MCPE/lav7/raknet"
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
	PlayerInputHead
	FullChunkDataHead
	SetDifficultyHead
	ChangeDimensionHead
	SetPlayerGametypeHead
	PlayerListHead
	TelemetryEventHead
	SpawnExperienceOrbHead
)

var packets = map[byte]Packet{
	LoginHead:         new(Login),
	BatchHead:         new(Batch),
	TextHead:          new(Text),
	MovePlayerHead:    new(MovePlayer),
	RemoveBlockHead:   new(RemoveBlock),
	FullChunkDataHead: new(FullChunkData),
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

type Batch struct {
	Payloads [][]byte
}

func (i Batch) Pid() byte { return BatchHead } // 0x92

func (i *Batch) Read(buf *buffer.Buffer) {
	i.Payloads = make([][]byte, 0)
	payload, err := util.DecodeDeflate(buf.Read(uint32(buf.ReadInt())))
	if err != nil {
		util.Debug("Error while decompressing Batch payload:", err)
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

const (
	MoveModeNormal byte = iota
	MoveModeReset
	MoveModeRotation
)

type MovePlayer struct {
	EntityID            uint64
	X, Y, Z             float32
	Yaw, BodyYaw, Pitch float32
	Mode                byte
	OnGround            bool
}

func (i MovePlayer) Pid() byte { return MovePlayerHead }

func (i *MovePlayer) Read(buf *buffer.Buffer) {
	buf.BatchRead(&i.EntityID, &i.X, &i.Y, &i.Z,
		&i.Yaw, &i.BodyYaw, &i.Pitch,
		&i.Mode, &i.OnGround)
}

func (i MovePlayer) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.BatchWrite(i.EntityID, i.X, i.Y, i.Z,
		i.Yaw, i.BodyYaw, i.Pitch,
		i.Mode, i.OnGround)
	return buf
}

type RemoveBlock struct {
	EntityID uint64
	X, Z     int32
	Y        byte
}

func (i RemoveBlock) Pid() byte { return RemoveBlockHead }

func (i *RemoveBlock) Read(buf *buffer.Buffer) {
	util.Debug(hex.Dump(buf.Payload), buf.Offset)
	var x, z uint32
	buf.BatchRead(&i.EntityID, &x, &z, &i.Y)
	i.X, i.Z = int32(x), int32(z)
	util.Debug(i.EntityID, i.X, i.Z, i.Y)
}

func (i *RemoveBlock) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.BatchWrite(i.EntityID, uint32(i.X), uint32(i.Z), i.Y)
	return buf
}

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
