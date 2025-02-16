package cache

import "log"

const (
	CacheTypeMemory  = "memory"
	CacheTypeLevelDB = "leveldb"
)

type Cache interface {
	Set(key string, value []byte) error
	Get(key string) ([]byte, error)
	Delete(key string) error
	GetStatus() Stat
	Close() error
	NewScanner() Scanner
}

var globalCache Cache

func InitCache(name string, ttl int) {
	var c Cache
	switch name {
	case CacheTypeMemory:
		c = newMCache(ttl)
	case CacheTypeLevelDB:
		c = newLevelDBCache(ttl)
	default:
		log.Panicln("unknown cache:", name)
	}
	log.Println("cache:", name, " is created")
	globalCache = c
}

func ShotdownCache() error {
	return globalCache.Close()
}

func GlobalCache() Cache {
	return globalCache
}
