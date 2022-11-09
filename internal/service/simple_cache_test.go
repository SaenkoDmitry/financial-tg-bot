package service

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_simpleCache_Add(t *testing.T) {
	c := &simpleCache{
		inner: make(map[string]*ValueData),
	}
	key1 := "CURRENCY_2022-01-02"
	key2 := "CURRENCY_2022-01-03"
	res1 := decimal.NewFromFloat(0.16)
	res2 := decimal.NewFromFloat(0.18)
	err1 := c.Add(key1, res1, 10*time.Millisecond)
	err2 := c.Add(key2, res2, 10*time.Millisecond)
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	v, ok := c.Get(key1)
	assert.Equal(t, v, res1)
	assert.Equal(t, ok, true)
	time.Sleep(20 * time.Millisecond)
	v2, ok2 := c.Get(key1)
	assert.Equal(t, v2, nil)
	assert.Equal(t, ok2, false)
}
