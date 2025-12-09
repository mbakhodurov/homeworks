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

func (s *OrderStorage) GetAllOrder() ([]*Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	orders := make([]*Order, 0, len(s.orders))

	for _, v := range s.orders {
		orders = append(orders, v)
	}
	return orders, nil
}
