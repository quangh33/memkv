package main

import (
	"flag"
	"fmt"
	"memkv/config"
	"memkv/server"
)

func init() {
	flag.StringVar(&config.Host, "host", "0.0.0.0", "host")
	flag.IntVar(&config.Port, "port", config.Port, "port")
	flag.Parse()
}

func main() {
	fmt.Println("starting memkv database ...")
	// server.RunSyncTcpServer()
	server.RunAsyncTCPServer()
}
