package routers

import (
	"jitD/controllers"

	"github.com/gin-gonic/gin"
)

func TestRoute(route *gin.Engine) {
	v1Stress := route.Group("v1/test/stress")
	v1Stress.GET("/", func(c *gin.Context) {
		controllers.GetTest(c, "TestStress")
	})
	v1Stress.PUT("/result/point/:point", controllers.CalTestStress)
	v1Stress.GET("/result", func(c *gin.Context) {
		controllers.GetTestResult(c, "TestStress")
	})

	v1Consult := route.Group("v1/test/consult")
	v1Consult.GET("/", func(c *gin.Context) {
		controllers.GetTest(c, "TestConsult")
	})
	v1Consult.PUT("/result/point/:point", controllers.CalTestStress)
	v1Consult.GET("/result", func(c *gin.Context) {
		controllers.GetTestResult(c, "TestConsult")
	})

}
