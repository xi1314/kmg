package kmgCache

import (
	"errors"
	"github.com/golang/groupcache/singleflight"
	"sync"
	"time"
)

var cacheExpire = errors.New("cache expire")

//会在缓存超时的时候,异步更新缓存,并且返回前一个数据.
type AsyncTtlCache struct {
	cache       map[string]TtlCacheEntry
	lock        sync.RWMutex
	singleGroup singleflight.Group
}

func NewAsyncCache() *AsyncTtlCache {
	return &AsyncTtlCache{
		cache: map[string]TtlCacheEntry{},
	}
}

//如果缓存不存在,会同步查询
//如果缓存过期,会异步查询,以便下次请求的时候有这个数据
func (s *AsyncTtlCache) DoWithTtl(key string, f func() (value interface{}, ttl uint32, err error)) (value interface{}, ttl uint32, err error) {
	entry, err := s.get(key)
	if err == nil {
		return entry.Value, entry.GetTtl(), nil
	}
	updateCache := func() (value interface{}, ttl uint32, err error) {
		//异步更新缓存
		entryi, err := s.singleGroup.Do(key, func() (interface{}, error) {
			value, ttl, err := f()
			if err != nil {
				return nil, err
			}
			timeout := time.Now().Add(time.Duration(ttl) * time.Second)
			return TtlCacheEntry{
				Value:   value,
				Timeout: timeout,
			}, nil
		})

		switch err {
		case nil:
			entryn := entryi.(TtlCacheEntry)
			ttl := entryn.GetTtl()
			if ttl > 0 {
				s.save(key, entryn)
			}
			return entryn.Value, ttl, nil
		default:
			return nil, 0, err
		}
	}
	switch err {
	case CacheMiss:
		return updateCache()
	case cacheExpire:
		go updateCache()
		return entry.Value, 0, nil
	default:
		return nil, 0, err
	}
}

func (s *AsyncTtlCache) save(key string, entry TtlCacheEntry) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.cache[key] = entry
	return
}

func (s *AsyncTtlCache) get(key string) (entry TtlCacheEntry, err error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	now := time.Now()
	entry, ok := s.cache[key]
	if !ok {
		return entry, CacheMiss
	}
	if now.After(entry.Timeout) {
		return entry, cacheExpire
	}
	return entry, nil
}

//要有个进程在一边进行gc,避免内存泄漏
func (s *AsyncTtlCache) GcThread() {
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
