package main

import (
	// "fmt"
	"flag"
	"log"
	"time"
)


// Conf - Shared config
type Conf struct {
	Connect        string
	Listen         string
	Key            string
	IsPipe         bool
	PipeBlockSize  int64
	Timeout        time.Duration
}

const (
	DefaultConfigFile = "~/.pin.toml"
	DefaultPipeBlockSize = 4
)


func expandConfigFile(path string) string {
	file, err := Expand(path)
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
	pipeBlockSize := flag.Int64("bsize", 4, "pipe block size")
	flag.Parse()

	conf := Config(configFile, timeout)
	conf.PipeBlockSize = *pipeBlockSize * ONE_M_BSIZE
	conf.IsPipe = *isPipe

	if *isServer {
		RunServer(conf)
	} else {
		if *isCopy {
			RunSender(conf, *filepath, *str)
		} else {
			RunReceiver(conf, *filepath)
		}
		//if *isPipe {
		//	if *filepath == "" {
		//		log.Fatal("please specify the filepath you want to transfer")
		//	}
		//	if *isCopy {
		//		if !FileExists(*filepath) {
		//			log.Fatal("file not exist")
		//		}
		//		RunPipeCopy(conf, *filepath)
		//	} else {
		//		if FileExists(*filepath) {
		//			log.Fatal("file already exist")
		//		}
		//		RunPipePaste(conf, *filepath)
		//	}
		//} else {
		//	if *isCopy {
		//		RunCopy(conf, *filepath, *str)
		//	} else {
		//		RunPaste(conf)
		//	}
		//}


	}
	
}
