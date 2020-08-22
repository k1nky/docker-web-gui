package main

import (
	"flag"
	"os"
)

var (
	Cfg Config
)

type Config struct {
	DockerHost string
	Listen     string
}

func (cfg *Config) ParseArgs() {
	flag.StringVar(&cfg.DockerHost, "d", "unix:///var/run/docker.sock", "Defines the docker host location")
	flag.StringVar(&cfg.Listen, "l", ":8000", "Defines a listen interface and port")
	if v, ok := os.LookupEnv("LISTEN"); ok {
		cfg.Listen = v
	}
	flag.Parse()
}
