package repositories

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"firebase.google.com/go/auth"
	"github.com/Bualoi-s-Dev/backend/dto"
	"github.com/gin-gonic/gin"
)

type FirebaseRepository struct {
	authClient *auth.Client
}

func NewFirebaseRepository(authClient *auth.Client) *FirebaseRepository {
	return &FirebaseRepository{authClient: authClient}
}

// Register user using email & password
func (repo *FirebaseRepository) RegisterUser(c *gin.Context, req dto.AuthUserCredentials) (*auth.UserRecord, error) {
	params := (&auth.UserToCreate{}).
		Email(req.Email).
		Password(req.Password)

	return repo.authClient.CreateUser(c, params)
}

// Login user using Firebase Custom Token
func (repo *FirebaseRepository) LoginUser(c *gin.Context, req dto.AuthUserCredentials) (string, error) {
	url := "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=" + os.Getenv("FIREBASE_WEB_API_KEY")

	payload, err := json.Marshal(map[string]string{
		"email":             req.Email,
		"password":          req.Password,
		"returnSecureToken": "true",
	})
	if err != nil {
		return "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Decode the response
	var loginResp struct {
		IDToken string `json:"idToken"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return "", err
	}

	// Handle login errors
	if loginResp.IDToken == "" {
		return "", errors.New("failed to get ID token")
	}

	return loginResp.IDToken, nil
}
