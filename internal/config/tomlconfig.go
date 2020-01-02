package config

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
)

type tomlConfig struct {
	Connect     string
	Listen      string
	Key         string
}

func getTomlConfig(configFile *string) tomlConfig {
	tomlData, err := ioutil.ReadFile(expandConfigFile(*configFile))

	var tomlConf tomlConfig
	if _, err = toml.Decode(string(tomlData), &tomlConf); err != nil {
		log.Fatal(err)
	}
	return tomlConf
}
