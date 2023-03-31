package routers

import (
	controllers "jitD/controllers"

	"github.com/gin-gonic/gin"
)

func BookMarkRoutes(route *gin.Engine) {
	v1 := route.Group("v1/posts/bookmark/")

	// get post in book mark
	v1.PUT("/add/:post_id", controllers.AddBookmark)

	// remove on book mark
	v1.PUT("/remove/:post_id", controllers.Remove)

	// add to book mark
	v1.GET("/", controllers.GetBookmarks)

}
