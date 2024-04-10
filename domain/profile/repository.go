package profile

import "inMemoryCache/aggregate"

type ProfileRepository interface {
	Get(uuid string) (aggregate.Profile, error)
	Add(aggregate.Profile)
	Delete(aggregate.Profile)
}
