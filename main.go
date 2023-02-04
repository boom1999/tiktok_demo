package main

import (
	"github.com/gin-gonic/gin"
	"tiktok_demo/config"
	"tiktok_demo/repository"
	"tiktok_demo/routes"
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
}
