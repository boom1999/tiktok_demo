package routes

import (
	"github.com/gin-gonic/gin"
	"tiktok_demo/controller"
	"tiktok_demo/middleware/jwt"
)

func Relation(r *gin.RouterGroup) {
	relation := r.Group("/relation")
	{
		relation.POST("/action/", jwt.Auth(), controller.RelationAction)
		relation.GET("/follow/list/", jwt.Auth(), controller.GetFollowingList)
		relation.GET("/follower/list/", jwt.Auth(), controller.GetFollowersList)
		relation.GET("/friend/list/", jwt.Auth())
	}
}
