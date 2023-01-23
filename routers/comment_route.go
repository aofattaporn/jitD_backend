package routers

import (
	controllers "jitD/controllers"

	"github.com/gin-gonic/gin"
)

func CommentRoutes(route *gin.Engine) {
	v1 := route.Group("v1/comment")
	v1.GET("/", controllers.GetAllComment)
	v1.GET("/:id", controllers.GetCommentById)
	v1.POST("/", controllers.CreateComment)
	v1.POST("/:id", controllers.UpdateComment)
	v1.DELETE("/:id", controllers.DeleteComment)
}
