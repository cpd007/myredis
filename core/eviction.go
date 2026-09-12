package core

import (
	"github.com/cpd007/myredis/config"
	"github.com/cpd007/myredis/constants"
)

func evict() {
	// TODO: add different eviction strategy
	switch config.EvictionCfg.EvictionStrategy {
	case constants.EvictFist:
		evictFirst()
	}
}

// delete the first key that is found
func evictFirst() {
	for k := range store {
		Del(k)
		return
	}
}
