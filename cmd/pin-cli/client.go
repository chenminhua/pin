package main



import (
	"bufio"
	"fmt"
	"github.com/chenminhua/pin/internal/protocol"
	"io"
	"log"
	"net"
	"github.com/chenminhua/pin/internal/config"
)

type Client struct {
	conn net.Conn
	conf config.Conf
	reader *bufio.Reader
	writer *bufio.Writer
}

func connect(conf config.Conf) *Client {
	conn, err := net.DialTimeout("tcp", conf.Connect, conf.Timeout)
	if err != nil {
		log.Fatal(fmt.Sprintf("unable to connect %v", conf.Connect))
	}
	reader, writer := bufio.NewReader(conn), bufio.NewWriter(conn)
	return &Client{conn, conf, reader, writer}
}

func connectWithoutTimeout(conf config.Conf) *Client {
	conn, err := net.Dial("tcp", conf.Connect)
	if err != nil {
		log.Fatal(fmt.Sprintf("unable to connect %v", conf.Connect))
	}
	reader, writer := bufio.NewReader(conn), bufio.NewWriter(conn)
	return &Client{conn, conf, reader, writer}
}

func (c *Client) send(header *protocol.Header, content []byte) {
	var err error
	_, err = c.writer.Write(header.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	if content != nil {
		_, err = c.writer.Write(content)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = c.writer.Flush()
	if err != nil {
		log.Fatal(err)
	}
}

func (c *Client) read(contentLength uint32) []byte {
	buf := make([]byte, contentLength)
	_, err := io.ReadFull(c.reader, buf)
	if err != nil {
		log.Fatal(err)
	}
	return buf
}


