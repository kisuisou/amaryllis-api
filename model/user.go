package model

import "time"

type User struct {
	ID           string `gorm:"primaryKey"`
	PasswordHash string
	CreatedAt    time.Time
	IsAdmin      bool
	Books        []Book `gorm:"many2many:user_books;"`
}
