package models

import (
	"time"

	"gorm.io/gorm"
)

type TokenBlacklist struct {
	gorm.Model
	Token     string `gorm:"uniqueIndex"`
	ExpiresAt time.Time
}
