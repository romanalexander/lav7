package dummy

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/L7-MCPE/lav7/types"
	"github.com/L7-MCPE/lav7/util/buffer"
)

type Provider struct {
	Name string

	ChunkMap   map[[2]int32]*types.Chunk
	ChunkMutex *sync.Mutex
	Generator  func(int32, int32) *types.Chunk
}

func (pv Provider) ChunkExists(cx, cz int32) bool {
	_, ok := pv.ChunkMap[[2]int32{cx, cz}]
	return ok
}

func (pv Provider) GetChunk(cx, cz int32, create bool) (chk *types.Chunk) {
	pv.ChunkMutex.Lock()
	defer pv.ChunkMutex.Unlock()
	if c, ok := pv.ChunkMap[[2]int32{cx, cz}]; ok {
		return c
	}
	if path, err := filepath.Abs("levels/" + pv.Name + "/" + strconv.Itoa(int(cx)) + "_" + strconv.Itoa(int(cz)) + ".raw"); err != nil {
		goto crt
	} else if _, err := os.Stat(path); os.IsNotExist(err) {
		goto crt
	} else {
		f, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
		if err != nil {
			goto crt
		}
		b := make([]byte, 83200)
		_, err = f.Read(b)
		if err != nil {
			goto crt
		}
		chk = new(types.Chunk)
		buf := buffer.FromBytes(b)
		chk.Read(buf)
		pv.SetChunk(cx, cz, false, chk)
		return
	}
crt:
	if create {
		chk = pv.Generator(cx, cz)
		pv.SetChunk(cx, cz, false, chk)
		return
	}
	return
}

func (pv Provider) SetChunk(cx, cz int32, force bool, c *types.Chunk) {
	pv.ChunkMutex.Lock()
	defer pv.ChunkMutex.Unlock()
	if _, ok := pv.ChunkMap[[2]int32{cx, cz}]; !force && ok {
		panic("Tried to overwrite existing chunk!")
	}
	if c.Mutex() == nil {
		panic("Nil mutex: chunk may have been uninitialized!")
	}
	pv.ChunkMap[[2]int32{cx, cz}] = c
}

func (pv Provider) Save(name string) error {
	pv.ChunkMutex.Lock()
	defer pv.ChunkMutex.Unlock()
	for k, c := range pv.ChunkMap {
		path, _ := filepath.Abs("levels/" + name + "/" + strconv.Itoa(int(k[0])) + "_" + strconv.Itoa(int(k[1])) + ".raw")
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			log.Println("Error while creating dir:", err)
		}
		buf, err := c.Write()
		if err != nil {
			log.Println("Error while writing chunk to file:", err)
			continue
		}
		if err := ioutil.WriteFile(path, buf.Done(), 0644); err != nil {
			log.Println("Error while saving:", err)
			continue
		}
	}
	return nil
}
