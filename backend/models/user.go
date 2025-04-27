package models

import (
	"gorm.io/gorm"
)

// 添加用戶模型
type User struct {
	gorm.Model
	Username string `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Password string `json:"password" gorm:"size:100;not null"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50"`
    Password string `json:"password" binding:"required,min=6,max=50"`
}

type LoginResponse struct {
    Token string `json:"token"`
}
