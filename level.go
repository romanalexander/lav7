package lav7

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/L7-MCPE/lav7/format"
	"github.com/L7-MCPE/lav7/proto"
	"github.com/L7-MCPE/lav7/types"
	"github.com/L7-MCPE/lav7/util"
	"github.com/L7-MCPE/lav7/util/vector"
)

const tickDuration = time.Millisecond * 50

var numWorkers = runtime.NumCPU()

type genRequest struct {
	done   chan<- struct{}
	cx, cz int32
}

// Level is a struct for processing MCPE worlds.
type Level struct {
	format.Provider
	Name string

	ChunkMap   map[[2]int32]*types.Chunk
	ChunkMutex util.Locker
	genTask    chan genRequest
	Gen        func(int32, int32) *types.Chunk

	CleanQueue map[[2]int32]struct{}

	Ticker *time.Ticker
	Stop   chan struct{}
}

// Init initializes the level.
func (lv *Level) Init(pv format.Provider) {
	lv.Provider = pv
	lv.ChunkMap = make(map[[2]int32]*types.Chunk)
	lv.ChunkMutex = util.NewMutex()
	lv.Ticker = time.NewTicker(tickDuration)
	lv.Stop = make(chan struct{}, 1)
	lv.genTask = make(chan genRequest, 512)
	lv.CleanQueue = make(map[[2]int32]struct{})
	pv.Init(lv.Name)
	log.Printf("* level: generating %d workers for chunk gen", numWorkers)
	for i := 0; i < numWorkers; i++ {
		go lv.genWorker()
	}
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

func (lv *Level) genWorker() {
	for task := range lv.genTask {
		c := lv.Gen(task.cx, task.cz)
		lv.ChunkMutex.Lock()
		if _, ok := lv.ChunkMap[[2]int32{task.cx, task.cz}]; ok {
			lv.ChunkMutex.Unlock()
			task.done <- struct{}{}
			continue
		}
		lv.SetChunk(task.cx, task.cz, c)
		lv.ChunkMutex.Unlock()
		task.done <- struct{}{}
	}
}

// OnUseItem handles UseItemPacket and determines position to update block position.
func (lv *Level) OnUseItem(p *Player, x, y, z int32, face byte, item *types.Item) {
	if !item.IsBlock() {
		return
	}
	switch face {
	case vector.SideDown:
		y--
	case vector.SideUp:
		y++
	case vector.SideNorth:
		z--
	case vector.SideSouth:
		z++
	case vector.SideWest:
		x--
	case vector.SideEast:
		x++
	case 255:
		return
	}
	if y > 127 {
		return
	}
	if f := lv.GetBlock(x, y, z); f == 0 {
		lv.Set(x, y, z, item.Block())
		records := []proto.BlockRecord{
			{
				X: uint32(x),
				Y: byte(y),
				Z: uint32(z),
				Block: types.Block{
					ID:   byte(item.ID),
					Meta: byte(0),
				},
				Flags: proto.UpdateAllPriority,
			},
		}
		lv.updateSides(x, y, z, &records)
		BroadcastPacket(&proto.UpdateBlock{
			BlockRecords: records,
		})
		p.SendMessage(fmt.Sprintf("Face: %d", face))
	} else {
		p.SendMessage(fmt.Sprintf("Block %d(%s) already exists on x:%d, y:%d, z: %d", f, types.ID(f), x, y, z))
		p.SendPacket(&proto.UpdateBlock{
			BlockRecords: []proto.BlockRecord{
				{
					X:     uint32(x),
					Y:     byte(y),
					Z:     uint32(z),
					Block: lv.Get(x, y, z),
					Flags: proto.UpdateAllPriority,
				},
			},
		})
	}
}

func (lv *Level) updateSides(x, y, z int32, record *[]proto.BlockRecord) {
	// TODO
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
	var err error
	if c, ok := lv.ChunkMap[[2]int32{cx, cz}]; ok {
		return c
	} else if path, ok := lv.Loadable(cx, cz); ok {
		if path == "" {
			goto fallback
		}
		var c *types.Chunk
		c, err = lv.LoadChunk(cx, cz, path)
		if err != nil {
			goto fallback
		}
		lv.SetChunk(cx, cz, c)
		return c
	}
	return nil
fallback:
	if err != nil {
		log.Println("Error while loading chunk:", err)
	} else {
		log.Println("An error occurred while loading chunk.")
	}
	log.Println("Using empty chunk anyway.")
	c := new(types.Chunk)
	*c = types.FallbackChunk
	lv.SetChunk(cx, cz, c)
	return c
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
	done := make(chan struct{}, 1)
	go func(done chan<- struct{}) {
		lv.genTask <- genRequest{
			cx:   cx,
			cz:   cz,
			done: done,
		}
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

// Clean unloads all 'unused' chunks from memory.
func (lv *Level) Clean() (cnt int) {
	lv.ChunkMutex.Lock()
	defer lv.ChunkMutex.Unlock()
	cnt = len(lv.CleanQueue)
	for k := range lv.CleanQueue {
		lv.UnloadChunk(k[0], k[1], true)
	}
	return
}

// Save saves all loaded chunks on memory.
func (lv *Level) Save() {
	lv.ChunkMutex.Lock()
	defer lv.ChunkMutex.Unlock()
	if err := lv.SaveAll(lv.ChunkMap); err != nil {
		log.Println("Error while saving level:", err)
	}
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
	c.Mutex().Lock()
	defer c.Mutex().Unlock()
	return types.Block{
		ID:   c.GetBlock(byte(x&0xf), byte(y), byte(z&0xf)),
		Meta: c.GetBlockMeta(byte(x&0xf), byte(y), byte(z&0xf)),
	}
}

// Set sets block to given types.Block struct on given coordinates.
func (lv Level) Set(x, y, z int32, block types.Block) {
	c := lv.GetChunk(x>>4, z>>4)
	c.Mutex().Lock()
	defer c.Mutex().Unlock()
	c.SetBlock(byte(x&0xf), byte(y), byte(z&0xf), block.ID)
	c.SetBlockMeta(byte(x&0xf), byte(y), byte(z&0xf), block.Meta)
}
