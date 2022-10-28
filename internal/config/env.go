package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
)

func ReadENVParams(conf *Config) {
	envConf := New()
	err := env.Parse(envConf)
	if err != nil {
		fmt.Println("Could not fetch server ENV params")
		panic(err)
	}
	if envConf.RunAddress != "" {
		conf.RunAddress = envConf.RunAddress
	}

	if envConf.DatabaseURI != "" {
		conf.DatabaseURI = envConf.DatabaseURI
	}

	if envConf.AccuralSystemAddress != "" {
		conf.AccuralSystemAddress = envConf.AccuralSystemAddress
	}

}
