package db

import (
	"encoding/json"
	"time"

	"github.com/patrickmn/go-cache"
)

type CacheService struct {
	cache *cache.Cache
}

var (
	cacheServices *CacheService
	maxCacheTime  = time.Hour * 12
)

type CacheServiceInterface interface {
	Get(key string) string
	Set(key string, value string)
	SetCustom(key string, value string, time time.Duration)
	SetInterfaceNotTimeLimit(key string, val interface{})
	GetInterface(key string) []byte
	SetInterface(key string, val interface{})
	Clear(key string)
}

func (s *CacheService) Get(key string) string {
	context, ok := s.cache.Get(key)
	if !ok {
		return ""
	}
	return context.(string)
}

func (s *CacheService) Set(key string, value string) {
	s.cache.Set(key, value, maxCacheTime)
}

func (s *CacheService) SetInterface(key string, value interface{}) {
	byte, _ := json.Marshal(&value)
	s.cache.Set(key, string(byte), maxCacheTime)
}

func (s *CacheService) GetInterface(key string) []byte {
	context, ok := s.cache.Get(key)
	if !ok || context == "" {
		return nil
	}

	return []byte(context.(string))
}

func (s *CacheService) SetCustom(key string, value string, time time.Duration) {
	s.cache.Set(key, value, time)
}

func (s *CacheService) SetInterfaceNotTimeLimit(key string, value interface{}) {
	bytes, _ := json.Marshal(&value)
	s.cache.Set(key, string(bytes), 0)
}

func (s *CacheService) Clear(key string) {
	s.cache.Delete(key)
}

func GetCache() CacheServiceInterface {
	if cacheServices == nil {
		cacheServices = &CacheService{cache: cache.New(time.Hour*12, time.Hour*1)}
	}
	return cacheServices
}
