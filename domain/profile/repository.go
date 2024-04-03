package profile

import "inMemoryCache/aggregate"

type ProfileRepository interface {
	Get(string) *aggregate.Profile
	Add(*aggregate.Profile)
	Delete(aggregate.Profile)
}
