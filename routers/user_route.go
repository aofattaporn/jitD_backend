package routers

import(
	"github.com/gin-gonic/gin"

	controllers "jitD/controllers"
)

func UserRoute(route *gin.Engine) {
	v1 := route.Group("/")
	v1.GET("/", controllers.GetAllUser)
	v1.GET("getUser/", controllers.GetUserID)

}

