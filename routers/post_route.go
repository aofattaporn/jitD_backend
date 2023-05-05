package routers

import (
	controllers "jitD/controllers"

	"github.com/gin-gonic/gin"
)

func PostRoutes(route *gin.Engine) {
	v1 := route.Group("v1/posts")

	v1.GET("/", controllers.GetAllNewPost)
	v1.GET("/homepage", controllers.GetAllPostHomePage)
	v1.GET("/id", controllers.GetMyPost)
	v1.GET("/byLike", controllers.GetPostByLikeIndividual)
	v1.GET("/keyword/:keyword", controllers.GetPostByKeyword)
	v1.GET("/category/:category", controllers.GetPostByCategorry)

	v1.POST("/", controllers.CreatePost)
	v1.DELETE("/:post_id", controllers.DeleteMyPost)
	v1.PUT("/:post_id", controllers.UpdatePost)

}
