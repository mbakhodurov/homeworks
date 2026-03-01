package order

import (
	"database/sql"

	def "github.com/mbakhodurov/homeworks/week4/order/internal/repository"
)

type repository struct {
	db *sql.DB
}

var _ def.OrderRepository = (*repository)(nil)

func NewRepository(db *sql.DB) *repository {
	return &repository{
		db: db,
	}
}
