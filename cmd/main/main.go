package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"github.com/Bualoi-s-Dev/backend/configs"
	"github.com/Bualoi-s-Dev/backend/controllers"
	"github.com/Bualoi-s-Dev/backend/middleware"
	"github.com/Bualoi-s-Dev/backend/models"
	repositories "github.com/Bualoi-s-Dev/backend/repositories/database"
	s3 "github.com/Bualoi-s-Dev/backend/repositories/s3"
	"github.com/Bualoi-s-Dev/backend/routes"
	"github.com/Bualoi-s-Dev/backend/services"
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

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	configs.LoadEnv()
	client := configs.ConnectMongoDB().Database("PhotoMatch")
	authClient := middleware.InitializeFirebaseAuth()

	// Validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		models.RegisterCustomValidators(v)
	}

	// Init
	packageRepo := repositories.NewPackageRepository(client.Collection("Package"))
	userRepo := repositories.NewUserRepository(client.Collection("User"))
	s3Repo := s3.NewS3Repository()

	userService := services.NewUserService(userRepo)
	s3Service := services.NewS3Service(s3Repo)
	packageService := services.NewPackageService(packageRepo, s3Service)

	packageController := controllers.NewPackageController(packageService, s3Service)
	userController := controllers.NewUserController(userService)
	s3Controller := controllers.NewS3Controller(s3Service)

	// Swagger
	routes.SwaggerRoutes(r)

	// Add routes
	r.Use(middleware.FirebaseAuthMiddleware(authClient, client.Collection("User"), userService))

	routes.PackageRoutes(r, packageController)
	routes.UserRoutes(r, userController)
	routes.S3Routes(r, s3Controller)

	r.Run()
}
