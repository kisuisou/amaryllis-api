package model

import (
	"time"
)

type UserBooks struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	UserID    string
	BookISBN  string
	IsRead    bool
	CreatedAt time.Time
}
