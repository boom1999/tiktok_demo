package routes

import (
	"github.com/gin-gonic/gin"
	"tiktok_demo/controller"
	"tiktok_demo/middleware/jwt"
)

func Message(r *gin.RouterGroup) {
	message := r.Group("/message")
	{
		message.GET("/chat/", jwt.Auth(), controller.MessageChat)
		message.POST("/action/", jwt.Auth(), controller.MessageAction)
	}
}
