package routes

import (
	"github.com/gin-gonic/gin"
	"tiktok_demo/controller"
	"tiktok_demo/middleware/jwt"
)

func Comment(r *gin.RouterGroup) {
	comment := r.Group("/comment")
	{
		comment.POST("/action/", jwt.Auth(), controller.CommentAction)
		comment.GET("/list/", jwt.AuthWithoutLogin(), controller.CommentList)
	}
}
