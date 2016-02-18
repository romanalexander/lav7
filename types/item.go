package types

import (
	"bytes"

	"github.com/L7-MCPE/lav7/util/buffer"
	"github.com/minero/minero-go/proto/nbt"
)

// Item contains item data for each container slots.
type Item struct {
	ID       ID
	Meta     uint16
	Amount   byte
	Compound *nbt.Compound
}

func (i *Item) Read(buf *buffer.Buffer) {
	i.ID = ID(buf.ReadShort())
	if i.ID == 0 {
		return
	}
	i.Amount = buf.ReadByte()
	i.Meta = buf.ReadShort()
	length := uint32(buf.ReadShort())
	if length > 0 {
		compound := bytes.NewBuffer(buf.Read(length))
		i.Compound = new(nbt.Compound)
		i.Compound.ReadFrom(compound)
	}
}

func (i Item) Write() []byte {
	if i.ID == 0 {
		return []byte{0, 0}
	}
	buf := new(buffer.Buffer)
	buf.WriteShort(uint16(i.ID))
	buf.WriteByte(i.Amount)
	buf.WriteShort(i.Meta)
	compound := new(bytes.Buffer)
	i.Compound.WriteTo(compound)
	buf.WriteShort(uint16(compound.Len()))
	buf.Write(compound.Bytes())
	return buf.Done()
}

func (i Item) Block() Block {
	return Block{
		ID:   i.ID.Block(),
		Meta: byte(i.Meta),
	}
}
