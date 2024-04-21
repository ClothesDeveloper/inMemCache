package entity

import (
	"time"
)

type Order struct {
	UUID      string
	Value     any
	CreatedAt time.Time
	UpdatedAt time.Time
}
