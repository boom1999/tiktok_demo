package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"tiktok_demo/config"
	"tiktok_demo/middleware/minio"
	"tiktok_demo/middleware/rabbitmq"
	"tiktok_demo/middleware/redis"
	"tiktok_demo/repository"
	"tiktok_demo/routes"
	"tiktok_demo/util"
)

func main() {
	Init()

	r := gin.Default()
	routes.CollectRoutes(r)
	err := r.Run(":8080")
	if err != nil {
		log.Println("Start failed.")
	}
}

func Init() {
	config.LoadConfig()
	util.LogConfig()
	repository.InitDataBase()
	minio.InitMinio()
	redis.InitRedis()
	rabbitmq.InitRabbitMQ()
	rabbitmq.InitFollowRabbitMQ()
	rabbitmq.InitLikeRabbitMQ()
}
