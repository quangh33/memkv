package main

import (
	"flag"
	"fmt"
	"memkv/config"
	"memkv/server"
)

func init() {
	flag.StringVar(&config.Host, "host", "0.0.0.0", "host")
	flag.IntVar(&config.Port, "port", 8080, "port")
	flag.Parse()
}

func main() {
	fmt.Println("starting memkv database ...")
	server.RunSyncTcpServer()
}
