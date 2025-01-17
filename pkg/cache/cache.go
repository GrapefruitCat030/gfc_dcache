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
}

var globalCache Cache

func InitCache(name string) {
	var c Cache
	switch name {
	case CacheTypeMemory:
		c = newMCache()
	case CacheTypeLevelDB:
		c = newLevelDBCache()
	default:
		log.Panicln("unknown cache:", name)
	}
	log.Println("cache:", name, " is created")
	globalCache = c
}

func GlobalCache() Cache {
	return globalCache
}
