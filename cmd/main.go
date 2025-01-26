package main

import (
	"github.com/gin-gonic/gin"

	"github.com/Bualoi-s-Dev/backend/configs"
	"github.com/Bualoi-s-Dev/backend/controllers"
	"github.com/Bualoi-s-Dev/backend/docs"
	repositories "github.com/Bualoi-s-Dev/backend/repositories/database"
	"github.com/Bualoi-s-Dev/backend/routes"
	"github.com/Bualoi-s-Dev/backend/services"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title PhotoMatch API
// @version 1.0
// @description This is a sample API to demonstrate Swagger with Gin.
// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com
func main() {
	r := gin.Default()

	configs.LoadEnv()
	client := configs.ConnectMongoDB().Database("PhotoMatch")

	// Init
	packageRepo := repositories.NewPackageRepository(client.Collection("Package"))

	packageService := services.NewPackageService(packageRepo)

	packageController := controllers.NewPackageController(packageService)

	s3Service := services.NewS3Service()
	s3Controller := controllers.NewS3Controller(s3Service)

	// Swagger UI route
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Add routes
	routes.PackageRoutes(r, packageController)
	routes.S3Routes(r, s3Controller)

	r.Run()
}
