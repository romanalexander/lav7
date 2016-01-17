package dummy

var DummyChunk = new(Chunk)

func InitDummyChunk() {
	for x := byte(0); x < 16; x++ {
		for z := byte(0); z < 16; z++ {
			DummyChunk.SetHeightMap(x, z, 60)
			DummyChunk.SetBiomeColor(x, z, x, 128, z)
			for y := byte(0); y <= 60; y++ {
				DummyChunk.SetBlock(x, y, z, 3)
			}
			for y := byte(60); y < 128; y++ {
				DummyChunk.SetBlockSkyLight(x, y, z, 15)
				DummyChunk.SetBlockLight(x, y, z, 15)
			}
			DummyChunk.SetBlock(x, 60, z, 2)
		}
	}
}

type Chunk struct {
	blockData    [16 * 16 * 128]byte
	metaData     [16 * 16 * 64]byte // Nibbles
	lightData    [16 * 16 * 64]byte // Nibbles
	skyLightData [16 * 16 * 64]byte // Nibbles
	heightMap    [16 * 16]byte
	biomeData    [16 * 16 * 4]byte // Uints
}

// GetBlock implements level.Chunk interface.
func (c Chunk) GetBlock(x, y, z byte) byte {
	return c.blockData[uint16(y)<<8|uint16(z)<<4|uint16(x)]
}

// SetBlock implements level.Chunk interface.
func (c *Chunk) SetBlock(x, y, z, id byte) {
	c.blockData[uint16(y)<<8|uint16(z)<<4|uint16(x)] = id
}

// GetBlockMeta implements level.Chunk interface.
func (c Chunk) GetBlockMeta(x, y, z byte) byte {
	if x&1 == 0 {
		return c.metaData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] & 0x0f
	}
	return c.metaData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] >> 4
}

// SetBlockMeta implements level.Chunk interface.
func (c *Chunk) SetBlockMeta(x, y, z, id byte) {
	b := c.metaData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1]
	if x&1 == 0 {
		c.metaData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] = (b & 0xf0) | (id & 0x0f)
	} else {
		c.metaData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] = (id&0xf)<<4 | (b & 0x0f)
	}
}

// GetBlockLight implements level.Chunk interface.
func (c Chunk) GetBlockLight(x, y, z byte) byte {
	if x&1 == 0 {
		return c.lightData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] & 0x0f
	}
	return c.lightData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] >> 4
}

// SetBlockLight implements level.Chunk interface.
func (c *Chunk) SetBlockLight(x, y, z, id byte) {
	b := c.lightData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1]
	if x&1 == 0 {
		c.lightData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] = (b & 0xf0) | (id & 0x0f)
	} else {
		c.lightData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] = (id&0xf)<<4 | (b & 0x0f)
	}
}

// GetBlockSkyLight implements level.Chunk interface.
func (c Chunk) GetBlockSkyLight(x, y, z byte) byte {
	if x&1 == 0 {
		return c.skyLightData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] & 0x0f
	}
	return c.skyLightData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] >> 4
}

// SetBlockSkyLight implements level.Chunk interface.
func (c *Chunk) SetBlockSkyLight(x, y, z, id byte) {
	b := c.skyLightData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1]
	if x&1 == 0 {
		c.skyLightData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] = (b & 0xf0) | (id & 0x0f)
	} else {
		c.skyLightData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] = (id&0xf)<<4 | (b & 0x0f)
	}
}

// GetHeightMap implements level.Chunk interface.
func (c Chunk) GetHeightMap(x, z byte) byte {
	return c.heightMap[uint16(z)<<4|uint16(x)]
}

// SetHeightMap implements level.Chunk interface.
func (c *Chunk) SetHeightMap(x, z, h byte) {
	c.heightMap[uint16(z)<<4|uint16(x)] = h
}

// GetBiomeID implements level.Chunk interface.
func (c Chunk) GetBiomeID(x, z byte) byte {
	return c.biomeData[uint16(z)<<6|uint16(x)<<2]
}

// SetBiomeID implements level.Chunk interface.
func (c *Chunk) SetBiomeID(x, z, id byte) {
	c.biomeData[uint16(z)<<6|uint16(x)<<2] = id
}

// GetBiomeColor implements level.Chunk interface.
func (c Chunk) GetBiomeColor(x, z byte) (r, g, b byte) {
	rgb := c.biomeData[uint16(z)<<6|uint16(x)<<2+1 : uint16(z)<<6|uint16(x)<<2+4]
	return rgb[0], rgb[1], rgb[2]
}

// SetBiomeColor implements level.Chunk interface.
func (c *Chunk) SetBiomeColor(x, z, r, g, b byte) {
	offset := uint16(z)<<6 | uint16(x)<<2
	c.biomeData[offset+1], c.biomeData[offset+2], c.biomeData[offset+3] = r, g, b
}

// FullChunkData returns full chunk payload for FullChunkDataPacket.
func (c Chunk) FullChunkData() []byte {
	a := append(c.blockData[:], c.metaData[:]...)     // Block ID, Block Metadata
	b := append(c.skyLightData[:], c.lightData[:]...) // SkyLight, Light
	c_ := append(c.heightMap[:], c.biomeData[:]...)   // Height Map, Biome colors
	d := []byte{0, 0, 0, 0}                           // Extra data: length 0
	// No tile entity NBT fields
	e := append(a, append(b, append(c_, d...)...)...) // Seems dirty :\
	return e
}
