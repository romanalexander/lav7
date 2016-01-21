package util

import (
	"encoding/hex"
	"fmt"
)

// EOFError represents the buffer needs more bytes than needed for read operations.
type EOFError struct {
	BufLen    int
	BufOffset uint32
	Needed    uint32
	Buf       []byte
}

func (err EOFError) Error() string {
	dump := hex.Dump(err.Buf)
	return fmt.Sprintf("EOF: Total %d, Offset %d, Needed %d", err.BufLen, err.BufOffset, err.Needed) + "\n" + dump[:len(dump)-1]
}

// LargeBufferError represents the buffer's size is larger than max capacity defined.
type LargeBufferError struct {
	OldCap int
	Append int
}

func (err LargeBufferError) Error() string {
	return fmt.Sprintf("Buffer too large: Expected capacity %d+%d is bigger than %d", err.OldCap, err.Append, 1024*1024*256)
}

// StringOverflowError represents the given string is too long for write
type StringOverflowError struct {
	Length int
}

func (err StringOverflowError) Error() string {
	return fmt.Sprintf("String too long: Given string is %d characters long, it overflows uint16(65535)", err.Length)
}
