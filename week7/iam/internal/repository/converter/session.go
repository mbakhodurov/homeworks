package converter

import (
	"encoding/json"
	"time"

	"github.com/mbakhodurov/homeworks/week7/iam/internal/model"
	repoModel "github.com/mbakhodurov/homeworks/week7/iam/internal/repository/model"
	"github.com/samber/lo"
)

func SessionAndUserFromRedisView(view repoModel.SessionRedisView) (model.Session, model.User) {
	var notificationMethods []model.NotificationMethod
	if view.NotificationMethod != "" {
		err := json.Unmarshal([]byte(view.NotificationMethod), &notificationMethods)
		if err != nil {
			notificationMethods = []model.NotificationMethod{}
		}
	}

	var userUpdatedAt *time.Time
	if view.UserUpdatedAt != nil {
		tmp := time.Unix(0, *view.UserUpdatedAt)
		userUpdatedAt = &tmp
	}

	user := model.User{
		UserUUID: view.UserUUID,
		UserInfo: model.UserInfo{
			Login:               view.Login,
			Email:               view.Email,
			NotificationMethods: notificationMethods,
		},
		CreatedAt: time.Unix(0, view.UserCreatedAt),
		UpdatedAt: userUpdatedAt,
	}

	var updatedAt *time.Time
	if view.UpdatedAtNs != nil {
		tmp := time.Unix(0, *view.UpdatedAtNs)
		updatedAt = &tmp
	}

	session := model.Session{
		UUID:      view.UUID,
		CreatedAt: time.Unix(0, view.CreatedAtNs),
		UpdatedAt: updatedAt,
		ExpiresAt: time.Unix(0, view.ExpiresAt),
	}

	return session, user
}

func SessionAndUserToRedisView(session model.Session, user model.User) repoModel.SessionRedisView {
	var updatedAt *int64
	if session.UpdatedAt != nil {
		updatedAt = lo.ToPtr(session.UpdatedAt.UnixNano())
	}

	var userUpdatedAt *int64
	if user.UpdatedAt != nil {
		userUpdatedAt = lo.ToPtr(user.UpdatedAt.UnixNano())
	}

	return repoModel.SessionRedisView{
		UserUUID:           user.UserUUID,
		Login:              user.UserInfo.Login,
		Email:              user.UserInfo.Email,
		NotificationMethod: serializeNotificationMethods(user.UserInfo.NotificationMethods),
		UserCreatedAt:      user.CreatedAt.UnixNano(),
		UserUpdatedAt:      userUpdatedAt,
		UUID:               session.UUID,
		CreatedAtNs:        session.CreatedAt.UnixNano(),
		UpdatedAtNs:        updatedAt,
		ExpiresAt:          session.ExpiresAt.UnixNano(),
	}
}

func serializeNotificationMethods(methods []model.NotificationMethod) string {
	serialized, err := json.Marshal(methods)
	if err != nil {
		return ""
	}
	return string(serialized)
}
