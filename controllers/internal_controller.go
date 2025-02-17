package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Bualoi-s-Dev/backend/dto"
	"github.com/Bualoi-s-Dev/backend/services"
	"github.com/gin-gonic/gin"
)

type InternalController struct {
	FirebaseService *services.FirebaseService
	S3Service       *services.S3Service
}

func NewInternalController(firebaseService *services.FirebaseService, s3Service *services.S3Service) *InternalController {
	return &InternalController{FirebaseService: firebaseService, S3Service: s3Service}
}

// <--- S3 --->

func (s *InternalController) UploadProfileImage(c *gin.Context) {
	file, fileError := c.FormFile("image")
	if fileError != nil {
		fmt.Println(fileError)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image is required."})
	}

	key := "profileImage/" + c.PostForm("username")

	imgKey, uploadErr := s.S3Service.UploadFile(file, key)
	url := os.Getenv("S3_PUBLIC_URL") + imgKey
	if uploadErr != nil {
		log.Println(uploadErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image."})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Upload image successfully: " + url})
	}
}

type ImageUploadRequest struct {
	ImageBase64 string `json:"image" binding:"required"`
	Username    string `json:"username" binding:"required"`
}

func (s *InternalController) UploadBase64ProfileImage(c *gin.Context) {
	var req ImageUploadRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request, " + err.Error()})
		return
	}

	key := "profileImage/" + req.Username

	imgKey, uploadErr := s.S3Service.UploadBase64([]byte(req.ImageBase64), key)
	url := os.Getenv("S3_PUBLIC_URL") + imgKey
	if uploadErr != nil {
		log.Println(uploadErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image."})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Upload image successfully: " + url})
	}
}

func (s *InternalController) RemoveProfileImage(c *gin.Context) {
	username := "profileImage/" + c.PostForm("username")

	deleteErr := s.S3Service.DeleteObject(username)
	if deleteErr != nil {
		log.Println(deleteErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove image."})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Remove image successfully"})
	}
}

// <--- Firebase --->

func (ctrl *InternalController) Login(c *gin.Context) {
	var req dto.AuthUserCredentials
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := ctrl.FirebaseService.Login(c, req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to login, " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (ctrl *InternalController) Register(c *gin.Context) {
	var req dto.AuthUserCredentials
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ctrl.FirebaseService.Register(c, req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to register, " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
