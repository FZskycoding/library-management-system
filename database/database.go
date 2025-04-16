package database

import (
	"library-sys/config"
	"library-sys/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB 初始化資料庫連接
func InitDB(config *config.Config) error {
	var err error

	// 使用配置建立資料庫連接
	DB, err = gorm.Open(postgres.Open(config.Database.GetDSN()), &gorm.Config{})
	if err != nil {
		return err
	}

	// 自動遷移資料庫結構
	err = DB.AutoMigrate(&models.Book{})
	if err != nil {
		return err
	}

	return nil
}

// GetDB 返回資料庫連接
func GetDB() *gorm.DB {
	return DB
}
