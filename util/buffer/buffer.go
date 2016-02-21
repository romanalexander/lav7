// Package buffer provides simple functions for processing with go's internal bytes.Buffer struct.
package buffer

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math"
	"net"
)

// Overflow is an error indicates the reader could not read as you requested.
type Overflow struct {
	Need   int
	Got    int
	Buffer *bytes.Buffer
}

// Error implements the error interface.
func (e Overflow) Error() string {
	return fmt.Sprintf("Overflow: Needed %d, got %d", e.Need, e.Got)
}

// StringOverflowError represents the given string is too long for write
type StringOverflow struct {
	Length int
}

// Error implements the error interface.
func (err StringOverflow) Error() string {
	return fmt.Sprintf("String too long: Given string is %d characters long, it overflows uint16(65535)", err.Length)
}

// Read reads n bytes of data from buf. If buf returns smaller slice than n, returns OverFlow.
func Read(buf *bytes.Buffer, n int) (b []byte, err error) {
	if b = buf.Next(n); n != len(b) {
		err = Overflow{
			Need:   n,
			Got:    len(b),
			Buffer: buf,
		}
		return
	}
	return
}

// ReadAny reads appropriate type from given reference value.
func ReadAny(buf *bytes.Buffer, p interface{}) {
	switch p.(type) {
	case *bool:
		*p.(*bool) = ReadBool(buf)
	case *byte:
		*p.(*byte) = ReadByte(buf)
	case *uint16:
		*p.(*uint16) = ReadShort(buf)
	case *uint32:
		*p.(*uint32) = ReadInt(buf)
	case *uint64:
		*p.(*uint64) = ReadLong(buf)
	case *float32:
		*p.(*float32) = ReadFloat(buf)
	case *float64:
		*p.(*float64) = ReadDouble(buf)
	case *string:
		*p.(*string) = ReadString(buf)
	case *net.UDPAddr:
		var addr *net.UDPAddr
		addr = ReadAddress(buf)
		*p.(*net.UDPAddr) = *addr
	case **net.UDPAddr:
		*p.(**net.UDPAddr) = ReadAddress(buf)
	case byte, uint16, uint32,
		uint64, float32, float64, string, net.UDPAddr:
		panic("ReadAny requires reference type")
	default:
		panic("Unsupported type for ReadAny")
	}
}

// BatchRead batches ReadAny from given reference pointers.
func BatchRead(buf *bytes.Buffer, p ...interface{}) {
	for _, pp := range p {
		ReadAny(buf, pp)
	}
}

// ReadBool reads boolean from buffer.
func ReadBool(buf *bytes.Buffer) bool {
	b, err := Read(buf, 1)
	if err != nil {
		panic(err)
	}
	return b[0] > 0
}

// ReadByte reads unsigned byte from buffer.
func ReadByte(buf *bytes.Buffer) byte {
	b, err := Read(buf, 1)
	if err != nil {
		panic(err)
	}
	return b[0]
}

// ReadShort reads unsigned short from buffer.
func ReadShort(buf *bytes.Buffer) uint16 {
	b, err := Read(buf, 2)
	if err != nil {
		panic(err)
	}
	return uint16(b[0])<<8 | uint16(b[1])
}

// ReadLShort reads unsigned little-endian short from buffer.
func ReadLShort(buf *bytes.Buffer) uint16 {
	b, err := Read(buf, 2)
	if err != nil {
		panic(err)
	}
	return uint16(b[1])<<8 | uint16(b[0])
}

// ReadInt reads unsigned int from buffer.
func ReadInt(buf *bytes.Buffer) uint32 {
	b, err := Read(buf, 4)
	if err != nil {
		panic(err)
	}
	return uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
}

// ReadLInt reads unsigned little-endian int from buffer.
func ReadLInt(buf *bytes.Buffer) uint32 {
	b, err := Read(buf, 4)
	if err != nil {
		panic(err)
	}
	return uint32(b[3])<<24 | uint32(b[2])<<16 | uint32(b[1])<<8 | uint32(b[0])
}

// ReadLong reads unsigned long from buffer.
func ReadLong(buf *bytes.Buffer) uint64 {
	b, err := Read(buf, 8)
	if err != nil {
		panic(err)
	}
	return uint64(b[0])<<56 | uint64(b[1])<<48 |
		uint64(b[2])<<40 | uint64(b[3])<<32 |
		uint64(b[4])<<24 | uint64(b[5])<<16 |
		uint64(b[6])<<8 | uint64(b[7])
}

// ReadLLong reads unsigned little-endian long from buffer.
func ReadLLong(buf *bytes.Buffer) uint64 {
	b, err := Read(buf, 8)
	if err != nil {
		panic(err)
	}
	return uint64(b[7])<<56 | uint64(b[6])<<48 |
		uint64(b[5])<<40 | uint64(b[4])<<32 |
		uint64(b[3])<<24 | uint64(b[2])<<16 |
		uint64(b[1])<<8 | uint64(b[0])
}

// ReadFloat reads 32-bit float from buffer.
func ReadFloat(buf *bytes.Buffer) float32 {
	r := ReadInt(buf)
	return math.Float32frombits(r)
}

// ReadDouble reads 64-bit float from buffer.
func ReadDouble(buf *bytes.Buffer) float64 {
	r := ReadLong(buf)
	return math.Float64frombits(r)
}

// ReadTriad reads unsigned 3-bytes triad from buffer.
func ReadTriad(buf *bytes.Buffer) uint32 {
	b, err := Read(buf, 3)
	if err != nil {
		panic(err)
	}
	return uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2])
}

// ReadLTriad reads unsigned little-endian 3-bytes triad from buffer.
func ReadLTriad(buf *bytes.Buffer) uint32 {
	b, err := Read(buf, 3)
	if err != nil {
		panic(err)
	}
	return uint32(b[2])<<16 | uint32(b[1])<<8 | uint32(b[0])
}

// ReadString reads string from buffer.
func ReadString(buf *bytes.Buffer) (str string) {
	b, err := Read(buf, int(ReadShort(buf)))
	if err != nil {
		panic(err)
	}
	return string(b)
}

// ReadAddress reads IP address/port from buffer.
func ReadAddress(buf *bytes.Buffer) (addr *net.UDPAddr) {
	v := ReadByte(buf)
	if v != 4 {
		panic(fmt.Sprintf("ReadAddress got unsupported IP version %d", v))
	}
	b, err := Read(buf, 4)
	if err != nil {
		panic(err)
	}
	p := ReadShort(buf)
	return &net.UDPAddr{
		IP:   append([]byte{b[0] ^ 0xff}, b[1]^0xff, b[2]^0xff, b[3]^0xff),
		Port: int(p),
	}
}

// Write writes given byte array to buffer.
func Write(buf *bytes.Buffer, b []byte) error {
	n, err := buf.Write(b)
	if err == nil && n != len(b) {
		err = Overflow{
			Need:   len(b),
			Got:    n,
			Buffer: buf,
		}
	}
	return err
}

// WriteAny writes appropriate type from given interface{} value to buffer.
func WriteAny(buf *bytes.Buffer, p interface{}) {
	switch p.(type) {
	case bool:
		WriteBool(buf, p.(bool))
	case byte:
		WriteByte(buf, p.(byte))
	case uint16:
		WriteShort(buf, p.(uint16))
	case uint32:
		WriteInt(buf, p.(uint32))
	case uint64:
		WriteLong(buf, p.(uint64))
	case float32:
		WriteFloat(buf, p.(float32))
	case float64:
		WriteDouble(buf, p.(float64))
	case string:
		WriteString(buf, p.(string))
	case []byte:
		Write(buf, p.([]byte))
	case *bool:
		WriteBool(buf, *p.(*bool))
	case *byte:
		WriteByte(buf, *p.(*byte))
	case *uint16:
		WriteShort(buf, *p.(*uint16))
	case *uint32:
		WriteInt(buf, *p.(*uint32))
	case *uint64:
		WriteLong(buf, *p.(*uint64))
	case *float32:
		WriteFloat(buf, *p.(*float32))
	case *float64:
		WriteDouble(buf, *p.(*float64))
	case *string:
		WriteString(buf, *p.(*string))
	case *[]byte:
		Write(buf, *p.(*[]byte))
	case *net.UDPAddr:
		WriteAddress(buf, p.(*net.UDPAddr))
	}
}

// BatchWrite batches WriteAny from given values.
func BatchWrite(buf *bytes.Buffer, p ...interface{}) {
	for _, pp := range p {
		WriteAny(buf, pp)
	}
}

// WriteBool writes boolean to buffer.
func WriteBool(buf *bytes.Buffer, n bool) {
	WriteByte(buf, func() byte {
		if n {
			return 1
		} else {
			return 0
		}
	}())
}

// WriteByte writes unsigned byte to buffer.
func WriteByte(buf *bytes.Buffer, n byte) {
	if err := Write(buf, []byte{n}); err != nil {
		panic(err)
	}
}

// WriteShort writes unsigned short to buffer.
func WriteShort(buf *bytes.Buffer, n uint16) {
	if err := Write(buf, []byte{byte(n >> 8), byte(n)}); err != nil {
		panic(err)
	}
}

// WriteLShort writes unsigned little-endian short to buffer.
func WriteLShort(buf *bytes.Buffer, n uint16) {
	if err := Write(buf, []byte{byte(n), byte(n >> 8)}); err != nil {
		panic(err)
	}
}

// WriteInt writes unsigned int to buffer.
func WriteInt(buf *bytes.Buffer, n uint32) {
	if err := Write(buf, []byte{byte(n >> 24), byte(n >> 16), byte(n >> 8), byte(n)}); err != nil {
		panic(err)
	}
}

// WriteLInt writes unsigned little-endian int to buffer.
func WriteLInt(buf *bytes.Buffer, n uint32) {
	if err := Write(buf, []byte{byte(n), byte(n >> 8), byte(n >> 16), byte(n >> 24)}); err != nil {
		panic(err)
	}
}

// WriteLong writes unsigned long to buffer.
func WriteLong(buf *bytes.Buffer, n uint64) {
	if err := Write(buf, []byte{
		byte(n >> 56), byte(n >> 48),
		byte(n >> 40), byte(n >> 32),
		byte(n >> 24), byte(n >> 16),
		byte(n >> 8), byte(n),
	}); err != nil {
		panic(err)
	}
}

// WriteLLong writes unsigned little-endian long to buffer.
func WriteLLong(buf *bytes.Buffer, n uint64) {
	if err := Write(buf, []byte{
		byte(n), byte(n >> 8),
		byte(n >> 16), byte(n >> 24),
		byte(n >> 32), byte(n >> 40),
		byte(n >> 48), byte(56),
	}); err != nil {
		panic(err)
	}
}

// WriteFloat writes 32-bit float to buffer.
func WriteFloat(buf *bytes.Buffer, f float32) {
	WriteInt(buf, math.Float32bits(f))
}

// WriteDouble writes 64-bit float to buffer.
func WriteDouble(buf *bytes.Buffer, f float64) {
	WriteLong(buf, math.Float64bits(f))
}

// WriteTriad writes unsigned 3-bytes triad to buffer.
func WriteTriad(buf *bytes.Buffer, n uint32) {
	if err := Write(buf, []byte{byte(n >> 16), byte(n >> 8), byte(n)}); err != nil {
		panic(err)
	}
}

// WriteLTriad writes unsigned little-endian 3-bytes triad to buffer.
func WriteLTriad(buf *bytes.Buffer, n uint32) error {
	return Write(buf, []byte{byte(n), byte(n >> 8), byte(n >> 16)})
}

// WriteString writes string to buffer.
func WriteString(buf *bytes.Buffer, s string) {
	if len(s) > 65535 {
		panic(StringOverflow{
			Length: len(s),
		})
	}
	WriteShort(buf, uint16(len(s)))
	Write(buf, []byte(s))
}

// WriteAddress writes net.UDPAddr address to buffer.
func WriteAddress(buf *bytes.Buffer, i *net.UDPAddr) {
	WriteByte(buf, 4)
	for _, v := range i.IP.To4() {
		WriteByte(buf, v^0xff)
	}
	WriteShort(buf, uint16(i.Port))
}

// Dump prints hexdump for given buffer.
func Dump(buf *bytes.Buffer) {
	fmt.Print(hex.Dump(buf.Bytes()))
}
