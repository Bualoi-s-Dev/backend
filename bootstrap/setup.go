package bootstrap

import (
	"github.com/Bualoi-s-Dev/backend/configs"
	"github.com/Bualoi-s-Dev/backend/controllers"
	"github.com/Bualoi-s-Dev/backend/middleware"
	"github.com/Bualoi-s-Dev/backend/routes"
	"github.com/Bualoi-s-Dev/backend/services"
	validators "github.com/Bualoi-s-Dev/backend/validator"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"

	database "github.com/Bualoi-s-Dev/backend/repositories/database"
	firebase "github.com/Bualoi-s-Dev/backend/repositories/firebase"
	s3 "github.com/Bualoi-s-Dev/backend/repositories/s3"
)

type ServerRepositories struct {
	packageRepo     *database.PackageRepository
	userRepo        *database.UserRepository
	appointmentRepo *database.AppointmentRepository
	s3Repo          *s3.S3Repository
	firebaseRepo    *firebase.FirebaseRepository
}

type ServerServices struct {
	packageService     *services.PackageService
	userService        *services.UserService
	appointmentService *services.AppointmentService
	s3Service          *services.S3Service
	firebaseService    *services.FirebaseService
}

func SetupServer(client *mongo.Database) (*gin.Engine, *ServerRepositories, *ServerServices) {
	r := gin.Default()
	r.RemoveExtraSlash = true

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://frontend-2gn.pages.dev/"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	authClient := configs.InitializeFirebaseAuth()

	// Validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validators.RegisterCustomValidators(v)
	}
	packageCollection := client.Collection("Package")
	subpackageCollection := client.Collection("Subpackage")
	userCollection := client.Collection("User")
	appointmentCollection := client.Collection("Appointment")
	busyTimeCollection := client.Collection("BusyTime")

	packageRepo := database.NewPackageRepository(packageCollection)
	subpackageRepo := database.NewSubpackageRepository(subpackageCollection)
	userRepo := database.NewUserRepository(userCollection)
	appointmentRepo := database.NewAppointmentRepository(appointmentCollection, busyTimeCollection)
	busyTimeRepo := database.NewBusyTimeRepository(busyTimeCollection)
	s3Repo := s3.NewS3Repository()
	firebaseRepo := firebase.NewFirebaseRepository(authClient)

	s3Service := services.NewS3Service(s3Repo)
	firebaseService := services.NewFirebaseService(firebaseRepo)
	appointmentService := services.NewAppointmentService(appointmentRepo, packageRepo, busyTimeRepo)
	subpackageService := services.NewSubpackageService(subpackageRepo)
	packageService := services.NewPackageService(packageRepo, s3Service, subpackageService)
	userService := services.NewUserService(userRepo, s3Service, packageService, authClient)
	busyTimeService := services.NewBusyTimeService(busyTimeRepo, subpackageRepo, packageRepo)

	packageController := controllers.NewPackageController(packageService, s3Service, userService)
	subPackageController := controllers.NewSubpackageController(subpackageService, packageService)

	appointmentController := controllers.NewAppointmentController(appointmentService, busyTimeService)
	userController := controllers.NewUserController(userService, s3Service, busyTimeService)
	BusyTimeController := controllers.NewBusyTimeController(busyTimeService)
	internalController := controllers.NewInternalController(firebaseService, s3Service)

	serverRepositories := &ServerRepositories{
		packageRepo:     packageRepo,
		userRepo:        userRepo,
		appointmentRepo: appointmentRepo,
		s3Repo:          s3Repo,
		firebaseRepo:    firebaseRepo,
	}
	serverServices := &ServerServices{
		packageService:     packageService,
		userService:        userService,
		appointmentService: appointmentService,
		s3Service:          s3Service,
		firebaseService:    firebaseService,
	}

	rateLimiter := middleware.NewRateLimiter(20, 5)
	r.Use(rateLimiter.RateLimitMiddleware())

	// Swagger
	routes.SwaggerRoutes(r)

	// Add routes
	routes.InternalRoutes(r, internalController)

	r.Use(middleware.FirebaseAuthMiddleware(authClient, client.Collection("User"), userService))

	routes.PackageRoutes(r, packageController)
	routes.SubpackageRoutes(r, subPackageController)
	routes.UserRoutes(r, userController)
	routes.AppointmentRoutes(r, appointmentController)
	routes.BusyTimeRoutes(r, BusyTimeController)

	return r, serverRepositories, serverServices
}
