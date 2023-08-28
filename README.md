# A simple in-memory key-value database

## How to run
```
  cd cmd/memkv
  go run main.go
  # on another terminal
  redis-cli -p 8081
```
## Supported features
- Single-threaded architecture
- Multiplexing IO for Linux using epoll
- [RESP protocol](https://redis.io/docs/reference/protocol-spec/)
- Graceful shutdown
- Simple eviction
- Commands:
  - PING
  - SET
  - GET
  - TTL
  - DEL
  - EXPIRE
  - INCR

## WIP
- Geohash commands: 
  - [ ] GEOADD
  - [ ] GEODIST
  - [ ] GEOHASH
  - [ ] GEOSEARCH
## Todo
- Approx LRU eviction
- Bloom filter commands
- Hyperloglog
- Count-min sketch
- ZSET using Skiplist
- Morris counter