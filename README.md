# A simple in-memory key-value database

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
- Skiplist
- Morris Counter