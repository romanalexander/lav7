package types

import (
	"bytes"

	"github.com/L7-MCPE/lav7/util"
	"github.com/L7-MCPE/lav7/util/buffer"
)

// ChunkDelivery is a type for passing full chunk data to players.
type ChunkDelivery struct {
	X, Z  int32
	Chunk *Chunk
}

// Chunk contains block data for each MCPE level chunks.
// Each chunk holds 16*16*128 blocks, and consumes at least 83200 bytes of memory.
//
// A zero value for Chunk is a valid value.
type Chunk struct {
	BlockData    [16 * 16 * 128]byte
	MetaData     [16 * 16 * 64]byte // Nibbles
	LightData    [16 * 16 * 64]byte // Nibbles
	SkyLightData [16 * 16 * 64]byte // Nibbles
	HeightMap    [16 * 16]byte
	BiomeData    [16 * 16 * 4]byte // Uints

	Refs    uint64
	RWMutex util.RWLocker
}

// FallbackChunk is a chunk to be returned if level provider fails to load chunk from file.
var FallbackChunk = *new(Chunk)

// CopyFrom gets everything from given chunk, and writes to the chunk instance.
// Mutex is not shared with given chunk. You don't need to RLock the copying chunk.
func (c *Chunk) CopyFrom(chunk Chunk) {
	chunk.Mutex().RLock()
	defer chunk.Mutex().RUnlock()
	copy(c.BlockData[:], chunk.BlockData[:])
	copy(c.MetaData[:], chunk.MetaData[:])
	copy(c.LightData[:], chunk.LightData[:])
	copy(c.SkyLightData[:], chunk.SkyLightData[:])
	copy(c.HeightMap[:], chunk.HeightMap[:])
	copy(c.BiomeData[:], chunk.BiomeData[:])
}

// GetBlock returns block ID at given coordinates.
func (c Chunk) GetBlock(x, y, z byte) byte {
	return c.BlockData[uint16(y)<<8|uint16(z)<<4|uint16(x)]
}

// SetBlock sets block ID at given coordinates.
func (c *Chunk) SetBlock(x, y, z, id byte) {
	c.BlockData[uint16(y)<<8|uint16(z)<<4|uint16(x)] = id
}

// GetBlockMeta returns block meta at given coordinates.
func (c Chunk) GetBlockMeta(x, y, z byte) byte {
	if x&1 == 0 {
		return c.MetaData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] & 0x0f
	}
	return c.MetaData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] >> 4
}

// SetBlockMeta sets block meta at given coordinates.
func (c *Chunk) SetBlockMeta(x, y, z, id byte) {
	b := c.MetaData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1]
	if x&1 == 0 {
		c.MetaData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] = (b & 0xf0) | (id & 0x0f)
	} else {
		c.MetaData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] = (id&0xf)<<4 | (b & 0x0f)
	}
}

// GetBlockLight returns block light level at given coordinates.
func (c Chunk) GetBlockLight(x, y, z byte) byte {
	if x&1 == 0 {
		return c.LightData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] & 0x0f
	}
	return c.LightData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] >> 4
}

// SetBlockLight sets block light level at given coordinates.
func (c *Chunk) SetBlockLight(x, y, z, id byte) {
	b := c.LightData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1]
	if x&1 == 0 {
		c.LightData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] = (b & 0xf0) | (id & 0x0f)
	} else {
		c.LightData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] = (id&0xf)<<4 | (b & 0x0f)
	}
}

// GetBlockSkyLight returns sky light level at given coordinates.
func (c Chunk) GetBlockSkyLight(x, y, z byte) byte {
	if x&1 == 0 {
		return c.SkyLightData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] & 0x0f
	}
	return c.SkyLightData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] >> 4
}

// SetBlockSkyLight sets sky light level at given coordinates.
func (c *Chunk) SetBlockSkyLight(x, y, z, id byte) {
	b := c.SkyLightData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1]
	if x&1 == 0 {
		c.SkyLightData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] = (b & 0xf0) | (id & 0x0f)
	} else {
		c.SkyLightData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] = (id&0xf)<<4 | (b & 0x0f)
	}
}

// GetHeightMap returns highest block height on given X-Z coordinates.
func (c Chunk) GetHeightMap(x, z byte) byte {
	return c.HeightMap[uint16(z)<<4|uint16(x)]
}

// SetHeightMap saves highest block height on given X-Z coordinates.
func (c *Chunk) SetHeightMap(x, z, h byte) {
	c.HeightMap[uint16(z)<<4|uint16(x)] = h
}

// GetBiomeID returns biome ID on given X-Z coordinates.
func (c Chunk) GetBiomeID(x, z byte) byte {
	return c.BiomeData[uint16(z)<<6|uint16(x)<<2]
}

// SetBiomeID sets biome ID on given X-Z coordinates.
func (c *Chunk) SetBiomeID(x, z, id byte) {
	c.BiomeData[uint16(z)<<6|uint16(x)<<2] = id
}

// GetBiomeColor returns biome color on given X-Z coordinates.
func (c Chunk) GetBiomeColor(x, z byte) (r, g, b byte) {
	rgb := c.BiomeData[uint16(z)<<6|uint16(x)<<2+1 : uint16(z)<<6|uint16(x)<<2+4]
	return rgb[0], rgb[1], rgb[2]
}

// SetBiomeColor sets biome color on given X-Z coordinates.
func (c *Chunk) SetBiomeColor(x, z, r, g, b byte) {
	offset := uint16(z)<<6 | uint16(x)<<2
	c.BiomeData[offset+1], c.BiomeData[offset+2], c.BiomeData[offset+3] = r, g, b
}

// PopulateHeight populates chunk's block height map.
func (c *Chunk) PopulateHeight() {
	for x := byte(0); x < 16; x++ {
		for z := byte(0); z < 16; z++ {
			for y := byte(127); y > 0; y-- {
				if c.GetBlock(x, y, z) != 0 {
					c.SetHeightMap(x, z, y)
					break
				}
			}
		}
	}
}

// Mutex returns chunk's RW mutex.
func (c *Chunk) Mutex() util.RWLocker {
	if c.RWMutex == nil {
		c.RWMutex = util.NewRWMutex()
	}
	return c.RWMutex
}

// FullChunkData returns full chunk payload for FullChunkDataPacket. Order is layered.
func (c Chunk) FullChunkData() []byte {
	buf := bytes.NewBuffer(append(c.BlockData[:], c.MetaData[:]...)) // Block ID, Block MetaData
	buffer.Write(buf, append(c.SkyLightData[:], c.LightData[:]...))  // SkyLight, Light
	buffer.Write(buf, append(c.HeightMap[:], c.BiomeData[:]...))     // Height Map, Biome colors
	buffer.Write(buf, []byte{0, 0, 0, 0})                            // Extra data: NBT length 0
	// No tile entity NBT fields
	return buf.Bytes()
}
