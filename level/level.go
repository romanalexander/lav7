// Package level implements MCPE world components and defines associated interfaces.
package level

import (
	"log"

	"github.com/L7-MCPE/lav7/types"
)

const (
	DayTime     = 0
	SunsetTime  = 12000
	NightTime   = 14000
	SunriseTime = 23000
	FullTime    = 24000
)

type Level struct {
	LevelProvider
	Name string
}

func (lv *Level) Init(pv LevelProvider) {
	lv.LevelProvider = pv
}

func (lv *Level) Provider() LevelProvider {
	return lv.LevelProvider
}

// OnUseItem handles UseItemPacket and determines position to update block position.
// Note: Value of x, y, z could be changed
//
// Face direction:
//
// 0: Down  (Y-)
// 1: Up    (Y+)
// 2: North (Z-)
// 3: South (Z+)
// 4: West  (X-)
// 5: East  (X+)
func (lv *Level) OnUseItem(x, y, z *int32, face byte, item *types.Item) (canceled bool) {
	px, py, pz := *x, *y, *z
	switch face {
	case 0:
		py--
	case 1:
		py++
	case 2:
		pz--
	case 3:
		pz++
	case 4:
		px--
	case 5:
		px++
	}
	if f := lv.GetBlock(*x, *y, *z); f == 0 {
		lv.Set(*x, *y, *z, item.Block())
		*x, *y, *z = px, py, pz
	} else {
		log.Println("Block already exists:", f)
		canceled = true
	}
	return
}

func (lv Level) GetBlock(x, y, z int32) byte {
	c := lv.GetChunk(x>>4, z>>4, true)
	c.Mutex().RLock()
	defer c.Mutex().RUnlock()
	return c.GetBlock(byte(x&0xf), byte(y), byte(z&0xf))
}

func (lv *Level) SetBlock(x, y, z int32, b byte) {
	c := lv.GetChunk(x>>4, z>>4, true)
	c.Mutex().Lock()
	defer c.Mutex().Unlock()
	c.SetBlock(byte(x&0xf), byte(y), byte(z&0xf), b)
}

func (lv Level) GetBlockMeta(x, y, z int32) byte {
	c := lv.GetChunk(x>>4, z>>4, true)
	c.Mutex().RLock()
	defer c.Mutex().RUnlock()
	return c.GetBlockMeta(byte(x&0xf), byte(y), byte(z&0xf))
}

func (lv *Level) SetBlockMeta(x, y, z int32, b byte) {
	c := lv.GetChunk(x>>4, z>>4, true)
	c.Mutex().Lock()
	defer c.Mutex().Unlock()
	c.SetBlockMeta(byte(x&0xf), byte(y), byte(z&0xf), b)
}

func (lv Level) Get(x, y, z int32) *types.Block {
	c := lv.GetChunk(x>>4, z>>4, true)
	return &types.Block{
		ID:   c.GetBlock(byte(x&0xf), byte(y), byte(z&0xf)),
		Meta: c.GetBlockMeta(byte(x&0xf), byte(y), byte(z&0xf)),
	}
}

func (lv Level) Set(x, y, z int32, block *types.Block) {
	c := lv.GetChunk(x>>4, z>>4, true)
	c.SetBlock(byte(x&0xf), byte(y), byte(z&0xf), block.ID)
	c.SetBlockMeta(byte(x&0xf), byte(y), byte(z&0xf), block.Meta)
}

func (lv Level) GetTime() uint16 {
	return 0
}

func (lv Level) SetTime(t uint16) {}
