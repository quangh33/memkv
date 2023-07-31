package server

import (
	"fmt"
	"io"
	"log"
	"memkv/config"
	"memkv/core"
	"net"
	"strconv"
)

func readCommand(con net.Conn) (*core.MemKVCmd, error) {
	var buf []byte = make([]byte, 512)
	n, err := con.Read(buf)
	if err != nil {
		return nil, err
	}

	return core.ParseCmd(buf[:n])
}

func responseError(err error, c net.Conn) {
	c.Write([]byte(fmt.Sprintf("-%s%s", err, core.CRLF)))
}

func response(cmd *core.MemKVCmd, con net.Conn) {
	err := core.EvalAndResponse(cmd, con)
	if err != nil {
		responseError(err, con)
	}
}

func RunSyncTcpServer() {
	log.Println("starting a sync TCP server on", config.Host, config.Port)
	listener, err := net.Listen("tcp", config.Host+":"+strconv.Itoa(config.Port))
	if err != nil {
		panic(err)
	}

	for {
		con, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		log.Println("a new client connected with address:", con.RemoteAddr())

		for {
			cmd, err := readCommand(con)
			if err != nil {
				con.Close()
				if err == io.EOF {
					break
				}
				log.Println("err", err)
			}
			log.Println("command:", cmd)
			response(cmd, con)
		}
	}
}
