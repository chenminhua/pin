package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"
)

type Client struct {
	conn net.Conn
	conf Conf
	reader *bufio.Reader
	writer *bufio.Writer
}

func connect(conf Conf) *Client {
	conn, err := net.DialTimeout("tcp", conf.Connect, conf.Timeout)
	if err != nil {
		log.Fatal(fmt.Sprintf("unable to connect %v", conf.Connect))
	}
	conn.SetDeadline(time.Now().Add(conf.Timeout))
	reader, writer := bufio.NewReader(conn), bufio.NewWriter(conn)
	return &Client{conn, conf, reader, writer}
}

func (c *Client) send(header *Header, content []byte) {
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

func (c *Client) read(contentLength uint32) []byte {
	buf := make([]byte, contentLength)
	_, err := io.ReadFull(c.reader, buf)
	if err != nil {
		log.Fatal(err)
	}
	return buf
}

// 将小文件的二进制流上传到服务器
// 默认从file中读，
// 如果file为空，则看有没有传字符串，
// 如果字符串为空则读os.Stdin
func RunCopy(conf Conf, filepath string, str string) {
	client := connect(conf)
	defer client.conn.Close()

	var contentBytes []byte
	var err error

	if filepath != "" {
		// 读文件
		contentBytes, err = ioutil.ReadFile(filepath)
		if err != nil {
			log.Fatal(err)
		}
	} else if str != "" {
		// 读字符串
		contentBytes = []byte(str)
	} else {
		// 从os.Stdin读数据
		var contentBuffer bytes.Buffer
		_, err = contentBuffer.ReadFrom(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		contentBytes = contentBuffer.Bytes()
	}

	header := CopyOpHeader(conf.Key, len(contentBytes))

	client.send(header, contentBytes)
}


func RunPaste(conf Conf) {
	client := connect(conf)
	conn, writer, reader := client.conn, client.writer, client.reader
	defer conn.Close()

	// 发送paste请求
	h := PasteOpHeader(conf.Key, 0)
	writer.Write(h.Bytes())
	writer.Flush()

	header, err := GetHeader(reader)
	if err != nil {
		log.Fatal(err)
	}
	if header.OpCode == 'E' {
		errMsg := client.read(header.ContentLen)
		log.Fatal(string(errMsg))
		return
	}
	content := client.read(header.ContentLen)
	binary.Write(os.Stdout, binary.LittleEndian, content)
}

func RunPipeCopy(conf Conf) {
	//client := connect(conf)
	//conn, writer, reader := client.conn, client.writer, client.reader
	//defer conn.Close()
	//
	//h := PipePasteOpHeader(conf.Key)
}

func RunPipePaste(conf Conf) {

}