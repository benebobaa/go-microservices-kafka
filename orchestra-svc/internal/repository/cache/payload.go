package cache

import (
	"log"
	"sync"
)

type PayloadCacher struct {
	data  map[string]map[string]any
	mutex sync.RWMutex
}

func NewPayloadCache() *PayloadCacher {
	return &PayloadCacher{
		data:  make(map[string]map[string]any),
		mutex: sync.RWMutex{},
	}
}

func (c *PayloadCacher) Set(key string, value map[string]any) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data[key] = value
}

func (c *PayloadCacher) Get(key string) (map[string]any, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	value, ok := c.data[key]
	log.Println("cache is found: ", ok)
	return value, ok
}

func (c *PayloadCacher) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.data, key)
}

func (c *PayloadCacher) GetAll() map[string]map[string]any {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.data
}

func (c *PayloadCacher) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data = make(map[string]map[string]any)
}
