package raknet

import (
	"bytes"
	"net"
	"time"

	"github.com/L7-MCPE/lav7/util/buffer"
)

// Packet is a struct which contains binary buffer, address, and send time.
type Packet struct {
	*bytes.Buffer
	Address *net.UDPAddr
}

// NewPacket creates new packet with given packet id.
func NewPacket(pid byte) Packet {
	return Packet{bytes.NewBuffer([]byte{pid}), new(net.UDPAddr)}
}

// EncapsulatedPacket is a struct, containing more values for decoding/encoding encapsualted packets.
type EncapsulatedPacket struct {
	*bytes.Buffer
	Reliability  byte
	HasSplit     bool
	MessageIndex uint32 // LE Triad
	OrderIndex   uint32 // LE Triad
	OrderChannel byte
	SplitCount   uint32
	SplitID      uint16
	SplitIndex   uint32
}

// NewEncapsulated returns decoded EncapsulatedPacket struct from given binary.
// Do NOT set buf with *Packet struct. It could cause panic.
func NewEncapsulated(buf *bytes.Buffer) (ep *EncapsulatedPacket) {
	ep = new(EncapsulatedPacket)
	flags := buffer.ReadByte(buf)
	ep.Reliability = flags >> 5
	ep.HasSplit = (flags>>4)&1 > 0
	l := uint32(buffer.ReadShort(buf))
	length := l >> 3
	if l%8 != 0 {
		length++
	}
	if ep.Reliability > 0 {
		if ep.Reliability >= 2 && ep.Reliability != 5 {
			ep.MessageIndex = buffer.ReadLTriad(buf)
		}
		if ep.Reliability <= 4 && ep.Reliability != 2 {
			ep.OrderIndex = buffer.ReadLTriad(buf)
			ep.OrderChannel = buffer.ReadByte(buf)
		}
	}
	if ep.HasSplit {
		ep.SplitCount = buffer.ReadInt(buf)
		ep.SplitID = buffer.ReadShort(buf)
		ep.SplitIndex = buffer.ReadInt(buf)
	}
	b, err := buffer.Read(buf, int(length))
	if err != nil {
		panic(err.Error())
	}
	ep.Buffer = bytes.NewBuffer(b)
	return
}

// TotalLen returns total binary length of EncapsulatedPacket.
func (ep *EncapsulatedPacket) TotalLen() int {
	return 3 + ep.Buffer.Len() + func() int {
		return func() int {
			if ep.Reliability >= 2 && ep.Reliability != 5 {
				return 3
			}
			return 0
		}() + func() int {
			if ep.Reliability != 0 && ep.Reliability <= 4 && ep.Reliability != 2 {
				return 4
			}
			return 0
		}()
	}() + func() int {
		if ep.HasSplit {
			return 10
		}
		return 0
	}()
}

// Bytes returns encoded binary from EncapsulatedPacket struct options.
func (ep *EncapsulatedPacket) Bytes() (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	buffer.WriteByte(buf, ep.Reliability<<5|func() byte {
		if ep.HasSplit {
			return 1 << 4
		}
		return 0
	}())
	buffer.WriteShort(buf, uint16(ep.Len())<<3)
	if ep.Reliability > 0 {
		buffer.Write(buf, func() []byte {
			buf := new(bytes.Buffer)
			if ep.Reliability >= 2 && ep.Reliability != 5 {
				buffer.WriteLTriad(buf, ep.MessageIndex)
			}
			if ep.Reliability <= 4 && ep.Reliability != 2 {
				buffer.WriteLTriad(buf, ep.OrderIndex)
				buffer.WriteByte(buf, ep.OrderChannel)
			}
			return buf.Bytes()
		}())
	}
	if ep.HasSplit {
		buffer.WriteInt(buf, ep.SplitCount)
		buffer.WriteShort(buf, ep.SplitID)
		buffer.WriteInt(buf, ep.SplitIndex)
	}
	b := ep.Buffer.Bytes()
	buffer.Write(buf, b)
	return
}

// DataPacket is a packet struct, containing Raknet data packet fields.
type DataPacket struct {
	*bytes.Buffer
	Head      byte
	SendTime  time.Time
	SeqNumber uint32 // LE Triad
	Packets   []*EncapsulatedPacket
}

// Decode decodes buffer to struct fields and decapsulates all packets.
func (dp *DataPacket) Decode() {
	// dp.Head = buffer.ReadByte(dp.Buffer)
	dp.SeqNumber = buffer.ReadLTriad(dp.Buffer)
	for dp.Buffer.Len() > 0 {
		ep := NewEncapsulated(dp.Buffer)
		dp.Packets = append(dp.Packets, ep)
	}
	return
}

// Len returns total buffer length of data packet.
func (dp *DataPacket) TotalLen() int {
	length := 4
	for _, d := range dp.Packets {
		length += d.TotalLen()
	}
	return length
}

// Encode encodes fields and packets to buffer.
func (dp *DataPacket) Encode() {
	dp.Buffer = new(bytes.Buffer)
	buffer.WriteByte(dp.Buffer, dp.Head)
	buffer.WriteLTriad(dp.Buffer, dp.SeqNumber)
	for _, ep := range dp.Packets {
		buffer.Write(dp.Buffer, ep.Bytes().Bytes())
	}
	return
}
