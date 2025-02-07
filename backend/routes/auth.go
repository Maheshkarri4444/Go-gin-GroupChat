package routes

import (
	"github.com/gin-gonic/gin"
	controller "githun.com/Maheshkarri4444/group-chat/controllers"
)

func AuthRoutes(r *gin.RouterGroup) {
	r.POST("/login", controller.Login)
	r.POST("/signup", controller.SignUp)
}
