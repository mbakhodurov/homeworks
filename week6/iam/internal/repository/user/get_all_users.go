package user

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/mbakhodurov/homeworks/week6/iam/internal/model"
	"github.com/mbakhodurov/homeworks/week6/iam/internal/repository/converter"
	repomodel "github.com/mbakhodurov/homeworks/week6/iam/internal/repository/model"
)

func (r *repository) GetAllUsers(ctx context.Context) ([]model.User, error) {
	builder := sq.Select("user_uuid", "info", "created_at", "updated_at", "deleted_at", "password_hash").
		From("users").
		Where(sq.Eq{"deleted_at": nil}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v\n", err)
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u repomodel.User
		err = rows.Scan(
			&u.UserUUID,
			&u.UserInfo,
			&u.CreatedAt,
			&u.UpdatedAt,
			&u.DeletedAt,
			&u.Password,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, converter.UserRepoModelToModel(u))
	}
	if len(users) == 0 {
		return nil, model.ErrUserNotFound
	}

	return users, nil
}
