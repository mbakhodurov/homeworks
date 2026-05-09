package v1

import (
	"context"
	"fmt"

	"github.com/mbakhodurov/homeworks/week6/iam/internal/converter"
	"github.com/mbakhodurov/homeworks/week6/platform/pkg/logger"
	user_v1 "github.com/mbakhodurov/homeworks/week6/shared/pkg/proto/user/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) Register(ctx context.Context, req *user_v1.RegisterRequest) (*user_v1.RegisterResponse, error) {
	fmt.Println("GetInfo:", req.Info)
	userInfo := converter.UserInfoProtoToModel(req.GetInfo())
	userUUID, err := a.userService.Register(ctx, userInfo, req.GetPassword())
	if err != nil {
		logger.Error(ctx, "error while registering user",
			zap.Error(err),
		)
		return nil, status.Errorf(codes.Internal, "internal error while registering user")
	}

	return &user_v1.RegisterResponse{
		UserUuid: userUUID,
	}, nil
}
