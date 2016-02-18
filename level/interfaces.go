package level

import "github.com/L7-MCPE/lav7/types"

// LevelProvider is a interface for managing chunks and their save formats.

type LevelProvider interface {
	Init(string)
	Loadable(int32, int32) (string, bool)
	LoadChunk(int32, int32, string) (*types.Chunk, error)
	WriteChunk(int32, int32, *types.Chunk) error
	SaveAll(map[[2]int32]*types.Chunk) error
}
