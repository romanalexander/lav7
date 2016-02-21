/*
 Package dummy provides an example for writing level format provider.

 This format saves each chunk in separate file, writing raw block id/meta/height map/skylight data, etc.
 Use of this format in production server is not recommended.
*/
package dummy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/L7-MCPE/lav7"
	"github.com/L7-MCPE/lav7/types"
	"github.com/L7-MCPE/lav7/util/buffer"
)

func init() {
	lav7.RegisterProvider(new(Dummy))
}

type Dummy struct {
	Name string
}

func (dm *Dummy) Init(name string) {
	dm.Name = name
}

func (dm *Dummy) Loadable(cx, cz int32) (path string, ok bool) {
	var err error
	path, err = filepath.Abs("levels/" + dm.Name + "/" + strconv.Itoa(int(cx)) + "_" + strconv.Itoa(int(cz)) + ".raw")
	if err != nil {
		ok = false
		return
	}
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		ok = false
		return
	}
	fi, err := f.Stat()
	if err != nil {
		ok = false
		return
	}
	ok = !fi.IsDir() && fi.Size() >= 83200
	return
}

func (dm *Dummy) LoadChunk(cx, cz int32, path string) (chunk *types.Chunk, err error) {
	var f *os.File
	f, err = os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return
	}
	b := make([]byte, 83200)
	_, err = f.Read(b)
	if err != nil {
		return
	}
	chunk = new(types.Chunk)
	buf := bytes.NewBuffer(b)
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

func (dm *Dummy) WriteChunk(cx, cz int32, c *types.Chunk) error {
	path, _ := filepath.Abs("levels/" + dm.Name + "/" + strconv.Itoa(int(cx)) + "_" + strconv.Itoa(int(cz)) + ".raw")
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	c.Mutex().Lock()
	defer c.Mutex().Unlock()
	buffer.BatchWrite(buf, c.BlockData[:], c.MetaData[:], c.LightData[:], c.SkyLightData[:], c.HeightMap[:], c.BiomeData[:])
	if err := ioutil.WriteFile(path, buf.Bytes(), 0644); err != nil {
		return err
	}
	return nil
}

func (dm *Dummy) SaveAll(chunks map[[2]int32]*types.Chunk) error {
	errstr := ""
	for k, c := range chunks {
		if err := dm.WriteChunk(k[0], k[1], c); err != nil {
			fmt.Sprintln(errstr, err.Error())
		}
	}
	if errstr == "" {
		return fmt.Errorf(errstr)
	}
	return nil
}
