package main

import "log"

type Pipe struct {
	sendClient *SClient
	receiveClient *SClient
}

var PIPE_BLOCK_SIZE int64 = 1024 * 1024

func (p *Pipe) checkAndRun() {
	if pipe.sendClient == nil {
		log.Print("sender is not ready")
		return
	}
	if pipe.receiveClient == nil {
		log.Print("receiver is not ready")
		return
	}
	log.Print("sender and receiver both ready")
	// if sender and receiver both ready, told sender to send
	h := PipeCopyOpHeader(pipe.sendClient.conf.Key)
	pipe.sendClient.send(h, nil)
	var buf []byte
	for {
		h1, err := GetHeader(p.sendClient.reader)
		if err != nil {
			log.Print(err)
		}
		buf = pipe.sendClient.read(h1.ContentLen)
		log.Print(len(buf))
		h2 := PipeTransferOpHeader(pipe.receiveClient.conf.Key, int(h1.ContentLen))
		pipe.receiveClient.send(h2, buf)
		if int64(h2.ContentLen) < PIPE_BLOCK_SIZE {
			return
		}
	}

}


