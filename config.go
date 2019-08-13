package main

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
	"time"
)

const (
	DefaultListen = "0.0.0.0:7788"
	DefaultConnect = "127.0.0.1:7788"
	DefaultKey = "iorDDkjFaMAJp8HNxwAWoyNKqLGTmG87"
)

type tomlConfig struct {
	Connect     string
	Listen      string
	Key         string
}

func Config(configFile *string, timeout *uint) Conf {

	tomlData, err := ioutil.ReadFile(expandConfigFile(*configFile))

	var tomlConf tomlConfig
	if _, err = toml.Decode(string(tomlData), &tomlConf); err != nil {
		log.Fatal(err)
	}

	var conf Conf
	conf.Connect = DefaultConnect
	if tomlConf.Connect != "" {
		conf.Connect = tomlConf.Connect
	}
	conf.Listen = DefaultListen
	if tomlConf.Listen != "" {
		conf.Listen = tomlConf.Listen
	}
	conf.Key = DefaultKey
	if tomlConf.Key != "" {
		conf.Key = tomlConf.Key
	}
	conf.Timeout = time.Duration(*timeout) * time.Second

	return conf
}
