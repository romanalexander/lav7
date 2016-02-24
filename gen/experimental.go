package gen

import (
	"math"

	"github.com/L7-MCPE/lav7/types"
)

func init() {
	RegisterGenerator(new(ExperimentalGenerator))
}

// ExperimentalGenerator generates normal worlds.
// Be careful to use this. This is just a concept.
type ExperimentalGenerator struct {
	Seed int64
}

// Init implements gen.Generator interface.
func (eg *ExperimentalGenerator) Init() {
	eg.Seed = 10
}

// Gen implements gen.Generator interface.
func (eg *ExperimentalGenerator) Gen(cx, cz int32) *types.Chunk {
	chunk := new(types.Chunk)
	chunk.Mutex().Lock()
	defer chunk.Mutex().Unlock()

	rcx, rcz := cx<<4, cz<<4
	var stage [18][18]float64
	for x := -1; x < 17; x++ {
		n1 := noise1dx(float64(rcx) + float64(x) - 20)
		n2 := noise1d(float64(rcx) + float64(x) + 20)
		for z := -1; z < 17; z++ {
			h1 := noise1d(n1 + float64(rcz) + float64(z))
			h2 := noise1dx(n2 - float64(rcz) - float64(z))
			h3 := noise1d(-float64(x) - float64(rcx) + float64(z) + float64(rcz))
			h4 := noise1dx(-float64(z) - float64(rcz) + float64(x) + float64(rcx))

			stage[x+1][z+1] = ((h1 + h2 + h3 + h4) / 4) + 20
		}
	}
	for x := byte(0); x < 16; x++ {
		for z := byte(0); z < 16; z++ {
			hh := byte((stage[x+1][z+1]*.4+
				stage[x+1][z]*.32+
				stage[x+1][z+2]*.32+
				stage[x+2][z]*.15+
				stage[x+2][z+1]*.32+
				stage[x][z+2]*.15+
				stage[x][z]*.15+
				stage[x][z+1]*.32+
				stage[x][z+2]*.15)/1.64) & 0x7f
			for y := byte(0); y < hh; y++ {
				chunk.SetBlock(x, y, z, byte(types.Dirt))
			}
			chunk.SetBlock(x, hh, z, byte(types.Grass))
			chunk.SetHeightMap(x, z, hh)
			chunk.SetBiomeColor(x, z, 0, 160, 0)
		}
	}
	return chunk
}

func noise1d(n float64) float64 {
	n += 35
	h := float64(0)
	for j := .25; j < 3.6; j *= 1.8 {
		h += wave(float64(n)+j, 0.02*j, 11.8/j, 0)
	}
	return h
}

func noise1dx(n float64) float64 {
	n -= 25
	h := float64(0)
	for j := .20; j < 3.6; j *= 1.5 {
		h += wave(float64(n)+j, 0.015*j, 13.2/j, 0)
	}
	return h
}

func wave(n, multx, multy, add float64) float64 {
	return math.Sin(n*multx)*math.Cos(n*multx)*multy + add
}
