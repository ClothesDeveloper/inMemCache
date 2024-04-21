package aggregate

import (
	"github.com/google/uuid"
	"inMemoryCache/entity"
)

type Profile struct {
	UUID   string
	Name   string
	Orders []*entity.Order
}

func NewProfile(name string) (Profile, error) {
	return Profile{
		UUID:   uuid.New().String(),
		Name:   name,
		Orders: make([]*entity.Order, 0),
	}, nil
}

func (p *Profile) IsOrderInList(order entity.Order) bool {
	if p.Orders == nil {
		return false
	}

	for _, orderInProfile := range p.Orders {
		if order == *orderInProfile {
			return true
		}
	}

	return false
}
