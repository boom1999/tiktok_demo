package routes

import (
	"github.com/gin-gonic/gin"
	"tiktok_demo/controller"
	"tiktok_demo/middleware/jwt"
)

func CollectRoutes(r *gin.Engine) *gin.Engine {
	tiktok := r.Group("/douyin")
	{
		tiktok.GET("/feed/", jwt.AuthWithoutLogin(), controller.Feed)
		User(tiktok)
		Publish(tiktok)
		Favorite(tiktok)
		Comment(tiktok)
		Relation(tiktok)
		Message(tiktok)
	}
	return r
}
