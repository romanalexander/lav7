// Package vilan implements simple world format for lav7 servers.
//
// Vilan splits worlds into section files, containing 16(4*4) chunks.
// Filename format is 'section.sectionX.sectionZ.v'.
//
// Each sections' chunk structures are same as dummy format.
// There are no tile entity/NBT support for now.
package vilan

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/L7-MCPE/lav7"
	"github.com/L7-MCPE/lav7/types"
	"github.com/L7-MCPE/lav7/util/buffer"
	"github.com/L7-MCPE/lav7/util/try"
)

func init() {
	lav7.RegisterProvider(new(Vilan))
}

// Vilan is a improved version of Dummy, grouping 16 chunks into a single section.
type Vilan struct {
	name string
}

// Init implements format.Provider interface.
func (v *Vilan) Init(name string) {
	v.name = name
}

// Loadable implements format.Provider interface.
func (v *Vilan) Loadable(cx, cz int32) (path string, ok bool) {
	sectionX, sectionZ := cx>>2, cz>>2
	path = fmt.Sprintf("levels/%s/section.%d.%d.v", v.name, sectionX, sectionZ)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		log.Println("Error while creating directory:", err)
		return "", false
	}
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return path, false
		}
		log.Println("Error while opening chunk section:", err)
		return "", false
	}
	defer file.Close()

	var chunkstat uint16
	if err := try.Safe(func() {
		chunkstat = buffer.ReadShort(file)
	}); err != nil {
		log.Println("Error while reading chunk status byte:", err)
		return "", false
	}

	ok = (chunkstat>>(byte(cx&3)<<2|byte(cz&3)))&1 == 1
	return
}

// LoadChunk implements format.Provider interface.
func (v *Vilan) LoadChunk(cx, cz int32, path string) (chunk *types.Chunk, err error) {
	sectionX, sectionZ := cx>>2, cz>>2
	path = fmt.Sprintf("levels/%s/section.%d.%d.v", v.name, sectionX, sectionZ)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	pos := 2 + int64(byte(cx&3)<<2|byte(cz&3))*83200
	fbuf := make([]byte, 83200)
	_, err = file.ReadAt(fbuf, pos)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(fbuf)
	chunk = new(types.Chunk)
	chunk.Mutex().Lock()
	copy(chunk.BlockData[:], buf.Next(16*16*128))
	copy(chunk.MetaData[:], buf.Next(16*16*64))
	copy(chunk.LightData[:], buf.Next(16*16*64))
	copy(chunk.SkyLightData[:], buf.Next(16*16*64))
	copy(chunk.HeightMap[:], buf.Next(16*16))
	copy(chunk.BiomeData[:], buf.Next(16*16*4))
	chunk.Mutex().Unlock()
	return
}

// WriteChunk implements format.Provider interface.
func (v *Vilan) WriteChunk(cx, cz int32, chunk *types.Chunk) error {
	sectionX, sectionZ := cx>>2, cz>>2
	path := fmt.Sprintf("levels/%s/section.%d.%d.v", v.name, sectionX, sectionZ)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	fstat, err := file.Stat()
	if err != nil {
		return err
	}

	fbuf := make([]byte, 2)
	if fstat.Size() < 16*83200+2 {
		file.WriteAt(make([]byte, 16*83200+2), 0)
	} else if _, err := file.Read(fbuf); err != nil {
		return err
	}

	offset := byte(cx&3)<<2 | byte(cz&3)
	if offset >= 8 {
		fbuf[0] |= 1 << (offset - 8)
	} else {
		fbuf[1] |= 1 << offset
	}
	if _, err := file.WriteAt(fbuf, 0); err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	chunk.Mutex().Lock()
	defer chunk.Mutex().Unlock()
	buffer.BatchWrite(buf, chunk.BlockData[:], chunk.MetaData[:], chunk.LightData[:], chunk.SkyLightData[:], chunk.HeightMap[:], chunk.BiomeData[:])

	pos := 2 + int64(byte(cx&3)<<2|byte(cz&3))*83200
	_, err = file.WriteAt(buf.Bytes(), pos)
	return err
}

// SaveAll implements format.Provider interface.
func (v *Vilan) SaveAll(chunks map[[2]int32]*types.Chunk) error {
	errstr := ""
	for k, c := range chunks {
		if err := v.WriteChunk(k[0], k[1], c); err != nil {
			errstr += err.Error()
		}
	}
	if errstr != "" {
		return fmt.Errorf(errstr)
	}
	return nil
}
