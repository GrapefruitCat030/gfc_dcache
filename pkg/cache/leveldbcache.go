package cache

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

const (
	levelDBCachePath = "/tmp/leveldb"
)

const (
	BATCH_SIZE = 100
)

type pair struct {
	key   string
	value []byte
}

type cacheValue struct {
	Value    []byte    `json:"value"`
	ExpireAt time.Time `json:"expireAt"`
}

type LevelDBCache struct {
	db         *leveldb.DB
	batch      *leveldb.Batch
	batchCh    chan *pair
	batchMutex sync.Mutex
	ttl        time.Duration
	done       chan struct{}
}

func (lc *LevelDBCache) Set(key string, value []byte) error {
	cacheVal := cacheValue{Value: value, ExpireAt: time.Now().Add(lc.ttl)}
	encodeVal, err := json.Marshal(cacheVal)
	if err != nil {
		return err
	}
	lc.batchCh <- &pair{key: key, value: encodeVal}
	return nil
}

func (lc *LevelDBCache) Get(key string) ([]byte, error) {
	encodeVal, err := lc.db.Get([]byte(key), nil)
	if err == leveldb.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var cacheVal cacheValue
	if err := json.Unmarshal(encodeVal, &cacheVal); err != nil {
		return nil, err
	}
	if time.Now().After(cacheVal.ExpireAt) {
		lc.Delete(key)
		return nil, fmt.Errorf("key %s is expired", key)
	}
	return cacheVal.Value, nil
}

func (lc *LevelDBCache) Delete(key string) error {
	return lc.db.Delete([]byte(key), nil)
}

func (lc *LevelDBCache) GetStatus() Stat {
	iter := lc.db.NewIterator(nil, nil)
	defer iter.Release()
	s := Stat{}
	for iter.Next() {
		var cacheVal cacheValue
		if err := json.Unmarshal(iter.Value(), &cacheVal); err != nil {
			continue
		}
		if time.Now().After(cacheVal.ExpireAt) {
			lc.Delete(string(iter.Key()))
			continue
		}
		s.add(string(iter.Key()), cacheVal.Value)
	}
	return s
}

func (lc *LevelDBCache) Close() error {
	close(lc.done)
	time.Sleep(1 * time.Second)
	return lc.db.Close()
}

func (lc *LevelDBCache) NewScanner() Scanner {
	iter := lc.db.NewIterator(nil, nil)
	return &LevelDBCacheScanner{iter: iter}
}

func (lc *LevelDBCache) flushBatch() error {
	lc.batchMutex.Lock()
	defer lc.batchMutex.Unlock()
	if lc.batch.Len() == 0 {
		return nil
	}
	if err := lc.db.Write(lc.batch, nil); err != nil {
		return err
	}
	lc.batch.Reset()
	return nil
}

func (lc *LevelDBCache) writeFunc() {
	t := time.NewTicker(1 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-lc.done:
			for len(lc.batchCh) > 0 {
				p := <-lc.batchCh
				lc.batch.Put([]byte(p.key), p.value)
			}
			lc.flushBatch()
			return
		case <-t.C:
			if err := lc.flushBatch(); err != nil {
				fmt.Println("flush batch error", err)
			}
		case p := <-lc.batchCh:
			lc.batch.Put([]byte(p.key), p.value)
			if lc.batch.Len() >= BATCH_SIZE {
				if err := lc.flushBatch(); err != nil {
					fmt.Println("flush batch error", err)
				}
			}
		}
	}
}

func newLevelDBCache(ttl int) *LevelDBCache {
	db, err := leveldb.OpenFile(levelDBCachePath, nil)
	if err != nil {
		panic(err)
	}
	lc := &LevelDBCache{
		db:      db,
		batch:   new(leveldb.Batch),
		batchCh: make(chan *pair, 5000),
		done:    make(chan struct{}),
		ttl:     time.Duration(ttl) * time.Second,
	}
	go lc.writeFunc()
	return lc
}

// LevelDBCacheScanner is the scanner for LevelDBCache
type LevelDBCacheScanner struct {
	iter iterator.Iterator
}

func (lcs *LevelDBCacheScanner) Scan() bool {
	return lcs.iter.Next()
}

func (lcs *LevelDBCacheScanner) Key() string {
	return string(lcs.iter.Key())
}

func (lcs *LevelDBCacheScanner) Value() []byte {
	return lcs.iter.Value()
}

func (lcs *LevelDBCacheScanner) Close() {
	lcs.iter.Release()
}
