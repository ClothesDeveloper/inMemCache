package memory

import (
	"errors"
	"inMemoryCache/aggregate"
	"inMemoryCache/cache"
	"inMemoryCache/cache/memory"
	"inMemoryCache/domain/profile"
	"time"
)

const profileCacheTTL = 2 * time.Second

var (
	ErrProfileNotFound = errors.New("profile not found")
)

type Repository struct {
	cache cache.Cache
}

func NewRepository() profile.ProfileRepository {
	return &Repository{
		cache: memory.New(),
	}
}

func (r *Repository) Add(profile aggregate.Profile) {
	r.cache.Set(profile.UUID, profile, profileCacheTTL)
}

func (r *Repository) Get(uuid string) (aggregate.Profile, error) {
	cacheVal, err := r.cache.Get(uuid)
	if err != nil {
		return aggregate.Profile{}, ErrProfileNotFound
	}
	return cacheVal, nil
}

func (r *Repository) Delete(profile aggregate.Profile) {
	r.cache.Set(profile.UUID, profile, profileCacheTTL)
}
