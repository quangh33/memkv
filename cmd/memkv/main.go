package main

import (
	"flag"
	"fmt"
	"memkv/internal/config"
	"memkv/internal/server"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func init() {
	flag.StringVar(&config.Host, "host", "0.0.0.0", "host")
	flag.IntVar(&config.Port, "port", config.Port, "port")
	flag.Parse()
}

func main() {
	fmt.Println("starting memkv database ...")
	var signals chan os.Signal = make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)
	var wg sync.WaitGroup
	wg.Add(2)

	go server.RunAsyncTCPServer(&wg)
	go server.WaitForSignal(&wg, signals)

	wg.Wait()
}
