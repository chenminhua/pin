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
	reader, writer := bufio.NewReader(conn), bufio.NewWriter(conn)
	return &Client{conn, conf, reader, writer}
}

func connectWithoutTimeout(conf Conf) *Client {
	conn, err := net.Dial("tcp", conf.Connect)
	if err != nil {
		log.Fatal(fmt.Sprintf("unable to connect %v", conf.Connect))
	}
	reader, writer := bufio.NewReader(conn), bufio.NewWriter(conn)
	return &Client{conn, conf, reader, writer}
}

func (c *Client) send(header *Header, content []byte) {
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
		if !FileExists(filepath) {
			log.Fatal("file not exist")
		}
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
	defer client.conn.Close()

	// 发送paste请求
	client.send(PasteOpHeader(conf.Key, 0), nil)

	header, err := GetHeader(client.reader)
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

func RunPipeCopy(conf Conf, filepath string) {
	client := connectWithoutTimeout(conf)
	defer client.conn.Close()
	// 发送 pipe copy请求
	client.send(PipeCopyOpHeader(conf.Key), nil)
	for {
		header, err := GetHeader(client.reader)
		if err != nil {
			log.Fatal(err)
		}

		if header.OpCode == 'E' {
			errMsg := client.read(header.ContentLen)
			log.Fatal(string(errMsg))
			return
		}

		if header.OpCode == 'w' {
			log.Print("waiting for the receiver")
			// 表示现在没有receiver，你需要等待
		}

		if header.OpCode == 'c' {
			// 表示你可以写了
			file, err := os.Open(filepath)
			if err != nil {
				println(err)
			}
			defer file.Close()
			buf := make([]byte, PIPE_BLOCK_SIZE)
			// todo progress
			var offset int64 = 0
			for {
				n, err := file.ReadAt(buf, offset)
				if err != nil && err != io.EOF {
					log.Print("ERROR happened while reading file ",err)
				}

				client.send(PipeTransferOpHeader(conf.Key, n), buf)
				offset += PIPE_BLOCK_SIZE
				if int64(n) < PIPE_BLOCK_SIZE {
					return
				}
			}
		}
	}
}

func RunPipePaste(conf Conf, filepath string) {
	client := connectWithoutTimeout(conf)
	defer client.conn.Close()
	client.send(PipePasteOpHeader(conf.Key), nil)

	// 发送 pipe paste请求
	client.send(PipePasteOpHeader(conf.Key), nil)

	var offset int64 = 0

	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		println(err)
	}
	defer file.Close()

	for {
		header, err := GetHeader(client.reader)
		if err != nil {
			log.Fatal(err)
		}

		if header.OpCode == 'E' {
			errMsg := client.read(header.ContentLen)
			log.Fatal(string(errMsg))
			return
		}

		if header.OpCode == 't' {
			// 表示你可以写了
			log.Print(header.ContentLen)
			buf := client.read(header.ContentLen)
			_, err = file.WriteAt(buf, offset)
			offset += int64(header.ContentLen)
			if err != nil {
				log.Print(err)
			}
			if int64(header.ContentLen) < PIPE_BLOCK_SIZE {
				log.Print("return")
				return
			}
		}
	}
}