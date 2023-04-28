package routers

import (
	controllers "jitD/controllers"

	"github.com/gin-gonic/gin"
)

func ReportRoutes(route *gin.Engine) {
	v1 := route.Group("v1/report")

	v1.GET("/", controllers.GetAllReport)
	v1.POST("/", controllers.AddReport)
	v1.DELETE("/:report_id", controllers.UpdateReport)
	v1.DELETE("/:report_id", controllers.DeleteReport)
}
