package kmgCache

import (
	"errors"
	"github.com/bronze1man/kmg/kmgMath"
	"github.com/golang/groupcache/singleflight"
	"sync"
	"time"
)

var CacheMiss = errors.New("cache miss")

type TtlCacheEntry struct {
	Value   interface{}
	Timeout time.Time
}

func (entry TtlCacheEntry) GetTtl() uint32 {
	ttlDur := entry.Timeout.Sub(time.Now())
	if ttlDur < 0 {
		ttlDur = 0
	}
	return uint32(kmgMath.FloorToInt(ttlDur.Seconds()))
}

//请调用 NewTtlCache() 初始化
type TtlCache struct {
	cache       map[string]TtlCacheEntry
	lock        sync.RWMutex
	singleGroup singleflight.Group
}

func NewTtlCache() *TtlCache {
	return &TtlCache{
		cache: map[string]TtlCacheEntry{},
	}
}

//如果err不是空,则不会把数据保存在缓存里面,但是会返回另外2项.
func (s *TtlCache) DoWithTtl(key string, f func() (value interface{}, ttl uint32, err error)) (value interface{}, ttl uint32, err error) {
	entry, err := s.get(key)
	if err == nil {
		return entry.Value, entry.GetTtl(), nil
	}
	if err != CacheMiss {
		return
	}
	entryi, err := s.singleGroup.Do(key, func() (interface{}, error) {
		value, ttl, err := f()
		timeout := time.Now().Add(time.Duration(ttl) * time.Second)
		return TtlCacheEntry{
			Value:   value,
			Timeout: timeout,
		}, err
	})
	entry = entryi.(TtlCacheEntry)
	ttl = entry.GetTtl()
	if err == nil && ttl > 0 {
		s.save(key, entry)
	}
	return entry.Value, ttl, nil
}

func (s *TtlCache) save(key string, entry TtlCacheEntry) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.cache[key] = entry
	return
}

func (s *TtlCache) get(key string) (entry TtlCacheEntry, err error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	now := time.Now()
	entry, ok := s.cache[key]
	if !ok {
		return entry, CacheMiss
	}
	if now.After(entry.Timeout) {
		return entry, CacheMiss
	}
	return entry, nil
}

//要有个进程在一边进行gc,避免内存泄漏
func (s *TtlCache) GcThread() {
	for {
		time.Sleep(time.Hour)
		s.lock.Lock()
		now := time.Now()
		for key, entry := range s.cache {
			if now.After(entry.Timeout) {
				delete(s.cache, key)
			}
		}
		s.lock.Unlock()
	}
}
