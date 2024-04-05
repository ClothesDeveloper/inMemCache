package memory

import (
	"inMemoryCache/aggregate"
	"inMemoryCache/cache"
	"inMemoryCache/cache/memory"
	"inMemoryCache/domain/profile"
	"time"
)

const profileCacheTTL = 2 * time.Second

type Repository struct {
	cache cache.Cache
}

func NewRepository() profile.ProfileRepository {
	return &Repository{
		cache: memory.New(),
	}
}

func (r *Repository) Add(profile *aggregate.Profile) {
	r.cache.Set(profile.UUID, profile, profileCacheTTL)
}

func (r *Repository) Get(uuid string) *aggregate.Profile {
	return r.cache.Get(uuid)
}

func (r *Repository) Delete(profile aggregate.Profile) {
	r.cache.Set(profile.UUID, nil, profileCacheTTL)
}
