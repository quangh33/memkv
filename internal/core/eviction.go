package core

import "memkv/internal/config"

func evictFirst() {
	for k := range keyValueStore {
		Del(k)
		return
	}
}

func evict() {
	switch config.EvictStrategy {
	case config.EvictFirst:
		evictFirst()
	default:
		evictFirst()
	}
}
