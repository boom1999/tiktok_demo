package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"tiktok_demo/config"
	"tiktok_demo/middleware/minio"
	"tiktok_demo/repository"
	"tiktok_demo/routes"
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
	repository.InitDataBase()
	minio.InitMinio()
}
