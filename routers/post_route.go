package routers

import (
	"github.com/gin-gonic/gin"
)

func PostRoutes(route *gin.Engine) {
	v1 := route.Group("/")
	v1.GET("getAll/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "get all",
		})
	})

	v1.GET("getById/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "get by id",
		})
	})
}
