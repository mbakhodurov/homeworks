package model

import (
	"time"
)

type User struct {
	UserUUID  string     `json:"user_uuid"`
	UserInfo  UserInfo   `json:"user_info"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	Password  string     `json:"password_hash"`
}

type UserInfo struct {
	Login               string               `json:"login"`
	Email               string               `json:"email"`
	NotificationMethods []NotificationMethod `json:"notificationMethods"`
}

type NotificationMethod struct {
	ProviderName string `json:"provider_name"` // Провайдер: telegram, email, push и т.д.
	Target       string `json:"target"`        // Адрес/идентификатор назначения (email, чат-id)
}
