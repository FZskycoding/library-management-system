package main

import (
"library-sys/config"
"library-sys/database"
"library-sys/routes"
"library-sys/services"
"log"
"time"

"github.com/gin-contrib/cors"
"github.com/gin-gonic/gin"
"github.com/joho/godotenv"
)

func main() {
// 載入 .env 檔案
if err := godotenv.Load(); err != nil {
    log.Fatal("Error loading .env file:", err)
}

	// 載入配置
	cfg := config.NewConfig()

	// 初始化資料庫
	err := database.InitDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db := database.GetDB()

	//創建服務
	authService := services.CreateAuthService(db, cfg)
	bookService := services.CreateBookService(db)

// 創建 Gin 引擎
r := gin.Default()

// 配置CORS
r.Use(cors.New(cors.Config{
	AllowOrigins:     []string{"http://localhost:3000"},
	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
	ExposeHeaders:    []string{"Content-Length"},
	AllowCredentials: true,
	MaxAge:           12 * time.Hour,
}))

// 安排圖書館服務台
	routes.SetupBookRouters(
		r,
		authService, //處理會員相關的服務
		bookService, //處理書籍相關的服務
	)
	// 啟動服務器
	r.Run(":8080")
}
