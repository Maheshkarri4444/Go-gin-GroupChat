package route

import (
	"github.com/gin-gonic/gin"
	controller "githun.com/Maheshkarri4444/group-chat/controllers"
	"githun.com/Maheshkarri4444/group-chat/middleware"
)

func MessageRoutes(r *gin.RouterGroup) {
	r.GET("/messages", middleware.AuthMiddleware(), controller.GetMessages)
	r.POST("/send", middleware.AuthMiddleware(), controller.SendMessage)
	// r.DELETE("/delete")
}
