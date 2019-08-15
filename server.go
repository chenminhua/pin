package main

import (
	"bufio"
	"io"
	"log"
	"net"
)

// client connection on server end
type SClient struct {
	conn net.Conn
	conf Conf
	reader *bufio.Reader
	writer *bufio.Writer
}

func NewSClient(conn net.Conn, conf Conf) *SClient {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	return &SClient{conn, conf, reader, writer}
}

func (c *SClient) returnException(msg string) {
	// 1.log
	// 2.return to client
	h := ErrReplyHeader(c.conf.Key, len(msg))
	c.send(h, []byte(msg))
}

func (c *SClient) send(header *Header, content []byte) {
	var err error
	_, err = c.writer.Write(header.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	_, err = c.writer.Write(content)
	if err != nil {
		log.Fatal(err)
	}
	err = c.writer.Flush()
	if err != nil {
		log.Fatal(err)
	}
}

func (c *SClient) read(contentLength uint32) []byte {
	buf := make([]byte, contentLength)
	_, err := io.ReadFull(c.reader, buf)
	if err != nil {
		log.Fatal(err)
	}
	return buf
}

var storedContent StoredContent
var pipe *Pipe = &Pipe{nil, []*SClient{}}

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

func (c *SClient) handleNormalCmd(header *Header) {
	if header.OpCode == byte('C') {
		store(c.reader, header.ContentLen)
	}
	if header.OpCode == byte('P') {
		h := PasteOpHeader(c.conf.Key, len(storedContent.content))
		c.send(h, storedContent.content)
	}
	defer c.conn.Close()

}

func (c *SClient) handlePipeCmd(header *Header) {
	// someone try to paste something from pipe channel
	if header.OpCode == byte('p') {
		// todo thread-safe??
		log.Print("new receiver try to join the pipe")
		pipe.receiveClients = append(pipe.receiveClients, c)
	}
	if header.OpCode == byte('c') {
		log.Print("new sender try to join the pipe")
		if pipe.sendClient == nil {
			pipe.sendClient = c
			log.Print("new sender joined the pipe")
		} else {
			log.Print("new sender failed to join the pipe, pipe occupied by other sender")
			c.returnException("pipe occupied by other sender")
		}
	}
}

func (c *SClient) handle() {

	header, err := GetHeader(c.reader)
	if err != nil {
		c.returnException("Wrong Header")
		return
	}
	if string(header.Key) != c.conf.Key {
		c.returnException("Wrong key")
		return
	}
	if header.OpCode == byte('p') || header.OpCode == byte('c') {
		c.handlePipeCmd(header)
	} else {
		c.handleNormalCmd(header)
	}
}

func RunServer(conf Conf) {
	// go handleSignals()
	listener, err := net.Listen("tcp", conf.Listen)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		client := NewSClient(conn, conf)
		client.handle()
	}
}