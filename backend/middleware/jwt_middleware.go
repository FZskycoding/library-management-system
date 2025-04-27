package middleware

import (
    "library-sys/services"
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
)

func JWTAuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 從 Header 獲取 Token
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
            c.Abort()
            return
        }

        // 檢查 Token 格式
        parts := strings.SplitN(authHeader, " ", 2)
        if !(len(parts) == 2 && parts[0] == "Bearer") {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
            c.Abort()
            return
        }

        // 驗證 Token
        claims, err := authService.ValidateToken(parts[1])
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        // 將用戶信息存儲到上下文
        c.Set("userID", claims.UserID)
        c.Set("username", claims.Username)

        c.Next()
    }
}
