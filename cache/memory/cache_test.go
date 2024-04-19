package memory

import (
	"gotest.tools/v3/assert"
	"inMemoryCache/aggregate"
	"inMemoryCache/entity"
	"log"
	"strconv"
	"testing"
	"time"
)

func Test_Cache_Expires_Properly(t *testing.T) {
	cache := New()

	orders := []*entity.Order{
		{
			UUID:      "123",
			Value:     111,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	profile, err := aggregate.NewProfile("testName")
	profile.Orders = orders
	if err != nil {
		log.Println(err)
	}
	cache.Set(profile.UUID, profile, 1*time.Second)
	time.Sleep(time.Second * 2)

	_, err = cache.Get(profile.UUID)

	assert.Equal(t, ErrValueNotFound, err)
}

func Test_Cleanup(t *testing.T) {
	cache := New()

	fakeOrders := []*entity.Order{
		{
			UUID:      "123",
			Value:     111,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	testCases := []struct {
		orders      []*entity.Order
		profileName string
		duration    time.Duration
	}{
		{
			orders:      fakeOrders,
			profileName: "profile1",
			duration:    3 * time.Second,
		},
		{
			orders:      fakeOrders,
			profileName: "profile2",
			duration:    10 * time.Second,
		},
		{
			orders:      fakeOrders,
			profileName: "profile3",
			duration:    3 * time.Second,
		},
	}

	for _, testCase := range testCases {
		profile, _ := aggregate.NewProfile(testCase.profileName)
		cache.Set(profile.UUID, profile, testCase.duration)
	}

	length := len(cache.elements)
	assert.Equal(t, length, 3)

	time.Sleep(3500 * time.Millisecond)
	lengthAfterCleanup := len(cache.elements)
	assert.Equal(t, lengthAfterCleanup, 1)
}

func Test_Concurrent_Sets_Not_Allowed(t *testing.T) {
	cache := New()

	fakeOrders := []*entity.Order{
		{
			UUID:      "123",
			Value:     111,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for i := 0; i < 10; i++ {
		i := i
		go func() {
			profile, _ := aggregate.NewProfile("test" + strconv.Itoa(i))
			profile.Orders = fakeOrders
			cache.Set(profile.UUID, profile, 2*time.Second)
		}()
	}

	time.Sleep(1 * time.Second)
}

func Test_Profile_IsNotModifiedInCache(t *testing.T) {
	cache := New()

	fakeOrders := []*entity.Order{
		{
			UUID:      "123",
			Value:     111,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			UUID:      "4546",
			Value:     233,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	profile, err := aggregate.NewProfile("dude")
	if err != nil {
		log.Println(err)
	}

	profile.Orders = fakeOrders

	cache.Set(profile.UUID, profile, 5*time.Second)
	profile.Name = "chopo"

	profileFromCache, _ := cache.Get(profile.UUID)
	assert.Equal(t, profileFromCache.Name, "dude")
}

func Test_ProfileOrders_NotModifiedInCache(t *testing.T) {
	cache := New()
	fakeOrders := []*entity.Order{
		{
			UUID:      "123",
			Value:     111,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			UUID:      "4546",
			Value:     233,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	profile, err := aggregate.NewProfile("dude")
	if err != nil {
		log.Println(err)
	}

	profile.Orders = fakeOrders

	cache.Set(profile.UUID, profile, 5*time.Second)
	fakeOrders[0].Value = 2

	profileFromCache, _ := cache.Get(profile.UUID)

	assert.Equal(t, profileFromCache.Orders[0].Value, 111)
}
