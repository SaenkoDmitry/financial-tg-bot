package service

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

type simpleCache struct {
	inner             map[string]*ValueData
	defaultExpiration time.Duration
	mutex             sync.RWMutex
}

type ValueData struct {
	data    interface{}
	expires int64
}

func (v *ValueData) expired(time int64) bool {
	if v.expires == 0 {
		return false
	}
	return time > v.expires
}

func NewSimpleCache(ctx context.Context, defaultExpiration, cleanupInterval time.Duration) *simpleCache {
	c := &simpleCache{
		inner:             make(map[string]*ValueData),
		defaultExpiration: defaultExpiration,
	}
	go func() {
		t := time.NewTicker(cleanupInterval)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				fmt.Println("graceful shutdown")
				break
			case <-t.C:
				c.mutex.Lock()
				for k, v := range c.inner {
					if v.expired(time.Now().UnixNano()) {
						log.Printf("%v has expires at %d", k, time.Now().UnixNano())
						delete(c.inner, k)
					}
				}
				c.mutex.Unlock()
			}
		}
	}()
	return c
}

func (c *simpleCache) Add(k string, value interface{}, d time.Duration) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.inner[k] = &ValueData{
		data:    value,
		expires: time.Now().Add(d).UnixNano(),
	}
	return nil
}

func (c *simpleCache) Get(k string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if value, ok := c.inner[k]; ok {
		if value.expired(time.Now().UnixNano()) {
			return nil, false
		}
		return value.data, true
	}
	return nil, false
}
