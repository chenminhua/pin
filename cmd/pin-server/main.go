package main

import (
	// "fmt"
	"flag"
	"github.com/chenminhua/pin/internal/config"
)

const (
	DefaultConfigFile = "~/.pin.toml"
	DefaultPipeBlockSize = 4
)




func main() {
	isPipe := flag.Bool("pipe", false, "pipe")
	timeout := flag.Uint("timeout", 10, "connection timeout (seconds)")
	configFile := flag.String("config", DefaultConfigFile, "configuration file")
	pipeBlockSize := flag.Int64("bsize", 4, "pipe block size")
	flag.Parse()

	conf := config.Config(configFile, timeout)
	conf.PipeBlockSize = *pipeBlockSize * ONE_M_BSIZE
	conf.IsPipe = *isPipe

	RunServer(conf)




}
