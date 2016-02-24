package gen

import (
	"math"
	"math/rand"

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

	rd := rand.New(rand.NewSource(eg.Seed)).Float64() * 256
	rcx, rcz := cx<<4, cz<<4
	var stage [18][18]float64
	for x := -1; x < 17; x++ {
		rx := rd + float64(rcx) + float64(x)
		n1 := noise1dx(rx)
		n2 := noise1d(rx)
		for z := -1; z < 17; z++ {
			rz := rd + float64(rcz) + float64(z)
			stage[x+1][z+1] = noise16(n1, n2, rx, rz)
		}
	}
	for x := byte(0); x < 16; x++ {
		for z := byte(0); z < 16; z++ {
			hh := stabilized(x, z, &stage) + 32
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

func stabilized(x, z byte, hmap *[18][18]float64) byte {
	h := byte((hmap[x+1][z+1]*.4+
		hmap[x+1][z]*.32+
		hmap[x+1][z+2]*.32+
		hmap[x+2][z]*.15+
		hmap[x+2][z+1]*.32+
		hmap[x][z+2]*.15+
		hmap[x][z]*.15+
		hmap[x][z+1]*.32+
		hmap[x][z+2]*.15)/1.64) & 0x7f
	if h > 0x3f {
		return 0x7f - h
	}
	return h
}

func noise16(n1, n2, rx, rz float64) float64 {
	h1 := noise1d(n1 + rx)
	h2 := noise1dx(n2 - rz)
	h3 := noise1d(n1 + rx - rz)
	h4 := noise1dx(n2 + rz - rx)
	h5 := noise1d(n1 + rx + 2*rz)
	h6 := noise1dx(n2 + rx - 2*rz)
	h7 := noise1d(n1 + 2*rx + rz)
	h8 := noise1dx(n2 + 2*rz - rx)
	return math.Abs(((h1+h2+h3+h4+h5+h6+h7+h8)/8)-35) + 20
}

func noise8(n1, n2, rx, rz float64) float64 {
	h1 := noise1d(n1 + rx)
	h2 := noise1dx(n2 - rz)
	h3 := noise1d(n1 + rx - rz)
	h4 := noise1dx(n2 + rz - rx)
	return math.Abs(((h1+h2+h3+h4)/4)-50) + 15
}

func noise1d(n float64) float64 {
	n += 35
	h := float64(0)
	for j, i := .25, 0; i < 4; j *= 1.8 {
		i++
		h += wave(float64(n)+j, 0.02*j, 12/j, 0)
	}
	return h
}

func noise1dx(n float64) float64 {
	n -= 25
	h := float64(0)
	for j, i := .20, 0; i < 4; j *= 1.5 {
		i++
		h += wave(float64(n)+j, 0.015*j, 16/j, 0)
	}
	return h
}

func wave(n, multx, multy, add float64) float64 {
	return math.Sin(n*multx)*math.Cos(n*multx)*multy + add
}
