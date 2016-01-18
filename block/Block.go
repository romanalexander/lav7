package block

import "fmt"

// IBlock is an interface to contain block element values, such as block ID, meta, etc.
type IBlock interface {
	GetID() byte
	SetID(byte)
	GetMeta() byte
	SetMeta(byte)
	Tick() error
	String() string
}

// Block is a generic struct of most 'nothing-special' blocks.
type Block struct {
	ID   byte
	Meta byte
}

func (b Block) GetID() byte {
	return b.ID
}

func (b *Block) SetID(id byte) {
	b.ID = id
}

func (b Block) GetMeta() byte {
	return b.Meta
}

func (b *Block) SetMeta(meta byte) {
	b.Meta = meta
}

func (b *Block) Tick() error { return nil }

func (b Block) String() string {
	return fmt.Sprintf("{Block ID: %d, Meta: %d}", b.ID, b.Meta)
}
