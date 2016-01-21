package dummy

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/L7-MCPE/lav7/level"
	"github.com/L7-MCPE/lav7/level/gen"
	"github.com/L7-MCPE/lav7/util"
	"github.com/L7-MCPE/lav7/util/buffer"
)

// Package dummy implements simple Level/Chunk interface and provides the example for implementing other formats.
// It will have independent saving format later, but use of it will be deprecated after another existing format port is completed.

type Level struct {
	ChunkMap   map[string]level.Chunk
	ChunkMutex *sync.Mutex
	Generator  gen.Generator
}

func (l *Level) Init() {
	l.ChunkMutex = new(sync.Mutex)
	l.ChunkMap = make(map[string]level.Chunk)
	l.Generator = new(gen.SampleGenerator)
}

func (l Level) ChunkExists(cx, cz int32) bool {
	_, ok := l.ChunkMap[strconv.Itoa(int(cx))+"_"+strconv.Itoa(int(cz))]
	return ok
}

func (l Level) GetChunk(cx, cz int32, create bool) level.Chunk {
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
		c := new(Chunk)
		c.RWMutex = new(sync.RWMutex)
		buf := buffer.FromBytes(b)
		copy(c.blockData[:], buf.Read(16*16*128))
		copy(c.metaData[:], buf.Read(16*16*64))
		copy(c.lightData[:], buf.Read(16*16*64))
		copy(c.skyLightData[:], buf.Read(16*16*64))
		copy(c.heightMap[:], buf.Read(16*16))
		copy(c.biomeData[:], buf.Read(16*16*4))
		return c
	}
crt:
	if create {
		bs := l.GetGenerator().Gen(cx, cz)
		c := &Chunk{
			RWMutex: new(sync.RWMutex),
		}
		c.BlockChunk(bs)
		l.ChunkMutex.Lock()
		defer l.ChunkMutex.Unlock()
		l.SetChunk(cx, cz, c)
		return c
	}
	return nil
}

func (l Level) SetChunk(x, z int32, c level.Chunk) {
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

func (l Level) GetGenerator() gen.Generator {
	return l.Generator
}

func (l Level) GetName() string {
	return "Dummy"
}

func (l Level) Save() error {
	for k, c := range l.ChunkMap {
		c := c.(*Chunk)
		path, _ := filepath.Abs("levels/" + l.GetName() + "/" + k + ".raw")
		if err := os.MkdirAll(filepath.Dir(path), 0644); err != nil {
			util.Debug("Error while creating dir:", err)
		}
		util.Debug("Mkdir", filepath.Dir(path))
		if err := ioutil.WriteFile(path, c.write().Done(), 0644); err != nil {
			util.Debug("Error while saving:", err)
			continue
		}
	}
	return nil
}
