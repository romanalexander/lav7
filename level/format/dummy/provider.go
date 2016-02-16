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

func (pv *Provider) Init(gen func(int32, int32) *types.Chunk, name string) {
	pv.ChunkMap = make(map[[2]int32]*types.Chunk)
	pv.ChunkMutex = new(sync.Mutex)
	pv.Generator = gen
	pv.Name = name
}

func (pv *Provider) ChunkExists(cx, cz int32) bool {
	pv.ChunkMutex.Lock()
	_, ok := pv.ChunkMap[[2]int32{cx, cz}]
	pv.ChunkMutex.Unlock()
	return ok
}

func (pv *Provider) GetChunk(cx, cz int32, create bool) (chk *types.Chunk) {
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
		chk.Mutex().Lock()
		copy(chk.BlockData[:], buf.Read(16*16*128))
		copy(chk.MetaData[:], buf.Read(16*16*64))
		copy(chk.LightData[:], buf.Read(16*16*64))
		copy(chk.SkyLightData[:], buf.Read(16*16*64))
		copy(chk.HeightMap[:], buf.Read(16*16))
		copy(chk.BiomeData[:], buf.Read(16*16*4))
		chk.Mutex().Unlock()
		pv.ChunkMap[[2]int32{cx, cz}] = chk
		return
	}
crt:
	if create {
		chk = pv.Generator(cx, cz)
		pv.ChunkMap[[2]int32{cx, cz}] = chk
		return
	}
	return
}

func (pv *Provider) SetChunk(cx, cz int32, force bool, c *types.Chunk) {
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

func (pv *Provider) Save() error {
	pv.ChunkMutex.Lock()
	defer pv.ChunkMutex.Unlock()
	for k, c := range pv.ChunkMap {
		path, _ := filepath.Abs("levels/" + pv.Name + "/" + strconv.Itoa(int(k[0])) + "_" + strconv.Itoa(int(k[1])) + ".raw")
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			log.Println("Error while creating dir:", err)
			continue
		}
		buf := new(buffer.Buffer)
		buf.BatchWrite(c.BlockData[:], c.MetaData[:], c.LightData[:], c.SkyLightData[:], c.HeightMap[:], c.BiomeData[:])
		if err := ioutil.WriteFile(path, buf.Done(), 0644); err != nil {
			log.Println("Error while saving:", err)
			continue
		}
	}
	return nil
}
