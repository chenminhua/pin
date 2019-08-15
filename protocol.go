package main

import (
	"bufio"
	"encoding/binary"
	"io"
)

const HeaderSize = 38
const ProtocolVersion = 1
const CopyOpCode = byte('C')
const PasteOpCode = byte('P')
const ErrReplyCode = byte('E')
const PipeCopyOpCode = byte('c')
const PipePasteOpCode = byte('p')


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

func GetHeader(reader *bufio.Reader) (*Header, error) {
	// todo 检查版本和 header长度
	var b = make([]byte, HeaderSize)
	_, err := io.ReadFull(reader, b)
	if err != nil {
		return nil, err
	}
	return &Header{b[0], b[1:33], b[33], binary.LittleEndian.Uint32(b[34:38])}, nil
}

// header，用于客户端发起Copy请求的报文
func CopyOpHeader(key string, contentLength int) *Header {
	return &Header{1, []byte(key), CopyOpCode, uint32(contentLength)}
}

// header, 用于客户端发起Paste请求的报文，或者服务端返回数据给客户端来paste时的报文
func PasteOpHeader(key string, contentLength int) *Header {
	return &Header{ProtocolVersion, []byte(key), PasteOpCode, uint32(contentLength)}
}

// 服务端遇到错误，返回给客户端的报文的header
func ErrReplyHeader(key string, contentLength int) *Header{
	return &Header{ProtocolVersion, []byte(key), ErrReplyCode, uint32(contentLength)}
}

func PipePasteOpHeader(key string) *Header {
	return &Header{ProtocolVersion, []byte(key), PipePasteOpCode, 0}
}

func PipeCopyOpHeader(key string) *Header {
	return &Header{ProtocolVersion, []byte(key), PipeCopyOpCode, 0}
}