package app

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/mbakhodurov/homeworks/week6/payment/internal/config"
	"github.com/mbakhodurov/homeworks/week6/payment/internal/interceptor"
	"github.com/mbakhodurov/homeworks/week6/platform/pkg/closer"
	"github.com/mbakhodurov/homeworks/week6/platform/pkg/grpc/health"
	"github.com/mbakhodurov/homeworks/week6/platform/pkg/logger"
	grpcMiddleware "github.com/mbakhodurov/homeworks/week6/platform/pkg/middleware/grpc"
	payment_v1 "github.com/mbakhodurov/homeworks/week6/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type App struct {
	grpcServer *grpc.Server
	di         *diContainer
	listener   net.Listener
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
		a.initDi,
		a.initLogger,
		a.initCloser,
		a.initListeners,
		a.initGrpcServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initListeners(ctx context.Context) error {
	listener, err := net.Listen("tcp", config.AppConfig().PaymentGRPC.Address())
	if err != nil {
		return err
	}

	closer.AddNamed("TCP listener", func(ctx context.Context) error {
		lerr := listener.Close()
		if lerr != nil && !errors.Is(lerr, net.ErrClosed) {
			return lerr
		}

		return nil
	})

	a.listener = listener
	return nil
}

func (a *App) initDi(ctx context.Context) error {
	a.di = NewDiContainer()
	return nil
}

func (a *App) initLogger(ctx context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
	)
}

func (a *App) initCloser(ctx context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initGrpcServer(ctx context.Context) error {
	authInterceptor := grpcMiddleware.NewAuthInterceptor(a.di.IAMClient(ctx))

	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.ChainUnaryInterceptor(
			interceptor.LoggerInterceptor(),
			authInterceptor.Unary(),
		),
	)

	closer.AddNamed("gRPC server", func(ctx context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	reflection.Register(a.grpcServer)
	health.RegisterService(a.grpcServer)
	payment_v1.RegisterPaymentServiceServer(a.grpcServer, a.di.PaymentV1Api(ctx))
	return nil
}

func (a *App) runGRPCServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 gRPC PaymentService server listening on %s", config.AppConfig().PaymentGRPC.Address()))

	err := a.grpcServer.Serve(a.listener)
	if err != nil {
		return err
	}

	return nil
}
