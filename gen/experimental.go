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
	for x := byte(0); x < 16; x++ {
		nx := noise1dx(int(rcx) + int(x) - 20)
		ny := noise1d(int(rcx) + int(x) + 20)
		for z := byte(0); z < 16; z++ {
			h1 := byte(noise1d(nx + int(rcz) + int(z)))
			//log.Println(hh)
			if h1 >= 128 {
				h1 -= 128
			}

			h2 := byte(noise1d(ny + int(rcz) + int(z)))
			//log.Println(hh)
			if h2 >= 128 {
				h2 -= 128
			}

			hh := (h1 + h2) / 2
			for y := byte(0); y < hh; y++ {
				chunk.SetBlock(x, y, z, byte(types.Dirt))
			}
			chunk.SetBlock(x, hh, z, byte(types.Grass))
			chunk.SetBiomeColor(x, z, 0, 200, 0)
		}
	}
	chunk.PopulateHeight()
	return chunk
}

func noise1d(n int) int {
	n += 35
	h := float64(0)
	for j := .25; j < 20; j *= 1.8 {
		h += wave(float64(n)+j, 0.02*j, 12/j, 6)
	}
	h = math.Abs(h)
	return int(h)
}

func noise1dx(n int) int {
	n -= 25
	h := float64(0)
	for j := .20; j < 20; j *= 1.5 {
		h += wave(float64(n)+j, 0.015*j, 15/j, 5)
	}
	h = math.Abs(h)
	return int(h)
}

func wave(n, multx, multy, add float64) float64 {
	return math.Sin(n*multx)*math.Cos(n*multx)*multy + add
}
