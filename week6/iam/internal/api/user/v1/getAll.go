package v1

import (
	"context"
	"errors"

	"github.com/mbakhodurov/homeworks/week6/iam/internal/converter"
	"github.com/mbakhodurov/homeworks/week6/iam/internal/model"
	"github.com/mbakhodurov/homeworks/week6/platform/pkg/logger"
	common_v1 "github.com/mbakhodurov/homeworks/week6/shared/pkg/proto/common/v1"
	user_v1 "github.com/mbakhodurov/homeworks/week6/shared/pkg/proto/user/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) GetAllUser(ctx context.Context, req *user_v1.GetAllUserRequest) (*user_v1.GetAllUserResponse, error) {
	res, err := a.userService.GetAll(ctx)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			logger.Error(ctx, "error while getting user",
				zap.Error(err),
			)
			return nil, status.Errorf(codes.NotFound, "users not found")
		}
		logger.Error(ctx, "error while getting user",
			zap.Error(err),
		)
		return nil, status.Errorf(codes.Internal, "internal error while getting user")
	}

	userDTOList := make([]*common_v1.User, 0, len(res))
	for _, v := range res {
		dto := converter.UserToProto(v)
		userDTOList = append(userDTOList, dto)
	}

	return &user_v1.GetAllUserResponse{
		Users: userDTOList,
	}, nil
}
