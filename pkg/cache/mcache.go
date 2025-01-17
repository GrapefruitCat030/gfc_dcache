package cache

import (
	"sync"
)

type MCache struct {
	mcache map[string][]byte
	status Stat
	mtx    sync.RWMutex
}

func (mc *MCache) Set(key string, value []byte) error {
	mc.mtx.Lock()
	defer mc.mtx.Unlock()
	if tmp, ok := mc.mcache[key]; ok {
		mc.status.del(key, tmp)
	}
	mc.mcache[key] = value
	mc.status.add(key, value)
	return nil
}

func (mc *MCache) Get(key string) ([]byte, error) {
	mc.mtx.RLock() // use RLock to allow multiple readers
	defer mc.mtx.RUnlock()
	return mc.mcache[key], nil
}

func (mc *MCache) Delete(key string) error {
	mc.mtx.Lock()
	defer mc.mtx.Unlock()
	if _, ok := mc.mcache[key]; ok {
		mc.status.del(key, mc.mcache[key])
		delete(mc.mcache, key)
	}
	return nil
}

func (mc *MCache) GetStatus() Stat {
	return mc.status
}

func (mc *MCache) Close() error {
	mc.mtx.Lock()
	defer mc.mtx.Unlock()
	return nil
}

func newMCache() *MCache {
	return &MCache{
		mcache: make(map[string][]byte),
		status: Stat{},
		mtx:    sync.RWMutex{},
	}
}
