package memory

import (
	"context"
	"fmt"
	"inMemoryCache/aggregate"
	"sync"
	"time"
)

type CacheInMemory struct {
	ctx       context.Context
	cancelCtx context.CancelFunc
	elements  map[string]*CacheElement
	mx        sync.Mutex
}

type CacheElement struct {
	profile              *aggregate.Profile
	cancelClearCacheFunc func()
}

func New() *CacheInMemory {
	ctx, cancel := context.WithCancel(context.Background())
	return &CacheInMemory{
		ctx:       ctx,
		elements:  make(map[string]*CacheElement),
		cancelCtx: cancel,
	}
}

func (c *CacheInMemory) Cleanup() {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.cancelCtx()
	ctx, cancel := context.WithCancel(context.Background())
	c.ctx = ctx
	c.cancelCtx = cancel

	c.elements = make(map[string]*CacheElement)
}

func (c *CacheInMemory) Get(uuid string) *aggregate.Profile {
	c.mx.Lock()
	defer c.mx.Unlock()
	if element, ok := c.elements[uuid]; ok {
		return element.profile
	}

	return nil
}

func (c *CacheInMemory) Set(uuid string, profile *aggregate.Profile, duration time.Duration) {
	c.mx.Lock()
	defer c.mx.Unlock()
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

	//goroutine drops cache by uuid
	go func(ctx context.Context) {
		<-time.After(duration * time.Second)

		delete(c.elements, uuid)
		fmt.Println(fmt.Sprintf("Dropped cache for uuid: %v", uuid))
	}(ctx)
}
