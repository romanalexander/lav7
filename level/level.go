package level

import "github.com/L7-MCPE/lav7/level/gen"

const (
	DayTime     = 0
	SunsetTime  = 12000
	NightTime   = 14000
	SunriseTime = 23000
	FullTime    = 24000
)

// Level is a interface for storing block chunks, entities, tile objects, etc.
type Level interface {
	Init()
	GetChunk(int32, int32, bool) Chunk
	SetChunk(int32, int32, Chunk)
	GetBlock(int32, int32, int32) byte
	SetBlock(int32, byte, int32, byte)
	GetBlockMeta(int32, int32, int32) byte
	SetBlockMeta(int32, int32, int32, byte)
	GetGenerator() gen.Generator
	GetName() string
	Save() error
}
