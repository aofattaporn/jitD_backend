package routers

import (
	configs "jitD/configs"
	controllers "jitD/controllers"

	"github.com/gin-gonic/gin"
)

func PostRoutes(route *gin.Engine) {
	v1 := route.Group("v1/posts")
	v1.POST("/:id", controllers.CreatePost)
	v1.GET("/", configs.Verify, controllers.GetAllPost)
}
