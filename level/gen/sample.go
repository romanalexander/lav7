package gen

import "github.com/L7-MCPE/lav7/block"

type SampleGenerator struct{}

func (s SampleGenerator) Gen(x, z int32) (bs [16 * 16 * 128]block.IBlock) {
	for x := 0; x < 16; x++ {
		for z := 0; z < 16; z++ {
			for y := 0; y < 60; y++ {
				bs[y<<8|z<<4|x] = &block.Block{
					ID: 2,
				}
			}
		}
	}
	return
}
