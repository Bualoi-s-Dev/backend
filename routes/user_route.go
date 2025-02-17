package routes

import (
	"github.com/Bualoi-s-Dev/backend/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, userController *controllers.UserController) {
	userRoutes := router.Group("/user")
	{
		userRoutes.GET("/me", userController.GetUserJWT)

		userRoutes.GET("/profile", userController.GetUserProfile)
		userRoutes.GET("/profile/:id", userController.GetUserProfileByID)
		userRoutes.PATCH("/profile", userController.UpdateUserProfile)

		// userRoutes.PUT("/profile/showcase", userController.UpdateUserShowcasePackage)
	}
}
