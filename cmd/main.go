package main

import (
	"github.com/gin-gonic/gin"

	"github.com/Bualoi-s-Dev/backend/configs"
	"github.com/Bualoi-s-Dev/backend/controllers"
	repositories "github.com/Bualoi-s-Dev/backend/repositories/database"
	"github.com/Bualoi-s-Dev/backend/routes"
	"github.com/Bualoi-s-Dev/backend/services"
)

func main() {
	r := gin.Default()

	configs.LoadEnv()
	client := configs.ConnectMongoDB().Database("PhotoMatch")

	// Init
	packageRepo := repositories.NewPackageRepository(client.Collection("Package"))

	packageService := services.NewPackageService(packageRepo)

	packageController := controllers.NewPackageController(packageService)

	// Add routes
	routes.PackageRoutes(r, packageController)

	r.Run()
}
