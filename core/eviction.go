package core

import (
	"time"

	"github.com/IshankSharma2178/go-redis/internals/config"
)

// Evicts the first key it found while iterating the map
// TODO: Make it efficient by doing thorough sampling
func evictFirst() {
	for k := range store {
		Del(k)
		return
	}
}

func getCurrentClock() uint32 {
	return uint32(time.Now().Unix()) & 0x00FFFFFF
}

func getIdleTime(lastAccessedAt uint32) uint32 {
	c := getCurrentClock()
	if c >= lastAccessedAt {
		return c - lastAccessedAt
	}
	return (0x00FFFFFF - lastAccessedAt) + c
}

func populateEvictionPool() {
	sampleSize := 5
	for k := range store {
		ePool.Push(k, store[k].LastAccessedAt)
		sampleSize--
		if sampleSize == 0 {
			break
		}
	}
}

// TODO: no need to populate everytime. should populate
// only when the number of keys to evict is less than what we have in the pool
func evictAllkeysLRU() {
	populateEvictionPool()
	evictCount := int16(config.Cfg.EvictionRatio * float64(config.Cfg.KeysLimit))
	for i := 0; i < int(evictCount) && len(ePool.pool) > 0; i++ {
		item := ePool.Pop()
		if item == nil {
			return
		}
		Del(item.key)
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
	case "allkeys-lru":
		evictAllkeysLRU()
	}
}
