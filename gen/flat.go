package gen

import "github.com/L7-MCPE/lav7/types"

func init() {
	RegisterGenerator(new(FlatGenerator))
}

// FlatGenerator generates flat MCPE chunks.
type FlatGenerator struct {
	Cache types.Chunk
}

// Init implements gen.Generator interface.
func (s *FlatGenerator) Init() {
	chunk := new(types.Chunk)
	for x := byte(0); x < 16; x++ {
		for z := byte(0); z < 16; z++ {
			chunk.SetBlock(x, 0, z, byte(types.Bedrock))
			for y := byte(1); y < 5; y++ {
				chunk.SetBlock(x, y, z, byte(types.Dirt))
			}
			chunk.SetBlock(x, 5, z, byte(types.Grass))
			chunk.SetBiomeColor(x, z, 0, 200, 0)
		}
	}
	chunk.PopulateHeight()
	s.Cache = *chunk
}

// Gen implements gen.Generator interface.
func (s *FlatGenerator) Gen(x, z int32) *types.Chunk {
	chunk := new(types.Chunk)
	chunk.CopyFrom(s.Cache)

	return chunk
}
