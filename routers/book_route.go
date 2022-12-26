package routers

import (
	"github.com/gin-gonic/gin"
	controllers "jitD/controllers"
)

func BookRoutes(route *gin.Engine) {
	v1 := route.Group("v1/books")
	v1.GET("/", controllers.GetAllBook)
	v1.GET("/:id", controllers.GetBookById)
}
