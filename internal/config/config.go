package config

var Host = "0.0.0.0"
var Port = 8081
var MaxConnection = 20000
var KeyNumberLimit = 5000000

const (
	EvictFirst int = 0
	LRU            = 1
	LFU            = 2
)

var EvictStrategy = EvictFirst
var AOFFileName = "./memkv-master.aof"
