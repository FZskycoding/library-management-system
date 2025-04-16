package config

import (
	"fmt"
)

type Config struct {
	Database *DatabaseConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
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
