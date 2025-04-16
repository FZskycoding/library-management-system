package routes

import (
	"library-sys/controllers"
	"github.com/gin-gonic/gin"
)

func SetupBookRouters(router *gin.Engine) {
	lc := controllers.DefaultController() 

	v1 := router.Group("/api/v1")
	{
		v1.GET("/books", lc.GetAll)            //查詢所有書籍
		v1.GET("/books/:id", lc.GetByID)       //查詢特定書籍
		v1.POST("/books", lc.Create)           //新增書籍
		v1.PUT("/books/:id", lc.Update)        //更新書籍信息
		v1.DELETE("/books/:id", lc.Delete)     //刪除書籍
		v1.PUT("/books/:id/borrow", lc.Borrow) //借書
		v1.PUT("/books/:id/return", lc.Return) //還書
	}
}
