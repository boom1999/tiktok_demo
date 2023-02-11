package main

import (
	"tiktok_demo/config"
	"tiktok_demo/middleware/redis"
	"tiktok_demo/repository"
	"tiktok_demo/routes"
	"tiktok_demo/util"

	"github.com/gin-gonic/gin"
)

func main() {
	Init()

	r := gin.Default()
	routes.CollectRoutes(r)
	err := r.Run(":8080")
	if err != nil {
		return
	}
}

func Init() {
	config.LoadConfig()
	repository.Init()
	redis.InitRedis()
	util.InitFilter()

}
