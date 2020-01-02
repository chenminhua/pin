package main

import (
	"encoding/binary"
	"github.com/chenminhua/pin/internal/fs"
	"github.com/chenminhua/pin/internal/protocol"
	"github.com/chenminhua/pin/internal/config"
	"log"
	"os"
)

func RunReceiver(conf config.Conf, filepath string) {
	if conf.IsPipe {
		RunPipePaste(conf, filepath)
	} else {
		RunPaste(conf)
	}
}

func RunPaste(conf config.Conf) {
	client := connect(conf)
	defer client.conn.Close()

	// 发送paste请求
	client.send(protocol.PasteOpHeader(conf.Key, 0), nil)

	header, err := protocol.GetHeader(client.reader)
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


func RunPipePaste(conf config.Conf, filepath string) {

	if filepath == "" {
		log.Fatal("please specify the filepath you want to transfer")
	}
	if fs.FileExists(filepath) {
		log.Fatal("transfer file ", filepath, " already exists")
	}

	client := connectWithoutTimeout(conf)
	defer client.conn.Close()
	client.send(protocol.PipePasteOpHeader(conf.Key), nil)

	// 发送 pipe paste请求
	client.send(protocol.PipePasteOpHeader(conf.Key), nil)

	var offset int64 = 0

	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		println(err)
	}
	defer file.Close()

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

		if header.OpCode == protocol.PipeTransferOpCode || header.OpCode == protocol.PipeTransferLastOpCode {
			// 表示你可以写了
			log.Print("transfer at ", offset / ONE_M_BSIZE, "M")
			buf := client.read(header.ContentLen)
			_, err = file.WriteAt(buf, offset)
			offset += int64(header.ContentLen)
			if err != nil {
				log.Print(err)
			}
			// 收到最后一个包退出
			if header.OpCode == protocol.PipeTransferLastOpCode {
				log.Print("transfer finished")
				return
			}
		}
	}
}

