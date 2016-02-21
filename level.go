package lav7

import (
	"fmt"
	"log"
	"time"

	"github.com/L7-MCPE/lav7/format"
	"github.com/L7-MCPE/lav7/types"
	"github.com/L7-MCPE/lav7/util"
)

const tickDuration = time.Millisecond * 50

// Level is a struct for processing MCPE worlds.
type Level struct {
	format.Provider
	Name string

	ChunkMap   map[[2]int32]*types.Chunk
	TaskMap    map[[2]int32]struct{}
	ChunkMutex util.Locker
	Gen        func(int32, int32) *types.Chunk

	Ticker *time.Ticker
	Stop   chan struct{}
}

// Init initializes the level.
func (lv *Level) Init(pv format.Provider) {
	lv.Provider = pv
	lv.ChunkMap = make(map[[2]int32]*types.Chunk)
	lv.TaskMap = make(map[[2]int32]struct{})
	lv.ChunkMutex = util.NewMutex()
	lv.Ticker = time.NewTicker(tickDuration)
	lv.Stop = make(chan struct{}, 1)
	pv.Init(lv.Name)
}

// Process receives signals from two channels, Ticker.C and Stop.
func (lv *Level) Process() {
	for {
		select {
		case <-lv.Ticker.C:
			return
			//lv.tick()
		case <-lv.Stop:
			return
		}
	}
}

func (lv *Level) tick() {

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
	switch face {
	case 0:
		*y--
	case 1:
		*y++
	case 2:
		*z--
	case 3:
		*z++
	case 4:
		*x--
	case 5:
		*x++
	case 255:
		return
	}
	if f := lv.GetBlock(*x, *y, *z); f == 0 {
		lv.Set(*x, *y, *z, item.Block())
	} else {
		log.Printf("Block %d(%s) already exists on x:%d, y:%d, z: %d; cancelling. (face %d)", f, types.ID(f), *x, *y, *z, face)
		canceled = true
	}
	return
}

// ChunkExists returns if the chunk is loaded on the given chunk coordinates.
func (lv *Level) ChunkExists(cx, cz int32) bool {
	lv.ChunkMutex.Lock()
	_, ok := lv.ChunkMap[[2]int32{cx, cz}]
	lv.ChunkMutex.Unlock()
	return ok
}

// GetChunk returns *types.Chunk from ChunkMap with given chunk coordinates.
// If the chunk is not loaded, this will try to load a chunk from Provider.
//
// If Provider fails to load the chunk, this function will return nil.
func (lv *Level) GetChunk(cx, cz int32) *types.Chunk {
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
	}
	return nil
}

// SetChunk sets given chunk to chunk map.
// Callers should lock ChunkMutex before call.
func (lv *Level) SetChunk(cx, cz int32, c *types.Chunk) {
	// lv.ChunkMutex.Lock()
	// defer lv.ChunkMutex.Unlock()
	if _, ok := lv.ChunkMap[[2]int32{cx, cz}]; ok {
		panic("Tried to overwrite existing chunk!")
	}
	lv.ChunkMap[[2]int32{cx, cz}] = c
}

// CreateChunk generates chunk asynchronously.
func (lv *Level) CreateChunk(cx, cz int32) <-chan struct{} {
	lv.ChunkMutex.Lock()
	if _, ok := lv.TaskMap[[2]int32{cx, cz}]; ok {
		lv.ChunkMutex.Unlock()
		return nil
	}
	lv.TaskMap[[2]int32{cx, cz}] = struct{}{}
	lv.ChunkMutex.Unlock()
	done := make(chan struct{}, 1)
	go func(done chan<- struct{}) {
		c := lv.Gen(cx, cz)
		lv.ChunkMutex.Lock()
		lv.SetChunk(cx, cz, c)
		lv.ChunkMutex.Unlock()
		done <- struct{}{}
	}(done)
	return done
}

// UnloadChunk unloads chunk from memory.
// If save is given true, this will save the chunk before unload.
//
// Callers should lock ChunkMutex before call.
func (lv *Level) UnloadChunk(cx, cz int32, save bool) error {
	if c, ok := lv.ChunkMap[[2]int32{cx, cz}]; ok {
		delete(lv.ChunkMap, [2]int32{cx, cz})
		if save {
			return lv.WriteChunk(cx, cz, c)
		}
		return nil
	}
	return fmt.Errorf("Chunk %d:%d is not loaded", cx, cz)
}

// Save saves all loaded chunks on memory.
func (lv *Level) Save() {
	lv.ChunkMutex.Lock()
	defer lv.ChunkMutex.Unlock()
	lv.SaveAll(lv.ChunkMap)
}

// GetBlock returns block ID on given coordinates.
func (lv Level) GetBlock(x, y, z int32) byte {
	c := lv.GetChunk(x>>4, z>>4)
	c.Mutex().RLock()
	defer c.Mutex().RUnlock()
	return c.GetBlock(byte(x&0xf), byte(y), byte(z&0xf))
}

// SetBlock sets block ID on given coordinates.
func (lv *Level) SetBlock(x, y, z int32, b byte) {
	c := lv.GetChunk(x>>4, z>>4)
	c.Mutex().Lock()
	defer c.Mutex().Unlock()
	c.SetBlock(byte(x&0xf), byte(y), byte(z&0xf), b)
}

// GetBlockMeta returns block meta on given coordinates.
func (lv Level) GetBlockMeta(x, y, z int32) byte {
	c := lv.GetChunk(x>>4, z>>4)
	c.Mutex().RLock()
	defer c.Mutex().RUnlock()
	return c.GetBlockMeta(byte(x&0xf), byte(y), byte(z&0xf))
}

// SetBlockMeta sets block meta on given coordinates.
func (lv *Level) SetBlockMeta(x, y, z int32, b byte) {
	c := lv.GetChunk(x>>4, z>>4)
	c.Mutex().Lock()
	defer c.Mutex().Unlock()
	c.SetBlockMeta(byte(x&0xf), byte(y), byte(z&0xf), b)
}

// Get returns types.Block struct on given coordinates.
// The struct will contain block ID/meta.
func (lv Level) Get(x, y, z int32) types.Block {
	c := lv.GetChunk(x>>4, z>>4)
	return types.Block{
		ID:   c.GetBlock(byte(x&0xf), byte(y), byte(z&0xf)),
		Meta: c.GetBlockMeta(byte(x&0xf), byte(y), byte(z&0xf)),
	}
}

// Set sets block to given types.Block struct on given coordinates.
func (lv Level) Set(x, y, z int32, block types.Block) {
	c := lv.GetChunk(x>>4, z>>4)
	c.SetBlock(byte(x&0xf), byte(y), byte(z&0xf), block.ID)
	c.SetBlockMeta(byte(x&0xf), byte(y), byte(z&0xf), block.Meta)
}
