package routers

import(
	"github.com/gin-gonic/gin"
	controllers "jitD/controllers"
)

func UserRoute(route *gin.Engine) {
	v1 := route.Group("/")
	v1.GET("/", controllers.GetAllUser)
	v1.GET("/:id", controllers.GetUserById)
	v1.POST("create/", controllers.CreateUser)
	v1.DELETE("delete/", controllers.DeleteUser)
}

