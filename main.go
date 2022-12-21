package main

import (
	// "fmt"
	// "context"
	"github.com/gin-gonic/gin"

	// configs "jitD/configs"
	routes "jitD/routers"
	"net/http"
	"time"
)

func main() {

	// initail route
	router := gin.Default()

	// use middleware
	// router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// provide route
	routes.UserRoute(router)

	router.GET("/a", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "message: status Ok")
	})

	// configue on port 3000
	s := &http.Server{
		Addr:           ":3000",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
