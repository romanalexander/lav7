package dummy

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/L7-MCPE/lav7/types"
	"github.com/L7-MCPE/lav7/util/buffer"
)

type Provider struct {
	Name string
}

func (pv *Provider) Init(name string) {
	pv.Name = name
}

func (pv *Provider) Loadable(cx, cz int32) (path string, ok bool) {
	var err error
	path, err = filepath.Abs("levels/" + pv.Name + "/" + strconv.Itoa(int(cx)) + "_" + strconv.Itoa(int(cz)) + ".raw")
	if err != nil {
		ok = false
		return
	}
	f, err := os.Open(path)
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

func (pv *Provider) LoadChunk(cx, cz int32, path string) (chunk *types.Chunk, err error) {
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
	buf := buffer.FromBytes(b)
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

func (pv *Provider) WriteChunk(cx, cz int32, c *types.Chunk) error {
	path, _ := filepath.Abs("levels/" + pv.Name + "/" + strconv.Itoa(int(cx)) + "_" + strconv.Itoa(int(cz)) + ".raw")
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	buf := new(buffer.Buffer)
	buf.BatchWrite(c.BlockData[:], c.MetaData[:], c.LightData[:], c.SkyLightData[:], c.HeightMap[:], c.BiomeData[:])
	if err := ioutil.WriteFile(path, buf.Done(), 0644); err != nil {
		return err
	}
	return nil
}

func (pv *Provider) SaveAll(chunks map[[2]int32]*types.Chunk) error {
	errstr := ""
	for k, c := range chunks {
		if err := pv.WriteChunk(k[0], k[1], c); err != nil {
			fmt.Sprintln(errstr, err.Error())
		}
	}
	if errstr == "" {
		return fmt.Errorf(errstr)
	}
	return nil
}
