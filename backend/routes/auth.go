package routes

import "github.com/gin-gonic/gin"

func AuthRoutes(r *gin.RouterGroup) {
	r.POST("/login")
	r.POST("/signup")

}
