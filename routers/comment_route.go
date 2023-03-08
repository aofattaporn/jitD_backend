package routers

import (
	controllers "jitD/controllers"

	"github.com/gin-gonic/gin"
)

func CommentRoutes(route *gin.Engine) {
	v1 := route.Group("v1/comments")

	// test new struct
	v1.POST("post/:post_id", controllers.CreateComment)
	v1.GET("/post/:post_id", controllers.GetAllCommentByPostID)
	v1.PUT("/:comment_id/post/:post_id/", controllers.UpdateComment)
	v1.DELETE("/:comment_id/post/:post_id/", controllers.DeleteComment)

}
