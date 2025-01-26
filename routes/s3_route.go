package routes

import (
	"github.com/Bualoi-s-Dev/backend/controllers"
	"github.com/gin-gonic/gin"
)

func S3Routes(router *gin.Engine, ctrl *controllers.S3Controller) {
	router.POST("/upload/image", ctrl.UploadImage)
}
