package routers

import (
	controllers "jitD/controllers"

	"github.com/gin-gonic/gin"
)

func CommentRoutes(route *gin.Engine) {
	v1 := route.Group("v1/comment")
	v1.POST("/:post_id", controllers.CreateComment)
	v1.GET("/:id", controllers.GetAllComment)
	v1.GET("/", controllers.GetMyComment)
	v1.GET("/post/:post_id", controllers.GetCommentByPostID)
	// v1.POST("/:id", controllers.UpdateComment)
	// v1.DELETE("/:id", controllers.DeleteComment)
}
