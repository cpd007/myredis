package core

import (
	"time"

	"github.com/cpd007/myredis/config"
)

var store map[string]*Obj

type Obj struct {
	Value     any
	ExpiresAt int64
}

func InitializeStore() {
	store = make(map[string]*Obj)
}

func NewObj(value any, exDurationMs int64) *Obj {

	expiresAt := int64(-1)
	if exDurationMs > 0 {
		expiresAt = time.Now().UnixMilli() + exDurationMs
	}

	return &Obj{
		Value:     value,
		ExpiresAt: expiresAt,
	}
}

func Put(k string, v *Obj) {
	if _, exists := store[k]; !exists && len(store) >= config.EvictionCfg.MaxKeyLimit {
		evict()
	}
	store[k] = v
}

func Get(k string) *Obj {
	obj := store[k]
	if obj != nil && obj.ExpiresAt > 0 {
		if obj.ExpiresAt < time.Now().UnixMilli() {
			Del(k)
			return nil
		}
	}
	return obj
}

func Del(k string) bool {
	if _, exists := store[k]; !exists {
		return false
	}
	delete(store, k)
	return true
}
