package cache

import (
	"log"
	"sync"
)

type PayloadCache struct {
	data  map[string]map[string]any
	mutex sync.RWMutex
}

func NewPayloadCache() *PayloadCache {
	return &PayloadCache{
		data:  make(map[string]map[string]any),
		mutex: sync.RWMutex{},
	}
}

func (c *PayloadCache) Set(key string, value map[string]any) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data[key] = value
}

func (c *PayloadCache) Get(key string) (map[string]any, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	value, ok := c.data[key]
	log.Println("cache is found: ", ok)
	return value, ok
}

func (c *PayloadCache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.data, key)
}

func (c *PayloadCache) GetAll() map[string]map[string]any {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.data
}

func (c *PayloadCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data = make(map[string]map[string]any)
}
