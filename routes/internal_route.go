package routes

import (
	"os"

	"github.com/Bualoi-s-Dev/backend/controllers"
	"github.com/gin-gonic/gin"
)

func InternalRoutes(router *gin.Engine, ctrl *controllers.InternalController) {
	if os.Getenv("APP_MODE") == "development" {
		internalGroup := router.Group("/internal")

		s3Group := internalGroup.Group("/s3")
		{
			s3Group.POST("/upload/image", ctrl.UploadProfileImage)
			s3Group.POST("/upload/image/base64", ctrl.UploadBase64ProfileImage)
			s3Group.DELETE("/delete/image", ctrl.RemoveProfileImage)
		}
		firebaseGroup := internalGroup.Group("/firebase")
		{
			firebaseGroup.POST("/login", ctrl.Login)
			firebaseGroup.POST("/register", ctrl.Register)
		}

		internalGroup.GET("/health", ctrl.HealthCheck)
	}

}
