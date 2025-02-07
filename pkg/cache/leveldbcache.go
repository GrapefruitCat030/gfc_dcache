package cache

import (
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

type LevelDBCache struct {
	db         *leveldb.DB
	batch      *leveldb.Batch
	batchCh    chan *pair
	batchMutex sync.Mutex
	done       chan struct{}
}

func (lc *LevelDBCache) Set(key string, value []byte) error {
	lc.batchCh <- &pair{key: key, value: value}
	return nil
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

func newLevelDBCache() *LevelDBCache {
	db, err := leveldb.OpenFile(levelDBCachePath, nil)
	if err != nil {
		panic(err)
	}
	lc := &LevelDBCache{
		db:      db,
		batch:   new(leveldb.Batch),
		batchCh: make(chan *pair, 5000),
		done:    make(chan struct{}),
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
