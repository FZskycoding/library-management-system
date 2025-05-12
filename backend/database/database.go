package database

import (
	"library-sys/config"
	"library-sys/models"

	"golang.org/x/crypto/bcrypt"
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
	err = DB.AutoMigrate(&models.Book{}, &models.User{}, &models.TokenBlacklist{})
	if err != nil {
		return err
	}

	// 檢查並創建管理員帳號
	var adminUser models.User
	result := DB.Where("username = ?", config.Admin.Username).First(&adminUser)
	if result.RowsAffected == 0 {
		// 加密密碼
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(config.Admin.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		// 創建管理員用戶
		adminUser = models.User{
			Username: config.Admin.Username,
			Password: string(hashedPassword),
			IsAdmin:  true,
		}
		if err := DB.Create(&adminUser).Error; err != nil {
			return err
		}
	}

	return nil
}

// GetDB 返回資料庫連接
func GetDB() *gorm.DB {
	return DB
}
