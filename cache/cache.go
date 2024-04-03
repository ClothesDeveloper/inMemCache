package cache

import (
	"inMemoryCache/aggregate"
	"time"
)

// Cache describes cache that works only with Profiles
// Cleanup flushes all data
type Cache interface {
	Get(string) *aggregate.Profile
	Set(string, *aggregate.Profile, time.Duration)
	Cleanup()
}
