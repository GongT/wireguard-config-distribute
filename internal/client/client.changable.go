package client

import "sync"

type editableConfig struct {
	services     []string
	externalPort uint32
	_mutex       sync.Mutex
}

func (c *editableConfig) lock() {
	c._mutex.Lock()
}

func (c *editableConfig) unlock() {
	c._mutex.Unlock()
}
