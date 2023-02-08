package routers

import (
	controllers "jitD/controllers"

	"github.com/gin-gonic/gin"
)

func CommentRoutes(route *gin.Engine) {
	v1 := route.Group("v1/comments")
	v1.POST("/:post_id", controllers.CreateComment)
	v1.GET("/:id", controllers.GetAllComment)
	// v1.GET("/id", controllers.GetMyComment)
	// v1.GET("/post/:id", controllers.GetCommentByPostID)
	// v1.GET("/user/:id", controllers.GetCommentByUserID)
	// v1.GET("/:id", controllers.UpdateComment)
	// v1.GET("/like/:id", controllers.DeleteComment)
	// v1.GET("/:id", controllers.LikeComment)
	// v1.GET("/:id", controllers.DislikeComment)
	// v1.POST("/:id", controllers.UpdateComment)
	// v1.DELETE("/:id", controllers.DeleteComment)
}
