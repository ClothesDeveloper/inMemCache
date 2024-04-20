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

		return c.fromCacheElement(*element), nil
	}

	return aggregate.Profile{}, ErrValueNotFound
}

func (c *CacheInMemory) Set(uuid string, profile aggregate.Profile, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	newElement := c.toCacheElement(profile, duration)

	c.elements[uuid] = newElement
}

func (c *CacheInMemory) toCacheElement(profile aggregate.Profile, duration time.Duration) *CacheElement {
	orders := make([]*entity.Order, len(profile.Orders))
	for index, order := range profile.Orders {
		orderCopy := *order
		orders[index] = &orderCopy
	}

	profile.Orders = orders

	return &CacheElement{
		profile:   profile,
		expiresAt: time.Now().Add(duration),
	}
}

func (c *CacheInMemory) fromCacheElement(element CacheElement) aggregate.Profile {
	profile := element.profile

	orders := make([]*entity.Order, len(element.profile.Orders))
	for index, order := range element.profile.Orders {
		orderCopy := *order
		orders[index] = &orderCopy
	}

	profile.Orders = orders

	return profile
}
