package routers

import (
	controllers "jitD/controllers"

	"github.com/gin-gonic/gin"
)

func PostRoutes(route *gin.Engine) {
	v1 := route.Group("v1/posts")
	v1.POST("/", controllers.CreatePost)
	v1.GET("/", controllers.GetAllPost)
	v1.GET("/keyword/:keyword", controllers.GetPostByKeyword)
	v1.GET("/id", controllers.GetMyPost)
	v1.DELETE("/:post_id", controllers.DeleteMyPost)
	v1.PUT("/:post_id", controllers.UpdatePost)
}
