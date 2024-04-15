package memory

import (
	"context"
	"errors"
	"inMemoryCache/aggregate"
	"sync"
	"time"
)

var (
	ErrValueNotFound = errors.New("value not found")
)

type CacheInMemory struct {
	ctx               context.Context
	cancelCtx         context.CancelFunc
	expiredCachesChan chan string
	elements          map[string]*CacheElement
	mu                sync.RWMutex
}

type CacheElement struct {
	profile   aggregate.Profile
	expiresAt time.Time
}

func New() *CacheInMemory {
	ctx, cancel := context.WithCancel(context.Background())

	instance := &CacheInMemory{
		ctx:               ctx,
		elements:          make(map[string]*CacheElement),
		expiredCachesChan: make(chan string, 1),
		cancelCtx:         cancel,
	}

	go func(ctx context.Context) {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				instance.cleanup()
			}
		}
	}(ctx)

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
		return element.profile, nil
	}

	return aggregate.Profile{}, ErrValueNotFound
}

func (c *CacheInMemory) Set(uuid string, profile aggregate.Profile, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	newElement := &CacheElement{
		profile:   profile,
		expiresAt: time.Now().Add(duration),
	}

	c.elements[uuid] = newElement
}
