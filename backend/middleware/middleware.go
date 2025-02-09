package middleware

import (
	"github.com/gin-gonic/gin"
	"githun.com/Maheshkarri4444/group-chat/auth"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !auth.ValidateSession(c) {
			c.Abort()
			return
		}
		userid, _ := c.Cookie("userid")
		c.Set("userid", userid)
		c.Next()
	}

}
