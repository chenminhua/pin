package main

import (
	"github.com/chenminhua/pin/internal/protocol"
	"log"
)

type Pipe struct {
	sendClient *SClient
	receiveClient *SClient
}

var ONE_M_BSIZE int64 = 1024 * 1024

func (p *Pipe) checkAndRun() {
	if pipe.sendClient == nil {
		log.Print("sender is not ready")
		return
	}
	if pipe.receiveClient == nil {
		log.Print("receiver is not ready")
		return
	}
	p.run()
}

func (p *Pipe) run() {
	sclient, rclient := p.sendClient, p.receiveClient
	log.Print("sender and receiver both ready")
	// told sender to send
	sclient.send(protocol.PipeCopyOpHeader(pipe.sendClient.conf.Key), nil)
	var buf []byte
	for {
		h, err := protocol.GetHeader(p.sendClient.reader)
		if err != nil {
			log.Print(err)
		}
		buf = pipe.sendClient.read(h.ContentLen)
		rclient.send(h, buf)
		if h.OpCode == protocol.PipeTransferLastOpCode {
			sclient.conn.Close()
			rclient.conn.Close()
			p.sendClient = nil
			p.receiveClient = nil
			return
		}

	}
}


