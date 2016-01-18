package gen

import "github.com/L7-MCPE/lav7/block"

// Generator is an interface for MCPE map generator.
type Generator interface {
	Gen(int32, int32) [16 * 16 * 128]block.IBlock
}
