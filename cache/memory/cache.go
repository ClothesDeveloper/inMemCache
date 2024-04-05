package memory

import (
	"context"
	"fmt"
	"inMemoryCache/aggregate"
	"sync"
	"time"
)

type CacheInMemory struct {
	ctx               context.Context
	cancelCtx         context.CancelFunc
	expiredCachesChan chan string
	elements          map[string]*CacheElement
	rwMx              sync.RWMutex
}

type CacheElement struct {
	profile              *aggregate.Profile
	cancelClearCacheFunc func()
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
		fmt.Println("Cleanup instance is started")
		instance.Cleanup()
	}(ctx)

	return instance
}

// Cleanup removes expired caches
func (c *CacheInMemory) Cleanup() {
	for uuid := range c.expiredCachesChan {
		c.rwMx.Lock()

		delete(c.elements, uuid)
		fmt.Printf("Dropped cache for uuid: %v", uuid)
		c.rwMx.Unlock()
	}
}

func (c *CacheInMemory) Get(uuid string) *aggregate.Profile {
	c.rwMx.Lock()
	defer c.rwMx.Unlock()
	if element, ok := c.elements[uuid]; ok {
		return element.profile
	}

	return nil
}

func (c *CacheInMemory) Set(uuid string, profile *aggregate.Profile, duration time.Duration) {
	c.rwMx.Lock()
	defer c.rwMx.Unlock()
	element := c.elements[uuid]
	//cancelling previous goroutine, to prevent deleting new value
	if element != nil {
		element.cancelClearCacheFunc()
	}

	ctx, cancel := context.WithCancel(c.ctx)

	newElement := &CacheElement{
		profile:              profile,
		cancelClearCacheFunc: cancel,
	}
	c.elements[uuid] = newElement

	//goroutine sends expired cache to expired cache chan
	go func(ctx context.Context) {
		<-time.After(duration)
		c.expiredCachesChan <- uuid
	}(ctx)
}
