package config

var Host string = "0.0.0.0"
var Port int = 8081
var MaxConnection int = 20000
var KeyNumberLimit int = 5000000

const (
	EvictFirst int = 0
	LRU            = 1
	LFU            = 2
)

var EvictStrategy int = EvictFirst
var AOFFileName = "./memkv-master.aof"
