package user

import (
	"context"
	"log"

	"github.com/mbakhodurov/homeworks/week7/iam/internal/repository/converter"
	repomodel "github.com/mbakhodurov/homeworks/week7/iam/internal/repository/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/mbakhodurov/homeworks/week7/iam/internal/model"
	"golang.org/x/crypto/bcrypt"
)

func (r *repository) GetUserByLogin(ctx context.Context, login, password string) (model.User, error) {
	builder := sq.Select("user_uuid", "info", "created_at", "updated_at", "password_hash").
		From("users").
		Where("info->>'login' = ?", login).
		PlaceholderFormat(sq.Dollar).
		Limit(1)
	query, args, err := builder.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v\n", err)
		return model.User{}, err
	}

	var user repomodel.User
	err = r.db.QueryRowContext(ctx, query, args...).Scan(
		&user.UserUUID,
		&user.UserInfo,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Password,
	)
	if err != nil {
		log.Printf("failed to select user: %v\n", err)
		return model.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return model.User{}, repomodel.ErrUserNotFound
	}

	return converter.UserRepoModelToModel(user), nil
}
