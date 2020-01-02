package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"github.com/chenminhua/pin/internal/protocol"
	"github.com/chenminhua/pin/internal/config"
	"github.com/chenminhua/pin/internal/fs"
)

func RunSender(conf config.Conf, filepath string, str string) {
	if conf.IsPipe {
		RunPipeCopy(conf, filepath)
	} else {
		RunCopy(conf, filepath, str)
	}
}

// 将小文件的二进制流上传到服务器
// 默认从file中读，
// 如果file为空，则看有没有传字符串，
// 如果字符串为空则读os.Stdin
func RunCopy(conf config.Conf, filepath string, str string) {
	client := connect(conf)
	defer client.conn.Close()

	var contentBytes []byte
	var err error

	if filepath != "" {
		// 读文件
		if !fs.FileExists(filepath) {
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

	header := protocol.CopyOpHeader(conf.Key, len(contentBytes))

	client.send(header, contentBytes)
}


func RunPipeCopy(conf config.Conf, filepath string) {
	if filepath == "" {
		log.Fatal("please specify the filepath you want to transfer")
	}
	if !fs.FileExists(filepath) {
		log.Fatal("transfer file ", filepath, " not exists")
	}

	client := connectWithoutTimeout(conf)
	defer client.conn.Close()
	// 发送 pipe copy请求
	client.send(protocol.PipeCopyOpHeader(conf.Key), nil)
	for {
		header, err := protocol.GetHeader(client.reader)
		if err != nil {
			log.Fatal(err)
		}

		if header.OpCode == protocol.ErrReplyCode {
			errMsg := client.read(header.ContentLen)
			log.Fatal(string(errMsg))
			return
		}

		if header.OpCode == 'c' {
			// 表示你可以写了
			file, err := os.Open(filepath)
			if err != nil {
				println(err)
			}
			defer file.Close()
			buf := make([]byte, conf.PipeBlockSize)
			// todo progress
			var offset int64 = 0
			for {
				n, err := file.ReadAt(buf, offset)
				if err != nil && err != io.EOF {
					log.Print("ERROR happened while reading file ",err)
				}
				log.Print("transfer at ", offset / ONE_M_BSIZE, "M")

				if int64(n) < conf.PipeBlockSize || err == io.EOF {
					// send the last transfer packet
					client.send(protocol.PipeTransferLastOpHeader(conf.Key, n), buf)
					return
				} else {
					client.send(protocol.PipeTransferOpHeader(conf.Key, n), buf)
					offset += conf.PipeBlockSize
				}
			}
		}
	}
}
