package main

import(
	"library-sys/routes"
    "library-sys/config"
    "library-sys/database"
    "log"
	"github.com/gin-gonic/gin"

)

func main(){
    // 載入配置
    cfg := config.NewConfig()

    // 初始化資料庫
    err := database.InitDB(cfg)
    if err != nil{
        log.Fatal("Failed to connect to database:", err)
    }

    // 創建 Gin 引擎
	r := gin.Default()
    // 設置路由
	routes.SetupBookRouters(r)
    // 啟動服務器
	r.Run(":8080")
}
