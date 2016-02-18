// Package format provides MCPE level formats written for lav7.
package format

import "github.com/L7-MCPE/lav7/types"

// Provider is a interface for level formats.
type Provider interface {
	Init(string)
	Loadable(int32, int32) (string, bool)
	LoadChunk(int32, int32, string) (*types.Chunk, error)
	WriteChunk(int32, int32, *types.Chunk) error
	SaveAll(map[[2]int32]*types.Chunk) error
}
