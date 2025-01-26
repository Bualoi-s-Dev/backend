package controllers

import (
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

func (s *S3Controller) UploadImage(c *gin.Context) {
	file, fileError := c.FormFile("image")
	if fileError != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image is required."})
	}

	url, uploadErr := s.Service.UploadFile(file)
	if uploadErr != nil {
		log.Println(uploadErr)
	} else {
		log.Println("Upload image suceessfully:", url)
	}
}
