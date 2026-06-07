package model

import (
	"time"
)

type UserBooks struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	UserID    string
	BookID    uint
	IsRead    bool
	CreatedAt time.Time
}
