package main

import (
	"library-sys/config"
	"library-sys/database"
	"library-sys/routes"
	"library-sys/services"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
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
	// 安排圖書館服務台
	routes.SetupBookRouters(
		r,
		authService, //處理會員相關的服務
		bookService, //處理書籍相關的服務
	)
	// 啟動服務器
	r.Run(":8080")
}
