package mem

import (
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"gitlab.ozon.dev/dmitryssaenko/financial-tg-bot/internal/logger"
)

type Memcached struct {
	mc *memcache.Client
}

func NewMemcached(addresses ...string) *Memcached {
	mc := memcache.New(addresses...)
	return &Memcached{
		mc: mc,
	}
}

func (m *Memcached) Get(key string) (string, bool) {
	it, err := m.mc.Get(key)
	if err != nil {
		logger.Error("cannot extract value from cache")
		return "", false
	}
	return string(it.Value), true
}

func (m *Memcached) Add(key string, val string, d time.Duration) error {
	if err := m.mc.Set(&memcache.Item{Key: key, Value: []byte(val), Expiration: int32(time.Now().Add(d).Unix())}); err != nil {
		logger.Error("cannot save value to cache")
		return fmt.Errorf("cannot save value to cache")
	}
	return nil
}

func (m *Memcached) Delete(key string) error {
	if err := m.mc.Delete(key); err != nil {
		logger.Error("cannot remove value from cache")
		return fmt.Errorf("cannot remove value from cache")
	}
	return nil
}
