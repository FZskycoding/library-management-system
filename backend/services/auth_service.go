package services

import (
"errors"
"library-sys/config"
"library-sys/models"
"strings"
"sync"
"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
    serverStartTime time.Time
    once           sync.Once
)

//在執行main前會自動執行
func init() {
    once.Do(func() {
        serverStartTime = time.Now()
    })
}

type AuthService struct {
    db     *gorm.DB
    config *config.Config
}

//登入憑證
type Claims struct {
    UserID          uint      `json:"user_id"`
    Username        string    `json:"username"`
    IsAdmin         bool      `json:"is_admin"`
    ServerStartTime time.Time `json:"server_start_time"`
    jwt.RegisteredClaims
}

func CreateAuthService(db *gorm.DB, config *config.Config) *AuthService {
	return &AuthService{
		db:     db,
		config: config,
	}
}

// Register 註冊新用戶
func (s *AuthService) Register(req *models.RegisterRequest) error {
	// 移除首尾空格
	req.Username = strings.TrimSpace(req.Username)

	// 檢查是否包含空格
	if strings.Contains(req.Username, " ") {
		return errors.New("使用者名稱不能包含空格")
	}
	// 檢查用戶名是否已存在
	var existingUser models.User
	if err := s.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return errors.New("使用者名稱已被使用")
	}

	//加密密碼
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 創建新用戶
	user := &models.User{
		Username: req.Username,
		Password: string(hashedPassword),
	}

	return s.db.Create(user).Error

}

// Login 用戶登入
func (s *AuthService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	// 移除首尾空格
	req.Username = strings.TrimSpace(req.Username)
	// 檢查是否包含空格
	if strings.Contains(req.Username, " ") {
		return nil, errors.New("使用者名稱不能包含空格")
	}

	var user models.User
	if err := s.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return nil, errors.New("找不到此使用者，請註冊")
	}

	//驗證密碼
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("密碼錯誤")
	}

	// 生成 Token
	token, err := s.GenerateToken(&user)
	if err != nil {
		return nil, err
	}
	return &models.LoginResponse{
		Token: token,
	}, nil
}

// GenerateToken 生成 JWT token
func (s *AuthService) GenerateToken(user *models.User) (string, error) {
    // 設置 JWT 聲明
    claims := &Claims{
        UserID:          user.ID,
        Username:        user.Username,
        IsAdmin:         user.IsAdmin,
        ServerStartTime: serverStartTime,
        RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(s.config.JWT.ExpireHours))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// 創建 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 簽名 token
	return token.SignedString([]byte(s.config.JWT.SecretKey))
}

// ValidateToken 驗證 JWT token
func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	// 添加黑名單檢查
	var blacklistedToken models.TokenBlacklist
	if err := s.db.Where("token = ?", tokenString).First(&blacklistedToken).Error; err == nil {
		return nil, errors.New("此登入憑證已失效")
	}
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWT.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

if !token.Valid {
    return nil, errors.New("無效的登入憑證")
}

// 檢查伺服器啟動時間
if !claims.ServerStartTime.Equal(serverStartTime) {
    return nil, errors.New("伺服器已重啟，請重新登入")
}

return claims, nil
}

//用戶登出
func (s *AuthService) Logout(tokenString string) error {
	// 先驗證 token
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return err
	}

	// 將Token加入黑名單
	blacklist := &models.TokenBlacklist{
		Token:     tokenString,
		ExpiresAt: time.Unix(claims.ExpiresAt.Unix(), 0),
	}

	return s.db.Create(blacklist).Error
}
