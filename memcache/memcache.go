package memcache

import (
	"time"

	"github.com/DLag/session"
	"github.com/bradfitz/gomemcache/memcache"
)

const DefaultPrefix = "GOSESSION-"

type Store struct {
	prefix string
	client *memcache.Client
}

func NewStore(prefix string, servers ...string) *Store {
	return &Store{
		prefix: prefix,
		client: memcache.New(servers...),
	}
}

func NewDefaultStore(servers ...string) *Store {
	return NewStore(DefaultPrefix, servers...)
}

func NewDefaultManager(servers ...string) *session.Manager {
	return session.DefaultManager(NewDefaultStore(servers...))
}

func (s *Store) Get(key string) ([]byte, error) {
	i, err := s.client.Get(s.prefix + key)
	if err != nil {
		return nil, err
	}
	return i.Value, nil
}

func (s *Store) Set(key string, value []byte, ttl time.Duration) error {
	return s.client.Set(&memcache.Item{Key: key, Value: value, Expiration: int32(ttl.Seconds())})
}

func (s *Store) Delete(key string) error {
	return s.client.Delete(key)
}
