# A simple in-memory key-value database

## How to run
```
  cd cmd
  go run main.go
  # on another terminal
  redis-cli -p 8081
```
## Supported features
- Compatible with [Redis CLI](https://redis.io/docs/ui/cli/)
- Single-threaded architecture
- Multiplexing IO for Linux using epoll
- [RESP protocol](https://redis.io/docs/reference/protocol-spec/)
- Graceful shutdown
- Simple eviction mechanism
- Commands:
  - PING
  - SET
  - GET
  - TTL
  - DEL
  - EXPIRE
  - INCR
  - ZADD
  - ZRANK
  - ZREM
  - ZSCORE
  - ZCARD
  - GEOADD
  - GEODIST
## WIP
- Geohash commands: 
  - [x] GEOADD
  - [x] GEODIST
  - [ ] GEOHASH
  - [ ] GEOSEARCH
## Todo
- Bloom filter commands
- Hyperloglog
- Count-min sketch
- Morris counter
- Cuckoo filter
- Approx LRU eviction