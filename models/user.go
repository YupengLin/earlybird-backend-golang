package models

import "time"

type User struct {
	UserID          int64     `json:"user_id"`
	Username        *string   `json:"username"`
	LastLogin       time.Time `json:"last_login"`
	CreatedAt       time.Time `json:"created_at"`
	Register        *bool     `json:"register"`
	Guestname       *string   `json:"guestname"`
	Email           string    `json:"email"`
	Password        *string   `json:"password,omitempty"`
	PasswordVersion int64     `json:"password_version"`
	Thumbnail       *string   `json:"thumbnail"`
}
