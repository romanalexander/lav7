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
	_ // 0xbe is skipped
	FullChunkDataHead
	SetDifficultyHead
	_ // 0xc1 is skipped
	SetPlayerGametypeHead
	PlayerListHead
)

var packets = map[byte]Packet{
	LoginPacket:               new(Login),
	PlayStatusPacket:          new(PlayStatus),
	DisconnectPacket:          new(Disconnect),
	BatchPacket:               new(Batch),
	TextPacket:                new(Text),
	SetTimePacket:             new(SetTime),
	StartGamePacket:           new(StartGame),
	AddPlayerPacket:           new(AddPlayer),
	RemovePlayerPacket:        new(RemovePlayer),
	AddEntityPacket:           new(AddEntity),
	RemoveEntityPacket:        new(RemoveEntity),
	AddItemEntityPacket:       new(AddItemEntity),
	TakeItemEntityPacket:      new(TakeItemEntity),
	MoveEntityPacket:          new(MoveEntity),
	MovePlayerPacket:          new(MovePlayer),
	RemoveBlockPacket:         new(RemoveBlock),
	UpdateBlockPacket:         new(UpdateBlock),
	AddPaintingPacket:         new(AddPainting),
	ExplodePacket:             new(Explode),
	LevelEventPacket:          new(LevelEvent),
	BlockEventPacket:          new(BlockEvent),
	EntityEventPacket:         new(EntityEvent),
	MobEffectPacket:           new(MobEffect),
	UpdateAttributesPacket:    new(UpdateAttributes),
	MobEquipmentPacket:        new(MobEquipment),
	MobArmorEquipmentPacket:   new(MobArmorEquipment),
	InteractPacket:            new(Interact),
	UseItemPacket:             new(UseItem),
	PlayerActionPacket:        new(PlayerAction),
	HurtArmorPacket:           new(HurtArmor),
	SetEntityDataPacket:       new(SetEntityData),
	SetEntityMotionPacket:     new(SetEntityMotion),
	SetEntityLinkPacket:       new(SetEntityLink),
	SetHealthPacket:           new(SetHealth),
	SetSpawnPositionPacket:    new(SetSpawnPosition),
	AnimatePacket:             new(Animate),
	RespawnPacket:             new(Respawn),
	DropItemPacket:            new(DropItem),
	ContainerOpenPacket:       new(ContainerOpen),
	ContainerClosePacket:      new(ContainerClose),
	ContainerSetSlotPacket:    new(ContainerSetSlot),
	ContainerSetDataPacket:    new(ContainerSetData),
	ContainerSetContentPacket: new(ContainerSetContent),
	CraftingDataPacket:        new(CraftingData),
	CraftingEventPacket:       new(CraftingEvent),
	AdventureSettingsPacket:   new(AdventureSettings),
	BlockEntityDataPacket:     new(BlockEntityData),
	FullChunkDataPacket:       new(FullChunkData),
	SetDifficultyPacket:       new(SetDifficulty),
	SetPlayerGametypePacket:   new(SetPlayerGametype),
	PlayerListPacket:          new(PlayerList),
}

type Login struct {
}

// Pid implements Packet interface.
func (i *Login) Pid() byte { return LoginHead }

// Read implements Packet interface.
func (i *Login) Read(buf *buffer.Buffer) {
}

// Write implements Packet interface.
func (i *Login) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
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

// Pid implements Packet interface.
func (i *PlayStatus) Pid() byte { return PlayStatusHead }

// Read implements Packet interface.
func (i *PlayStatus) Read(buf *buffer.Buffer) {
	i.Int = buf.ReadStatus()
}

// Write implements Packet interface.
func (i *PlayStatus) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteInt(i.Status)
	return buf
}

type Disconnect struct {
	Message string
}

// Pid implements Packet interface.
func (i *Disconnect) Pid() byte { return DisconnectHead }

// Read implements Packet interface.
func (i *Disconnect) Read(buf *buffer.Buffer) {
	i.String = buf.ReadMessage()
}

// Write implements Packet interface.
func (i *Disconnect) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteString(i.Message)
	return buf
}

type Batch struct {
}

// Pid implements Packet interface.
func (i *Batch) Pid() byte { return BatchHead }

// Read implements Packet interface.
func (i *Batch) Read(buf *buffer.Buffer) {
	i.Payload = buf.Read(buf.ReadInt())
}

// Write implements Packet interface.
func (i *Batch) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteInt(len(i.Payload))
	return buf
}

// An exception was thrown while parsing/converting PocketMine-MP protocol.
// Please read original PHP code and port it manually.
// Exception: string index out of range
//
// 	const NETWORK_ID = Info::TEXT_PACKET;
//
// 	const TYPE_RAW = 0;
// 	const TYPE_CHAT = 1;
// 	const TYPE_TRANSLATION = 2;
// 	const TYPE_POPUP = 3;
// 	const TYPE_TIP = 4;
// 	const TYPE_SYSTEM = 5;
//
// 	public $type;
// 	public $source;
// 	public $message;
// 	public $parameters = [];
//
// 	public function decode(){
// 		$this->type = $this->getByte();
// 		switch($this->type){
// 			case self::TYPE_POPUP:
// 			case self::TYPE_CHAT:
// 				$this->source = $this->getString();
// 			case self::TYPE_RAW:
// 			case self::TYPE_TIP:
// 			case self::TYPE_SYSTEM:
// 				$this->message = $this->getString();
// 				break;
//
// 			case self::TYPE_TRANSLATION:
// 				$this->message = $this->getString();
// 				$count = $this->getByte();
// 				for($i = 0; $i < $count; ++$count){
// 					$this->parameters[] = $this->getString();
// 				}
// 		}
// 	}
//
// 	public function encode(){
// 		$this->reset();
// 		$this->putByte($this->type);
// 		switch($this->type){
// 			case self::TYPE_POPUP:
// 			case self::TYPE_CHAT:
// 				$this->putString($this->source);
// 			case self::TYPE_RAW:
// 			case self::TYPE_TIP:
// 			case self::TYPE_SYSTEM:
// 				$this->putString($this->message);
// 				break;
//
// 			case self::TYPE_TRANSLATION:
// 				$this->putString($this->message);
// 				$this->putByte(count($this->parameters));
// 				foreach($this->parameters as $p){
// 					$this->putString($p);
// 				}
// 		}
// 	}
//
//

// An exception was thrown while parsing/converting PocketMine-MP protocol.
// Please read original PHP code and port it manually.
// Exception: 'int' object is not subscriptable
//
// 	const NETWORK_ID = Info::SET_TIME_PACKET;
//
// 	public $time;
// 	public $started = true;
//
// 	public function decode(){
//
// 	}
//
// 	public function encode(){
// 		$this->reset();
// 		$this->putInt((int) (($this->time / Level::TIME_FULL) * 19200));
// 		$this->putByte($this->started ? 1 : 0);
// 	}
//
//

// An exception was thrown while parsing/converting PocketMine-MP protocol.
// Please read original PHP code and port it manually.
// Exception: 'int' object is not subscriptable
//
// 	const NETWORK_ID = Info::START_GAME_PACKET;
//
// 	public $seed;
// 	public $dimension;
// 	public $generator;
// 	public $gamemode;
// 	public $eid;
// 	public $spawnX;
// 	public $spawnY;
// 	public $spawnZ;
// 	public $x;
// 	public $y;
// 	public $z;
//
// 	public function decode(){
//
// 	}
//
// 	public function encode(){
// 		$this->reset();
// 		$this->putInt($this->seed);
// 		$this->putByte($this->dimension);
// 		$this->putInt($this->generator);
// 		$this->putInt($this->gamemode);
// 		$this->putLong($this->eid);
// 		$this->putInt($this->spawnX);
// 		$this->putInt($this->spawnY);
// 		$this->putInt($this->spawnZ);
// 		$this->putFloat($this->x);
// 		$this->putFloat($this->y);
// 		$this->putFloat($this->z);
// 		$this->putByte(0);
// 	}
//
//

// An exception was thrown while parsing/converting PocketMine-MP protocol.
// Please read original PHP code and port it manually.
// Exception: 'int' object is not subscriptable
//
// 	const NETWORK_ID = Info::ADD_PLAYER_PACKET;
//
// 	public $uuid;
// 	public $username;
// 	public $eid;
// 	public $x;
// 	public $y;
// 	public $z;
// 	public $speedX;
// 	public $speedY;
// 	public $speedZ;
// 	public $pitch;
// 	public $yaw;
// 	public $item;
// 	public $metadata;
//
// 	public function decode(){
//
// 	}
//
// 	public function encode(){
// 		$this->reset();
// 		$this->putUUID($this->uuid);
// 		$this->putString($this->username);
// 		$this->putLong($this->eid);
// 		$this->putFloat($this->x);
// 		$this->putFloat($this->y);
// 		$this->putFloat($this->z);
// 		$this->putFloat($this->speedX);
// 		$this->putFloat($this->speedY);
// 		$this->putFloat($this->speedZ);
// 		$this->putFloat($this->yaw);
// 		$this->putFloat($this->yaw); //TODO headrot
// 		$this->putFloat($this->pitch);
// 		$this->putSlot($this->item);
//
// 		$meta = Binary::writeMetadata($this->metadata);
// 		$this->put($meta);
// 	}
//
//

type RemovePlayer struct {
	Eid uint64
}

// Pid implements Packet interface.
func (i *RemovePlayer) Pid() byte { return RemovePlayerHead }

// Read implements Packet interface.
func (i *RemovePlayer) Read(buf *buffer.Buffer) {
	i.Long = buf.ReadEid()
	// Unexpected code:UUID ClientId
}

// Write implements Packet interface.
func (i *RemovePlayer) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteLong(i.Eid)
	// Unexpected code:UUID ClientId
	return buf
}

// An exception was thrown while parsing/converting PocketMine-MP protocol.
// Please read original PHP code and port it manually.
// Exception: string index out of range
//
// 	const NETWORK_ID = Info::ADD_ENTITY_PACKET;
//
// 	public $eid;
// 	public $type;
// 	public $x;
// 	public $y;
// 	public $z;
// 	public $speedX;
// 	public $speedY;
// 	public $speedZ;
// 	public $yaw;
// 	public $pitch;
// 	public $metadata;
// 	public $links = [];
//
// 	public function decode(){
//
// 	}
//
// 	public function encode(){
// 		$this->reset();
// 		$this->putLong($this->eid);
// 		$this->putInt($this->type);
// 		$this->putFloat($this->x);
// 		$this->putFloat($this->y);
// 		$this->putFloat($this->z);
// 		$this->putFloat($this->speedX);
// 		$this->putFloat($this->speedY);
// 		$this->putFloat($this->speedZ);
// 		$this->putFloat($this->yaw);
// 		$this->putFloat($this->pitch);
// 		$meta = Binary::writeMetadata($this->metadata);
// 		$this->put($meta);
// 		$this->putShort(count($this->links));
// 		foreach($this->links as $link){
// 			$this->putLong($link[0]);
// 			$this->putLong($link[1]);
// 			$this->putByte($link[2]);
// 		}
// 	}
//
//

type RemoveEntity struct {
	Eid uint64
}

// Pid implements Packet interface.
func (i *RemoveEntity) Pid() byte { return RemoveEntityHead }

// Read implements Packet interface.
func (i *RemoveEntity) Read(buf *buffer.Buffer) {
	i.Long = buf.ReadEid()
}

// Write implements Packet interface.
func (i *RemoveEntity) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteLong(i.Eid)
	return buf
}

type AddItemEntity struct {
	Eid    uint64
	X      float32
	Y      float32
	Z      float32
	SpeedX float32
	SpeedY float32
	SpeedZ float32
}

// Pid implements Packet interface.
func (i *AddItemEntity) Pid() byte { return AddItemEntityHead }

// Read implements Packet interface.
func (i *AddItemEntity) Read(buf *buffer.Buffer) {
	i.Long = buf.ReadEid()
	// Unexpected code:Slot Item
	i.Float = buf.ReadX()
	i.Float = buf.ReadY()
	i.Float = buf.ReadZ()
	i.Float = buf.ReadSpeedX()
	i.Float = buf.ReadSpeedY()
	i.Float = buf.ReadSpeedZ()
}

// Write implements Packet interface.
func (i *AddItemEntity) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteLong(i.Eid)
	// Unexpected code:Slot Item
	buf.WriteFloat(i.X)
	buf.WriteFloat(i.Y)
	buf.WriteFloat(i.Z)
	buf.WriteFloat(i.SpeedX)
	buf.WriteFloat(i.SpeedY)
	buf.WriteFloat(i.SpeedZ)
	return buf
}

type TakeItemEntity struct {
	Target uint64
	Eid    uint64
}

// Pid implements Packet interface.
func (i *TakeItemEntity) Pid() byte { return TakeItemEntityHead }

// Read implements Packet interface.
func (i *TakeItemEntity) Read(buf *buffer.Buffer) {
	i.Long = buf.ReadTarget()
	i.Long = buf.ReadEid()
}

// Write implements Packet interface.
func (i *TakeItemEntity) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteLong(i.Target)
	buf.WriteLong(i.Eid)
	return buf
}

// An exception was thrown while parsing/converting PocketMine-MP protocol.
// Please read original PHP code and port it manually.
// Exception: string index out of range
//
// 	const NETWORK_ID = Info::MOVE_ENTITY_PACKET;
//
//
// 	// eid, x, y, z, yaw, pitch
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
// 			$this->putLong($d[0]); //eid
// 			$this->putFloat($d[1]); //x
// 			$this->putFloat($d[2]); //y
// 			$this->putFloat($d[3]); //z
// 			$this->putFloat($d[4]); //yaw
// 			$this->putFloat($d[5]); //headYaw
// 			$this->putFloat($d[6]); //pitch
// 		}
// 	}
//
//

const (
	ModeNormal   byte = 0
	ModeReset    byte = 1
	ModeRotation byte = 2
)

type MovePlayer struct {
	Eid      uint64
	X        float32
	Y        float32
	Z        float32
	Yaw      float32
	BodyYaw  float32
	Pitch    float32
	Mode     byte
	OnGround byte
}

// Pid implements Packet interface.
func (i *MovePlayer) Pid() byte { return MovePlayerHead }

// Read implements Packet interface.
func (i *MovePlayer) Read(buf *buffer.Buffer) {
	i.Long = buf.ReadEid()
	i.Float = buf.ReadX()
	i.Float = buf.ReadY()
	i.Float = buf.ReadZ()
	i.Float = buf.ReadYaw()
	i.Float = buf.ReadBodyYaw()
	i.Float = buf.ReadPitch()
	i.Byte = buf.ReadMode()
	i.Byte = buf.ReadOnGround()
}

// Write implements Packet interface.
func (i *MovePlayer) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteLong(i.Eid)
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
}

// Pid implements Packet interface.
func (i *RemoveBlock) Pid() byte { return RemoveBlockHead }

// Read implements Packet interface.
func (i *RemoveBlock) Read(buf *buffer.Buffer) {
}

// Write implements Packet interface.
func (i *RemoveBlock) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	return buf
}

// An exception was thrown while parsing/converting PocketMine-MP protocol.
// Please read original PHP code and port it manually.
// Exception: string index out of range
//
// 	const NETWORK_ID = Info::UPDATE_BLOCK_PACKET;
//
// 	const FLAG_NONE      = 0b0000;
// 	const FLAG_NEIGHBORS = 0b0001;
//     const FLAG_NETWORK   = 0b0010;
// 	const FLAG_NOGRAPHIC = 0b0100;
// 	const FLAG_PRIORITY  = 0b1000;
//
// 	const FLAG_ALL = (self::FLAG_NEIGHBORS | self::FLAG_NETWORK);
// 	const FLAG_ALL_PRIORITY = (self::FLAG_ALL | self::FLAG_PRIORITY);
//
// 	public $records = []; //x, z, y, blockId, blockData, flags
//
// 	public function decode(){
//
// 	}
//
// 	public function encode(){
// 		$this->reset();
// 		$this->putInt(count($this->records));
// 		foreach($this->records as $r){
// 			$this->putInt($r[0]);
// 			$this->putInt($r[1]);
// 			$this->putByte($r[2]);
// 			$this->putByte($r[3]);
// 			$this->putByte(($r[5] << 4) | $r[4]);
// 		}
// 	}
//

type AddPainting struct {
	Eid       uint64
	X         uint32
	Y         uint32
	Z         uint32
	Direction uint32
	Title     string
}

// Pid implements Packet interface.
func (i *AddPainting) Pid() byte { return AddPaintingHead }

// Read implements Packet interface.
func (i *AddPainting) Read(buf *buffer.Buffer) {
	i.Long = buf.ReadEid()
	i.Int = buf.ReadX()
	i.Int = buf.ReadY()
	i.Int = buf.ReadZ()
	i.Int = buf.ReadDirection()
	i.String = buf.ReadTitle()
}

// Write implements Packet interface.
func (i *AddPainting) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteLong(i.Eid)
	buf.WriteInt(i.X)
	buf.WriteInt(i.Y)
	buf.WriteInt(i.Z)
	buf.WriteInt(i.Direction)
	buf.WriteString(i.Title)
	return buf
}

// An exception was thrown while parsing/converting PocketMine-MP protocol.
// Please read original PHP code and port it manually.
// Exception: string index out of range
//
// 	const NETWORK_ID = Info::EXPLODE_PACKET;
//
// 	public $x;
// 	public $y;
// 	public $z;
// 	public $radius;
// 	public $records = [];
//
// 	public function clean(){
// 		$this->records = [];
// 		return parent::clean();
// 	}
//
// 	public function decode(){
//
// 	}
//
// 	public function encode(){
// 		$this->reset();
// 		$this->putFloat($this->x);
// 		$this->putFloat($this->y);
// 		$this->putFloat($this->z);
// 		$this->putFloat($this->radius);
// 		$this->putInt(count($this->records));
// 		if(count($this->records) > 0){
// 			foreach($this->records as $record){
// 				$this->putByte($record->x);
// 				$this->putByte($record->y);
// 				$this->putByte($record->z);
// 			}
// 		}
// 	}
//

const (
	EventSoundClick            byte = 1000
	EventSoundClickFail        byte = 1001
	EventSoundShoot            byte = 1002
	EventSoundDoor             byte = 1003
	EventSoundFizz             byte = 1004
	EventSoundGhast            byte = 1007
	EventSoundGhastShoot       byte = 1008
	EventSoundBlazeShoot       byte = 1009
	EventSoundDoorBump         byte = 1010
	EventSoundDoorCrash        byte = 1012
	EventSoundBatFly           byte = 1015
	EventSoundZombieInfect     byte = 1016
	EventSoundZombieHeal       byte = 1017
	EventSoundEndermanTeleport byte = 1018
	EventSoundAnvilBreak       byte = 1020
	EventSoundAnvilUse         byte = 1021
	EventSoundAnvilFall        byte = 1022
	EventParticleShoot         byte = 2000
	EventParticleDestroy       byte = 2001
	EventParticleSplash        byte = 2002
	EventParticleEyeDespawn    byte = 2003
	EventParticleSpawn         byte = 2004
	EventStartRain             byte = 3001
	EventStartThunder          byte = 3002
	EventStopRain              byte = 3003
	EventStopThunder           byte = 3004
	EventSetData               byte = 4000
	EventPlayersSleeping       byte = 9800
)

type LevelEvent struct {
	Evid uint16
	X    float32
	Y    float32
	Z    float32
	Data uint32
}

// Pid implements Packet interface.
func (i *LevelEvent) Pid() byte { return LevelEventHead }

// Read implements Packet interface.
func (i *LevelEvent) Read(buf *buffer.Buffer) {
	i.Short = buf.ReadEvid()
	i.Float = buf.ReadX()
	i.Float = buf.ReadY()
	i.Float = buf.ReadZ()
	i.Int = buf.ReadData()
}

// Write implements Packet interface.
func (i *LevelEvent) Write() *buffer.Buffer {
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

// Pid implements Packet interface.
func (i *BlockEvent) Pid() byte { return BlockEventHead }

// Read implements Packet interface.
func (i *BlockEvent) Read(buf *buffer.Buffer) {
	i.Int = buf.ReadX()
	i.Int = buf.ReadY()
	i.Int = buf.ReadZ()
	i.Int = buf.ReadCase()
	i.Int = buf.ReadCase()
}

// Write implements Packet interface.
func (i *BlockEvent) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteInt(i.X)
	buf.WriteInt(i.Y)
	buf.WriteInt(i.Z)
	buf.WriteInt(i.Case)
	buf.WriteInt(i.Case)
	return buf
}

const (
	HurtAnimation     byte = 2
	DeathAnimation    byte = 3
	TameFail          byte = 6
	TameSuccess       byte = 7
	ShakeWet          byte = 8
	UseItem           byte = 9
	EatGrassAnimation byte = 10
	FishHookBubble    byte = 11
	FishHookPosition  byte = 12
	FishHookHook      byte = 13
	FishHookTease     byte = 14
	SquidInkCloud     byte = 15
	AmbientSound      byte = 16
	Respawn           byte = 17
)

type EntityEvent struct {
	Eid   uint64
	Event byte
}

// Pid implements Packet interface.
func (i *EntityEvent) Pid() byte { return EntityEventHead }

// Read implements Packet interface.
func (i *EntityEvent) Read(buf *buffer.Buffer) {
	i.Long = buf.ReadEid()
	i.Byte = buf.ReadEvent()
}

// Write implements Packet interface.
func (i *EntityEvent) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteLong(i.Eid)
	buf.WriteByte(i.Event)
	return buf
}

const (
	EventAdd    byte = 1
	EventModify byte = 2
	EventRemove byte = 3
)

type MobEffect struct {
	Eid       uint64
	EventId   byte
	EffectId  byte
	Amplifier byte
	Particles byte
	Duration  uint32
}

// Pid implements Packet interface.
func (i *MobEffect) Pid() byte { return MobEffectHead }

// Read implements Packet interface.
func (i *MobEffect) Read(buf *buffer.Buffer) {
	i.Long = buf.ReadEid()
	i.Byte = buf.ReadEventId()
	i.Byte = buf.ReadEffectId()
	i.Byte = buf.ReadAmplifier()
	i.Byte = buf.ReadParticles()
	i.Int = buf.ReadDuration()
}

// Write implements Packet interface.
func (i *MobEffect) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteLong(i.Eid)
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
	Eid          uint64
	Slot         byte
	SelectedSlot byte
}

// Pid implements Packet interface.
func (i *MobEquipment) Pid() byte { return MobEquipmentHead }

// Read implements Packet interface.
func (i *MobEquipment) Read(buf *buffer.Buffer) {
	i.Long = buf.ReadEid()
	// Unexpected code:Slot Item
	i.Byte = buf.ReadSlot()
	i.Byte = buf.ReadSelectedSlot()
}

// Write implements Packet interface.
func (i *MobEquipment) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteLong(i.Eid)
	// Unexpected code:Slot Item
	buf.WriteByte(i.Slot)
	buf.WriteByte(i.SelectedSlot)
	return buf
}

type MobArmorEquipment struct {
	Eid uint64
}

// Pid implements Packet interface.
func (i *MobArmorEquipment) Pid() byte { return MobArmorEquipmentHead }

// Read implements Packet interface.
func (i *MobArmorEquipment) Read(buf *buffer.Buffer) {
	i.Long = buf.ReadEid()
	// Unexpected code:Slot Slots
	// Unexpected code:Slot Slots
	// Unexpected code:Slot Slots
	// Unexpected code:Slot Slots
}

// Write implements Packet interface.
func (i *MobArmorEquipment) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteLong(i.Eid)
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

// Pid implements Packet interface.
func (i *Interact) Pid() byte { return InteractHead }

// Read implements Packet interface.
func (i *Interact) Read(buf *buffer.Buffer) {
	i.Byte = buf.ReadAction()
	i.Long = buf.ReadTarget()
}

// Write implements Packet interface.
func (i *Interact) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteByte(i.Action)
	buf.WriteLong(i.Target)
	return buf
}

type UseItem struct {
}

// Pid implements Packet interface.
func (i *UseItem) Pid() byte { return UseItemHead }

// Read implements Packet interface.
func (i *UseItem) Read(buf *buffer.Buffer) {
}

// Write implements Packet interface.
func (i *UseItem) Write() *buffer.Buffer {
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
	Eid    uint64
	Action uint32
	X      uint32
	Y      uint32
	Z      uint32
	Face   uint32
}

// Pid implements Packet interface.
func (i *PlayerAction) Pid() byte { return PlayerActionHead }

// Read implements Packet interface.
func (i *PlayerAction) Read(buf *buffer.Buffer) {
	i.Long = buf.ReadEid()
	i.Int = buf.ReadAction()
	i.Int = buf.ReadX()
	i.Int = buf.ReadY()
	i.Int = buf.ReadZ()
	i.Int = buf.ReadFace()
}

// Write implements Packet interface.
func (i *PlayerAction) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteLong(i.Eid)
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

// Pid implements Packet interface.
func (i *HurtArmor) Pid() byte { return HurtArmorHead }

// Read implements Packet interface.
func (i *HurtArmor) Read(buf *buffer.Buffer) {
	i.Byte = buf.ReadHealth()
}

// Write implements Packet interface.
func (i *HurtArmor) Write() *buffer.Buffer {
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
// 	public $eid;
// 	public $metadata;
//
// 	public function decode(){
//
// 	}
//
// 	public function encode(){
// 		$this->reset();
// 		$this->putLong($this->eid);
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
// 	// eid, motX, motY, motZ
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
// 			$this->putLong($d[0]); //eid
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

// Pid implements Packet interface.
func (i *SetEntityLink) Pid() byte { return SetEntityLinkHead }

// Read implements Packet interface.
func (i *SetEntityLink) Read(buf *buffer.Buffer) {
	i.Long = buf.ReadFrom()
	i.Long = buf.ReadTo()
	i.Byte = buf.ReadType()
}

// Write implements Packet interface.
func (i *SetEntityLink) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteLong(i.From)
	buf.WriteLong(i.To)
	buf.WriteByte(i.Type)
	return buf
}

type SetHealth struct {
	Health uint32
}

// Pid implements Packet interface.
func (i *SetHealth) Pid() byte { return SetHealthHead }

// Read implements Packet interface.
func (i *SetHealth) Read(buf *buffer.Buffer) {
	i.Int = buf.ReadHealth()
}

// Write implements Packet interface.
func (i *SetHealth) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteInt(i.Health)
	return buf
}

type SetSpawnPosition struct {
	X uint32
	Y uint32
	Z uint32
}

// Pid implements Packet interface.
func (i *SetSpawnPosition) Pid() byte { return SetSpawnPositionHead }

// Read implements Packet interface.
func (i *SetSpawnPosition) Read(buf *buffer.Buffer) {
	i.Int = buf.ReadX()
	i.Int = buf.ReadY()
	i.Int = buf.ReadZ()
}

// Write implements Packet interface.
func (i *SetSpawnPosition) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteInt(i.X)
	buf.WriteInt(i.Y)
	buf.WriteInt(i.Z)
	return buf
}

type Animate struct {
	Action byte
	Eid    uint64
}

// Pid implements Packet interface.
func (i *Animate) Pid() byte { return AnimateHead }

// Read implements Packet interface.
func (i *Animate) Read(buf *buffer.Buffer) {
	i.Byte = buf.ReadAction()
	i.Long = buf.ReadEid()
}

// Write implements Packet interface.
func (i *Animate) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteByte(i.Action)
	buf.WriteLong(i.Eid)
	return buf
}

type Respawn struct {
	X float32
	Y float32
	Z float32
}

// Pid implements Packet interface.
func (i *Respawn) Pid() byte { return RespawnHead }

// Read implements Packet interface.
func (i *Respawn) Read(buf *buffer.Buffer) {
	i.Float = buf.ReadX()
	i.Float = buf.ReadY()
	i.Float = buf.ReadZ()
}

// Write implements Packet interface.
func (i *Respawn) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteFloat(i.X)
	buf.WriteFloat(i.Y)
	buf.WriteFloat(i.Z)
	return buf
}

type DropItem struct {
}

// Pid implements Packet interface.
func (i *DropItem) Pid() byte { return DropItemHead }

// Read implements Packet interface.
func (i *DropItem) Read(buf *buffer.Buffer) {
}

// Write implements Packet interface.
func (i *DropItem) Write() *buffer.Buffer {
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

// Pid implements Packet interface.
func (i *ContainerOpen) Pid() byte { return ContainerOpenHead }

// Read implements Packet interface.
func (i *ContainerOpen) Read(buf *buffer.Buffer) {
	i.Byte = buf.ReadWindowid()
	i.Byte = buf.ReadType()
	i.Short = buf.ReadSlots()
	i.Int = buf.ReadX()
	i.Int = buf.ReadY()
	i.Int = buf.ReadZ()
}

// Write implements Packet interface.
func (i *ContainerOpen) Write() *buffer.Buffer {
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

// Pid implements Packet interface.
func (i *ContainerClose) Pid() byte { return ContainerCloseHead }

// Read implements Packet interface.
func (i *ContainerClose) Read(buf *buffer.Buffer) {
	i.Byte = buf.ReadWindowid()
}

// Write implements Packet interface.
func (i *ContainerClose) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteByte(i.Windowid)
	return buf
}

type ContainerSetSlot struct {
	Windowid   byte
	Slot       uint16
	HotbarSlot uint16
}

// Pid implements Packet interface.
func (i *ContainerSetSlot) Pid() byte { return ContainerSetSlotHead }

// Read implements Packet interface.
func (i *ContainerSetSlot) Read(buf *buffer.Buffer) {
	i.Byte = buf.ReadWindowid()
	i.Short = buf.ReadSlot()
	i.Short = buf.ReadHotbarSlot()
	// Unexpected code:Slot Item
}

// Write implements Packet interface.
func (i *ContainerSetSlot) Write() *buffer.Buffer {
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

// Pid implements Packet interface.
func (i *ContainerSetData) Pid() byte { return ContainerSetDataHead }

// Read implements Packet interface.
func (i *ContainerSetData) Read(buf *buffer.Buffer) {
	i.Byte = buf.ReadWindowid()
	i.Short = buf.ReadProperty()
	i.Short = buf.ReadValue()
}

// Write implements Packet interface.
func (i *ContainerSetData) Write() *buffer.Buffer {
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

type CraftingEvent struct {
}

// Pid implements Packet interface.
func (i *CraftingEvent) Pid() byte { return CraftingEventHead }

// Read implements Packet interface.
func (i *CraftingEvent) Read(buf *buffer.Buffer) {
}

// Write implements Packet interface.
func (i *CraftingEvent) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	return buf
}

type AdventureSettings struct {
	Flags uint32
}

// Pid implements Packet interface.
func (i *AdventureSettings) Pid() byte { return AdventureSettingsHead }

// Read implements Packet interface.
func (i *AdventureSettings) Read(buf *buffer.Buffer) {
	i.Int = buf.ReadFlags()
}

// Write implements Packet interface.
func (i *AdventureSettings) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteInt(i.Flags)
	return buf
}

type BlockEntityData struct {
	X uint32
	Y uint32
	Z uint32
}

// Pid implements Packet interface.
func (i *BlockEntityData) Pid() byte { return BlockEntityDataHead }

// Read implements Packet interface.
func (i *BlockEntityData) Read(buf *buffer.Buffer) {
	i.Int = buf.ReadX()
	i.Int = buf.ReadY()
	i.Int = buf.ReadZ()
	// Unexpected code: Namedtag
}

// Write implements Packet interface.
func (i *BlockEntityData) Write() *buffer.Buffer {
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
	ChunkX uint32
	ChunkZ uint32
	Order  byte
}

// Pid implements Packet interface.
func (i *FullChunkData) Pid() byte { return FullChunkDataHead }

// Read implements Packet interface.
func (i *FullChunkData) Read(buf *buffer.Buffer) {
	i.Int = buf.ReadChunkX()
	i.Int = buf.ReadChunkZ()
	i.Byte = buf.ReadOrder()
	i.Data = buf.Read(buf.ReadInt())
}

// Write implements Packet interface.
func (i *FullChunkData) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteInt(i.ChunkX)
	buf.WriteInt(i.ChunkZ)
	buf.WriteByte(i.Order)
	buf.WriteInt(len(i.Data))
	return buf
}

type SetDifficulty struct {
	Difficulty uint32
}

// Pid implements Packet interface.
func (i *SetDifficulty) Pid() byte { return SetDifficultyHead }

// Read implements Packet interface.
func (i *SetDifficulty) Read(buf *buffer.Buffer) {
	i.Int = buf.ReadDifficulty()
}

// Write implements Packet interface.
func (i *SetDifficulty) Write() *buffer.Buffer {
	buf := new(buffer.Buffer)
	buf.WriteInt(i.Difficulty)
	return buf
}
