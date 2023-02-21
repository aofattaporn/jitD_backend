package routers

import (
	controllers "jitD/controllers"

	"github.com/gin-gonic/gin"
)

func CommentRoutes(route *gin.Engine) {
	v1 := route.Group("v1/comments")

	// test new struct
	v1.POST("post/:post_id", controllers.NewCreateComment)
	v1.GET("/:comment_id/post/:post_id/", controllers.NewGetCommentByID)
	v1.GET("/post/:post_id", controllers.NewGetAllCommentByPostID)
	v1.PUT("/:comment_id/post/:post_id/", controllers.NewUpdateComment)
	v1.DELETE("/:comment_id/post/:post_id/", controllers.NewDeleteComment)

}
