package core

import "github.com/IshankSharma2178/go-redis/internals/config"

// Evicts the first key it found while iterating the map
// TODO: Make it efficient by doing thorough sampling
func evictFirst() {
	for k := range store {
		Del(k)
		return
	}
}
func evictAllkeysRandom() {
	evictCount := int64(config.Cfg.EvictionRatio * float64(config.Cfg.KeysLimit))
	// Iteration of Golang dictionary can be considered sas a random
	// because it depends on the hash of the inserted key
	for k := range store {
		Del(k)
		evictCount--
		if evictCount <= 0 {
			break
		}
	}
}

// TODO: Make the eviction strategy configuration driven
// TODO: Support multiple eviction strategies
func evict() {
	switch config.Cfg.EvictionStrategy {
	case "simple-first":
		evictFirst()
	case "allkeys-random":
		evictAllkeysRandom()
	}
}
