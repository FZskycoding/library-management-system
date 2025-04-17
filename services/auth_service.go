package services

import (
	"errors"
	"library-sys/config"
	"library-sys/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 這是會員服務處的專業人員，負責：
type AuthService struct {
	db     *gorm.DB       // 可以查詢會員資料庫
	config *config.Config // 知道如何驗證會員證
}


type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
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
	// 檢查用戶名是否已存在
	var existingUser models.User
	if err := s.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return errors.New("username already exists")
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
	var user models.User
	if err := s.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	//驗證密碼
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid password")
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
		UserID:   user.ID,
		Username: user.Username,
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
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWT.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
