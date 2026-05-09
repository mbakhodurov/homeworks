package app

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/mbakhodurov/homeworks/week6/iam/internal/config"
	"github.com/mbakhodurov/homeworks/week6/platform/pkg/closer"
	"github.com/mbakhodurov/homeworks/week6/platform/pkg/grpc/health"
	"github.com/mbakhodurov/homeworks/week6/platform/pkg/logger"
	auth_v1 "github.com/mbakhodurov/homeworks/week6/shared/pkg/proto/auth/v1"
	user_v1 "github.com/mbakhodurov/homeworks/week6/shared/pkg/proto/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type App struct {
	diContainer *diContainer
	grpcServer  *grpc.Server
	listener    net.Listener
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	return a.runGRPCServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initDatabase,
		a.initRedis,
		a.initListener,
		a.initGRPCServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initRedis(ctx context.Context) error {
	_ = a.diContainer.RedisClient(ctx)
	logger.Info(ctx, "🔴 Redis initialized")
	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
	)
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initDatabase(ctx context.Context) error {
	_ = a.diContainer.PgxConn(ctx)

	migratorRunner := a.diContainer.MigrationRunner(ctx)
	err := migratorRunner.Up(ctx)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	logger.Info(ctx, "🗄️ Database initialized and migrations applied")
	return nil
}

func (a *App) initListener(_ context.Context) error {
	listener, err := net.Listen("tcp", config.AppConfig().IAMGRPC.Address())
	if err != nil {
		return err
	}
	closer.AddNamed("TCP listener", func(ctx context.Context) error {
		err = listener.Close()
		if err != nil && !errors.Is(err, net.ErrClosed) {
			return err
		}

		return nil
	})

	a.listener = listener
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	closer.AddNamed("gRPC server", func(ctx context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	reflection.Register(a.grpcServer)

	health.RegisterService(a.grpcServer)

	user_v1.RegisterUserServiceServer(a.grpcServer, a.diContainer.UserV1Api(ctx))
	auth_v1.RegisterAuthServiceServer(a.grpcServer, a.diContainer.AuthV1API(ctx))

	return nil
}

func (a *App) runGRPCServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 gRPC IAMService server listening on %s", config.AppConfig().IAMGRPC.Address()))

	err := a.grpcServer.Serve(a.listener)
	if err != nil {
		return err
	}

	return nil
}
