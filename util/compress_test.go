package util

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestCompress(t *testing.T) {
	bufs := [][]byte{
		[]byte("azerdfaksjdfkljl"),
		[]byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"),
		[]byte("\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0a\x0b\x0c\x0d\x0e\x0f"),
	}

	for _, buf := range bufs {
		c := EncodeDeflate(buf)
		if d, err := DecodeDeflate(c); err != nil {
			t.Error("Error while decoding DEFLATE:", err)
			return
		} else {
			if !bytes.Equal(d, buf) {
				t.Errorf("Decoded result mismatch!\n%s\n%s", hex.Dump(buf), hex.Dump(d))
				return
			}
		}
	}
}
