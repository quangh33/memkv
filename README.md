# A simple in-memory key-value database

## How to run
```
  cd cmd/memkv
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

## WIP
- Geohash commands: 
  - [ ] GEOADD
  - [ ] GEODIST
  - [ ] GEOHASH
  - [ ] GEOSEARCH
- Ordered set using Skiplist
## Todo
- Approx LRU eviction
- Bloom filter commands
- Hyperloglog
- Count-min sketch
- Morris counter