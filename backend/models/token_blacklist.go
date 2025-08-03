package models

import (
	"time"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

//登入憑證
type Claims struct {
    UserID          uint      `json:"user_id"`
    Username        string    `json:"username"`
    IsAdmin         bool      `json:"is_admin"`
    ServerStartTime time.Time `json:"server_start_time"`
    jwt.RegisteredClaims
}

type TokenBlacklist struct {
	gorm.Model
	Token     string `gorm:"uniqueIndex"`
	ExpiresAt time.Time
}
