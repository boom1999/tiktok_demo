package routes

import (
	"github.com/gin-gonic/gin"
	"tiktok_demo/controller"
	"tiktok_demo/middleware/jwt"
)

func Publish(r *gin.RouterGroup) {
	publish := r.Group("/publish")
	{
		publish.POST("/action/", jwt.Auth(), controller.Publish)
		publish.GET("/list/", jwt.Auth(), controller.PublishList)
	}
}
