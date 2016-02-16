package level

import "github.com/L7-MCPE/lav7/types"

// Generator is an interface for MCPE map generator.
type Generator interface {
	Init()
	Gen(int32, int32) *types.Chunk
}
