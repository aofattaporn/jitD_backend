package main

import (
	configs "jitD/configs"
	"jitD/controllers"
	routes "jitD/routers"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

type Data struct {
	value string
}

func main() {

	// initail route
	router := gin.Default()

	controllers.ReccomendPost()
	// use middleware
	router.Use(CORSMiddleware())
	router.Use(configs.Verify)

	routes.UserRoute(router)
	routes.QuestRoute(router)
	routes.PostRoutes(router)
	routes.CommentRoutes(router)
	routes.LikeRoutes(router)
	routes.BookMarkRoutes(router)
	routes.TestRoute(router)
	// routes.ReportRoutes(router)

	routes.MockingRoute(router)

	// configue on port 3000
	router.Run("0.0.0.0:3000")
}
