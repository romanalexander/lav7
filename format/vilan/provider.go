package vilan

import (
	"log"
	"os"

	"github.com/L7-MCPE/lav7"
	"github.com/L7-MCPE/lav7/types"
)

func init() {
	lav7.RegisterProvider(new(Vilan))
}

type Vilan struct {
	name string
	file *os.File
}

func (v *Vilan) Init(name string) {
	v.name = name
	var err error
	v.file, err = os.OpenFile("levels/"+name+"/chunk.dat", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal("Error while opening level data:", err)
	}
}

func (v *Vilan) Loadable(cx, cz int32) (path string, ok bool) {
	return
}

func (v *Vilan) LoadChunk(cx, cz int32, path string) (chunk *types.Chunk, err error) {
	return
}

func (v *Vilan) WriteChunk(cx, cz int32, chunk *types.Chunk) error {
	return nil
}

func (v *Vilan) SaveAll(chunks map[[2]int32]*types.Chunk) error {
	return nil
}
