package order

import (
	"database/sql"

	def "github.com/mbakhodurov/homeworks/week7/order/internal/repository"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	return &repository{
		db: db,
	}
}
