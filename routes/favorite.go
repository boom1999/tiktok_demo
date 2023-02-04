package routes

import (
	"github.com/gin-gonic/gin"
	"tiktok_demo/middleware/jwt"
)

func Favorite(r *gin.RouterGroup) {
	favorite := r.Group("/favorite")
	{
		favorite.POST("/action/", jwt.Auth())
		favorite.GET("/list/", jwt.Auth())
	}
}
