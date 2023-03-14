package routers

import (
	"jitD/controllers"

	"github.com/gin-gonic/gin"
)

func TestRoute(route *gin.Engine) {
	v1Stress := route.Group("v1/test/stress")

	v1Stress.POST("/", controllers.CreateSetTestStress)
	v1Stress.GET("/", controllers.GetTestStress)

}
