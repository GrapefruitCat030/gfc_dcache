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

func (mc *MCache) NewScanner() Scanner {
	mc.mtx.RLock()
	defer mc.mtx.RUnlock()
	keys := make([]string, 0, len(mc.mcache))
	for k := range mc.mcache {
		keys = append(keys, k)
	}
	return &MCacheScanner{
		cache: mc,
		keys:  keys,
		idx:   -1, // start from -1, so the first call of Scan will move to 0
	}
}

func newMCache() *MCache {
	return &MCache{
		mcache: make(map[string][]byte),
		status: Stat{},
		mtx:    sync.RWMutex{},
	}
}

type MCacheScanner struct {
	cache *MCache
	keys  []string
	idx   int
	mtx   sync.RWMutex
}

func (mcs *MCacheScanner) Scan() bool {
	mcs.mtx.RLock()
	defer mcs.mtx.RUnlock()
	mcs.idx++
	return mcs.idx < len(mcs.keys)
}

func (mcs *MCacheScanner) Key() string {
	mcs.mtx.RLock()
	defer mcs.mtx.RUnlock()
	if mcs.idx < 0 || mcs.idx >= len(mcs.keys) {
		return ""
	}
	return mcs.keys[mcs.idx]
}

func (mcs *MCacheScanner) Value() []byte {
	mcs.mtx.RLock()
	defer mcs.mtx.RUnlock()
	if mcs.idx < 0 || mcs.idx >= len(mcs.keys) {
		return nil
	}
	return mcs.cache.mcache[mcs.keys[mcs.idx]]
}

func (mcs *MCacheScanner) Close() {
	// do nothing
}
