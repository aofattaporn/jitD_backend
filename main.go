package main

import (
	// "fmt"
	// "context"
	"github.com/gin-gonic/gin"

	// configs "jitD/configs"
	routes "jitD/routers"
	// "net/http"
	// "time"
)

func main() {

	// initail route
	router := gin.Default()

	// use middleware
	// router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// provide route
	routes.UserRoute(router)

	// configue on port 3000
	router.Run("localhost:3000")
}
