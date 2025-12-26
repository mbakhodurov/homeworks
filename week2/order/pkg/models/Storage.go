package models

import (
	"sync"
)

type OrderStorage struct {
	mu sync.RWMutex

	orders map[string]*Order
}

type NewOrderStorage struct {
}
