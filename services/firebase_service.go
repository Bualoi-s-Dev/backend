package services

import (
	"firebase.google.com/go/auth"
	"github.com/Bualoi-s-Dev/backend/dto"
	repositories "github.com/Bualoi-s-Dev/backend/repositories/firebase"
	"github.com/gin-gonic/gin"
)

type FirebaseService struct {
	Repo *repositories.FirebaseRepository
}

func NewFirebaseService(repo *repositories.FirebaseRepository) *FirebaseService {
	return &FirebaseService{Repo: repo}
}

func (s *FirebaseService) Login(c *gin.Context, req dto.AuthUserCredentials) (string, error) {
	return s.Repo.LoginUser(c, req)
}

func (s *FirebaseService) Register(c *gin.Context, req dto.AuthUserCredentials) (*auth.UserRecord, error) {
	return s.Repo.RegisterUser(c, req)
}
