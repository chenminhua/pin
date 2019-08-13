package main

import (
	"bufio"
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
		log.Print(err)
		return
	}
	storedContent.Lock()
	storedContent.content = contentBuf
	storedContent.Unlock()
}

func paste(writer *bufio.Writer) {
	writer.Write(storedContent.content)
	writer.Flush()
}

func handleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	var headerBuf = make([]byte, HeaderSize)
	_, err := io.ReadFull(reader, headerBuf)
	header := GetHeader(headerBuf)
	if err != nil {
		log.Print(err)
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
		handleConnection(conn)
	}
}