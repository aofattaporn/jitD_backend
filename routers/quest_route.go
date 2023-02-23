package routers

import (
	controllers "jitD/controllers"

	"github.com/gin-gonic/gin"
)

func QuestRoute(route *gin.Engine) {
	v1 := route.Group("v1/quest")

	v1.PUT("/updatePoint", controllers.UpdateProgressQuest)
	v1.PUT("/getPoint", controllers.GetPointFromQuest)

}
