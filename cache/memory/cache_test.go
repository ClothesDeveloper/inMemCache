package memory

import (
	"fmt"
	"gotest.tools/v3/assert"
	"inMemoryCache/aggregate"
	"inMemoryCache/entity"
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
		fmt.Println("error on creating new profile")
	}
	cache.Set(profile.UUID, &profile, 1*time.Second)
	time.Sleep(time.Second * 2)

	got := cache.Get(profile.UUID)
	var want *aggregate.Profile

	assert.Equal(t, got, want)
}

//func Test_Cleanup(t *testing.T) {
//	cache := New()
//
//	fakeOrders := []*entity.Order{
//		{
//			UUID:      "123",
//			Value:     111,
//			CreatedAt: time.Now(),
//			UpdatedAt: time.Now(),
//		},
//	}
//
//	testCases := []struct {
//		orders      []*entity.Order
//		profileName string
//		duration    uint64
//	}{
//		{
//			orders:      fakeOrders,
//			profileName: "profile1",
//			duration:    3,
//		},
//		{
//			orders:      fakeOrders,
//			profileName: "profile2",
//			duration:    3,
//		},
//		{
//			orders:      fakeOrders,
//			profileName: "profile3",
//			duration:    3,
//		},
//	}
//
//	for _, testCase := range testCases {
//		profile, _ := aggregate.NewProfile(testCase.profileName)
//		cache.Set(profile.UUID, &profile, time.Duration(testCase.duration))
//	}
//
//	length := len(cache.elements)
//	assert.Equal(t, length, 3)
//
//	cache.Cleanup()
//	lengthAfterCleanup := len(cache.elements)
//	assert.Equal(t, lengthAfterCleanup, 0)
//}

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
			cache.Set(profile.UUID, &profile, 2*time.Second)
		}()
	}

	time.Sleep(1 * time.Second)
}
