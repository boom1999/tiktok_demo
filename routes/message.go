package routes

import (
	"github.com/gin-gonic/gin"
	"tiktok_demo/middleware/jwt"
)

func Message(r *gin.RouterGroup) {
	message := r.Group("/message")
	{
		message.POST("/chat/", jwt.Auth())
		message.GET("/action/", jwt.Auth())
	}
}
