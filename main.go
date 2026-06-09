package main

import (
	"log"

	"github.com/IshankSharma2178/go-redis/internals/config"
	"github.com/IshankSharma2178/go-redis/server"
)

func main() {
	config.SetupFlags()
	log.Println("rolling the dice")
	server.RunAsyncTCPServer()
}
