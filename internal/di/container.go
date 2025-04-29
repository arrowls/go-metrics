package di

import (
	"fmt"
	"sync"
)

type ContainerInterface interface {
	Add(key string, value any) error
	Get(key string) any
}

type Container struct {
	mu   sync.RWMutex
	deps map[string]any
}

func NewContainer() ContainerInterface {
	return &Container{
		deps: make(map[string]any),
	}
}

func (c *Container) Add(key string, value any) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.deps[key]; ok {
		return fmt.Errorf("container already has a dependency %s", key)
	}

	c.deps[key] = value

	return nil
}

func (c *Container) Get(key string) any {
	c.mu.RLock()
	defer c.mu.RUnlock()

	d, ok := c.deps[key]

	if !ok {
		return nil
	}
	return d
}
