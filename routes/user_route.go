package routes

import (
	"github.com/Bualoi-s-Dev/backend/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, userController *controllers.UserController) {
	userRoutes := router.Group("/user")
	{
		userRoutes.GET("/me", userController.GetUserProfile)
		userRoutes.PUT("/me", userController.UpdateUserProfile)

		userRoutes.GET("/profile", userController.GetUserProfilePic)
		userRoutes.PUT("/profile", userController.UpdateUserProfilePic)
	}
}
