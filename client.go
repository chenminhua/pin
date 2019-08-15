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

// 将小文件的二进制流上传到服务器
// 默认从file中读，
// 如果file为空，则看有没有传字符串，
// 如果字符串为空则读os.Stdin
func RunCopy(conf Conf, filepath string, str string) {
	client := connect(conf)
	conn, writer := client.conn, client.writer
	defer conn.Close()

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

	contentLength := uint32(len(contentBytes))
	log.Print(contentLength)
	header := Header{1, []byte(conf.Key), byte('C'), contentLength}
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
	h := Header{1, []byte(conf.Key), byte('P'), 0}
	writer.Write(h.Bytes())
	writer.Flush()
	// read
	headerBuf := make([]byte, HeaderSize)
	if _, err := io.ReadFull(reader, headerBuf); err != nil {
		log.Fatal(err)
	}
	header := GetHeader(headerBuf)
	if header.OpCode == 'E' {
		errmsgBuf := make([]byte, header.ContentLen)
		if _, err := io.ReadFull(reader, errmsgBuf); err != nil {
			log.Fatal(err)
		}
		log.Fatal(string(errmsgBuf))
		return
	}
	contentBuf := make([]byte, header.ContentLen)
	if _, err := io.ReadFull(reader, contentBuf); err != nil {
		log.Fatal(err)
	}
	binary.Write(os.Stdout, binary.LittleEndian, contentBuf)
}

func RunPipeCopy(conf Conf) {
	//client := connect(conf)
	//conn, writer, reader := client.conn, client.writer, client.reader
	//defer conn.Close()
	//
	//h := Header{}
}

func RunPipePaste(conf Conf) {

}