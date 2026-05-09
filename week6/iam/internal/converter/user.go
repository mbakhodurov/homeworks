package converter

import (
	"github.com/mbakhodurov/homeworks/week6/iam/internal/model"
	common_v1 "github.com/mbakhodurov/homeworks/week6/shared/pkg/proto/common/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func UserInfoProtoToModel(user *common_v1.UserInfo) model.UserInfo {
	return model.UserInfo{
		Login:               user.GetLogin(),
		Email:               user.GetEmail(),
		NotificationMethods: notificationMethodToModel(user.GetNotificationMethods()),
	}
}
func notificationMethodToModel(methods []*common_v1.NotificationMethod) []model.NotificationMethod {
	result := make([]model.NotificationMethod, 0, len(methods))
	for _, method := range methods {
		result = append(result, model.NotificationMethod{
			ProviderName: method.GetProviderName(),
			Target:       method.GetTarget(),
		})
	}
	return result
}

func UserToProto(user model.User) *common_v1.User {
	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt != nil {
		updatedAt = timestamppb.New(*user.UpdatedAt)
	}

	return &common_v1.User{
		Uuid: user.UserUUID,
		Info: &common_v1.UserInfo{
			Login:               user.UserInfo.Login,
			Email:               user.UserInfo.Email,
			NotificationMethods: notificationMethodToProto(user.UserInfo.NotificationMethods),
		},
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

func notificationMethodToProto(methods []model.NotificationMethod) []*common_v1.NotificationMethod {
	result := make([]*common_v1.NotificationMethod, len(methods))
	for _, method := range methods {
		result = append(result, &common_v1.NotificationMethod{
			ProviderName: method.ProviderName,
			Target:       method.Target,
		})
	}

	return result
}
