package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/IshankSharma2178/go-redis/core"
	"github.com/IshankSharma2178/go-redis/internals/config"
)

func toArrayString(ai []interface{}) ([]string, error) {
	as := make([]string, len(ai))
	for i := range ai {
		as[i] = ai[i].(string)
	}
	return as, nil
}

func readCommands(c io.ReadWriter) (core.RedisCmds, error) {
	buf := make([]byte, 512)
	n, err := c.Read(buf[:])
	if err != nil {
		return nil, err
	}
	values, err := core.Decode(buf[:n])
	if err != nil {
		return nil, err
	}

	arrayValues, ok := values.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected RESP value type %T", values)
	}

	var cmds []*core.RedisCmd = make([]*core.RedisCmd, 0)
	for _, value := range arrayValues {
		tokens, err := toArrayString(value.([]interface{}))
		if err != nil {
			return nil, err
		}
		cmds = append(cmds, &core.RedisCmd{
			Cmd:  strings.ToUpper(tokens[0]),
			Args: tokens[1:],
		})
	}
	return cmds, nil
}

func respondError(err error, c io.ReadWriter) {
	c.Write([]byte(fmt.Sprintf("-%s\r\n", err))) // adding '-' becoz error need to be transformed into resp before sending
}

func respond(cmds core.RedisCmds, c io.ReadWriter) {
	core.EvalAndRespond(cmds, c)
}

func RunSyncTCPServer() {
	log.Println("starting a synchronousTCP server on", config.Cfg.Host, config.Cfg.Port)

	var con_clients int = 0 // hold the number of concurrent client at a specific moment

	lsnr, err := net.Listen("tcp", config.Cfg.Host+":"+strconv.Itoa(config.Cfg.Port))
	if err != nil {
		log.Println("err", err)
		return
	}
	for {
		c, err := lsnr.Accept() // blocking call : waiting for the new client to connect
		if err != nil {
			log.Println("err", err)
			return
		}

		con_clients += 1
		log.Println("new client connected on ", c.RemoteAddr(), " concurrent clients:", con_clients)

		for {
			cmd, err := readCommands(c)
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
			respond(cmd, c)
		}
	}
}
