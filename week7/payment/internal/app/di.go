package app

import (
	"context"
	"fmt"

	v1 "github.com/mbakhodurov/homeworks/week7/payment/internal/api/payment/v1"
	"github.com/mbakhodurov/homeworks/week7/payment/internal/config"
	"github.com/mbakhodurov/homeworks/week7/payment/internal/service"
	paymentService "github.com/mbakhodurov/homeworks/week7/payment/internal/service/payment"
	"github.com/mbakhodurov/homeworks/week7/platform/pkg/closer"
	auth_v1 "github.com/mbakhodurov/homeworks/week7/shared/pkg/proto/auth/v1"
	payment_v1 "github.com/mbakhodurov/homeworks/week7/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type diContainer struct {
	iamGRPCConn *grpc.ClientConn
	iamClient   auth_v1.AuthServiceClient

	paymentV1Api   payment_v1.PaymentServiceServer
	paymentService service.PaymentService
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (di *diContainer) IAMClient(ctx context.Context) auth_v1.AuthServiceClient {
	if di.iamClient == nil {
		di.iamClient = auth_v1.NewAuthServiceClient(di.IAMGRPCConn(ctx))
	}
	return di.iamClient
}

func (di *diContainer) IAMGRPCConn(_ context.Context) *grpc.ClientConn {
	if di.iamGRPCConn == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().IAMGRPCClient.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Sprintf("💥 failed to connect to IAM service: %v", err))
		}

		closer.AddNamed("IAM gRPC connection", func(ctx context.Context) error {
			return conn.Close()
		})

		di.iamGRPCConn = conn
	}
	return di.iamGRPCConn
}

func (di *diContainer) PaymentV1Api(ctx context.Context) payment_v1.PaymentServiceServer {
	if di.paymentV1Api == nil {
		di.paymentV1Api = v1.NewApi(di.PaymentService(ctx))
	}
	return di.paymentV1Api
}

func (di *diContainer) PaymentService(ctx context.Context) service.PaymentService {
	if di.paymentService == nil {
		di.paymentService = paymentService.NewService()
	}
	return di.paymentService
}
