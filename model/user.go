package model

import "time"

type User struct {
	UserID       string `gorm:"primaryKey"`
	PasswordHash string
	CreatedAt    time.Time
	IsAdmin      bool
}
