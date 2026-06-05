package main

import (
	"log"

	"github.com/IshankSharma2178/redis-internals/TCP-Echo-Server/internals/config"
	"github.com/IshankSharma2178/redis-internals/TCP-Echo-Server/server"
)

func main() {
	config.SetupFlags()
	log.Println("rolling the dice")
	server.RunSyncTCPServer()
}
