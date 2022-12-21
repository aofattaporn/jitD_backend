package routers

import (
	"github.com/gin-gonic/gin"
	controllers "jitD/controllers"
)

func PostRoutes(route *gin.Engine) {
	v1 := route.Group("/")
	v1.GET("getAll/", controllers.GetAllUser)
	v1.GET("getById/", controllers.GetAllUser)
}
