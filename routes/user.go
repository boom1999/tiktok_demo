package routes

import (
	"github.com/gin-gonic/gin"
	"tiktok_demo/middleware/jwt"
)

func User(r *gin.RouterGroup) {
	user := r.Group("/user")
	{
		user.POST("/register/")
		user.POST("/login/")
		user.GET("/", jwt.Auth())
	}
}
