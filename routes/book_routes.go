package routes

import (
	"library-sys/controllers"

	"github.com/gin-gonic/gin"
)

func SetupBookRouters(router *gin.Engine) {
	LibraryController := controllers.LibraryController{}

	v1 := router.Group("/api/v1")
	{
		v1.GET("/books", LibraryController.GetAll)            //查詢所有書籍
		v1.GET("/books/:id", LibraryController.GetByID)       //查詢特定書籍
		v1.POST("/books", LibraryController.Create)           //新增書籍
		v1.PUT("/books/:id", LibraryController.Update)        //更新書籍信息
		v1.DELETE("/books/:id", LibraryController.Delete)     //刪除書籍
		v1.PUT("/books/:id/borrow", LibraryController.Borrow) //借書
		v1.PUT("/books/:id/return", LibraryController.Return) //還書
	}
}
