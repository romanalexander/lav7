// Package level implements MCPE world components and defines associated interfaces.
package level

import (
	"fmt"
	"log"
	"sync"

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

	ChunkMap   map[[2]int32]*types.Chunk
	ChunkMutex *sync.Mutex
	Gen        func(int32, int32) *types.Chunk
}

// Init initializes the level.
func (lv *Level) Init(pv LevelProvider) {
	lv.LevelProvider = pv
	lv.ChunkMap = make(map[[2]int32]*types.Chunk)
	lv.ChunkMutex = new(sync.Mutex)
	pv.Init(lv.Name)
}

// OnUseItem handles UseItemPacket and determines position to update block position.
// Note: Value of x, y, z could be changed
//
// Face direction:
//
// `0: Down  (Y-)
// 1: Up    (Y+)
// 2: North (Z-)
// 3: South (Z+)
// 4: West  (X-)
// 5: East  (X+)`
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

func (lv *Level) ChunkExists(cx, cz int32) bool {
	lv.ChunkMutex.Lock()
	_, ok := lv.ChunkMap[[2]int32{cx, cz}]
	lv.ChunkMutex.Unlock()
	return ok
}

func (lv *Level) GetChunk(cx, cz int32, create bool) *types.Chunk {
	lv.ChunkMutex.Lock()
	defer lv.ChunkMutex.Unlock()
	if c, ok := lv.ChunkMap[[2]int32{cx, cz}]; ok {
		return c
	} else if path, ok := lv.Loadable(cx, cz); ok {
		c, err := lv.LoadChunk(cx, cz, path)
		if err != nil {
			log.Println("Error while loading chunk:", err)
			log.Println("Using empty chunk anyway.")
			c = new(types.Chunk)
			*c = types.FallbackChunk
		}
		lv.SetChunk(cx, cz, c)
		return c
	} else {
		c := lv.Gen(cx, cz)
		lv.SetChunk(cx, cz, c)
		return c
	}
}

func (lv *Level) SetChunk(cx, cz int32, c *types.Chunk) { // Should lock ChunkMutex before call
	// lv.ChunkMutex.Lock()
	// defer lv.ChunkMutex.Unlock()
	if _, ok := lv.ChunkMap[[2]int32{cx, cz}]; ok {
		panic("Tried to overwrite existing chunk!")
	}
	lv.ChunkMap[[2]int32{cx, cz}] = c
}

func (lv *Level) UnloadChunk(cx, cz int32, save bool) error { // Should lock ChunkMutex before call
	if c, ok := lv.ChunkMap[[2]int32{cx, cz}]; ok {
		delete(lv.ChunkMap, [2]int32{cx, cz})
		if save {
			return lv.WriteChunk(cx, cz, c)
		}
		return nil
	}
	return fmt.Errorf("Chunk %d:%d is not loaded", cx, cz)
}

func (lv *Level) Save() {
	lv.ChunkMutex.Lock()
	defer lv.ChunkMutex.Unlock()
	lv.SaveAll(lv.ChunkMap)
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
