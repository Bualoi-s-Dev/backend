package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/Bualoi-s-Dev/backend/models"
	"github.com/Bualoi-s-Dev/backend/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"firebase.google.com/go/auth"
)

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
		user, err := userService.FindUser(ctx, email)
		if err != nil {
			log.Printf("[INFO] User not found: %s. Creating new user...", email)

			// Create new user
			newUser := models.User{
				ID:       primitive.NewObjectID(),
				Email:    email,
				Name:     "",
				Gender:   "",
				Profile:  "",
				Phone:    "",
				Location: "",
				Role:     models.Guest,
			}

			if err := userService.CreateUser(ctx, &newUser); err != nil {
				log.Printf("[ERROR] Failed to create user: %v", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
				return
			}

			user = &newUser

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

func GetUserRoleFromContext(c *gin.Context) models.UserRole {
	user := GetUserFromContext(c)
	if user == nil {
		log.Println("[ERROR] Role not found in context")
		return models.Guest
	}
	return user.Role
}