package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

var storedContent StoredContent

type StoredContent struct {
	sync.RWMutex
	content []byte
}

func store(reader *bufio.Reader, contentLen uint32) {
	contentBuf := make([]byte, contentLen)
	_,err := io.ReadFull(reader, contentBuf)
	if err != nil {
		log.Fatal(err)
		return
	}
	storedContent.Lock()
	storedContent.content = contentBuf
	storedContent.Unlock()
}

func paste(writer *bufio.Writer) {
	h := Header{1, nil, 'P',
		uint32(len(storedContent.content))}
	writer.Write(h.Bytes())
	writer.Write(storedContent.content)
	writer.Flush()
}

func handleConnection(conf Conf, conn net.Conn) {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	var headerBuf = make([]byte, HeaderSize)
	_, err := io.ReadFull(reader, headerBuf)
	header := GetHeader(headerBuf)
	if err != nil {
		log.Print(err)
		return
	}
	if string(header.Key) != conf.Key {
		fmt.Println("fefwe")
		errMsg := []byte("wrong key")
		h := Header{1, nil, 'E',
			uint32(len(errMsg))}
		writer.Write(h.Bytes())
		writer.Write(errMsg)
		writer.Flush()
		return
	}
	if header.OpCode == byte('C') {
		store(reader, header.ContentLen)
	}
	if header.OpCode == byte('P') {
		paste(writer)
	}

	defer conn.Close()
}

func RunServer(conf Conf) {
	// go handleSignals()
	listen, err := net.Listen("tcp", conf.Listen)
	if err != nil {
		log.Fatal(err)
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
		}
		handleConnection(conf, conn)
	}
}