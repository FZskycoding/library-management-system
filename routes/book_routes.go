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
	public := router.Group("/api/v1")
	{
		public.POST("/register", authController.Register)
		public.POST("/login", authController.Login)
		public.GET("/books", bookController.GetAll)      //查詢所有書籍
		public.GET("/books/:id", bookController.GetByID) //查詢特定書籍
	}

	// 需要認證的路由
	protected := router.Group("/api/v1").Use(middleware.JWTAuthMiddleware(authService))
	{
		protected.POST("/books", bookController.Create)           //新增書籍
		protected.PUT("/books/:id", bookController.Update)        //更新書籍信息
		protected.DELETE("/books/:id", bookController.Delete)     //刪除書籍
		protected.PUT("/books/:id/borrow", bookController.Borrow) //借書
		protected.PUT("/books/:id/return", bookController.Return) //還書
		// 用戶相關路由
		protected.GET("/me", authController.GetCurrentUser)
		protected.POST("/logout", authController.Logout)

	}
}
