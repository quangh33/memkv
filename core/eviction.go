package core

import "memkv/config"

func evictFirst() {
	for k := range store {
		delete(store, k)
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
