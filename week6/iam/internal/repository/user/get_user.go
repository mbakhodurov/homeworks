package user

import (
	"context"
	"database/sql"
	"errors"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/mbakhodurov/homeworks/week6/iam/internal/model"
	"github.com/mbakhodurov/homeworks/week6/iam/internal/repository/converter"
	repomodel "github.com/mbakhodurov/homeworks/week6/iam/internal/repository/model"
)

func (r *repository) GetUser(ctx context.Context, userUUID string) (model.User, error) {
	builder := sq.Select("user_uuid", "info", "created_at", "updated_at", "deleted_at", "password_hash").
		From("users").
		Where(sq.Eq{"user_uuid": userUUID}).
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
		&user.DeletedAt,
		&user.Password,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, model.ErrUserNotFound
		}
		return model.User{}, err
	}

	return converter.UserRepoModelToModel(user), nil
	// return model.User{
	// 	UserUUID: user.UserUUID,
	// 	UserInfo: model.UserInfo{
	// 		Login: user.UserInfo.Login,
	// 		Email: user.UserInfo.Email,
	// 	},
	// 	CreatedAt: user.CreatedAt,
	// 	UpdatedAt: user.UpdatedAt,
	// 	DeletedAt: user.DeletedAt,
	// }, nil
}
