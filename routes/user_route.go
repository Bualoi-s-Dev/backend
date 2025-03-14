package routes

import (
	"github.com/Bualoi-s-Dev/backend/controllers"
	"github.com/Bualoi-s-Dev/backend/middleware"
	"github.com/Bualoi-s-Dev/backend/models"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, userController *controllers.UserController) {
	userRoutes := router.Group("/user")
	{
		userRoutes.GET("/me", userController.GetUserJWT)

		userRoutes.GET("/profile", userController.GetUserProfile)
		userRoutes.GET("/profile/:id", userController.GetUserProfileByID)
		userRoutes.PATCH("/profile", userController.UpdateUserProfile)

		userRoutes.GET("/photographers", userController.GetPhotographers)

		userRoutes.POST("/busytime", middleware.AllowRoles(models.Photographer), userController.CreateUserBusyTime)
		userRoutes.DELETE("/busytime/:busyTimeId", middleware.AllowRoles(models.Photographer), userController.DeleteUserBusyTime)
	}
}
