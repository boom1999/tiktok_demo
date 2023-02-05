package routes

import (
	"github.com/gin-gonic/gin"
	"tiktok_demo/controller"
	"tiktok_demo/middleware/jwt"
)

func User(r *gin.RouterGroup) {
	user := r.Group("/user")
	{
		user.POST("/register/", controller.Register)
		user.POST("/login/", controller.Login)
		user.GET("/", jwt.Auth(), controller.GetUserInfo)
	}
}
