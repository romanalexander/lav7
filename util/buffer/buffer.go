package buffer

import (
	"fmt"
	"math"
	"net"

	"github.com/L7-MCPE/lav7/util"
)

// Buffer contains packet payload, serves binary functions for efficiency
type Buffer struct {
	Payload []byte
	Offset  uint32
}

// FromBytes converts byte array to Buffer struct
func FromBytes(buf []byte) *Buffer {
	return &Buffer{
		Payload: buf,
		Offset:  0,
	}
}

// Append appends given buffer to payload.
func (buf *Buffer) Append(b *Buffer) {
	buf.Payload = append(buf.Payload, b.Payload...)
}

// Done resets entire buffer and returns payload.
func (buf *Buffer) Done() (b []byte) {
	b = buf.Payload
	buf.Payload = make([]byte, 0)
	buf.Offset = 0
	return
}

// Head returns the first byte of packet payload.
func (buf *Buffer) Head() byte {
	if len(buf.Payload) > 0 {
		return buf.Payload[0]
	}
	return 0
}

// Require returns true if packet has needed unread bytes from buffer, otherwise false.
func (buf *Buffer) Require(need uint32) bool {
	if uint32(len(buf.Payload)) < buf.Offset+need {
		return false
	}
	return true
}

// Read reads given byte from buffer.
func (buf *Buffer) Read(length uint32) (r []byte) {
	if length == 0 {
		if buf.Require(1) {
			return buf.Payload[buf.Offset:]
		}
		panic(util.EOFError{
			BufLen:    buf.Len(),
			BufOffset: buf.Offset,
			Needed:    0,
			Buf:       buf.Payload,
		}.Error())
	}
	if !buf.Require(length) {
		panic(util.EOFError{
			BufLen:    buf.Len(),
			BufOffset: buf.Offset,
			Needed:    length,
			Buf:       buf.Payload,
		}.Error())
	}
	r = buf.Payload[buf.Offset : buf.Offset+length]
	buf.Offset += length
	return
}

// ReadAny reads appropriate type from given reference value.
func (buf *Buffer) ReadAny(p interface{}) {
	switch p.(type) {
	case *bool:
		*p.(*bool) = buf.ReadBool()
	case *byte:
		*p.(*byte) = buf.ReadByte()
	case *uint16:
		*p.(*uint16) = buf.ReadShort()
	case *uint32:
		*p.(*uint32) = buf.ReadInt()
	case *uint64:
		*p.(*uint64) = buf.ReadLong()
	case *float32:
		*p.(*float32) = buf.ReadFloat()
	case *float64:
		*p.(*float64) = buf.ReadDouble()
	case *string:
		*p.(*string) = buf.ReadString()
	case *net.UDPAddr:
		*p.(*net.UDPAddr) = *buf.ReadAddress()
	case byte, uint16, uint32,
		uint64, float32, float64, string:
		panic("ReadAny requires reference type")
	}
}

// BatchRead batches ReadAny from given reference pointers.
func (buf *Buffer) BatchRead(p ...interface{}) {
	for _, pp := range p {
		buf.ReadAny(pp)
	}
}

// ReadBool reads boolean from buffer.
func (buf *Buffer) ReadBool() bool {
	return buf.ReadByte() > 0
}

// ReadByte reads unsigned byte from buffer.
func (buf *Buffer) ReadByte() byte {
	b := buf.Read(1)
	return b[0]
}

// ReadShort reads unsigned short from buffer.
func (buf *Buffer) ReadShort() uint16 {
	b := buf.Read(2)
	return uint16(b[0])<<8 | uint16(b[1])
}

// ReadLShort reads unsigned little-endian short from buffer.
func (buf *Buffer) ReadLShort() uint16 {
	b := buf.Read(2)
	return uint16(b[1])<<8 | uint16(b[0])
}

// ReadInt reads unsigned int from buffer.
func (buf *Buffer) ReadInt() uint32 {
	b := buf.Read(4)
	return uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
}

// ReadLInt reads unsigned little-endian int from buffer.
func (buf *Buffer) ReadLInt() uint32 {
	b := buf.Read(4)
	return uint32(b[3])<<24 | uint32(b[2])<<16 | uint32(b[1])<<8 | uint32(b[0])
}

// ReadLong reads unsigned long from buffer.
func (buf *Buffer) ReadLong() uint64 {
	b := buf.Read(8)
	return uint64(b[0])<<56 | uint64(b[1])<<48 |
		uint64(b[2])<<40 | uint64(b[3])<<32 |
		uint64(b[4])<<24 | uint64(b[5])<<16 |
		uint64(b[6])<<8 | uint64(b[7])
}

// ReadLLong reads unsigned little-endian long from buffer.
func (buf *Buffer) ReadLLong() uint64 {
	b := buf.Read(8)
	return uint64(b[7])<<56 | uint64(b[6])<<48 |
		uint64(b[5])<<40 | uint64(b[4])<<32 |
		uint64(b[3])<<24 | uint64(b[2])<<16 |
		uint64(b[1])<<8 | uint64(b[0])
}

// ReadFloat reads 32-bit float from buffer.
func (buf *Buffer) ReadFloat() float32 {
	return math.Float32frombits(buf.ReadInt())
}

// ReadDouble reads 64-bit float from buffer.
func (buf *Buffer) ReadDouble() float64 {
	return math.Float64frombits(buf.ReadLong())
}

// ReadTriad reads unsigned 3-bytes triad from buffer.
func (buf *Buffer) ReadTriad() uint32 {
	b := buf.Read(3)
	return uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2])
}

// ReadLTriad reads unsigned little-endian 3-bytes triad from buffer.
func (buf *Buffer) ReadLTriad() uint32 {
	b := buf.Read(3)
	return uint32(b[2])<<16 | uint32(b[1])<<8 | uint32(b[0])
}

// ReadString reads string from buffer.
func (buf *Buffer) ReadString() (str string) {
	var l uint16
	l = buf.ReadShort()
	var b []byte
	b = buf.Read(uint32(l))
	return string(b)
}

// ReadAddress reads IP address/port from buffer.
func (buf *Buffer) ReadAddress() (addr *net.UDPAddr) {
	var v byte
	v = buf.ReadByte()
	if v != 4 {
		panic(fmt.Sprintf("Expected IPv4, got IP version %d", v))
	}
	b := buf.Read(4)
	p := buf.ReadShort()
	return &net.UDPAddr{
		IP:   append([]byte{b[0] ^ 0xff}, b[1]^0xff, b[2]^0xff, b[3]^0xff),
		Port: int(p),
	}
}

// Write writes given byte array to buffer.
func (buf *Buffer) Write(b []byte) error {
	if buf.Len()+len(b) > 1024*1024*256 {
		return util.LargeBufferError{
			OldCap: buf.Len(),
			Append: len(b),
		}
	}
	buf.Payload = append(buf.Payload, b...)
	return nil
}

// WriteAny writes appropriate type from given interface{} value to buffer.
func (buf *Buffer) WriteAny(p interface{}) {
	switch p.(type) {
	case bool:
		buf.WriteBool(p.(bool))
	case byte:
		buf.WriteByte(p.(byte))
	case uint16:
		buf.WriteShort(p.(uint16))
	case uint32:
		buf.WriteInt(p.(uint32))
	case uint64:
		buf.WriteLong(p.(uint64))
	case float32:
		buf.WriteFloat(p.(float32))
	case float64:
		buf.WriteDouble(p.(float64))
	case string:
		buf.WriteString(p.(string))
	case []byte:
		buf.Write(p.([]byte))
	case *bool:
		buf.WriteBool(*p.(*bool))
	case *byte:
		buf.WriteByte(*p.(*byte))
	case *uint16:
		buf.WriteShort(*p.(*uint16))
	case *uint32:
		buf.WriteInt(*p.(*uint32))
	case *uint64:
		buf.WriteLong(*p.(*uint64))
	case *float32:
		buf.WriteFloat(*p.(*float32))
	case *float64:
		buf.WriteDouble(*p.(*float64))
	case *string:
		buf.WriteString(*p.(*string))
	case *[]byte:
		buf.Write(*p.(*[]byte))
	case *net.UDPAddr:
		buf.WriteAddress(p.(*net.UDPAddr))
	}
}

// BatchWrite batches WriteAny from given values.
func (buf *Buffer) BatchWrite(p ...interface{}) {
	for _, pp := range p {
		buf.WriteAny(pp)
	}
}

// WriteBool writes boolean to buffer.
func (buf *Buffer) WriteBool(n bool) error {
	return buf.WriteByte(func() byte {
		if n {
			return 1
		} else {
			return 0
		}
	}())
}

// WriteByte writes unsigned byte to buffer.
func (buf *Buffer) WriteByte(n byte) error {
	return buf.Write([]byte{n})
}

// WriteShort writes unsigned short to buffer.
func (buf *Buffer) WriteShort(n uint16) error {
	return buf.Write([]byte{byte(n >> 8), byte(n)})
}

// WriteLShort writes unsigned little-endian short to buffer.
func (buf *Buffer) WriteLShort(n uint16) error {
	return buf.Write([]byte{byte(n), byte(n >> 8)})
}

// WriteInt writes unsigned int to buffer.
func (buf *Buffer) WriteInt(n uint32) error {
	return buf.Write([]byte{byte(n >> 24), byte(n >> 16), byte(n >> 8), byte(n)})
}

// WriteLInt writes unsigned little-endian int to buffer.
func (buf *Buffer) WriteLInt(n uint32) error {
	return buf.Write([]byte{byte(n), byte(n >> 8), byte(n >> 16), byte(n >> 24)})
}

// WriteLong writes unsigned long to buffer.
func (buf *Buffer) WriteLong(n uint64) error {
	return buf.Write([]byte{
		byte(n >> 56), byte(n >> 48),
		byte(n >> 40), byte(n >> 32),
		byte(n >> 24), byte(n >> 16),
		byte(n >> 8), byte(n),
	})
}

// WriteLLong writes unsigned little-endian long to buffer.
func (buf *Buffer) WriteLLong(n uint64) error {
	return buf.Write([]byte{
		byte(n), byte(n >> 8),
		byte(n >> 16), byte(n >> 24),
		byte(n >> 32), byte(n >> 40),
		byte(n >> 48), byte(56),
	})
}

// WriteFloat writes 32-bit float to buffer.
func (buf *Buffer) WriteFloat(f float32) error {
	return buf.WriteInt(math.Float32bits(f))
}

// WriteDouble writes 64-bit float to buffer.
func (buf *Buffer) WriteDouble(f float64) error {
	return buf.WriteLong(math.Float64bits(f))
}

// WriteTriad writes unsigned 3-bytes triad to buffer.
func (buf *Buffer) WriteTriad(n uint32) error {
	return buf.Write([]byte{byte(n >> 16), byte(n >> 8), byte(n)})
}

// WriteLTriad writes unsigned little-endian 3-bytes triad to buffer.
func (buf *Buffer) WriteLTriad(n uint32) error {
	return buf.Write([]byte{byte(n), byte(n >> 8), byte(n >> 16)})
}

// WriteString writes string to buffer
func (buf *Buffer) WriteString(s string) (err error) {
	if len(s) > 65535 {
		return util.StringOverflowError{
			Length: len(s),
		}
	}
	if err = buf.WriteShort(uint16(len(s))); err != nil {
		return
	}
	return buf.Write([]byte(s))
}

// WriteAddress writes net.UDPAddr address to buffer.
func (buf *Buffer) WriteAddress(i *net.UDPAddr) (err error) {
	if err = buf.WriteByte(4); err != nil {
		return
	}
	for _, v := range i.IP.To4() {
		if err = buf.WriteByte(v ^ 0xff); err != nil {
			return
		}
	}
	return buf.WriteShort(uint16(i.Port))
}

// Len returns the number of the bytes of the entire buffer.
func (buf *Buffer) Len() int {
	return len(buf.Payload)
}
