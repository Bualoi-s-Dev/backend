package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Bualoi-s-Dev/backend/models"
	"github.com/Bualoi-s-Dev/backend/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"firebase.google.com/go/auth"
	//firebase "firebase.google.com/go/v4"
	//"google.golang.org/api/option"
)

func FirebaseAuthMiddleware(authClient *auth.Client, userCollection *mongo.Collection, userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip authentication for Stripe webhooks
		if c.Request.URL.Path == "/payment/webhook" {
			c.Next() // Allow webhook requests
			return
		}

		// Skip authentication for Provider checking
		if c.Request.URL.Path == "/user/provider" {
			c.Next()
			return
		}

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
		user, err := userService.FindUser(ctx, email)
		if err != nil {
			log.Printf("[INFO] User not found: %s. Creating new user...", email)

			// Create new user
			newUser := models.NewUser(email)

			if err := userService.CreateUser(ctx, newUser); err != nil {
				log.Printf("[ERROR] Failed to create user: %v", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
				return
			}

			user = newUser
		}
		c.Set("user", user)
		c.Next()
	}
}

func GetUserFromContext(c *gin.Context) *models.User {
	user, exists := c.Get("user")
	if !exists {
		log.Println("[ERROR] User not found in context")
		return nil
	}
	return user.(*models.User)
}

func GetUserRoleFromContext(c *gin.Context) models.UserRole {
	user := GetUserFromContext(c)
	if user == nil {
		log.Println("[ERROR] Role not found in context")
		return models.Guest
	}
	return user.Role
}

func AllowRoles(userService *services.UserService, allowRoles ...models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := GetUserRoleFromContext(c)
		for _, allowRole := range allowRoles {
			if role == allowRole {
				c.Next()
				return
			}
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("%s cannot access this endpoint", role)})
		c.Abort()
	}
}

func CheckProviderByEmail(authClient *auth.Client, email string) ([]string, error) {
	ctx := context.Background()

	user, err := authClient.GetUserByEmail(ctx, email)
	if err != nil {
		if auth.IsUserNotFound(err) || strings.Contains(err.Error(), "no user") {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	var providers []string
	for _, info := range user.ProviderUserInfo {
		providers = append(providers, info.ProviderID)
	}

	return providers, nil
}
