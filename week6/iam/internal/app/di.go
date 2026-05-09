package app

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/mbakhodurov/homeworks/week6/iam/internal/config"
	"github.com/mbakhodurov/homeworks/week6/iam/internal/repository"
	userRepository "github.com/mbakhodurov/homeworks/week6/iam/internal/repository/user"
	userService "github.com/mbakhodurov/homeworks/week6/iam/internal/service/user"
	"github.com/mbakhodurov/homeworks/week6/platform/pkg/cache"
	"github.com/mbakhodurov/homeworks/week6/platform/pkg/cache/redis"
	"github.com/mbakhodurov/homeworks/week6/platform/pkg/logger"
	pgMigrator "github.com/mbakhodurov/homeworks/week6/platform/pkg/migrator/pg"

	"github.com/mbakhodurov/homeworks/week6/platform/pkg/closer"
	"github.com/mbakhodurov/homeworks/week6/platform/pkg/migrator"

	redigo "github.com/gomodule/redigo/redis"
	authV1API "github.com/mbakhodurov/homeworks/week6/iam/internal/api/auth/v1"
	userV1API "github.com/mbakhodurov/homeworks/week6/iam/internal/api/user/v1"
	sessionRepository "github.com/mbakhodurov/homeworks/week6/iam/internal/repository/session"
	"github.com/mbakhodurov/homeworks/week6/iam/internal/service"
	authService "github.com/mbakhodurov/homeworks/week6/iam/internal/service/auth"
	auth_v1 "github.com/mbakhodurov/homeworks/week6/shared/pkg/proto/auth/v1"
	user_v1 "github.com/mbakhodurov/homeworks/week6/shared/pkg/proto/user/v1"
)

type diContainer struct {
	migrationRunner migrator.Migrator

	sqlDb          *sql.DB
	pgxConn        *pgx.Conn
	userRepository repository.UserRepository
	userService    service.UserService
	userV1Api      user_v1.UserServiceServer

	authV1Api         auth_v1.AuthServiceServer
	authService       service.AuthService
	sessionRepository repository.SessionRepository
	redisClient       cache.RedisClient
	redisPool         *redigo.Pool
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) RedisPool() *redigo.Pool {
	if d.redisPool == nil {
		d.redisPool = &redigo.Pool{
			MaxIdle:     config.AppConfig().Redis.MaxIdle(),
			IdleTimeout: config.AppConfig().Redis.IdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", config.AppConfig().Redis.Address())
			},
		}

		closer.AddNamed("Redis pool", func(ctx context.Context) error {
			return d.redisPool.Close()
		})
	}

	return d.redisPool
}

func (di *diContainer) RedisClient(ctx context.Context) cache.RedisClient {
	if di.redisClient == nil {
		di.redisClient = redis.NewClient(di.RedisPool(), logger.Logger(), config.AppConfig().Redis.ConnectionTimeout())
	}
	return di.redisClient
}

func (di *diContainer) SessionRepository(ctx context.Context) repository.SessionRepository {
	if di.sessionRepository == nil {
		di.sessionRepository = sessionRepository.NewRepository(di.RedisClient(ctx))
	}
	return di.sessionRepository
}

func (d *diContainer) AuthService(ctx context.Context) service.AuthService {
	if d.authService == nil {
		d.authService = authService.NewService(d.SessionRepository(ctx), d.UserRepository(ctx), config.AppConfig().Redis.CacheTTL())
	}

	return d.authService
}

func (di *diContainer) AuthV1API(ctx context.Context) auth_v1.AuthServiceServer {
	if di.authV1Api == nil {
		di.authV1Api = authV1API.NewAPI(di.AuthService(ctx))
	}
	return di.authV1Api
}

func (di *diContainer) MigrationRunner(ctx context.Context) migrator.Migrator {
	if di.migrationRunner == nil {
		di.migrationRunner = pgMigrator.NewMigrator(di.SqlDB(ctx), config.AppConfig().Postgres.MigrationDir())
	}
	return di.migrationRunner
}

func (d *diContainer) SqlDB(ctx context.Context) *sql.DB {
	if d.sqlDb == nil {
		d.sqlDb = stdlib.OpenDB(*d.PgxConn(ctx).Config().Copy())
	}

	return d.sqlDb
}

func (d *diContainer) PgxConn(ctx context.Context) *pgx.Conn {
	if d.pgxConn == nil {
		conn, err := pgx.Connect(ctx, config.AppConfig().Postgres.URL())
		if err != nil {
			panic(fmt.Sprintf("💥 failed to connect to database: %v", err))
		}

		err = conn.Ping(ctx)
		if err != nil {
			panic(fmt.Sprintf("💥 failed to ping database: %v", err))
		}

		closer.AddNamed("PostgresSQL connection", func(ctx context.Context) error {
			return conn.Close(ctx)
		})

		d.pgxConn = conn
	}

	return d.pgxConn
}

func (di *diContainer) UserRepository(ctx context.Context) repository.UserRepository {
	if di.userRepository == nil {
		di.userRepository = userRepository.NewRepository(di.SqlDB(ctx))
	}
	return di.userRepository
}

func (di *diContainer) UserService(ctx context.Context) service.UserService {
	if di.userService == nil {
		di.userService = userService.NewService(di.UserRepository(ctx))
	}
	return di.userService
}

func (di *diContainer) UserV1Api(ctx context.Context) user_v1.UserServiceServer {
	if di.userV1Api == nil {
		di.userV1Api = userV1API.NewAPI(di.UserService(ctx))
	}
	return di.userV1Api
}
