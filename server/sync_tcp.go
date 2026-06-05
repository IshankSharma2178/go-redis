package server

import (
	"io"
	"log"
	"net"
	"strconv"

	"github.com/IshankSharma2178/redis-internals/TCP-Echo-Server/internals/config"
)

func readCommand(c net.Conn) (string, error) {
	buf := make([]byte, 512)
	n, err := c.Read(buf[:])
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}

func respond(cmd string, c net.Conn) error {
	_, err := c.Write([]byte(cmd))
	return err
}

func RunSyncTCPServer() {
	log.Println("starting a synchronousTCP server on", config.Cfg.Host, config.Cfg.Port)

	var con_clients int = 0 // hold the number of concurrent client at a specific moment

	lsnr, err := net.Listen("tcp", config.Cfg.Host+":"+strconv.Itoa(config.Cfg.Port))
	if err != nil {
		panic(err)
	}
	for {
		c, err := lsnr.Accept() // blocking call : waiting for the new client to connect
		if err != nil {
			panic(err)
		}

		con_clients += 1
		log.Println("new client connected on ", c.RemoteAddr(), " concurrent clients:", con_clients)

		for {
			cmd, err := readCommand(c)
			if err != nil {
				c.Close()
				con_clients -= 1
				log.Println("client disconnected", c.RemoteAddr(), " concurrent clients", con_clients)
				if err == io.EOF {
					break
				}
				log.Println("err", err)
			}
			log.Println("command", cmd)
			if err = respond(cmd, c); err != nil {
				log.Print("err write: ", err)
			}
		}

	}
}
