package bootstrap

import (
	"os"

	"github.com/Bualoi-s-Dev/backend/configs"
	"github.com/Bualoi-s-Dev/backend/controllers"
	"github.com/Bualoi-s-Dev/backend/middleware"
	"github.com/Bualoi-s-Dev/backend/routes"
	"github.com/Bualoi-s-Dev/backend/services"
	validators "github.com/Bualoi-s-Dev/backend/validator"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stripe/stripe-go/v81"
	"go.mongodb.org/mongo-driver/mongo"

	database "github.com/Bualoi-s-Dev/backend/repositories/database"
	firebase "github.com/Bualoi-s-Dev/backend/repositories/firebase"
	s3 "github.com/Bualoi-s-Dev/backend/repositories/s3"
	stripeRepo "github.com/Bualoi-s-Dev/backend/repositories/stripe"
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

	r.Use(configs.EnableCORS())

	authClient := configs.InitializeFirebaseAuth()
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	// Validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validators.RegisterCustomValidators(v)
	}
	packageCollection := client.Collection("Package")
	subpackageCollection := client.Collection("Subpackage")
	userCollection := client.Collection("User")
	appointmentCollection := client.Collection("Appointment")
	busyTimeCollection := client.Collection("BusyTime")
	paymentCollection := client.Collection("Payment")
	ratingCollection := client.Collection("Rating")

	packageRepo := database.NewPackageRepository(packageCollection)
	subpackageRepo := database.NewSubpackageRepository(subpackageCollection)
	userRepo := database.NewUserRepository(userCollection)
	appointmentRepo := database.NewAppointmentRepository(appointmentCollection, busyTimeCollection)
	busyTimeRepo := database.NewBusyTimeRepository(busyTimeCollection)
	s3Repo := s3.NewS3Repository()
	firebaseRepo := firebase.NewFirebaseRepository(authClient)
	paymentRepo := database.NewPaymentRepository(paymentCollection, appointmentCollection)
	stripeRepo := stripeRepo.NewStripeRepository()
	ratingRepo := database.NewRatingRepository(ratingCollection)

	s3Service := services.NewS3Service(s3Repo)
	firebaseService := services.NewFirebaseService(firebaseRepo)
	subpackageService := services.NewSubpackageService(subpackageRepo, packageRepo, busyTimeRepo)
	packageService := services.NewPackageService(packageRepo, s3Service, subpackageService, userRepo)
	ratingService := services.NewRatingService(ratingRepo)
	userService := services.NewUserService(userRepo, s3Service, packageService, subpackageService, authClient, ratingService)
	busyTimeService := services.NewBusyTimeService(busyTimeRepo, subpackageRepo, packageRepo)
	paymentService := services.NewPaymentService(paymentRepo, userRepo, appointmentRepo, stripeRepo)
	appointmentService := services.NewAppointmentService(appointmentRepo, packageRepo, subpackageRepo, busyTimeRepo, userRepo, paymentService)

	packageController := controllers.NewPackageController(packageService, s3Service, userService)
	subPackageController := controllers.NewSubpackageController(subpackageService, packageService)
	appointmentController := controllers.NewAppointmentController(appointmentService, busyTimeService)
	userController := controllers.NewUserController(userService, s3Service, busyTimeService)
	BusyTimeController := controllers.NewBusyTimeController(busyTimeService)
	internalController := controllers.NewInternalController(firebaseService, s3Service)
	paymentController := controllers.NewPaymentController(paymentService, appointmentService, packageService)
	RatingController := controllers.NewRatingController(ratingService, userService)

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
	routes.UserRoutes(r, userController, RatingController)
	routes.AppointmentRoutes(r, appointmentController)
	routes.BusyTimeRoutes(r, BusyTimeController)
	routes.PaymentRoutes(r, paymentController)

	return r, serverRepositories, serverServices
}
