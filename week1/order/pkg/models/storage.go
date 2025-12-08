package models

import (
	"errors"
	"sync"
)

var ErrNotFound = errors.New("not found")

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*Order
}

func NewStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*Order),
	}
}

func (s *OrderStorage) CreateOrder(o *Order) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[o.OrderUUID] = o
}

func (s *OrderStorage) GetOrderByUUID(uuid string) (*Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[uuid]
	if !ok {
		return nil, ErrNotFound
	}
	return order, nil
}
