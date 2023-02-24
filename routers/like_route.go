package routers

import (
	controllers "jitD/controllers"

	"github.com/gin-gonic/gin"
)

func LikeRoutes(route *gin.Engine) {
	v1Like := route.Group("v1/like")
	v1Like.PUT("/post/:post_id", controllers.LikePost)
	v1Like.PUT("/comment/:comment_id/post/:post_id", controllers.LikeComment)

	v1UnLike := route.Group("v1/unlike")
	// v1UnLike.PUT("/post/:post_id", controllers.UnlikePost)
	v1UnLike.PUT("/comment/:comment_id/post/:post_id", controllers.UnLikeComment)

}
