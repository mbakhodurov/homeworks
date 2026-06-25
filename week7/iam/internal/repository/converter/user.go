package converter

import (
	"github.com/mbakhodurov/homeworks/week7/iam/internal/model"
	repoModel "github.com/mbakhodurov/homeworks/week7/iam/internal/repository/model"
)

func UserInfoToRepo(from model.UserInfo) repoModel.UserInfo {
	methods := make([]repoModel.NotificationMethod, 0, len(from.NotificationMethods))

	for _, m := range from.NotificationMethods {
		methods = append(methods, repoModel.NotificationMethod{
			ProviderName: m.ProviderName,
			Target:       m.Target,
		})
	}

	return repoModel.UserInfo{
		Login:               from.Login,
		Email:               from.Email,
		NotificationMethods: methods,
	}
}

func UserInfoRepoToModel(from repoModel.UserInfo) model.UserInfo {
	methods := make([]model.NotificationMethod, 0, len(from.NotificationMethods))

	for _, m := range from.NotificationMethods {
		methods = append(methods, model.NotificationMethod{
			ProviderName: m.ProviderName,
			Target:       m.Target,
		})
	}

	return model.UserInfo{
		Login:               from.Login,
		Email:               from.Email,
		NotificationMethods: methods,
	}
}

func UserRepoModelToModel(from repoModel.User) model.User {

	return model.User{
		UserUUID:  from.UserUUID,
		UserInfo:  UserInfoRepoToModel(from.UserInfo),
		CreatedAt: from.CreatedAt,
		UpdatedAt: from.UpdatedAt,
		DeletedAt: from.DeletedAt,
		// Password:  from.Password,
	}
}
