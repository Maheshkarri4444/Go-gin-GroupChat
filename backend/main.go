package main

import (
	"github.com/gin-gonic/gin"
	"githun.com/Maheshkarri4444/group-chat/routes"
)

func main() {
	router := gin.Default()

	authRoutes := router.Group("/api/auth")
	{
		routes.AuthRoutes(authRoutes)
	}

	router.Run(":4000")

}
