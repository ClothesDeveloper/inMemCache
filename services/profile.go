package services

import (
	"errors"
	"inMemoryCache/aggregate"
	"inMemoryCache/domain/profile"
	"inMemoryCache/domain/profile/memory"
	"inMemoryCache/entity"
)

type ProfileService struct {
	profileRepository profile.ProfileRepository
}

func NewProfileService() *ProfileService {
	return &ProfileService{
		profileRepository: memory.NewRepository(),
	}
}

func (ps *ProfileService) AddOrder(profile *aggregate.Profile, order entity.Order) error {
	profileInfo := ps.profileRepository.Get(profile.UUID)
	if profileInfo == nil {
		profileInfo = profile
	}

	if profileInfo.IsOrderInList(order) {
		return errors.New("order is in list already")
	}
	profileInfo.Orders = append(profileInfo.Orders, &order)
	ps.profileRepository.Add(profileInfo)
	return nil
}

func (ps *ProfileService) UpdateOrder(profileUUID string, order entity.Order) error {
	profileInfo := ps.profileRepository.Get(profileUUID)
	if profileInfo == nil {
		return errors.New("profile is not found")
	}
	profileOrders := profileInfo.Orders
	if len(profileOrders) == 0 {
		return errors.New("there are no orders in profile")
	}

	for index, o := range profileOrders {
		if o.UUID == order.UUID {
			profileOrders[index] = &order
			return nil
		}
	}

	return errors.New("didn't find order to update")
}

func (ps *ProfileService) DeleteOrder(profileUUID string, order entity.Order) {
	profile := ps.profileRepository.Get(profileUUID)
	savedOrders := profile.Orders
	for index, savedOrder := range savedOrders {
		if order == *savedOrder {
			profile.Orders = append(savedOrders[:index], savedOrders[index+1:]...)
		}
	}

	ps.profileRepository.Add(profile)
}

func (ps *ProfileService) GetOrdersList(profileUUID string) ([]*entity.Order, error) {
	p := ps.profileRepository.Get(profileUUID)
	if p == nil {
		return []*entity.Order{}, errors.New("profile does not have orders")
	}

	return p.Orders, nil
}
