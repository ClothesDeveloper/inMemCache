package memory

import (
	"errors"
	"inMemoryCache/aggregate"
	"inMemoryCache/entity"
	"sync"
	"time"
)

var (
	ErrValueNotFound = errors.New("value not found")
)

type CacheInMemory struct {
	expiredCachesChan chan string
	elements          map[string]*CacheElement
	mu                sync.RWMutex
}

type CacheElement struct {
	profile   aggregate.Profile
	expiresAt time.Time
}

func New() *CacheInMemory {
	instance := &CacheInMemory{
		elements:          make(map[string]*CacheElement),
		expiredCachesChan: make(chan string, 1),
	}

	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				instance.cleanup()
			}
		}
	}()

	return instance
}

// Cleanup removes expired caches
func (c *CacheInMemory) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for uuid, element := range c.elements {
		if element.expiresAt.Before(time.Now()) {
			delete(c.elements, uuid)
		}
	}
}

func (c *CacheInMemory) Get(uuid string) (aggregate.Profile, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if element, ok := c.elements[uuid]; ok {
		if element.expiresAt.Before(time.Now()) {
			return aggregate.Profile{}, ErrValueNotFound
		}

		return element.profile, nil
	}

	return aggregate.Profile{}, ErrValueNotFound
}

func (c *CacheInMemory) Set(uuid string, profile aggregate.Profile, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	orders := make([]*entity.Order, len(profile.Orders))
	for index, order := range profile.Orders {
		orderCopy := *order
		orders[index] = &orderCopy
	}

	profile.Orders = orders

	newElement := &CacheElement{
		profile:   profile,
		expiresAt: time.Now().Add(duration),
	}

	c.elements[uuid] = newElement
}
