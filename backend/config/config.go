package config

import (
	"fmt"
	"os"
	"strconv"
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

// getRequiredEnv 從環境變數獲取值，如果不存在則報錯
func getRequiredEnv(key string) string {
    value := os.Getenv(key)
    if value == "" {
        panic(fmt.Sprintf("必要的環境變數 %s 未設置", key))
    }
    return value
}

// getRequiredEnvInt 從環境變數獲取整數值，如果不存在或無法轉換則報錯
func getRequiredEnvInt(key string) int {
    value := getRequiredEnv(key)
    intValue, err := strconv.Atoi(value)
    if err != nil {
        panic(fmt.Sprintf("環境變數 %s 必須是有效的整數", key))
    }
    return intValue
}

// NewConfig 創建新的配置，從環境變數讀取必要的設定
func NewConfig() *Config {
	return &Config{
Database: &DatabaseConfig{
Host:     getRequiredEnv("DB_HOST"),
Port:     getRequiredEnv("DB_PORT"),
User:     getRequiredEnv("DB_USER"),
Password: getRequiredEnv("DB_PASSWORD"),
DBName:   getRequiredEnv("DB_NAME"),
SSLMode:  getRequiredEnv("DB_SSLMODE"),
},
JWT: &JWTConfig{
SecretKey:    getRequiredEnv("JWT_SECRET_KEY"),
ExpireHours:  getRequiredEnvInt("JWT_EXPIRE_HOURS"),
RefreshHours: getRequiredEnvInt("JWT_REFRESH_HOURS"),
		},
	}
}

// GetDSN 生成資料庫連接字串
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
