package config

import "flag"

type Config struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

var Cfg = Config{
	Host: "0.0.0.0",
	Port: 7379,
}

func SetupFlags() {
	flag.StringVar(&Cfg.Host, "host", Cfg.Host, "host for the server")
	flag.IntVar(&Cfg.Port, "port", Cfg.Port, "port for the server")
	flag.Parse()
}
