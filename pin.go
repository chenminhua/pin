package main

import (
	// "fmt"
	"flag"
	"time"
)

// Conf - Shared config
type Conf struct {
	Connect        string
	Listen         string
	Timeout        time.Duration
}

const (
	DefaultConfigFile = "~/.pin.toml"
)

func main() {
	isCopy := flag.Bool("copy", false, "copy sth to server")
	isServer := flag.Bool("server", false, "start a server")
	timeout := flag.Uint("timeout", 10, "connection timeout (seconds)")
	flag.Parse()

	var conf Conf
	conf.Connect = "127.0.0.1:7788"
	conf.Listen = "0.0.0.0:7788"
	conf.Timeout = time.Duration(*timeout) * time.Second

	if *isServer {
		RunServer(conf)
	} else {
		if (*isCopy) {
			RunCopy(conf)
		} else {
			RunPaste(conf)
		}

	}
	
}
