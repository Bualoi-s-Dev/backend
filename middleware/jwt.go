package middleware

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Bualoi-s-Dev/backend/models"
	"github.com/Bualoi-s-Dev/backend/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

func InitializeFirebaseAuth() *auth.Client {
	opt := option.WithCredentialsFile(os.Getenv("FIREBASE_CREDENTIALS"))
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Error initializing Firebase: %v", err)
	}
	authClient, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("Error getting Firebase Auth client: %v", err)
	}
	return authClient
}

func FirebaseAuthMiddleware(authClient *auth.Client, userCollection *mongo.Collection, userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Println("[ERROR] Missing Authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			log.Println("[ERROR] Invalid token format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}

		token, err := authClient.VerifyIDToken(context.Background(), tokenString)
		if err != nil {
			log.Printf("[ERROR] Invalid JWT token: %v\n", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Extract email from the verified token
		email, ok := token.Claims["email"].(string)
		if !ok || email == "" {
			log.Println("[ERROR] Email not found in token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Email not found in token"})
			return
		}

		// Check if user exists by email
		ctx := context.Background()
		user, err := userService.GetUserProfile(ctx, email)
		if err != nil {
			log.Printf("[INFO] User not found: %s. Creating new user...", email)

			// Create new user
			newUser := models.User{
				Email:    email,
				Name:     "",
				Gender:   "",
				Profile:  "",
				Phone:    "",
				Location: "",
			}

			if err := userService.CreateUser(ctx, &newUser); err != nil {
				log.Printf("[ERROR] Failed to create user: %v", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
				return
			}

			user = &newUser
		}

		// Store user in context
		c.Set("user", user)
		c.Next()
	}
}
