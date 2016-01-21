package util

import (
	"bytes"
	"compress/zlib"
	"io"
)

// DecodeDeflate returns decompressed data of given byte slice.
func DecodeDeflate(b []byte) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewBuffer(b))
	if err != nil {
		return make([]byte, 0), err
	}
	output := new(bytes.Buffer)
	io.Copy(output, r)
	r.Close()
	return output.Bytes(), nil
}

// EncodeDeflate returns compressed data of given byte slice.
func EncodeDeflate(b []byte) []byte {
	o := new(bytes.Buffer)
	w := zlib.NewWriter(o)
	w.Write(b)
	w.Close()
	return o.Bytes()
}
