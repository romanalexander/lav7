package types

import (
	"bytes"
	"io"

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

// Read reads item data from io.Reader interface.
func (i *Item) Read(buf io.Reader) {
	i.ID = ID(buffer.ReadShort(buf))
	if i.ID == 0 {
		return
	}
	i.Amount = buffer.ReadByte(buf)
	i.Meta = buffer.ReadShort(buf)
	length := uint32(buffer.ReadShort(buf))
	if length > 0 {
		b, _ := buffer.Read(buf, int(length))
		compound := bytes.NewBuffer(b)
		i.Compound = new(nbt.Compound)
		i.Compound.ReadFrom(compound)
	}
}

// Write returns byte slice with item data.
func (i Item) Write() []byte {
	if i.ID == 0 {
		return []byte{0, 0}
	}
	buf := new(bytes.Buffer)
	buffer.WriteShort(buf, uint16(i.ID))
	buffer.WriteByte(buf, i.Amount)
	buffer.WriteShort(buf, i.Meta)
	compound := new(bytes.Buffer)
	i.Compound.WriteTo(compound)
	buffer.WriteShort(buf, uint16(compound.Len()))
	buf.Write(compound.Bytes())
	return buf.Bytes()
}

// Block converts the item to block struct.
// If ID is not a block ID, it panics.
func (i Item) Block() Block {
	return Block{
		ID:   i.ID.Block(),
		Meta: byte(i.Meta),
	}
}
