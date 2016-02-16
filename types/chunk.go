package types

import "sync"

type ChunkDelivery struct {
	X, Z  int32
	Chunk *Chunk
}

// Chunk contains block data for each MCPE level chunks.
// Each chunk holds 16*16*128 blocks.
type Chunk struct {
	BlockData    [16 * 16 * 128]byte
	MetaData     [16 * 16 * 64]byte // Nibbles
	LightData    [16 * 16 * 64]byte // Nibbles
	SkyLightData [16 * 16 * 64]byte // Nibbles
	HeightMap    [16 * 16]byte
	BiomeData    [16 * 16 * 4]byte // Uints
	mutex        *sync.RWMutex
}

func (c *Chunk) Init() {
	c.mutex = new(sync.RWMutex)
}

// GetBlock implements level.Chunk interface.
func (c Chunk) GetBlock(x, y, z byte) byte {
	return c.BlockData[uint16(y)<<8|uint16(z)<<4|uint16(x)]
}

// SetBlock implements level.Chunk interface.
func (c *Chunk) SetBlock(x, y, z, id byte) {
	c.BlockData[uint16(y)<<8|uint16(z)<<4|uint16(x)] = id
}

// GetBlockMeta implements level.Chunk interface.
func (c Chunk) GetBlockMeta(x, y, z byte) byte {
	if x&1 == 0 {
		return c.MetaData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] & 0x0f
	}
	return c.MetaData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] >> 4
}

// SetBlockMeta implements level.Chunk interface.
func (c *Chunk) SetBlockMeta(x, y, z, id byte) {
	b := c.MetaData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1]
	if x&1 == 0 {
		c.MetaData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] = (b & 0xf0) | (id & 0x0f)
	} else {
		c.MetaData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] = (id&0xf)<<4 | (b & 0x0f)
	}
}

// GetBlockLight implements level.Chunk interface.
func (c Chunk) GetBlockLight(x, y, z byte) byte {
	if x&1 == 0 {
		return c.LightData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] & 0x0f
	}
	return c.LightData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] >> 4
}

// SetBlockLight implements level.Chunk interface.
func (c *Chunk) SetBlockLight(x, y, z, id byte) {
	b := c.LightData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1]
	if x&1 == 0 {
		c.LightData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] = (b & 0xf0) | (id & 0x0f)
	} else {
		c.LightData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] = (id&0xf)<<4 | (b & 0x0f)
	}
}

// GetBlockSkyLight implements level.Chunk interface.
func (c Chunk) GetBlockSkyLight(x, y, z byte) byte {
	if x&1 == 0 {
		return c.SkyLightData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] & 0x0f
	}
	return c.SkyLightData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] >> 4
}

// SetBlockSkyLight implements level.Chunk interface.
func (c *Chunk) SetBlockSkyLight(x, y, z, id byte) {
	b := c.SkyLightData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1]
	if x&1 == 0 {
		c.SkyLightData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] = (b & 0xf0) | (id & 0x0f)
	} else {
		c.SkyLightData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] = (id&0xf)<<4 | (b & 0x0f)
	}
}

// GetHeightMap implements level.Chunk interface.
func (c Chunk) GetHeightMap(x, z byte) byte {
	return c.HeightMap[uint16(z)<<4|uint16(x)]
}

// SetHeightMap implements level.Chunk interface.
func (c *Chunk) SetHeightMap(x, z, h byte) {
	c.HeightMap[uint16(z)<<4|uint16(x)] = h
}

// GetBiomeID implements level.Chunk interface.
func (c Chunk) GetBiomeID(x, z byte) byte {
	return c.BiomeData[uint16(z)<<6|uint16(x)<<2]
}

// SetBiomeID implements level.Chunk interface.
func (c *Chunk) SetBiomeID(x, z, id byte) {
	c.BiomeData[uint16(z)<<6|uint16(x)<<2] = id
}

// GetBiomeColor implements level.Chunk interface.
func (c Chunk) GetBiomeColor(x, z byte) (r, g, b byte) {
	rgb := c.BiomeData[uint16(z)<<6|uint16(x)<<2+1 : uint16(z)<<6|uint16(x)<<2+4]
	return rgb[0], rgb[1], rgb[2]
}

// SetBiomeColor implements level.Chunk interface.
func (c *Chunk) SetBiomeColor(x, z, r, g, b byte) {
	offset := uint16(z)<<6 | uint16(x)<<2
	c.BiomeData[offset+1], c.BiomeData[offset+2], c.BiomeData[offset+3] = r, g, b
}

// PopulateHeight implements level.Chunk interface.
func (c *Chunk) PopulateHeight() {
	for x := byte(0); x < 16; x++ {
		for z := byte(0); z < 16; z++ {
			for y := byte(127); y >= 0; y-- {
				if c.GetBlock(x, y, z) != 0 {
					c.SetHeightMap(x, z, y)
					break
				}
			}
		}
	}
}

// Mutex implements level.Chunk interface.
func (c *Chunk) Mutex() *sync.RWMutex {
	return c.mutex
}

// FullChunkData returns full chunk payload for FullChunkDataPacket.
func (c Chunk) FullChunkData() []byte {
	a := append(c.BlockData[:], c.MetaData[:]...)     // Block ID, Block MetaData
	b := append(c.SkyLightData[:], c.LightData[:]...) // SkyLight, Light
	c_ := append(c.HeightMap[:], c.BiomeData[:]...)   // Height Map, Biome colors
	d := []byte{0, 0, 0, 0}                           // Extra data: length 0
	// No tile entity NBT fields
	return append(a, append(b, append(c_, d...)...)...) // Seems dirty :\
}
