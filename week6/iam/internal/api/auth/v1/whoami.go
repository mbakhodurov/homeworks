package v1

import (
	"context"

	"github.com/mbakhodurov/homeworks/week6/iam/internal/converter"
	auth_v1 "github.com/mbakhodurov/homeworks/week6/shared/pkg/proto/auth/v1"
)

func (a *api) Whoami(ctx context.Context, req *auth_v1.WhoamiRequest) (*auth_v1.WhoamiResponse, error) {
	session, user, err := a.authService.Whoami(ctx, req.GetSessionUuid())
	if err != nil {
		return nil, err
	}

	return &auth_v1.WhoamiResponse{
		Session: converter.SessionToProto(session),
		User:    converter.UserToProto(user),
	}, nil
}
