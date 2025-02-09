package route

import (
	"github.com/gin-gonic/gin"
	controller "githun.com/Maheshkarri4444/group-chat/controllers"
	"githun.com/Maheshkarri4444/group-chat/middleware"
)

func AuthRoutes(r *gin.RouterGroup) {
	r.POST("/login", controller.Login)
	r.POST("/signup", controller.SignUp)
	r.POST("/logout", controller.Logout)
	r.GET("/check-auth", middleware.AuthMiddleware(), controller.CheckAuth)
}
