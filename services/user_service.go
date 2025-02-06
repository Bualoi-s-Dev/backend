package services

import (
	"context"

	"github.com/Bualoi-s-Dev/backend/models"
	repositories "github.com/Bualoi-s-Dev/backend/repositories/database"
)

type UserService struct {
	Repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) GetUserProfile(ctx context.Context, email string) (*models.User, error) {
	return s.Repo.GetUserByEmail(ctx, email)
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
	return s.Repo.CreateUser(ctx, user)
}

func (s *UserService) UpdateUser(ctx context.Context, email string, updates map[string]interface{}) error {
	return s.Repo.UpdateUser(ctx, email, updates)
}

func (s *UserService) GetUserFromJWT(ctx context.Context, c *gin.Context, authClient *auth.Client) (*models.User, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil, &gin.Error{Err: http.ErrNoCookie, Type: gin.ErrorTypePublic}
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" {
		return nil, &gin.Error{Err: http.ErrNoCookie, Type: gin.ErrorTypePublic}
	}

	token, err := authClient.VerifyIDToken(ctx, tokenString)
	if err != nil {
		return nil, err
	}

	email, ok := token.Claims["email"].(string)
	if !ok || email == "" {
		return nil, &gin.Error{Err: http.ErrNoCookie, Type: gin.ErrorTypePublic}
	}

	return s.GetUserProfile(ctx, email)
}