package core

import "memkv/internal/config"

func evictFirst() {
	for k := range dictStore {
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
