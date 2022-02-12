package crand

import (
	"crypto/rand"
	"encoding/binary"
)

func Uint64() uint64 {
	return binary.BigEndian.Uint64(Bytes(8))
}

func Bytes(count int) []byte {
	buf := make([]byte, count)
	n, err := rand.Read(buf)
	if err != nil {
		panic("crypto random failed: " + err.Error())
	}
	if n != count {
		panic("didn't read enough data to be random")
	}
	return buf
}
