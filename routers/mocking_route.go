package routers

import (
	"jitD/controllers"

	"github.com/gin-gonic/gin"
)

func MockingRoute(route *gin.Engine) {
	mock := route.Group("v1/mock")

	mock.POST("/testStress", controllers.CreateSetTestStress)
	mock.POST("/testConsult", controllers.CreateSetTestConsult)

}
