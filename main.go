package main

import (
	// "fmt"
	// "context"
	// "net/http"

	routes "jitD/routers"

	"github.com/gin-gonic/gin"
	//"github.com/gin-gonic/gin"
)

func main() {

	// initail route
	router := gin.Default()

	// use middleware
	// router.Use(gin.Logger())
	// router.Use(gin.Recovery())

	// provide route
	routes.BookRoutes(router)
	routes.UserRoute(router)
	routes.CommentRoutes(router)
	routes.PostRoutes(router)

	// checking token
	routes.PostRoutes(router)

	// configue on port 3000
	router.Run("0.0.0.0:3000")
}
