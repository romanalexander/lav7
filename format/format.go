// Package format provides MCPE level formats written for lav7.
package format

import "github.com/L7-MCPE/lav7/types"

// Provider is a interface for level formats.
type Provider interface {
	Init(string)                          // Level name: usually used for file directories
	Loadable(int32, int32) (string, bool) // Path: path to file, Ok: if the chunk is saved on the file
	LoadChunk(int32, int32, string) (*types.Chunk, error)
	WriteChunk(int32, int32, *types.Chunk) error
	SaveAll(map[[2]int32]*types.Chunk) error
}
