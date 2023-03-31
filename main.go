package main

import (
	configs "jitD/configs"
	routes "jitD/routers"

	"github.com/gin-gonic/gin"
)

func main() {

	// initail route
	router := gin.Default()

	// use middleware
	router.Use(configs.Verify)

	routes.UserRoute(router)
	routes.QuestRoute(router)
	routes.PostRoutes(router)
	routes.CommentRoutes(router)
	routes.LikeRoutes(router)
	routes.BookMarkRoutes(router)
	routes.TestRoute(router)
	routes.MockingRoute(router)

	// configue on port 3000
	router.Run("0.0.0.0:3000")
}
