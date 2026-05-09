package v1

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mbakhodurov/homeworks/week6/iam/internal/model"
	"github.com/mbakhodurov/homeworks/week6/platform/pkg/logger"
	auth_v1 "github.com/mbakhodurov/homeworks/week6/shared/pkg/proto/auth/v1"
)

func (a *api) Login(ctx context.Context, req *auth_v1.LoginRequest) (*auth_v1.LoginResponse, error) {
	sessionUUID, err := a.authService.Login(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		if errors.Is(err, model.ErrInvalidCredentials) {
			logger.Error(ctx, "invalid credentials", zap.Error(err))
			return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
		}
		if errors.Is(err, model.ErrSessionBadRequest) {
			logger.Error(ctx, "bad request during login", zap.Error(err))
			return nil, status.Errorf(codes.InvalidArgument, "bad request during login")
		}
		logger.Error(ctx, "failed to login", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal error during login")
	}

	return &auth_v1.LoginResponse{
		SessionUuid: sessionUUID,
	}, nil
}
