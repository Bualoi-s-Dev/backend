package main

import (
	"github.com/gin-gonic/gin"

	"github.com/Bualoi-s-Dev/backend/configs"
	"github.com/Bualoi-s-Dev/backend/routes"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	configs.LoadEnv()

	configs.ConnectMongoDB()

	r.Run()
}
