package routes

import (
	"tiktok_demo/controller"
	"tiktok_demo/middleware/jwt"

	"github.com/gin-gonic/gin"
)

func Favorite(r *gin.RouterGroup) {
	favorite := r.Group("/favorite")
	{
		favorite.POST("/action/", jwt.Auth(), controller.FavoriteAction)
		favorite.GET("/list/", jwt.Auth(), controller.GetFavouriteList)
	}
}
