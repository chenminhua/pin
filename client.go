package main

import (
	"encoding/binary"
	"io"
	"net"
	"log"
	"fmt"
	"time"
	"bytes"
	"os"
	"bufio"
)

type ClientConnection struct {
	conn net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

func connect(conf Conf) *ClientConnection {
	conn, err := net.DialTimeout("tcp", conf.Connect, conf.Timeout)
	if err != nil {
		log.Fatal(fmt.Sprintf("unable to connect %v", conf.Connect))
	}
	conn.SetDeadline(time.Now().Add(conf.Timeout))
	reader, writer := bufio.NewReader(conn), bufio.NewWriter(conn)
	return &ClientConnection{conn, reader, writer}
}



func RunCopy(conf Conf) {
	client := connect(conf)
	conn, writer := client.conn, client.writer
	defer conn.Close()

	var content bytes.Buffer
	_, err := content.ReadFrom(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	contentBytes := content.Bytes()
	contentLength := uint32(len(contentBytes))
	log.Print(contentLength)
	header := Header{1, nil, byte('C'), contentLength}
	writer.Write(header.Bytes())
	writer.Write(contentBytes)
	// 写入socket
	err = client.writer.Flush()
	if err != nil {
		log.Fatal(err)
	}
}


func RunPaste(conf Conf) {
	client := connect(conf)
	conn, writer, reader := client.conn, client.writer, client.reader
	defer conn.Close()

	// 发送paste请求
	h := Header{1, nil, byte('P'), 0}
	writer.Write(h.Bytes())
	writer.Flush()
	// read
	headerBuf := make([]byte, HeaderSize)
	if _, err := io.ReadFull(reader, headerBuf); err != nil {
		log.Fatal(err)
	}
	header := GetHeader(headerBuf)
	contentBuf := make([]byte, header.ContentLen)
	if _, err := io.ReadFull(reader, contentBuf); err != nil {
		log.Fatal(err)
	}
	binary.Write(os.Stdout, binary.LittleEndian, contentBuf)
}