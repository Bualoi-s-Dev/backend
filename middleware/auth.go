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
)

func FirebaseAuthMiddleware(authClient *auth.Client, userCollection *mongo.Collection, userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip authentication for Stripe webhooks
		if c.Request.URL.Path == "/payment/webhook" {
			c.Next() // Allow webhook requests
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

			// Set Firebase Custom Claim for role
			err = authClient.SetCustomUserClaims(ctx, token.UID, map[string]interface{}{
				"role": string(models.Guest),
			})
			if err != nil {
				log.Printf("[ERROR] Failed to set custom claims: %v", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to set role"})
				return
			}
		}

		// Extract role from Firebase Custom Claims
		role, exists := token.Claims["role"].(string)
		if !exists {
			// Set default role if not exists
			role = string(models.Guest)
		}

		// Store user in context
		user.Role = models.UserRole(role)
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

func GetUserRoleFromContext(c *gin.Context, userService *services.UserService) models.UserRole {
	userRaw, exists := c.Get("user")
	if !exists {
		log.Println("[ERROR] User not found in context")
		return models.Guest
	}

	user, ok := userRaw.(*models.User)
	if !ok {
		log.Println("[ERROR] Invalid user type in context")
		return models.Guest
	}

	ctx := context.Background()
	// Re-fetch user from DB to ensure latest role
	latestUser, err := userService.FindUser(ctx, user.Email)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch user from DB: %v", err)
		return models.Guest
	}

	return latestUser.Role
}

func AllowRoles(userService *services.UserService, allowRoles ...models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := GetUserRoleFromContext(c, userService)
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
