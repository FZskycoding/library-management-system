package routes

import (
	"library-sys/controllers"
	"library-sys/middleware"
	"library-sys/services"

	"github.com/gin-gonic/gin"
)

func SetupBookRouters(router *gin.Engine, authService *services.AuthService, bookService *services.BookService) {
	// 創建控制器
	authController := controllers.NewAuthController(authService)
	bookController := controllers.NewLibraryController(bookService)

	// 不需要認證的路由
	public := router.Group("/")
	{
		public.POST("/register", authController.Register) //註冊
		public.POST("/login", authController.Login)       //登入
		public.GET("/books", bookController.GetAll)       //查詢所有書籍
		public.GET("/books/:id", bookController.GetByID)  //查詢特定書籍
	}

	// 需要認證的路由
	protected := router.Group("/", middleware.JWTAuthMiddleware(authService))
	{
		// 一般使用者可以使用的功能
		protected.PUT("/books/:id/borrow", bookController.Borrow) //借書
		protected.PUT("/books/:id/return", bookController.Return) //還書
		protected.GET("/me", authController.GetCurrentUser)       // 查詢userid
		protected.POST("/logout", authController.Logout)          // 登出
	}

	// 需要管理員權限的路由
	admin := router.Group("/", middleware.JWTAuthMiddleware(authService), middleware.AdminAuthMiddleware())
	{
		admin.POST("/books", bookController.Create)       //新增書籍
		admin.PUT("/books/:id", bookController.Update)    //更新書籍信息
		admin.DELETE("/books/:id", bookController.Delete) //刪除書籍
	}
}
