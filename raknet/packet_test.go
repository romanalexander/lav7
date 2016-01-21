package raknet

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/L7-MCPE/lav7/util/buffer"
)

func TestRequire(t *testing.T) {
	pk := buffer.FromBytes([]byte("\x00\x00\x00\x00\x00\x00"))
	if !pk.Require(6) {
		t.Errorf("Require test failed: payload 6, read 0, required 6, return value=false")
		return
	}
	pk.Read(4)
	if !pk.Require(2) {
		t.Errorf("Require test failed: payload 6, read 4, required 2, return value=false")
		return
	}
	if pk.Require(3) {
		t.Errorf("Require test failed: payload 6, read 4, required 3, return value=true")
		return
	}
	pk.Read(2)
	if pk.Require(1) {
		t.Errorf("Require test failed: payload 6, read 6, required 1, return value=true")
	}
}

type ReadCase struct {
	Total      uint64
	ReadBefore uint64
	ReadAfter  uint64
	ShouldErr  bool
}

func TestRead(t *testing.T) {
	cases := []ReadCase{
		{
			Total:     32,
			ReadAfter: 32,
		},
		{
			Total:      32,
			ReadBefore: 32,
			ReadAfter:  1,
			ShouldErr:  true,
		},
		{
			Total:     32,
			ReadAfter: 33,
			ShouldErr: true,
		},
	}
	for _, c := range cases {
		pk := buffer.FromBytes(make([]byte, c.Total))
		if pk.Head() != 0 {
			t.Error("Head test failed: Expected 0, got", pk.Head())
			return
		}
		pk.Read(c.ReadBefore)
		if _, err := pk.Read(c.ReadAfter); (err != nil) != c.ShouldErr {
			t.Error("Read test failed:", c, "err exists:", !c.ShouldErr)
			return
		}
	}
	pk := buffer.FromBytes([]byte("\xe8\x0f\x0d\xfd\x3f\xdd\xdd\x00\x0a\x00\xfd\xff\xfd\x00\x00\x64\x01\x04"))
	if n, err := pk.ReadByte(); err != nil || n != 232 {
		t.Error("ReadByte test failed:", pk, "Result:", n, "Expected: 232", "Error exists: ", err != nil)
	}
	if n, err := pk.ReadShort(); err != nil || n != 3853 {
		t.Error("ReadShort test failed:", pk, "Result:", n, "Expected: 3853", "Error exists: ", err != nil)
	}
	if n, err := pk.ReadInt(); err != nil || n != 4248821213 {
		t.Error("ReadInt test failed:", pk, "Result:", n, "Expected: 4248821213", "Error exists: ", err != nil)
	}
	if n, err := pk.ReadLong(); err != nil || n != 2815840688603136 {
		t.Error("ReadLong test failed:", pk, "Result:", n, "Expected: 2815840688603136", "Error exists: ", err != nil)
	}
	if n, err := pk.ReadLTriad(); err != nil || n != 262500 {
		t.Error("ReadLTriad test failed:", pk, "Result:", n, "Expected: 262500", "Error exists: ", err != nil)
	}
	pk = buffer.FromBytes(append([]byte{0x00, 0x0d}, []byte("Hello, 世界")...))
	if s, err := pk.ReadString(); err != nil || s != "Hello, 世界" {
		t.Error("ReadString tets failed:", pk, "Result:", s, "Expected: Hello, 世界", "Error exists: ", err != nil)
	}
}

func TestWrite(t *testing.T) {
	pk := new(buffer.Buffer)
	pk.WriteByte(4)
	pk.WriteShort(523)
	pk.WriteInt(153925)
	pk.WriteLong(539528483653)
	pk.WriteString("Hello, 世界")
	pk.WriteAddress(&net.UDPAddr{
		IP:   []byte{0x7f, 0x00, 0x00, 0x01},
		Port: 19132,
	})
	pk.WriteLTriad(564365)
	if b, err := pk.ReadByte(); err != nil || b != 4 {
		t.Error("Test failed: expected 4, got", b, "Error exists:", err != nil)
	}
	if b, err := pk.ReadShort(); err != nil || b != 523 {
		t.Error("Test failed: expected 523, got", b, "Error exists:", err != nil)
	}
	if b, err := pk.ReadInt(); err != nil || b != 153925 {
		t.Error("Test failed: expected 153925, got", b, "Error exists:", err != nil)
	}
	if b, err := pk.ReadLong(); err != nil || b != 539528483653 {
		t.Error("Test failed: expected 539528483653, got", b, "Error exists:", err != nil)
	}
	if s, err := pk.ReadString(); err != nil || s != "Hello, 世界" {
		t.Error("Test failed: expected Hello, 世界, got", s, "Error exists:", err != nil)
	}
	if a, err := pk.ReadAddress(); err != nil || a.String() != "127.0.0.1:19132" {
		t.Error("Test failed: expected 127.0.0.1:19132, got", a.String(), "Error exists:", err != nil)
	}
	if b, err := pk.ReadLTriad(); err != nil || b != 564365 {
		t.Error("Test failed: expected 564365, got", b, "Error exists:", err != nil)
	}
}

func TestEncapsulated(t *testing.T) {
	tests := []struct {
		Base64 string
		Length int
	}{
		{"kACQBAAAIAAAAgAAAAMACAAAAAE1MTUxNTE1MTUxNTE1MTUxNTE=", 38},
		{"MACQIAAAAgAAAAMACAAAAAE1MTUxNTE1MTUxNTE1MTUxNTE=", 35},
	}
	for _, test := range tests {
		var ep *EncapsulatedPacket
		var err error
		b, err := base64.StdEncoding.DecodeString(test.Base64)
		if err != nil {
			panic(fmt.Sprint("Error while decoding base64 payload:", err))
		}
		if ep, err = NewEncapsulated(buffer.FromBytes(b)); err != nil {
			t.Error("Error while creating new EncapsulatedPacket:", err)
			return
		}
		if ep.TotalLen() != test.Length {
			t.Error("EncapsulatedPacket length test failed:", ep.TotalLen(), "!=", test.Length, ep)
			return
		}
		var buf *buffer.Buffer
		if buf, err = ep.Bytes(); err != nil {
			t.Error("Error while encoding EncapsulatedPacket:", err)
			return
		}
		if !bytes.Equal(buf.Payload, b) {
			t.Error("EncapsulatedPacket test failed: mismatch after encode/decode")
			return
		}
	}
}

func TestDataPacket(t *testing.T) {
	dp := new(DataPacket)
	dp.Head = 4
	dp.SeqNumber = 3
	dp.SendTime = *new(time.Time)
	dp.Packets = make([]*EncapsulatedPacket, 0)
	b, _ := base64.StdEncoding.DecodeString("kACQBAAAIAAAAgAAAAMACAAAAAE1MTUxNTE1MTUxNTE1MTUxNTE=")
	ep, _ := NewEncapsulated(buffer.FromBytes(b))
	dp.Packets = append(dp.Packets, ep)
	b, _ = base64.StdEncoding.DecodeString("MACQIAAAAgAAAAMACAAAAAE1MTUxNTE1MTUxNTE1MTUxNTE=")
	ep, _ = NewEncapsulated(buffer.FromBytes(b))
	dp.Packets = append(dp.Packets, ep)
	if err := dp.Encode(); err != nil {
		t.Error("Error while encoding DataPacket:", err)
		return
	}
	dp = &DataPacket{
		Buffer:    buffer.FromBytes(dp.Done()),
		SeqNumber: 0,
		SendTime:  *new(time.Time),
		Packets:   make([]*EncapsulatedPacket, 0),
	}
	if err := dp.Decode(); err != nil {
		t.Error("Error while decoding DataPacket:", err)
		return
	}
	bb, _ := dp.Packets[1].Bytes()
	if !bytes.Equal(bb.Done(), b) {
		t.Error("DataPacket encode/decode mismatch!")
		return
	}
}
