package app

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/mbakhodurov/homeworks/week7/payment/internal/config"
	"github.com/mbakhodurov/homeworks/week7/payment/internal/interceptor"
	"github.com/mbakhodurov/homeworks/week7/platform/pkg/closer"
	"github.com/mbakhodurov/homeworks/week7/platform/pkg/grpc/health"
	"github.com/mbakhodurov/homeworks/week7/platform/pkg/logger"
	grpcMiddleware "github.com/mbakhodurov/homeworks/week7/platform/pkg/middleware/grpc"
	payment_v1 "github.com/mbakhodurov/homeworks/week7/shared/pkg/proto/payment/v1"
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
		a.initListener,
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

func (a *App) initListener(ctx context.Context) error {
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

func (a *App) initDI(ctx context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initLogger(ctx context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJSON(),
		config.AppConfig().Logger.EnableOTLP(),
		config.AppConfig().Logger.OtlpEndpoint(),
		config.AppConfig().Logger.ServiceName(),
		config.AppConfig().Logger.ServiceEnv(),
	)
}

func (a *App) initCloser(ctx context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initGrpcServer(ctx context.Context) error {
	// a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	authInterceptor := grpcMiddleware.NewAuthInterceptor(a.diContainer.IAMClient(ctx))
	// fmt.Println("fdsa")
	// str, _ := grpcMiddleware.GetSessionUUIDFromContext(ctx)
	// fmt.Println("str:", str)
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
	payment_v1.RegisterPaymentServiceServer(a.grpcServer, a.diContainer.PaymentV1Api(ctx))

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
