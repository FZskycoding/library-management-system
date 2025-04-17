package config

import (
	"fmt"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	SecretKey    string
	ExpireHours  int
	RefreshHours int
}
type Config struct {
	Database *DatabaseConfig
	JWT      *JWTConfig
}

// 創建默認配置
func NewConfig() *Config {
	return &Config{
		Database: &DatabaseConfig{
			Host:     "localhost",
			Port:     "5432",      // PostgreSQL 默認端口
			User:     "postgres",  // 你的資料庫用戶名
			Password: "fei080808", // 你的資料庫密碼
			DBName:   "library",   // 你的資料庫名稱
			SSLMode:  "disable",   // 開發環境可以禁用 SSL
		},
		JWT: &JWTConfig{
			SecretKey:    "your-secret-key-here", // 實際使用時應該使用環境變數
			ExpireHours:  24,                     // Token 有效期 24 小時
			RefreshHours: 72,                     // Refresh Token 有效期 72 小時
		},
	}
}

// 生成資料庫連接字串的方法
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		c.Host,
		c.User,
		c.Password,
		c.DBName,
		c.Port,
		c.SSLMode,
	)
}
