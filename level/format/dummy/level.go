package dummy

import (
	"strconv"
	"sync"

	"github.com/L7-MCPE/lav7/block"
	"github.com/L7-MCPE/lav7/level"
	"github.com/L7-MCPE/lav7/level/gen"
)

// Package dummy implements simple Level/Chunk interface and provides the example for implementing other formats.
// It will have independent saving format later, but use of it will be deprecated after another existing format port is completed.

type Level struct {
	ChunkMap  map[string]level.Chunk
	Generator gen.Generator
}

func (l *Level) Init() {
	l.ChunkMap = make(map[string]level.Chunk)
	l.Generator = new(gen.SampleGenerator)
}

func (l Level) ChunkExists(cx, cz int32) bool {
	_, ok := l.ChunkMap[strconv.Itoa(int(cx))+":"+strconv.Itoa(int(cz))]
	return ok
}

func (l Level) GetChunk(cx, cz int32, create bool) level.Chunk {
	if c, ok := l.ChunkMap[strconv.Itoa(int(cx))+":"+strconv.Itoa(int(cz))]; ok {
		return c
	}
	if create {
		bs := l.GetGenerator().Gen(cx, cz)
		c := &Chunk{
			RWMutex: new(sync.RWMutex),
		}
		c.ArrayChunk(bs)
		l.ChunkMap[strconv.Itoa(int(cx))+":"+strconv.Itoa(int(cz))] = c
		return c
	}
	return nil
}

func (l Level) SetChunk(x, z int32, c level.Chunk) {
	if _, ok := l.ChunkMap[strconv.Itoa(int(x))+":"+strconv.Itoa(int(z))]; ok {
		panic("Tried to overwrite existing chunk!")
	}
	if c.Mutex() == nil {
		panic("Nil mutex: chunk may have been uninitialized!")
	}
	l.ChunkMap[strconv.Itoa(int(x))+":"+strconv.Itoa(int(z))] = c
}

func (l Level) GetBlock(x, y, z int32) block.IBlock {
	c := l.GetChunk(x>>4, z>>4, true)
	c.Mutex().RLock()
	defer c.Mutex().RUnlock()
	b := &block.Block{
		ID:   c.GetBlock(byte(x&0xf), byte(y), byte(z&0xf)),
		Meta: c.GetBlockMeta(byte(x&0xf), byte(y), byte(z&0xf)),
	}
	return b
}

func (l *Level) SetBlock(x, y, z int32, b block.IBlock) {
	c := l.GetChunk(x>>4, z>>4, true)
	c.Mutex().Lock()
	defer c.Mutex().Unlock()
	c.SetBlock(byte(x&0x0f), byte(y), byte(z&0x0f), b.GetID())
	c.SetBlockMeta(byte(x&0x0f), byte(y), byte(z&0x0f), b.GetMeta())
}

func (l Level) GetGenerator() gen.Generator {
	return l.Generator
}

func (l Level) GetName() string {
	return "Dummy level"
}
