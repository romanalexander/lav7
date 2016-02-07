package level

import (
	"sync"

	"github.com/L7-MCPE/lav7/util/buffer"
)

// Chunk is a interface to access fixed amount of blocks.
type Chunk interface {
	GetBlock(byte, byte, byte) byte
	SetBlock(byte, byte, byte, byte)

	GetBlockMeta(byte, byte, byte) byte
	SetBlockMeta(byte, byte, byte, byte)

	GetBlockLight(byte, byte, byte) byte
	SetBlockLight(byte, byte, byte, byte)

	GetBlockSkyLight(byte, byte, byte) byte
	SetBlockSkyLight(byte, byte, byte, byte)

	GetHeightMap(byte, byte) byte
	SetHeightMap(byte, byte, byte)

	GetBiomeID(byte, byte) byte
	SetBiomeID(byte, byte, byte)

	GetBiomeColor(byte, byte) (byte, byte, byte)
	SetBiomeColor(byte, byte, byte, byte, byte)

	PopulateHeight()

	Mutex() *sync.RWMutex

	FullChunkData() []byte
	Read(*buffer.Buffer) error
	Write() (*buffer.Buffer, error)
}
