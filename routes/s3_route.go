package routes

import (
	"github.com/Bualoi-s-Dev/backend/controllers"
	"github.com/gin-gonic/gin"
)

func S3Routes(router *gin.Engine, ctrl *controllers.S3Controller) {
	group := router.Group("/s3/profile")
	{
		group.POST("/upload/image", ctrl.UploadProfileImage)
		group.POST("/upload/image/base64", ctrl.UploadBase64ProfileImage)
		group.DELETE("/delete/image", ctrl.RemoveProfileImage)
	}
}
