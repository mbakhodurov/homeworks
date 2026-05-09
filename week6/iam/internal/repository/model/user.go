package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
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

// Scan реализует интерфейс sql.Scanner для UserInfo
func (ui *UserInfo) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan UserInfo: value is not []byte")
	}

	return json.Unmarshal(bytes, ui)
}

// Value реализует интерфейс driver.Valuer для UserInfo
func (ui UserInfo) Value() (driver.Value, error) {
	b, err := json.Marshal(ui)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}
