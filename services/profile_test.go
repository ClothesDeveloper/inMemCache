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

	err := profileServ.AddOrder(&profile, order)
	if err != nil {
		fmt.Printf("Error on add order to profile occured %v", err)
		return
	}
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
		err := profileServ.AddOrder(&profile, order)
		fmt.Printf("Error on add order to profile occured %v", err)
	}()
	err := profileServ.AddOrder(&profile, order)
	if err != nil {
		fmt.Printf("Error on add order to profile occured %v", err)
		return
	}
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

	err := profileServ.AddOrder(&profile, order)
	if err != nil {
		fmt.Printf("Error on add order to profile occured %v", err)
		return
	}

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

	err := profileServ.AddOrder(&profile, order)
	if err != nil {
		fmt.Printf("Error on add order to profile occured %v", err)
		return
	}

	secondOrder := entity.Order{
		UUID:      "test2",
		Value:     "someV",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = profileServ.AddOrder(&profile, secondOrder)
	if err != nil {
		fmt.Printf("Error on add order to profile occured %v", err)
		return
	}

	profileServ.DeleteOrder(profile.UUID, order)

	userCreatedOrders, _ := profileServ.GetOrdersList(profile.UUID)

	assert.Equal(t, len(userCreatedOrders), 1)
	assert.Equal(t, *userCreatedOrders[0], secondOrder)
}
