package http

import (
	"context"
	"net/http"

	grpcAuth "github.com/mbakhodurov/homeworks/week6/platform/pkg/middleware/grpc"

	auth_v1 "github.com/mbakhodurov/homeworks/week6/shared/pkg/proto/auth/v1"
	common_v1 "github.com/mbakhodurov/homeworks/week6/shared/pkg/proto/common/v1"
)

const SessionUUIDHeader = "X-Session-Uuid"

// type IAMClient = auth_v1.AuthServiceClient

type AuthMiddleware struct {
	iamClient auth_v1.AuthServiceClient
}

// NewAuthMiddleware создает новый middleware аутентификации
func NewAuthMiddleware(iamClient auth_v1.AuthServiceClient) *AuthMiddleware {
	return &AuthMiddleware{
		iamClient: iamClient,
	}
}

func (m *AuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionUUID := r.Header.Get(SessionUUIDHeader)
		if sessionUUID == "" {
			writeErrorResponse(w, http.StatusUnauthorized, "MISSING_SESSION", "Authentication required")
			return
		}

		// Валидируем сессию через IAM сервис
		whoamiRes, err := m.iamClient.Whoami(r.Context(), &auth_v1.WhoamiRequest{
			SessionUuid: sessionUUID,
		})
		if err != nil {
			writeErrorResponse(w, http.StatusUnauthorized, "INVALID_SESSION", "Authentication failed")
			return
		}

		// Добавляем пользователя и session UUID в контекст используя функции из grpc middleware
		ctx := r.Context()
		ctx = grpcAuth.AddSessionUUIDToContext(ctx, sessionUUID)
		// Также добавляем пользователя в контекст
		ctx = context.WithValue(ctx, grpcAuth.GetUserContextKey(), whoamiRes.User)
		// ctx = grpcAuth.AddUserInfoToContext(ctx, whoamiRes.User)

		// Передаем управление следующему handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserFromContext извлекает пользователя из контекста
func GetUserFromContext(ctx context.Context) (*common_v1.User, bool) {
	return grpcAuth.GetUserFromContext(ctx)
}

// GetSessionUUIDFromContext извлекает session UUID из контекста
func GetSessionUUIDFromContext(ctx context.Context) (string, bool) {
	return grpcAuth.GetSessionUUIDFromContext(ctx)
}
