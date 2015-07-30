package kmgCache

import (
	"sync"
	"time"

	"github.com/golang/groupcache/singleflight"
)

//请调用 NewTtlCache() 初始化
type MemoryTtlCacheV2 struct {
	cache       map[string]ttlCacheEntry
	lock        sync.RWMutex
	singleGroup singleflight.Group
}

// TODO 写一个close,以便可以关闭缓存
func NewMemoryTtlCacheV2() *MemoryTtlCacheV2 {
	c := &MemoryTtlCacheV2{
		cache: map[string]ttlCacheEntry{},
	}
	go c.GcThread()
	return c
}

// Ttl 内存缓存实现第二版
// 1.使用key表示不同的项
// 2.f 返回err 表示本次无法获取到信息,该err会被返回给调用者,并且此时调用者的value是nil
// 3.使用singleGroup避免大量请求同时访问某个一个key
func (s *MemoryTtlCacheV2) Do(key string, f func() (value interface{}, ttl time.Duration, err error)) (value interface{}, err error) {
	entry, err := s.get(key)
	if err == nil {
		return entry.Value, nil
	}
	if err != CacheMiss {
		return
	}
	entryi, err := s.singleGroup.Do(key, func() (interface{}, error) {
		value, ttl, err := f()
		return ttlCacheEntry{
			Value:   value,
			Timeout: time.Now().Add(ttl),
		}, err
	})
	if err != nil {
		return nil, err
	}
	entry = entryi.(ttlCacheEntry)
	s.save(key, entry)
	return entry.Value, nil
}

func (s *MemoryTtlCacheV2) save(key string, entry ttlCacheEntry) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.cache[key] = entry
	return
}

func (s *MemoryTtlCacheV2) get(key string) (entry ttlCacheEntry, err error) {
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
func (s *MemoryTtlCacheV2) GcThread() {
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
