package gen

import (
	"crypto/rand"

	"github.com/L7-MCPE/lav7/types"
)

// SampleGenerator is a simple generator concept.
type SampleGenerator struct {
	Cache *types.Chunk
}

// Init implements gen.Generator interface.
func (s *SampleGenerator) Init() {
	chunk := new(types.Chunk)
	for x := byte(0); x < 16; x++ {
		for z := byte(0); z < 16; z++ {
			for y := byte(0); y < 60; y++ {
				chunk.SetBlock(x, y, z, 3)
			}
			chunk.SetBlock(x, 60, z, 2)
			// chunk.SetBiomeColor(x, z, x*16, x*z, z*16)
		}
	}
	chunk.PopulateHeight()
	s.Cache = chunk
}

// Gen implements gen.Generator interface.
func (s *SampleGenerator) Gen(x, z int32) *types.Chunk {
	chunk := new(types.Chunk)
	chunk.CopyFrom(s.Cache)

	for x := byte(0); x < 16; x++ {
		for z := byte(0); z < 16; z++ {
			var rgb [3]byte
			rand.Read(rgb[:])
			chunk.SetBiomeColor(x, z, rgb[0], rgb[1], rgb[2])
		}
	}

	return chunk
}
