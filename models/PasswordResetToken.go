package models

import (
	"time"

	"gorm.io/gorm"
)

type PasswordResetToken struct {
	gorm.Model
	UserID    uint
	Token     string `gorm:"uniqueIndex"`
	ExpiresAt time.Time
	UsedAt    *time.Time
}
