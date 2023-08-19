package core

import "memkv/config"

func evictFirst() {
	for k := range store {
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
