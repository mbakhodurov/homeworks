package order

import (
	"database/sql"

	def "github.com/mbakhodurov/homeworks/week3/order/internal/repository"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *repository {
	return &repository{
		db: db,
	}
}
