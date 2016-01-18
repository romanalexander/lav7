package level

import (
	"github.com/L7-MCPE/lav7/block"
	"github.com/L7-MCPE/lav7/level/gen"
)

// Level is a interface for storing block chunks, entities, tile objects, etc.
type Level interface {
	Init()
	GetChunk(int32, int32, bool) Chunk
	SetChunk(int32, int32, Chunk)
	GetBlock(int32, int32, int32) block.IBlock
	SetBlock(int32, int32, int32, block.IBlock)
	GetGenerator() gen.Generator
	GetName() string
}
