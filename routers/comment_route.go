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

	// test new struct
	v1.POST("/new/:post_id", controllers.NewCreateComment)
	v1.GET("/new/:post_id", controllers.NewGetComment)
	v1.GET("/new/all/:post_id", controllers.NewGetAllComment)
	v1.PUT("/new/:post_id/comment/:comment_id", controllers.NewUpdateComment)

}
