package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"github.com/Bualoi-s-Dev/backend/configs"
	"github.com/Bualoi-s-Dev/backend/controllers"
	"github.com/Bualoi-s-Dev/backend/docs"
	"github.com/Bualoi-s-Dev/backend/models"
	repositories "github.com/Bualoi-s-Dev/backend/repositories/database"
	s3 "github.com/Bualoi-s-Dev/backend/repositories/s3"
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
// @host      localhost:8080

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	r := gin.Default()

	configs.LoadEnv()
	client := configs.ConnectMongoDB().Database("PhotoMatch")

	// Validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		models.RegisterCustomValidators(v)
	}

	// Init
	packageRepo := repositories.NewPackageRepository(client.Collection("Package"))
	s3Repo := s3.NewS3Repository()

	packageService := services.NewPackageService(packageRepo)
	s3Service := services.NewS3Service(s3Repo)

	packageController := controllers.NewPackageController(packageService)
	s3Controller := controllers.NewS3Controller(s3Service)

	// Swagger UI route
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Add routes
	routes.PackageRoutes(r, packageController)
	routes.S3Routes(r, s3Controller)

	r.Run()
}
