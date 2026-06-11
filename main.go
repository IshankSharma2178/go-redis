package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/IshankSharma2178/go-redis/internals/config"
	"github.com/IshankSharma2178/go-redis/server"
)

func main() {
	config.SetupFlags()
	log.Println("rolling the dice")
	var sigs chan os.Signal = make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	var wg sync.WaitGroup
	wg.Add(2)

	go server.RunAsyncTCPServer(&wg)
	go server.WaitForSignal(&wg, sigs)

	wg.Wait()
}
