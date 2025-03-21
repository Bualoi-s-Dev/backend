package routes

import (
	"github.com/Bualoi-s-Dev/backend/controllers"
	"github.com/Bualoi-s-Dev/backend/middleware"
	"github.com/Bualoi-s-Dev/backend/models"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, userController *controllers.UserController, RatingController *controllers.RatingController) {
	userGroup := router.Group("/user")
	commonRoutes := userGroup.Group("", middleware.AllowRoles(models.Photographer, models.Customer))
	{
		commonRoutes.GET("/photographers", userController.GetPhotographers)
		commonRoutes.GET("/:photographerId/rating", RatingController.GetAllRatingsFromPhotographer)
		commonRoutes.GET("/:photographerId/rating/:ratingId", RatingController.GetRatingById)

	}
	customerRoutes := userGroup.Group("", middleware.AllowRoles(models.Customer))
	{
		customerRoutes.POST("/:photographerId/rating", RatingController.CreateRating)
		customerRoutes.PUT("/:photographerId/rating/:ratingId", RatingController.UpdateRating)
		customerRoutes.DELETE("/:photographerId/rating/:ratingId", RatingController.DeleteRating)
	}
	photographerRoutes := userGroup.Group("", middleware.AllowRoles(models.Photographer))
	{
		photographerRoutes.POST("/busytime", userController.CreateUserBusyTime)
		photographerRoutes.DELETE("/busytime/:busyTimeId", userController.DeleteUserBusyTime)
	}
	publicRoutes := userGroup.Group("")
	{
		publicRoutes.GET("/me", userController.GetUserJWT)

		publicRoutes.GET("/profile", userController.GetUserProfile)
		publicRoutes.GET("/profile/:id", userController.GetUserProfileByID)
		publicRoutes.PATCH("/profile", userController.UpdateUserProfile)
	}

}
