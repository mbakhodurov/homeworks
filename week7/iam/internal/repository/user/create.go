package user

import (
	"context"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/mbakhodurov/homeworks/week7/iam/internal/model"
	"github.com/mbakhodurov/homeworks/week7/iam/internal/repository/converter"
	repoModel "github.com/mbakhodurov/homeworks/week7/iam/internal/repository/model"
	"golang.org/x/crypto/bcrypt"
)

func (r *repository) Create(ctx context.Context, info model.UserInfo, password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := repoModel.User{
		UserUUID:  uuid.NewString(),
		UserInfo:  converter.UserInfoToRepo(info),
		CreatedAt: time.Now(),
		Password:  string(hashedPassword),
	}

	builder := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("user_uuid", "info", "created_at", "password_hash").
		Values(user.UserUUID, user.UserInfo, user.CreatedAt, user.Password).
		Suffix("RETURNING user_uuid")

	query, args, err := builder.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v\n", err)
		return "", err
	}

	var uuid string
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&uuid)
	if err != nil {
		return "", err
	}

	return uuid, nil
}
