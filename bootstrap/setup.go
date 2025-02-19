package bootstrap

import (
	"github.com/Bualoi-s-Dev/backend/controllers"
	"github.com/Bualoi-s-Dev/backend/middleware"
	"github.com/Bualoi-s-Dev/backend/models"
	"github.com/Bualoi-s-Dev/backend/routes"
	"github.com/Bualoi-s-Dev/backend/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"

	database "github.com/Bualoi-s-Dev/backend/repositories/database"
	firebase "github.com/Bualoi-s-Dev/backend/repositories/firebase"
	s3 "github.com/Bualoi-s-Dev/backend/repositories/s3"
)

func SetupServer(client *mongo.Database) *gin.Engine {
	r := gin.Default()
	r.RemoveExtraSlash = true

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	authClient := middleware.InitializeFirebaseAuth()

	// Validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		models.RegisterCustomValidators(v)
	}

	// Init
	packageRepo := database.NewPackageRepository(client.Collection("Package"))
	userRepo := database.NewUserRepository(client.Collection("User"))
	appointmentRepo := database.NewAppointmentRepository(client.Collection("Appointment"))
	s3Repo := s3.NewS3Repository()
	firebaseRepo := firebase.NewFirebaseRepository(authClient)

	s3Service := services.NewS3Service(s3Repo)
	firebaseService := services.NewFirebaseService(firebaseRepo)
	userService := services.NewUserService(userRepo, s3Service)
	appointmentService := services.NewAppointmentService(appointmentRepo)
	packageService := services.NewPackageService(packageRepo, s3Service)

	packageController := controllers.NewPackageController(packageService, s3Service, userService)
	userController := controllers.NewUserController(userService, s3Service)
	appointmentController := controllers.NewAppointmentController(appointmentService)
	internalController := controllers.NewInternalController(firebaseService, s3Service)

	// Swagger
	routes.SwaggerRoutes(r)

	// Add routes
	routes.InternalRoutes(r, internalController)

	r.Use(middleware.FirebaseAuthMiddleware(authClient, client.Collection("User"), userService))

	routes.PackageRoutes(r, packageController)
	routes.UserRoutes(r, userController)
	routes.AppointmentRoutes(r, appointmentController)

	return r
}
