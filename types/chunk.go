package types

import "sync"

type Chunk struct {
	blockData    [16 * 16 * 128]byte
	metaData     [16 * 16 * 64]byte // Nibbles
	lightData    [16 * 16 * 64]byte // Nibbles
	skyLightData [16 * 16 * 64]byte // Nibbles
	heightMap    [16 * 16]byte
	biomeData    [16 * 16 * 4]byte // Uints
	mutex        *sync.RWMutex
}

func (c Chunk) GetBlock(x, y, z byte) byte {
	return c.blockData[uint16(y)<<8|uint16(z)<<4|uint16(x)]
}

func (c *Chunk) SetBlock(x, y, z, id byte) {
	c.blockData[uint16(y)<<8|uint16(z)<<4|uint16(x)] = id
}

func (c Chunk) GetBlockMeta(x, y, z byte) byte {
	if x&1 == 0 {
		return c.metaData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] & 0x0f
	}
	return c.metaData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] >> 4
}

func (c *Chunk) SetBlockMeta(x, y, z, id byte) {
	b := c.metaData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1]
	if x&1 == 0 {
		c.metaData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] = (b & 0xf0) | (id & 0x0f)
	} else {
		c.metaData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] = (id&0xf)<<4 | (b & 0x0f)
	}
}

func (c Chunk) GetBlockLight(x, y, z byte) byte {
	if x&1 == 0 {
		return c.lightData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] & 0x0f
	}
	return c.lightData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] >> 4
}

func (c *Chunk) SetBlockLight(x, y, z, id byte) {
	b := c.lightData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1]
	if x&1 == 0 {
		c.lightData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] = (b & 0xf0) | (id & 0x0f)
	} else {
		c.lightData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] = (id&0xf)<<4 | (b & 0x0f)
	}
}

func (c Chunk) GetBlockSkyLight(x, y, z byte) byte {
	if x&1 == 0 {
		return c.skyLightData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] & 0x0f
	}
	return c.skyLightData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] >> 4
}

func (c *Chunk) SetBlockSkyLight(x, y, z, id byte) {
	b := c.skyLightData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1]
	if x&1 == 0 {
		c.skyLightData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] = (b & 0xf0) | (id & 0x0f)
	} else {
		c.skyLightData[uint16(y)<<7|uint16(z)<<3|uint16(x)>>1] = (id&0xf)<<4 | (b & 0x0f)
	}
}

func (c Chunk) GetHeightMap(x, z byte) byte {
	return c.heightMap[uint16(z)<<4|uint16(x)]
}

func (c *Chunk) SetHeightMap(x, z, h byte) {
	c.heightMap[uint16(z)<<4|uint16(x)] = h
}

func (c Chunk) GetBiomeID(x, z byte) byte {
	return c.biomeData[uint16(z)<<6|uint16(x)<<2]
}

func (c *Chunk) SetBiomeID(x, z, id byte) {
	c.biomeData[uint16(z)<<6|uint16(x)<<2] = id
}

func (c Chunk) GetBiomeColor(x, z byte) (r, g, b byte) {
	rgb := c.biomeData[uint16(z)<<6|uint16(x)<<2+1 : uint16(z)<<6|uint16(x)<<2+4]
	return rgb[0], rgb[1], rgb[2]
}

func (c *Chunk) SetBiomeColor(x, z, r, g, b byte) {
	offset := uint16(z)<<6 | uint16(x)<<2
	c.biomeData[offset+1], c.biomeData[offset+2], c.biomeData[offset+3] = r, g, b
}

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

func (c *Chunk) Mutex() *sync.RWMutex {
	return c.mutex
}

func (c *Chunk) FromGen([16][16][128]Block) {

}

// FullChunkData returns full chunk payload for FullChunkDataPacket.
func (c Chunk) FullChunkData() []byte {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	a := append(c.blockData[:], c.metaData[:]...)     // Block ID, Block Metadata
	b := append(c.skyLightData[:], c.lightData[:]...) // SkyLight, Light
	c_ := append(c.heightMap[:], c.biomeData[:]...)   // Height Map, Biome colors
	d := []byte{0, 0, 0, 0}                           // Extra data: length 0
	// No tile entity NBT fields
	e := append(a, append(b, append(c_, d...)...)...) // Seems dirty :\
	return e
}

/*

   func (c *Chunk) BlockChunk(bs [16 * 16 * 128]types.Block) {
   	c.mutex.Lock()
   	defer c.mutex.Unlock()
   	for i, b := range bs {
   		c.blockData[i] = b.ID
   		if i&0x01 == 0 {
   			c.metaData[i>>1] = (c.metaData[i>>1] & 0xf0) | (b.ID & 0x0f)
   		} else {
   			c.metaData[i>>1] = (b.Meta&0xf)<<4 | (c.metaData[i>>1] & 0x0f)
   		}
   	}
   	c.PopulateHeight()
   	c.beautifulize()
   }
*/

func (c *Chunk) beautifulize() {
	for x := byte(0); x < 16; x++ {
		for z := byte(0); z < 16; z++ {
			c.SetBiomeColor(x, z, x*16, x*z, z*16)
		}
	}
}
