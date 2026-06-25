package converter

import (
	"github.com/mbakhodurov/homeworks/week7/iam/internal/model"
	common_v1 "github.com/mbakhodurov/homeworks/week7/shared/pkg/proto/common/v1"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func SessionToProto(session model.Session) *common_v1.Session {
	return &common_v1.Session{
		Uuid:      session.UUID,
		CreatedAt: timestamppb.New(session.CreatedAt),
		UpdatedAt: timestamppb.New(lo.FromPtr(session.UpdatedAt)),
		ExpiresAt: timestamppb.New(session.ExpiresAt),
	}
}
