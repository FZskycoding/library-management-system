package main

import(
	"library-sys/routes"
	"library-sys/models"
	"github.com/gin-gonic/gin"

)

func initializeBooks() {
    // 檢查是否已有書籍
    if len(models.Libraries) == 0 {
        // 加入預設書籍
        defaultBooks := []models.Library{
            {
                ID:     1,
                Title:  "哈利波特：神秘的魔法石",
                Author: "J.K. 羅琳",
                ISBN:   "9573317249",
                Status: "available",
            },
            {
                ID:     2,
                Title:  "魔戒首部曲：魔戒現身",
                Author: "J.R.R. 托爾金",
                ISBN:   "9573271575",
                Status: "available",
            },
        }
        
        models.Libraries = append(models.Libraries, defaultBooks...)
    }
}

func main(){
	r := gin.Default()
	initializeBooks()
	routes.SetupBookRouters(r)
	r.Run(":8080")
}
