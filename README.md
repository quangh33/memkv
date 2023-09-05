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
  - SET, GET, DEL, TTL, EXPIRE, INCR
  - ZADD, ZRANK, ZREM, ZSCORE, ZCARD
  - SADD, SREM, SCARD, SMEMEBERS, SISMEMBER, SRAND, SPOP
  - GEOADD, GEODIST, GEOHASH, GEOSEARCH
## WIP
- Geohash commands: 
  - [x] GEOADD
  - [x] GEODIST
  - [x] GEOHASH
  - [x] GEOSEARCH
## Todo
- Bloom filter commands
- Hyperloglog
- Count-min sketch
- Morris counter
- Cuckoo filter
- Approx LRU eviction