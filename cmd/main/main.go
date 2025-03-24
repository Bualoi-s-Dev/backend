package main

import (
	"context"
	"fmt"

	"github.com/Bualoi-s-Dev/backend/bootstrap"
	"github.com/Bualoi-s-Dev/backend/configs"
)

// @title PhotoMatch API
// @version 1.0
// @description This is a sample API to demonstrate Swagger with Gin.
// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com
// @host      localhost:8080

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	databaseName := "PhotoMatch"
	configs.LoadEnv()
	if configs.GetEnv("APP_MODE") == "development" {
		databaseName = "PhotoMatch_Dev"
	}
	// Connect to MongoDB
	client := configs.ConnectMongoDB().Database(databaseName)

	// Setup server
	r, _, serverService := bootstrap.SetupServer(client)

	// Auto update go routine
	go bootstrap.AutoUpdate(context.TODO(), serverService)

	// Run server
	port := configs.GetEnv("PORT")
	if port == "" {
		fmt.Println("PORT is not set")
		port = "8080"
	}
	fmt.Println("Server is running on port: " + port)
	r.Run(":" + port)
}
