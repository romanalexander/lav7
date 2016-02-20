package gen

import (
	"github.com/L7-MCPE/lav7/types"
	"image"
)

type ImageGenerator struct {
	Cache         *types.Chunk
	Image         image.Image
	Width, Height int32
}

func (s *ImageGenerator) Init() {
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

func (s *ImageGenerator) Gen(x, z int32) *types.Chunk {
	chunk := new(types.Chunk)
	chunk.CopyFrom(s.Cache)
	blockX, blockZ := x<<4, z<<4
	imgStartX := int32(blockX/s.Width) * s.Width
	imgStartZ := int32(blockZ/s.Height) * s.Height

	if x < 0 {
		imgStartX -= s.Width
	}

	if z < 0 {
		imgStartZ -= s.Height
	}

	for x := byte(0); x < 16; x++ {
		for z := byte(0); z < 16; z++ {
			ix, iz := s.getImageXZ(blockX+int32(x), blockZ+int32(z), imgStartX, imgStartZ)
			rgb := s.Image.At(int(ix), int(iz))
			r, g, b, _ := rgb.RGBA()
			chunk.SetBiomeColor(x, z, byte(r>>8), byte(g>>8), byte(b>>8))
		}
	}

	return chunk
}

func (s *ImageGenerator) getImageXZ(bx, bz, isx, isz int32) (int32, int32) {
	diffX, diffZ := bx-isx, bz-isz
	diffX %= s.Width
	diffZ %= s.Height
	return diffX, diffZ
}
