package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Bualoi-s-Dev/backend/services"
	"github.com/gin-gonic/gin"
)

type S3Controller struct {
	Service *services.S3Service
}

func NewS3Controller(service *services.S3Service) *S3Controller {
	return &S3Controller{Service: service}
}

func (s *S3Controller) UploadProfileImage(c *gin.Context) {
	file, fileError := c.FormFile("image")
	if fileError != nil {
		fmt.Println(fileError)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image is required."})
	}
	username := "profileImage/" + c.PostForm("username")

	url, uploadErr := s.Service.UploadFile(file, username)
	if uploadErr != nil {
		log.Println(uploadErr)
		c.AbortWithError(500, uploadErr)
	} else {
		log.Println("Upload image suceessfully:", url)
	}
}
func (s *S3Controller) RemoveProfileImage(c *gin.Context) {
	username := "profileImage/" + c.PostForm("username")

	deleteErr := s.Service.DeleteObject(username)
	if deleteErr != nil {
		log.Println(deleteErr)
		c.AbortWithError(500, deleteErr)
	} else {
		log.Println("Removed image suceessfully:")
	}
}
