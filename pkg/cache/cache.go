package cache

import "log"

type Cache interface {
	Set(key string, value []byte) error
	Get(key string) ([]byte, error)
	Delete(key string) error
	GetStatus() Stat
}

var globalCache Cache

func InitCache(name string) {
	var c Cache
	switch name {
	case "memory":
		c = newMCache()
	default:
		log.Panicln("unknown cache:", name)
	}
	log.Println("cache:", name, " is created")
	globalCache = c
}

func GlobalCache() Cache {
	return globalCache
}
