package session

import (
	"github.com/bradfitz/gomemcache/memcache"
	"time"
)

type MemcacheStore struct {
	prefix string
	client *memcache.Client
}

func NewMemcacheStore(prefix string, servers ...string) (*MemcacheStore) {
	return &MemcacheStore{
		prefix: prefix,
		client: memcache.New(servers...),
	}
}

func (s *MemcacheStore) Get(key string) ([]byte, error) {
	i, err := s.client.Get(s.prefix+key)
	if err != nil {
		return nil, err
	}
	return i.Value, nil
}

func (s *MemcacheStore) Set(key string, value []byte, ttl time.Duration) error {
	return s.client.Set(&memcache.Item{Key: key, Value: value, Expiration: int32(ttl.Seconds())})
}

func (s *MemcacheStore) Delete(key string) error {
	return s.client.Delete(key)
}



