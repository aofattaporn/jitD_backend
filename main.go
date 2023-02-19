package main

import (
	// "fmt"
	// "context"
	// "net/http"

	configs "jitD/configs"
	routes "jitD/routers"

	"github.com/gin-gonic/gin"
	//"github.com/gin-gonic/gin"
)

func main() {

	// initail route
	router := gin.Default()

	// use middleware
	router.Use(configs.Verify)

	routes.UserRoute(router)
	routes.PostRoutes(router)
	routes.CommentRoutes(router)
	routes.LikeRoutes(router)

	// configue on port 3000
	router.Run("0.0.0.0:3000")
}
