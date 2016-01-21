package raknet

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net"
	"testing"

	"github.com/L7-MCPE/lav7/util/buffer"
)

func TestFieldChecks(t *testing.T) {
	tests := []struct {
		Map    map[string]interface{}
		Fields []string
		Expect bool
	}{
		{
			map[string]interface{}{
				"aa": 32,
			},
			[]string{"aa"},
			true,
		},
		{
			map[string]interface{}{
				"aa": 32,
			},
			[]string{"aa", "bb"},
			false,
		},
		{
			map[string]interface{}{
				"aa": 32,
				"az": false,
			},
			[]string{"aa"},
			false,
		},
		{
			map[string]interface{}{
				"a3":  32,
				"baz": "hello!",
			},
			[]string{"a3", "baz"},
			true,
		},
		{
			map[string]interface{}{
				"a3":  32,
				"baz": "hello!",
			},
			[]string{"a3", "foo"},
			false,
		},
	}
	for _, test := range tests {
		if test.Expect != checkFields(test.Map, test.Fields...) {
			t.Error("checkFields test failed:", test)
			return
		}
	}
}

func TestACK(t *testing.T) {
	cases := []struct {
		Base64 string
		Expect []uint32
	}{
		{"wAABAAEAAAYAAA==", []uint32{1, 2, 3, 4, 5, 6}},
		{"wAAEASAAAAGEAAAB1gAAAUQBAA==", []uint32{32, 132, 214, 324}},
		{"wAACAAEAAAMAAAEGAAA=", []uint32{1, 2, 3, 6}},
		{"wAABAAEAAAcAAA==", []uint32{1, 2, 3, 4, 5, 6, 7}},
		{"wAAA", []uint32{}},
	}

	for _, v := range cases {
		b, err := base64.StdEncoding.DecodeString(v.Base64)
		if err != nil {
			t.Error("Error while decoding base64 payload:", err)
			return
		}
		b = b[1:]
		var result []uint32
		result, err = DecodeAck(buffer.FromBytes(b))
		if err != nil {
			t.Error("Error while decoding ACK:", err)
			return
		}
		if fmt.Sprint(v.Expect) != fmt.Sprint(result) {
			t.Errorf("ACK decoding result mismatch: \n%v\n%v\ninput: %v", v.Expect, result, b)
		}
		var rb *buffer.Buffer
		rb, err = EncodeAck(v.Expect)
		if err != nil {
			t.Error("Error while encoding ACK:", err)
			return
		}
		if !bytes.Equal(rb.Payload, b) {
			t.Errorf("ACK encoding result mismatch: \n%v\n%v\ninput: %v", b, rb.Payload, v.Expect)
		}
	}
}

func packetB64(b64 string) Packet {
	b, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		panic(fmt.Sprint("Error while decoding base64 payload:", err))
	}
	return Packet{buffer.FromBytes(b[1:]), new(net.UDPAddr)}
}

func TestOCR1(t *testing.T) {
	tests := []struct {
		Base64   string
		Protocol byte
		MtuSize  int
	}{
		{"BQD//wD+/v7+/f39/RI0VnggAAAAAAAAAAAAAAAAAA==", 32, 31},
		{"BQD//wD+/v7+/f39/RI0VngCAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", 2, 42},
		{"BQD//wD+/v7+/f39/RI0VngAAA==", 0, 19},
	}
	handler := new(openConnectionRequest1)

	for _, test := range tests {
		f, err := handler.Read(packetB64(test.Base64).Buffer)
		if err != nil {
			t.Error("Error while reading packet:", err)
			return
		}
		if f["mtuSize"].(int) != test.MtuSize {
			t.Error("MtuSize mismatch:", f["mtuSize"].(int), test.MtuSize)
			return
		}
		if f["protocol"].(byte) != test.Protocol {
			t.Error("Protocol mismatch:", f["protocol"].(byte), test.Protocol)
			return
		}
		if pk, err := handler.Write(f); err != nil {
			t.Error("Error while writing packet:", err)
			return
		} else if !bytes.Equal(pk.Done()[1:], packetB64(test.Base64).Done()) {
			t.Errorf("Write test failed:\n%s\n\n%s", hex.Dump(pk.Done()[1:]), hex.Dump(packetB64(test.Base64).Done()))
			return
		}
	}
}
