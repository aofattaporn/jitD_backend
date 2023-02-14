package routers

import (
	controllers "jitD/controllers"

	"github.com/gin-gonic/gin"
)

func CommentRoutes(route *gin.Engine) {
	v1 := route.Group("v1/comments")
	v1.POST("/:post_id", controllers.CreateComment)
	v1.GET("/post/:post_id", controllers.GetCommentByPostID)
	v1.PUT("/:comment_id", controllers.UpdateComment)
	v1.DELETE("/:comment_id/post/:post_id", controllers.DeleteComment)

	// service that's not use iin fronent
	v1.GET("/", controllers.GetAllComment)
	v1.GET("/id", controllers.GetMyComment)
}
