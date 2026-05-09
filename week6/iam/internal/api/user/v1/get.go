package v1

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/mbakhodurov/homeworks/week6/iam/internal/converter"
	"github.com/mbakhodurov/homeworks/week6/iam/internal/model"
	"github.com/mbakhodurov/homeworks/week6/platform/pkg/logger"
	user_v1 "github.com/mbakhodurov/homeworks/week6/shared/pkg/proto/user/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) GetUser(ctx context.Context, req *user_v1.GetUserRequest) (*user_v1.GetUserResponse, error) {
	_, err := uuid.Parse(req.GetUserUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_uuid")
	}

	res, err := a.userService.GetUser(ctx, req.GetUserUuid())
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			logger.Error(ctx, "error while getting user",
				zap.String("user_uuid", req.GetUserUuid()),
				zap.Error(err),
			)
			return nil, status.Errorf(codes.NotFound, "user with UUID %s not found", req.GetUserUuid())
		}
		logger.Error(ctx, "error while getting user",
			zap.String("user_uuid", req.GetUserUuid()),
			zap.Error(err),
		)
		return nil, status.Errorf(codes.Internal, "internal error while getting user")
	}

	return &user_v1.GetUserResponse{
		User: converter.UserToProto(res),
	}, nil
}
