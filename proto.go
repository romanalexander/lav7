package lav7

import (
	"github.com/L7-MCPE/lav7/level"
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
	LoginSuccess      byte = 0
	LoginFailedClient byte = 1
	LoginFailedServer byte = 2
	PlayerSpawn       byte = 3
)

type PlayStatus struct {
	Status uint32
}

func (i *PlayStatus) Pid() byte { return PlayStatusHead }

func (i *PlayStatus) Read(buf *buffer.Buffer) {
	i.Status = buf.ReadStatus()
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
	i.Message = buf.ReadMessage()
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
	Yaw, HeadYaw, Pitch    float32
	Metadata               []byte
}

func (i AddPlayer) Pid() byte { return AddPlayerHead }

func (i *AddPlayer) Read(buf *buffer.Buffer) {
	copy(i.RawUUID[:], buf.Read(16))
	buf.BatchRead(&i.Username, &i.EntityID,
		&i.X, &i.Y, &i.Z,
		&i.SpeedX, &i.SpeedY, &i.SpeedZ,
		&i.Yaw, &i.HeadYaw, &i.Pitch)
	i.MetaData = i.Read(0)
}

func (i AddPlayer) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.BatchWrite(i.RawUUID[:], i.Username, i.EntityID,
		i.X, i.Y, i.Z,
		i.SpeedX, i.SpeedY, i.SpeedZ,
		i.Yaw, i.HeadYaw, i.Pitch, i.MetaData)
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
		&i.Yaw, &Pitch)
	i.MetaData = i.Read(0)
	// TODO
}

func (i AddEntity) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.BatchWrite(i.RawUUID[:], i.Username, i.EntityID,
		i.X, i.Y, i.Z,
		i.SpeedX, i.SpeedY, i.SpeedZ,
		i.Yaw, i.HeadRot)
	buf.Write(0x7f)
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
	// TODO: implement slot functions
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
	buf.BatchWrite(uint16(1), byte(1), uint16(1), uint16(0)) //TODO
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
	entityCnt = buf.ReadInt()
	i.EntityIDs = make([]uint64, entityCnt)
	i.EntityPos = make([][6]float32, entityCnt)
	for j := uint32(0); j < entityCnt; j++ {
		i.EntityIDs[j] = buf.ReadLong()
		for k := 0; k < 6; k++ {
			i.EntityPos[i] = buf.ReadFloat()
		}
	}
}

func (i MoveEntity) Write() *buffer.Buffer {
	if len(i.EntityIDs) != len(i.EntityPos) {
		panic("Entity data slice length mismatch")
	}
	buf := new(buffer.Buffer)
	buf.WriteInt(len(i.EntityIDs))
	for k, e := range i.EntityIDs {
		buf.WriteLong(e)
		for j := 0; j < 6; j++ {
			buf.WriteFloat(i.EntityPos[k][j])
		}
	}
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
	i.X = buf.Readint()
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
	X, Z            uint32
	Y               byte
	ID, Meta, Flags byte
}

type UpdateBlock struct {
	BlockRecords []BlockRecord
}

func (i UpdateBlock) Pid() byte { return UpdateBlockHead }

func (i *UpdateBlock) Read(buf *buffer.Buffer) {
	records := buf.ReadInt()
	i.BlockRecords = make([]BlockRecord, records)
	for k := 0; k < records; k++ {
		x := buf.ReadInt()
		z := buf.ReadInt()
		y := buf.ReadByte()
		id := buf.ReadByte()
		flagMeta := buf.ReadByte()
		i[k] = BlockRecord{X: x, Y: y, Z: z, ID: id, Meta: flagMeta & 0x0f, Flags: (flagMeta >> 4) & 0x0f}
	}
}

func (i UpdateBlock) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteInt(len(i.BlockRecords))
	for _, record := range i.BlockRecords {
		i.BatchWrite(record.X, record.Z, record.Y, record.ID, (record.Flags<<4 | record.Meta))
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

// TODO: Explode

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
	Evid uint16
	X    float32
	Y    float32
	Z    float32
	Data uint32
}

func (i LevelEvent) Pid() byte { return LevelEventHead }

func (i *LevelEvent) Read(buf *buffer.Buffer) {
	i.Evid = buf.ReadShort()
	i.X = buf.ReadFloat()
	i.Y = buf.ReadFloat()
	i.Z = buf.ReadFloat()
	i.Data = buf.ReadInt()
}

func (i LevelEvent) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteShort(i.Evid)
	buf.WriteFloat(i.X)
	buf.WriteFloat(i.Y)
	buf.WriteFloat(i.Z)
	buf.WriteInt(i.Data)
	return buf
}

type BlockEvent struct {
	X    uint32
	Y    uint32
	Z    uint32
	Case uint32
	Case uint32
}

func (i BlockEvent) Pid() byte { return BlockEventHead }

func (i *BlockEvent) Read(buf *buffer.Buffer) {
	i.X = buf.ReadInt()
	i.Y = buf.ReadInt()
	i.Z = buf.ReadInt()
	i.Case = buf.ReadInt()
	i.Case = buf.ReadInt()
}

func (i BlockEvent) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteInt(i.X)
	buf.WriteInt(i.Y)
	buf.WriteInt(i.Z)
	buf.WriteInt(i.Case)
	buf.WriteInt(i.Case)
	return buf
}

const (
	EventHurtAnimation     byte = 2
	EventDeathAnimation    byte = 3
	EventTameFail          byte = 6
	EventTameSuccess       byte = 7
	EventShakeWet          byte = 8
	EventUseItem           byte = 9
	EventEatGrassAnimation byte = 10
	EventFishHookBubble    byte = 11
	EventFishHookPosition  byte = 12
	EventFishHookHook      byte = 13
	EventFishHookTease     byte = 14
	EventSquidInkCloud     byte = 15
	EventAmbientSound      byte = 16
	EventRespawn           byte = 17
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
	EventAdd    byte = 1
	EventModify byte = 2
	EventRemove byte = 3
)

type MobEffect struct {
	EntityID  uint64
	EventId   byte
	EffectId  byte
	Amplifier byte
	Particles byte
	Duration  uint32
}

func (i MobEffect) Pid() byte { return MobEffectHead }

func (i *MobEffect) Read(buf *buffer.Buffer) {
	i.EntityID = buf.ReadLong()
	i.EventId = buf.ReadByte()
	i.EffectId = buf.ReadByte()
	i.Amplifier = buf.ReadByte()
	i.Particles = buf.ReadByte()
	i.Duration = buf.ReadInt()
}

func (i MobEffect) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteLong(i.EntityID)
	buf.WriteByte(i.EventId)
	buf.WriteByte(i.EffectId)
	buf.WriteByte(i.Amplifier)
	buf.WriteByte(i.Particles)
	buf.WriteInt(i.Duration)
	return buf
}

// An exception was thrown while parsing/converting PocketMine-MP protocol.
// Please read original PHP code and port it manually.
// Exception: string index out of range
//
// 	const NETWORK_ID = Info::UPDATE_ATTRIBUTES_PACKET;
//
//
// 	public $entityId;
// 	/** @var Attribute[] */
// 	public $entries = [];
//
// 	public function decode(){
//
// 	}
//
// 	public function encode(){
// 		$this->reset();
//
// 		$this->putLong($this->entityId);
//
// 		$this->putShort(count($this->entries));
//
// 		foreach($this->entries as $entry){
// 			$this->putFloat($entry->getMinValue());
// 			$this->putFloat($entry->getMaxValue());
// 			$this->putFloat($entry->getValue());
// 			$this->putString($entry->getName());
// 		}
// 	}
//
//

type MobEquipment struct {
	EntityID     uint64
	Slot         byte
	SelectedSlot byte
}

func (i MobEquipment) Pid() byte { return MobEquipmentHead }

func (i *MobEquipment) Read(buf *buffer.Buffer) {
	i.EntityID = buf.ReadLong()
	// Unexpected code:Slot Item
	i.Slot = buf.ReadByte()
	i.SelectedSlot = buf.ReadByte()
}

func (i MobEquipment) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteLong(i.EntityID)
	// Unexpected code:Slot Item
	buf.WriteByte(i.Slot)
	buf.WriteByte(i.SelectedSlot)
	return buf
}

type MobArmorEquipment struct {
	EntityID uint64
}

func (i MobArmorEquipment) Pid() byte { return MobArmorEquipmentHead }

func (i *MobArmorEquipment) Read(buf *buffer.Buffer) {
	i.EntityID = buf.ReadLong()
	// Unexpected code:Slot Slots
	// Unexpected code:Slot Slots
	// Unexpected code:Slot Slots
	// Unexpected code:Slot Slots
}

func (i MobArmorEquipment) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteLong(i.EntityID)
	// Unexpected code:Slot Slots
	// Unexpected code:Slot Slots
	// Unexpected code:Slot Slots
	// Unexpected code:Slot Slots
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

type UseItem struct{}

func (i UseItem) Pid() byte { return UseItemHead }

func (i *UseItem) Read(buf *buffer.Buffer) {}

func (i UseItem) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	return buf
}

const (
	ActionStartBreak      byte = 0
	ActionAbortBreak      byte = 1
	ActionStopBreak       byte = 2
	ActionReleaseItem     byte = 5
	ActionStopSleeping    byte = 6
	ActionRespawn         byte = 7
	ActionJump            byte = 8
	ActionStartSprint     byte = 9
	ActionStopSprint      byte = 10
	ActionStartSneak      byte = 11
	ActionStopSneak       byte = 12
	ActionDimensionChange byte = 13
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

// An exception was thrown while parsing/converting PocketMine-MP protocol.
// Please read original PHP code and port it manually.
// Exception: 'int' object is not subscriptable
//
// 	const NETWORK_ID = Info::SET_ENTITY_DATA_PACKET;
//
// 	public $EntityID;
// 	public $metadata;
//
// 	public function decode(){
//
// 	}
//
// 	public function encode(){
// 		$this->reset();
// 		$this->putLong($this->EntityID);
// 		$meta = Binary::writeMetadata($this->metadata);
// 		$this->put($meta);
// 	}
//
//

// An exception was thrown while parsing/converting PocketMine-MP protocol.
// Please read original PHP code and port it manually.
// Exception: string index out of range
//
// 	const NETWORK_ID = Info::SET_ENTITY_MOTION_PACKET;
//
//
// 	// EntityID, motX, motY, motZ
// 	/** @var array[] */
// 	public $entities = [];
//
// 	public function clean(){
// 		$this->entities = [];
// 		return parent::clean();
// 	}
//
// 	public function decode(){
//
// 	}
//
// 	public function encode(){
// 		$this->reset();
// 		$this->putInt(count($this->entities));
// 		foreach($this->entities as $d){
// 			$this->putLong($d[0]); //EntityID
// 			$this->putFloat($d[1]); //motX
// 			$this->putFloat($d[2]); //motY
// 			$this->putFloat($d[3]); //motZ
// 		}
// 	}
//
//

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

type DropItem struct{}

func (i DropItem) Pid() byte { return DropItemHead }

func (i *DropItem) Read(buf *buffer.Buffer) {}

func (i DropItem) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	return buf
}

type ContainerOpen struct {
	Windowid byte
	Type     byte
	Slots    uint16
	X        uint32
	Y        uint32
	Z        uint32
}

func (i ContainerOpen) Pid() byte { return ContainerOpenHead }

func (i *ContainerOpen) Read(buf *buffer.Buffer) {
	i.Windowid = buf.ReadByte()
	i.Type = buf.ReadByte()
	i.Slots = buf.ReadShort()
	i.X = buf.ReadInt()
	i.Y = buf.ReadInt()
	i.Z = buf.ReadInt()
}

func (i ContainerOpen) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteByte(i.Windowid)
	buf.WriteByte(i.Type)
	buf.WriteShort(i.Slots)
	buf.WriteInt(i.X)
	buf.WriteInt(i.Y)
	buf.WriteInt(i.Z)
	return buf
}

type ContainerClose struct {
	Windowid byte
}

func (i ContainerClose) Pid() byte { return ContainerCloseHead }

func (i *ContainerClose) Read(buf *buffer.Buffer) {
	i.Windowid = buf.ReadByte()
}

func (i ContainerClose) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteByte(i.Windowid)
	return buf
}

type ContainerSetSlot struct {
	Windowid   byte
	Slot       uint16
	HotbarSlot uint16
}

func (i ContainerSetSlot) Pid() byte { return ContainerSetSlotHead }

func (i *ContainerSetSlot) Read(buf *buffer.Buffer) {
	i.Windowid = buf.ReadByte()
	i.Slot = buf.ReadShort()
	i.HotbarSlot = buf.ReadShort()
	// Unexpected code:Slot Item
}

func (i ContainerSetSlot) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteByte(i.Windowid)
	buf.WriteShort(i.Slot)
	buf.WriteShort(i.HotbarSlot)
	// Unexpected code:Slot Item
	return buf
}

type ContainerSetData struct {
	Windowid byte
	Property uint16
	Value    uint16
}

func (i ContainerSetData) Pid() byte { return ContainerSetDataHead }

func (i *ContainerSetData) Read(buf *buffer.Buffer) {
	i.Windowid = buf.ReadByte()
	i.Property = buf.ReadShort()
	i.Value = buf.ReadShort()
}

func (i ContainerSetData) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteByte(i.Windowid)
	buf.WriteShort(i.Property)
	buf.WriteShort(i.Value)
	return buf
}

// An exception was thrown while parsing/converting PocketMine-MP protocol.
// Please read original PHP code and port it manually.
// Exception: string index out of range
//
// 	const NETWORK_ID = Info::CONTAINER_SET_CONTENT_PACKET;
//
// 	const SPECIAL_INVENTORY = 0;
// 	const SPECIAL_ARMOR = 0x78;
// 	const SPECIAL_CREATIVE = 0x79;
//
// 	public $windowid;
// 	public $slots = [];
// 	public $hotbar = [];
//
// 	public function clean(){
// 		$this->slots = [];
// 		$this->hotbar = [];
// 		return parent::clean();
// 	}
//
// 	public function decode(){
// 		$this->windowid = $this->getByte();
// 		$count = $this->getShort();
// 		for($s = 0; $s < $count and !$this->feof(); ++$s){
// 			$this->slots[$s] = $this->getSlot();
// 		}
// 		if($this->windowid === self::SPECIAL_INVENTORY){
// 			$count = $this->getShort();
// 			for($s = 0; $s < $count and !$this->feof(); ++$s){
// 				$this->hotbar[$s] = $this->getInt();
// 			}
// 		}
// 	}
//
// 	public function encode(){
// 		$this->reset();
// 		$this->putByte($this->windowid);
// 		$this->putShort(count($this->slots));
// 		foreach($this->slots as $slot){
// 			$this->putSlot($slot);
// 		}
// 		if($this->windowid === self::SPECIAL_INVENTORY and count($this->hotbar) > 0){
// 			$this->putShort(count($this->hotbar));
// 			foreach($this->hotbar as $slot){
// 				$this->putInt($slot);
// 			}
// 		}else{
// 			$this->putShort(0);
// 		}
// 	}
//

// An exception was thrown while parsing/converting PocketMine-MP protocol.
// Please read original PHP code and port it manually.
// Exception: string index out of range
//
// 	const NETWORK_ID = Info::CRAFTING_DATA_PACKET;
//
// 	const ENTRY_SHAPELESS = 0;
// 	const ENTRY_SHAPED = 1;
// 	const ENTRY_FURNACE = 2;
// 	const ENTRY_FURNACE_DATA = 3;
// 	const ENTRY_ENCHANT_LIST = 4;
//
// 	/** @var object[] */
// 	public $entries = [];
// 	public $cleanRecipes = false;
//
// 	private static function writeEntry($entry, BinaryStream $stream){
// 		if($entry instanceof ShapelessRecipe){
// 			return self::writeShapelessRecipe($entry, $stream);
// 		}elseif($entry instanceof ShapedRecipe){
// 			return self::writeShapedRecipe($entry, $stream);
// 		}elseif($entry instanceof FurnaceRecipe){
// 			return self::writeFurnaceRecipe($entry, $stream);
// 		}elseif($entry instanceof EnchantmentList){
// 			return self::writeEnchantList($entry, $stream);
// 		}
//
// 		return -1;
// 	}
//
// 	private static function writeShapelessRecipe(ShapelessRecipe $recipe, BinaryStream $stream){
// 		$stream->putInt($recipe->getIngredientCount());
// 		foreach($recipe->getIngredientList() as $item){
// 			$stream->putSlot($item);
// 		}
//
// 		$stream->putInt(1);
// 		$stream->putSlot($recipe->getResult());
//
// 		$stream->putUUID($recipe->getId());
//
// 		return CraftingDataPacket::ENTRY_SHAPELESS;
// 	}
//
// 	private static function writeShapedRecipe(ShapedRecipe $recipe, BinaryStream $stream){
// 		$stream->putInt($recipe->getWidth());
// 		$stream->putInt($recipe->getHeight());
//
// 		for($z = 0; $z < $recipe->getHeight(); ++$z){
// 			for($x = 0; $x < $recipe->getWidth(); ++$x){
// 				$stream->putSlot($recipe->getIngredient($x, $z));
// 			}
// 		}
//
// 		$stream->putInt(1);
// 		$stream->putSlot($recipe->getResult());
//
// 		$stream->putUUID($recipe->getId());
//
// 		return CraftingDataPacket::ENTRY_SHAPED;
// 	}
//
// 	private static function writeFurnaceRecipe(FurnaceRecipe $recipe, BinaryStream $stream){
// 		if($recipe->getInput()->getDamage() !== 0){ //Data recipe
// 			$stream->putInt(($recipe->getInput()->getId() << 16) | ($recipe->getInput()->getDamage()));
// 			$stream->putSlot($recipe->getResult());
//
// 			return CraftingDataPacket::ENTRY_FURNACE_DATA;
// 		}else{
// 			$stream->putInt($recipe->getInput()->getId());
// 			$stream->putSlot($recipe->getResult());
//
// 			return CraftingDataPacket::ENTRY_FURNACE;
// 		}
// 	}
//
// 	private static function writeEnchantList(EnchantmentList $list, BinaryStream $stream){
//
// 		$stream->putByte($list->getSize());
// 		for($i = 0; $i < $list->getSize(); ++$i){
// 			$entry = $list->getSlot($i);
// 			$stream->putInt($entry->getCost());
// 			$stream->putByte(count($entry->getEnchantments()));
// 			foreach($entry->getEnchantments() as $enchantment){
// 				$stream->putInt($enchantment->getId());
// 				$stream->putInt($enchantment->getLevel());
// 			}
// 			$stream->putString($entry->getRandomName());
// 		}
//
// 		return CraftingDataPacket::ENTRY_ENCHANT_LIST;
// 	}
//
// 	public function addShapelessRecipe(ShapelessRecipe $recipe){
// 		$this->entries[] = $recipe;
// 	}
//
// 	public function addShapedRecipe(ShapedRecipe $recipe){
// 		$this->entries[] = $recipe;
// 	}
//
// 	public function addFurnaceRecipe(FurnaceRecipe $recipe){
// 		$this->entries[] = $recipe;
// 	}
//
// 	public function addEnchantList(EnchantmentList $list){
// 		$this->entries[] = $list;
// 	}
//
// 	public function clean(){
// 		$this->entries = [];
// 		return parent::clean();
// 	}
//
// 	public function decode(){
//
// 	}
//
// 	public function encode(){
// 		$this->reset();
// 		$this->putInt(count($this->entries));
//
// 		$writer = new BinaryStream();
// 		foreach($this->entries as $d){
// 			$entryType = self::writeEntry($d, $writer);
// 			if($entryType >= 0){
// 				$this->putInt($entryType);
// 				$this->putInt(strlen($writer->getBuffer()));
// 				$this->put($writer->getBuffer());
// 			}else{
// 				$this->putInt(-1);
// 				$this->putInt(0);
// 			}
//
// 			$writer->reset();
// 		}
//
// 		$this->putByte($this->cleanRecipes ? 1 : 0);
// 	}
//
//

type CraftingEvent struct{}

func (i CraftingEvent) Pid() byte { return CraftingEventHead }

func (i *CraftingEvent) Read(buf *buffer.Buffer) {}

func (i CraftingEvent) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	return buf
}

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
	X uint32
	Y uint32
	Z uint32
}

func (i BlockEntityData) Pid() byte { return BlockEntityDataHead }

func (i *BlockEntityData) Read(buf *buffer.Buffer) {
	i.X = buf.ReadInt()
	i.Y = buf.ReadInt()
	i.Z = buf.ReadInt()
	// Unexpected code: Namedtag
}

func (i BlockEntityData) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteInt(i.X)
	buf.WriteInt(i.Y)
	buf.WriteInt(i.Z)
	// Unexpected code: Namedtag
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

// An exception was thrown while parsing/converting PocketMine-MP protocol.
// Please read original PHP code and port it manually.
// Exception: string index out of range
//
// 	const NETWORK_ID = Info::PLAYER_LIST_PACKET;
//
// 	const TYPE_ADD = 0;
// 	const TYPE_REMOVE = 1;
//
// 	//REMOVE: UUID, ADD: UUID, entity id, name, isSlim, skin
// 	/** @var array[] */
// 	public $entries = [];
// 	public $type;
//
// 	public function clean(){
// 		$this->entries = [];
// 		return parent::clean();
// 	}
//
// 	public function decode(){
//
// 	}
//
// 	public function encode(){
// 		$this->reset();
// 		$this->putByte($this->type);
// 		$this->putInt(count($this->entries));
// 		foreach($this->entries as $d){
// 			if($this->type === self::TYPE_ADD){
// 				$this->putUUID($d[0]);
// 				$this->putLong($d[1]);
// 				$this->putString($d[2]);
// 				$this->putString($d[3]);
// 				$this->putString($d[4]);
// 			}else{
// 				$this->putUUID($d[0]);
// 			}
// 		}
// 	}
//
//
