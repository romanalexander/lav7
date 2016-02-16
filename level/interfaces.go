package level

import "github.com/L7-MCPE/lav7/types"

// LevelProvider is a interface for managing chunks and their save formats.

type LevelProvider interface {
	Init(func(int32, int32) *types.Chunk, string)
	GetChunk(int32, int32, bool) *types.Chunk
	SetChunk(int32, int32, bool, *types.Chunk)
	ChunkExists(int32, int32) bool
	Save() error
}
