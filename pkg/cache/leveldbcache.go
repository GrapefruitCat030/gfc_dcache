package cache

import "github.com/syndtr/goleveldb/leveldb"

const (
	levelDBCachePath = "/tmp/leveldb"
)

type LevelDBCache struct {
	db     *leveldb.DB
	status Stat
}

func (lc *LevelDBCache) Set(key string, value []byte) error {
	if err := lc.db.Put([]byte(key), value, nil); err != nil {
		return err
	}
	lc.status.add(key, value)
	return nil
}

func (lc *LevelDBCache) Get(key string) ([]byte, error) {
	return lc.db.Get([]byte(key), nil)
}

func (lc *LevelDBCache) Delete(key string) error {
	if err := lc.db.Delete([]byte(key), nil); err != nil {
		return err
	}
	lc.status.del(key, nil)
	return nil
}

func (lc *LevelDBCache) GetStatus() Stat {
	return lc.status
}

func newLevelDBCache() *LevelDBCache {
	db, err := leveldb.OpenFile(levelDBCachePath, nil)
	if err != nil {
		panic(err)
	}
	return &LevelDBCache{
		db:     db,
		status: Stat{},
	}
}
