/*
 Package vilan implements simple world format for lav7 servers.

 Vilan splits worlds into section files, containing 16(4*4) chunks.
 Filename format is 'section.sectionX.sectionZ.v'.

 Each sections' chunk structures are same as dummy format.
 There are no tile entity/NBT support for now.
*/
package vilan

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/L7-MCPE/lav7"
	"github.com/L7-MCPE/lav7/types"
	"github.com/L7-MCPE/lav7/util/buffer"
)

func init() {
	lav7.RegisterProvider(new(Vilan))
}

type Vilan struct {
	name string
}

func (v *Vilan) Init(name string) {
	v.name = name
}

func (v *Vilan) Loadable(cx, cz int32) (path string, ok bool) {
	sectionX, sectionZ := cx>>2, cz>>2
	path = fmt.Sprintf("levels/%s/section.%d.%d.v", v.name, sectionX, sectionZ)
	file, err := os.Open(path)
	if err != nil {
		log.Println("Error while opening chunk section:", err)
		return "", false
	}
	defer file.Close()

	buf := make([]byte, 2)
	_, err = file.Read(buf)
	if err != nil {
		log.Println("Error while reading chunk status byte:", err)
		return "", false
	}

	chunkstat := uint16(buf[1])<<8 | uint16(buf[0])

	ok = (chunkstat>>(byte(cx&3)<<2|byte(cz&3)))&1 == 1
	return
}

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
	copy(chunk.BlockData[:], buf.Read(16*16*128))
	copy(chunk.MetaData[:], buf.Read(16*16*64))
	copy(chunk.LightData[:], buf.Read(16*16*64))
	copy(chunk.SkyLightData[:], buf.Read(16*16*64))
	copy(chunk.HeightMap[:], buf.Read(16*16))
	copy(chunk.BiomeData[:], buf.Read(16*16*4))
	chunk.Mutex().Unlock()
	return
}

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
	fbuf := make([]byte, 2)
	if _, err := file.Read(fbuf); err != nil {
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
	buf.BatchWrite(chunk.BlockData[:], chunk.MetaData[:], chunk.LightData[:], chunk.SkyLightData[:], chunk.HeightMap[:], chunk.BiomeData[:])

	pos := 2 + int64(byte(cx&3)<<2|byte(cz&3))*83200
	_, err = file.WriteAt(buf.Done(), pos)
	return err
}

func (v *Vilan) SaveAll(chunks map[[2]int32]*types.Chunk) error {
	errstr := ""
	for k, c := range chunks {
		if err := v.WriteChunk(k[0], k[1], c); err != nil {
			fmt.Sprintln(errstr, err.Error())
		}
	}
	if errstr == "" {
		return fmt.Errorf(errstr)
	}
	return nil
}
