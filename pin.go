package main

import (
	// "fmt"
	"flag"
	"github.com/mitchellh/go-homedir"
	"log"
	"time"
)


// Conf - Shared config
type Conf struct {
	Connect        string
	Listen         string
	Key            string
	Timeout        time.Duration
}

const (
	DefaultConfigFile = "~/.pin.toml"
)


func expandConfigFile(path string) string {
	file, err := homedir.Expand(path)
	if err != nil {
		log.Fatal(err)
	}
	return file
}


func main() {
	isCopy := flag.Bool("copy", false, "copy sth to server")
	isPipe := flag.Bool("pipe", false, "pipe")

	filepath := flag.String("f", "", "file")
	str := flag.String("s", "", "string")
	isServer := flag.Bool("server", false, "start a server")
	timeout := flag.Uint("timeout", 10, "connection timeout (seconds)")
	configFile := flag.String("config", DefaultConfigFile, "configuration file")
	flag.Parse()

	conf := Config(configFile, timeout)

	if *isServer {
		RunServer(conf)
	} else {
		if *isPipe {
			if *isCopy {
				RunPipeCopy(conf, "")
			} else {
				RunPipePaste(conf)
			}
		} else {
			if *isCopy {
				RunCopy(conf, *filepath, *str)
			} else {
				RunPaste(conf)
			}
		}


	}
	
}
