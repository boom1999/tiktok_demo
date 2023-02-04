package routes

import (
	"github.com/gin-gonic/gin"
	"tiktok_demo/middleware/jwt"
)

func Relation(r *gin.RouterGroup) {
	relation := r.Group("/relation")
	{
		relation.POST("/action/", jwt.Auth())
		relation.GET("/follow/list/", jwt.Auth())
		relation.GET("/follower/list/", jwt.Auth())
		relation.GET("/friend/list/", jwt.Auth())
	}
}
