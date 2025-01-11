package cache

import (
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	levelDBCachePath = "/tmp/leveldb"
)

type LevelDBCache struct {
	db *leveldb.DB
}

func (lc *LevelDBCache) Set(key string, value []byte) error {
	return lc.db.Put([]byte(key), value, nil)
}

func (lc *LevelDBCache) Get(key string) ([]byte, error) {
	return lc.db.Get([]byte(key), nil)
}

func (lc *LevelDBCache) Delete(key string) error {
	return lc.db.Delete([]byte(key), nil)
}

func (lc *LevelDBCache) GetStatus() Stat {
	iter := lc.db.NewIterator(nil, nil)
	defer iter.Release()
	s := Stat{}
	for iter.Next() {
		s.add(string(iter.Key()), iter.Value())
	}
	return s
}

func newLevelDBCache() *LevelDBCache {
	db, err := leveldb.OpenFile(levelDBCachePath, nil)
	if err != nil {
		panic(err)
	}
	return &LevelDBCache{db: db}
}
