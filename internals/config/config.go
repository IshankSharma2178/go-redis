package config

import "flag"

type Config struct {
	Host             string `yaml:"host"`
	Port             int    `yaml:"port"`
	KeysLimit        int    `yaml:"keys_limit"`
	EvictionStrategy string `yaml:"eviction_strategy"`
	AOFFile          string `yaml:"aof_file"`
}

var Cfg = Config{
	Host:             "0.0.0.0",
	Port:             7379,
	KeysLimit:        5,
	EvictionStrategy: "simple-first",
	AOFFile:          "./bin/appendonly.aof",
}

func SetupFlags() {
	flag.StringVar(&Cfg.Host, "host", Cfg.Host, "host for the server")
	flag.IntVar(&Cfg.Port, "port", Cfg.Port, "port for the server")
	flag.StringVar(&Cfg.AOFFile, "aof-file", Cfg.AOFFile, "path to the AOF file")
	flag.Parse()
}
