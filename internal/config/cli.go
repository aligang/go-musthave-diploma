package config

import (
	"flag"
)

func ReadCLIParams(conf *Config) {
	flag.StringVar(&conf.RunAddress, "a", "127.0.0.1:8080", "host to listen on")
	flag.StringVar(&conf.DatabaseURI, "d", "", "Database URI")
	flag.StringVar(&conf.RunAddress, "r", "127.0.0.1:9090", "URL of accural system")
	flag.Parse()
}
