package main

import (
	"flag"
	"github.com/chenminhua/pin/internal/config"
)

const (
	DefaultConfigFile = "~/.pin.toml"
	DefaultPipeBlockSize = 4
)

var ONE_M_BSIZE int64 = 1024 * 1024

func main() {
	isCopy := flag.Bool("copy", false, "copy sth to server")
	isPipe := flag.Bool("pipe", false, "pipe")
	filepath := flag.String("f", "", "file")
	str := flag.String("s", "", "string")
	timeout := flag.Uint("timeout", 10, "connection timeout (seconds)")
	configFile := flag.String("config", DefaultConfigFile, "configuration file")
	pipeBlockSize := flag.Int64("bsize", 4, "pipe block size")
	flag.Parse()

	conf := config.Config(configFile, timeout)
	conf.PipeBlockSize = *pipeBlockSize * ONE_M_BSIZE
	conf.IsPipe = *isPipe


	if *isCopy {
		RunSender(conf, *filepath, *str)
	} else {
		RunReceiver(conf, *filepath)
	}
}
