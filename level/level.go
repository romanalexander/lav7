// Package level implements MCPE world components and defines associated interfaces.
package level

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/L7-MCPE/lav7/types"
	"github.com/L7-MCPE/lav7/util/buffer"
)

const (
	DayTime     = 0
	SunsetTime  = 12000
	NightTime   = 14000
	SunriseTime = 23000
	FullTime    = 24000
)

type Level struct {
	ChunkMap     map[string]Chunk
	ChunkMutex   *sync.Mutex
	Generator    Generator
	ChunkCreator func() Chunk
}

func (l *Level) Init(gen Generator, cc func() Chunk) {
	l.ChunkMutex = new(sync.Mutex)
	l.ChunkMap = make(map[string]Chunk)
	l.Generator = gen
	l.ChunkCreator = cc
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
func (l *Level) OnUseItem(x, y, z *int32, face byte, item *types.Item) (canceled bool) {
	log.Println("OnTouch:", x, y, z, face)
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
	}
	if f := l.GetBlock(*x, *y, *z); f == 0 {
		log.Printf("SetBlock %d %d %d %v", *x, *y, *z, item.Block())
		l.Set(*x, *y, *z, item.Block())
	} else {
		log.Println("Block", f)
	}
	return
}

func (l Level) ChunkExists(cx, cz int32) bool {
	_, ok := l.ChunkMap[strconv.Itoa(int(cx))+"_"+strconv.Itoa(int(cz))]
	return ok
}

func (l Level) GetChunk(cx, cz int32, create bool) Chunk {
	l.ChunkMutex.Lock()
	defer l.ChunkMutex.Unlock()
	cname := strconv.Itoa(int(cx)) + "_" + strconv.Itoa(int(cz))
	if c, ok := l.ChunkMap[cname]; ok {
		return c
	}
	if path, err := filepath.Abs("levels/" + l.GetName() + "/" + cname + ".raw"); err != nil {
		goto crt
	} else if _, err := os.Stat(path); os.IsNotExist(err) {
		goto crt
	} else {
		f, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
		if err != nil {
			goto crt
		}
		b := make([]byte, 83200)
		_, err = f.Read(b)
		if err != nil {
			goto crt
		}
		c := l.ChunkCreator()
		buf := buffer.FromBytes(b)
		c.Read(buf)
		l.ChunkMap[cname] = c
		return c
	}
crt:
	if create {
		c := l.ChunkCreator()
		if err := l.Generator.Gen(cx, cz, c); err != nil {
			log.Println("Error while generating chunk:", err)
			c = l.ChunkCreator()
		}
		l.SetChunk(cx, cz, c)
		l.ChunkMap[cname] = c
		return c
	}
	return nil
}

func (l Level) SetChunk(x, z int32, c Chunk) {
	if _, ok := l.ChunkMap[strconv.Itoa(int(x))+"_"+strconv.Itoa(int(z))]; ok {
		panic("Tried to overwrite existing chunk!")
	}
	if c.Mutex() == nil {
		panic("Nil mutex: chunk may have been uninitialized!")
	}
	l.ChunkMap[strconv.Itoa(int(x))+"_"+strconv.Itoa(int(z))] = c
}

func (l Level) GetBlock(x, y, z int32) byte {
	c := l.GetChunk(x>>4, z>>4, true)
	c.Mutex().RLock()
	defer c.Mutex().RUnlock()
	return c.GetBlock(byte(x&0xf), byte(y), byte(z&0xf))
}

func (l *Level) SetBlock(x, y, z int32, b byte) {
	c := l.GetChunk(x>>4, z>>4, true)
	c.Mutex().Lock()
	defer c.Mutex().Unlock()
	c.SetBlock(byte(x&0x0f), byte(y), byte(z&0x0f), b)
}

func (l Level) GetBlockMeta(x, y, z int32) byte {
	c := l.GetChunk(x>>4, z>>4, true)
	c.Mutex().RLock()
	defer c.Mutex().RUnlock()
	return c.GetBlockMeta(byte(x&0xf), byte(y), byte(z&0xf))
}

func (l *Level) SetBlockMeta(x, y, z int32, b byte) {
	c := l.GetChunk(x>>4, z>>4, true)
	c.Mutex().Lock()
	defer c.Mutex().Unlock()
	c.SetBlockMeta(byte(x&0x0f), byte(y), byte(z&0x0f), b)
}

func (l Level) Get(x, y, z int32) *types.Block {
	return &types.Block{
		ID:   l.GetBlock(x, y, z),
		Meta: l.GetBlockMeta(x, y, z),
	}
}

func (l Level) Set(x, y, z int32, block *types.Block) {
	l.SetBlock(x, y, z, block.ID)
	l.SetBlockMeta(x, y, z, block.Meta)
}

func (l Level) GetTime() uint16 {
	return 0
}

func (l Level) SetTime(t uint16) {}

func (l Level) GetName() string {
	return "Dummy"
}

func (l Level) Save() error {
	for k, c := range l.ChunkMap {
		path, _ := filepath.Abs("levels/" + l.GetName() + "/" + k + ".raw")
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			log.Println("Error while creating dir:", err)
		}
		buf, err := c.Write()
		if err != nil {
			log.Println("Error while writing chunk to file:", err)
			continue
		}
		if err := ioutil.WriteFile(path, buf.Done(), 0644); err != nil {
			log.Println("Error while saving:", err)
			continue
		}
	}
	return nil
}
