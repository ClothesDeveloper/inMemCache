package services

import (
	"fmt"
	"gotest.tools/v3/assert"
	"inMemoryCache/aggregate"
	"inMemoryCache/domain/profile/memory"
	"inMemoryCache/entity"
	"testing"
	"time"
)

func Test_AddOrder(t *testing.T) {
	profileServ := NewProfileService()

	profile, _ := aggregate.NewProfile("someTestName")
	order := entity.Order{
		UUID:      "someName",
		Value:     "someVal",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	profileServ.AddOrder(&profile, order)
	orders, _ := profileServ.GetOrdersList(profile.UUID)

	var got entity.Order
	for _, o := range orders {
		got = *o
	}

	assert.Equal(t, got, order)
	assert.Equal(t, len(orders), 1)
}

func Test_Order_In_Profile(t *testing.T) {
	profileServ := NewProfileService()

	profile, _ := aggregate.NewProfile("someTestName")
	order := entity.Order{
		UUID:      "someName",
		Value:     "someVal",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	go func() {
		profileServ.AddOrder(&profile, order)
	}()
	profileServ.AddOrder(&profile, order)
	time.Sleep(time.Second * 1)
	orders, _ := profileServ.GetOrdersList(profile.UUID)

	var got entity.Order
	for _, o := range orders {
		got = *o
	}

	fmt.Println(len(orders))

	assert.Equal(t, got, order)
	assert.Equal(t, len(orders), 1)
}

func Test_AddOrder_WithExpiredCache(t *testing.T) {
	profileServ := &ProfileService{
		profileRepository: memory.NewRepository(),
	}

	profile, _ := aggregate.NewProfile("someTestName")
	order := entity.Order{
		UUID:      "someName",
		Value:     "someVal",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	profileServ.AddOrder(&profile, order)

	time.Sleep(3 * time.Second)
	orders, _ := profileServ.GetOrdersList(profile.UUID)

	assert.Equal(t, len(orders), 0)
}

func Test_Delete_Order(t *testing.T) {
	profileServ := NewProfileService()

	profile, _ := aggregate.NewProfile("dude")
	order := entity.Order{
		UUID:      "someName",
		Value:     "someVal",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	profileServ.AddOrder(&profile, order)

	secondOrder := entity.Order{
		UUID:      "test2",
		Value:     "someV",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	profileServ.AddOrder(&profile, secondOrder)

	profileServ.DeleteOrder(profile.UUID, order)

	userCreatedOrders, _ := profileServ.GetOrdersList(profile.UUID)

	assert.Equal(t, len(userCreatedOrders), 1)
	assert.Equal(t, *userCreatedOrders[0], secondOrder)
}

//func Test_Update_Get_Simultaneously(t *testing.T) {
//	TODO: try to simulate race condition
//}
