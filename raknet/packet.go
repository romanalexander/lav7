package raknet

import (
	"net"
	"time"

	"github.com/L7-MCPE/lav7/util"
	"github.com/L7-MCPE/lav7/util/buffer"
)

// Packet is a struct which contains binary buffer, address, and send time.
type Packet struct {
	*buffer.Buffer
	Address *net.UDPAddr
}

// NewPacket creates new packet with given packet id.
func NewPacket(pid byte) Packet {
	return Packet{buffer.FromBytes([]byte{pid}), new(net.UDPAddr)}
}

// EncapsulatedPacket is a struct, containing more values for decoding/encoding encapsualted packets.
type EncapsulatedPacket struct {
	*buffer.Buffer
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
func NewEncapsulated(buf *buffer.Buffer) (ep *EncapsulatedPacket) {
	ep = new(EncapsulatedPacket)
	ep.Buffer = new(buffer.Buffer)
	flags := buf.ReadByte()
	ep.Reliability = flags >> 5
	ep.HasSplit = (flags>>4)&1 > 0
	l := uint32(buf.ReadShort())
	length := l >> 3
	if l%8 != 0 {
		length++
	}
	util.Debug("Length", length)
	if ep.Reliability > 0 {
		if ep.Reliability >= 2 && ep.Reliability != 5 {
			ep.MessageIndex = buf.ReadLTriad()
		}
		if ep.Reliability <= 4 && ep.Reliability != 2 {
			ep.OrderIndex = buf.ReadLTriad()
			ep.OrderChannel = buf.ReadByte()
		}
	}
	if ep.HasSplit {
		ep.SplitCount = buf.ReadInt()
		ep.SplitID = buf.ReadShort()
		ep.SplitIndex = buf.ReadInt()
	}
	b := buf.Read(length)
	ep.Buffer = buffer.FromBytes(b)
	return
}

// TotalLen returns total binary length of EncapsulatedPacket
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
func (ep *EncapsulatedPacket) Bytes() (buf *buffer.Buffer) {
	buf = new(buffer.Buffer)
	buf.WriteByte(ep.Reliability<<5 | func() byte {
		if ep.HasSplit {
			return 1 << 4
		}
		return 0
	}())
	buf.WriteShort(uint16(len(ep.Payload)) << 3)
	if ep.Reliability > 0 {
		buf.Write(func() []byte {
			buf := new(buffer.Buffer)
			if ep.Reliability >= 2 && ep.Reliability != 5 {
				buf.WriteLTriad(ep.MessageIndex)
			}
			if ep.Reliability <= 4 && ep.Reliability != 2 {
				buf.WriteLTriad(ep.OrderIndex)
				buf.WriteByte(ep.OrderChannel)
			}
			return buf.Done()
		}())
	}
	if ep.HasSplit {
		buf.WriteInt(ep.SplitCount)
		buf.WriteShort(ep.SplitID)
		buf.WriteInt(ep.SplitIndex)
	}
	buf.Append(ep.Buffer)
	return
}

// DataPacket is a packet struct, containing Raknet data packet fields.
type DataPacket struct {
	*buffer.Buffer
	Head      byte
	SendTime  time.Time
	SeqNumber uint32 // LE Triad
	Packets   []*EncapsulatedPacket
}

// Decode decodes buffer to struct fields and decapsulates all packets.
func (dp *DataPacket) Decode() {
	dp.Offset = 0
	dp.Head = dp.ReadByte()
	dp.SeqNumber = dp.ReadLTriad()
	for dp.Require(1) {
		b := dp.Read(0)
		ep := NewEncapsulated(buffer.FromBytes(b))
		dp.Packets = append(dp.Packets, ep)
		dp.Offset += uint32(ep.TotalLen())
	}
	return
}

// Len returns total buffer length of data packet.
func (dp *DataPacket) Len() int {
	length := 4
	for _, d := range dp.Packets {
		length += d.TotalLen()
	}
	return length
}

// Encode encodes fields and packets to buffer.
func (dp *DataPacket) Encode() {
	dp.Buffer = new(buffer.Buffer)
	dp.WriteByte(dp.Head)
	dp.WriteLTriad(dp.SeqNumber)
	for _, ep := range dp.Packets {
		b := ep.Bytes()
		dp.Write(b.Done())
	}
	return
}
