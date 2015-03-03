package kmgCache

import (
	"errors"
	"github.com/golang/groupcache/singleflight"
	"sync"
	"time"
)

var cacheExpire = errors.New("cache expire")
var DoNotNeedCache = errors.New("do not need cache")

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
// 如果你把needcache设置为false,表示这个请求本次不进行缓存(应该是什么地方出现错误了,具体错误请自行处理)
// 1.如果缓存不存在,会同步查询
// 2.如果缓存过去,会异步查询,并返回旧的数据
// 3.如果不要求保存数据,不会把数据保存到缓存中.
// 		1. 如果缓存不存在,返回f.value f.ttl
func (s *AsyncTtlCache) DoWithTtl(key string, f func() (value interface{}, ttl uint32, canSave bool)) (value interface{}, ttl uint32) {
	entry, err := s.get(key)
	if err == nil {
		return entry.Value, entry.GetTtl()
	}
	updateCache := func() (value interface{}, ttl uint32) {
		//异步更新缓存
		entryi, err := s.singleGroup.Do(key, func() (out interface{}, err error) {
			value, ttl, canSave := f()
			timeout := time.Now().Add(time.Duration(ttl) * time.Second)
			out = TtlCacheEntry{
				Value:   value,
				Timeout: timeout,
			}
			if !canSave {
				err = DoNotNeedCache
			}
			return
		})
		entryn := entryi.(TtlCacheEntry)
		if err == nil {
			s.save(key, entryn) //ttl 是0 也存进去,下次可以异步刷新
		}
		ttl = entryn.GetTtl()
		return entryn.Value, ttl
	}
	switch err {
	case CacheMiss:
		value, ttl := updateCache()
		return value, ttl
	case cacheExpire:
		go updateCache()
		return entry.Value, 0
	default:
		return nil, 0
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

//里面数据的个数
func (s *AsyncTtlCache) Len() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.cache)
}
