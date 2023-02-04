package routes

import (
	"github.com/gin-gonic/gin"
	"tiktok_demo/middleware/jwt"
)

func Comment(r *gin.RouterGroup) {
	comment := r.Group("/comment")
	{
		comment.POST("/action/", jwt.Auth())
		comment.GET("/list/", jwt.AuthWithoutLogin())
	}
}
