package controllers

import (
	"library-sys/models"
	"library-sys/services"
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}


// Register 處理用戶註冊
func (ac *AuthController) Register(c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	if err := ac.authService.Register(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
	})
}

// Login 處理用戶登入
func (ac *AuthController) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	response, err := ac.authService.Login(&req) 

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Logout 處理用戶登出
func (ac *AuthController) Logout(c *gin.Context) {
    // 從 header 獲取 token
    token := c.GetHeader("Authorization")
    if token == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No token provided"})
        return
    }
    
    // 移除 "Bearer " 前綴
    token = strings.TrimPrefix(token, "Bearer ")
    
    // 調用服務層的登出方法
    if err := ac.authService.Logout(token); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Logout failed: " + err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}
// GetCurrentUser 獲取當前用戶信息
func (ac *AuthController) GetCurrentUser(c *gin.Context) {
	// 從上下文中獲取用戶信息
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Not authenticated",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
	})
}
