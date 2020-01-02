package config

import (
	"github.com/chenminhua/pin/internal/fs"
	"log"
	"time"
)

const (
	DefaultListen = "0.0.0.0:7788"
	DefaultConnect = "127.0.0.1:7788"
	DefaultKey = "iorDDkjFaMAJp8HNxwAWoyNKqLGTmG87"
)

/**
	Conf - Shared config
	Key是用于鉴权的，前后端在发送请求的时候需要带着
 */
type Conf struct {
	Connect        string
	Listen         string
	Key            string
	IsPipe         bool
	PipeBlockSize  int64
	Timeout        time.Duration
}

func Config(configFile *string, timeout *uint) Conf {

	tomlConf := getTomlConfig(configFile)

	var conf Conf
	conf.Connect = DefaultConnect
	conf.Listen = DefaultListen
	conf.Key = DefaultKey

	if tomlConf.Connect != "" {
		conf.Connect = tomlConf.Connect
	}

	if tomlConf.Listen != "" {
		conf.Listen = tomlConf.Listen
	}

	if tomlConf.Key != "" {
		conf.Key = tomlConf.Key
	}

	conf.Timeout = time.Duration(*timeout) * time.Second
	return conf
}

func expandConfigFile(path string) string {
	file, err := fs.Expand(path)
	if err != nil {
		log.Fatal(err)
	}
	return file
}
