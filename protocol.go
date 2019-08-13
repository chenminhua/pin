package main

import "encoding/binary"

const HeaderSize = 38

type Header struct {
	Version    uint8
	Key        []byte
	OpCode     uint8
	ContentLen uint32
}

func (h *Header) Bytes() []byte {
	res := make([]byte, HeaderSize)
	res[0] = byte(h.Version)
	copy(res[1:33], h.Key)
	res[33] = byte(h.OpCode)
	binary.LittleEndian.PutUint32(res[34:38], h.ContentLen)
	return res
}

func GetHeader(b []byte) *Header {
	// todo 检查版本和 header长度
	return &Header{b[0], b[1:33], b[33], binary.LittleEndian.Uint32(b[34:38])}
}