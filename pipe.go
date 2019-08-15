package main

import "log"

type Pipe struct {
	sendClient *SClient
	receiveClients []*SClient
}

func (p *Pipe) checkAndRun() {
	if (pipe.sendClient == nil) {
		log.Print("sender is not ready")
		return
	}
	if (len(pipe.receiveClients) == 0) {
		log.Print("receiver is not ready")
		return
	}

}


