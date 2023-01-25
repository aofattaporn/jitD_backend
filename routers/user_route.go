package routers

import (
	controllers "jitD/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoute(route *gin.Engine) {
	v1 := route.Group("v1/users")
	v1.GET("/", controllers.GetAllUser)
	v1.POST("/", controllers.CreateUser)
	v1.GET("/id", controllers.GetUserById)
	v1.PUT("/pet/id", controllers.NamingPet)

	// # -------------- unsuse ---------------
	v1.DELETE("/id", controllers.DeleteUser)
}
