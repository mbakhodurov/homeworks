package user

import (
	"database/sql"

	def "github.com/mbakhodurov/homeworks/week6/iam/internal/repository"
)

var _ def.UserRepository = (*repository)(nil)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	return &repository{
		db: db,
	}
}
