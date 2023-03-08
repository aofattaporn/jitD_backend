package routers

import (
	"jitD/controllers"

	"github.com/gin-gonic/gin"
)

func BookMarkRoutes(route *gin.Engine) {
	v1 := route.Group("v1/posts/bookmark/")

	// get post in book mark
	v1.PUT("/add", controllers.AddBookmark2)

	// remove on book mark
	v1.PUT("/remove", controllers.Remove2)

	// add to book mark
	v1.GET("/", controllers.Remove2)

}
