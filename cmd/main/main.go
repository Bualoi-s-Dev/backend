package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"

	"github.com/Bualoi-s-Dev/backend/configs"
	"github.com/Bualoi-s-Dev/backend/controllers"
	"github.com/Bualoi-s-Dev/backend/docs"
	"github.com/Bualoi-s-Dev/backend/middleware"
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

	// Init
	packageRepo := repositories.NewPackageRepository(client.Collection("Package"))
	userRepo := repositories.NewUserRepository(client.Collection("User"))
	s3Repo := s3.NewS3Repository()

	packageService := services.NewPackageService(packageRepo)
	userService := services.NewUserService(userRepo)
	s3Service := services.NewS3Service(s3Repo)

	packageController := controllers.NewPackageController(packageService)
	userController := controllers.NewUserController(userService)
	s3Controller := controllers.NewS3Controller(s3Service)

	// Swagger UI route
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.Use(middleware.FirebaseAuthMiddleware(authClient, client.Collection("User"), userService))

	// Add routes
	routes.PackageRoutes(r, packageController)
	routes.UserRoutes(r, userController)
	routes.S3Routes(r, s3Controller)

	r.Run()
}


