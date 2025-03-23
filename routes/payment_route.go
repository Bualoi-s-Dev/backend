package routes

import (
	"github.com/Bualoi-s-Dev/backend/controllers"
	"github.com/Bualoi-s-Dev/backend/middleware"
	"github.com/Bualoi-s-Dev/backend/models"
	"github.com/gin-gonic/gin"
)

func PaymentRoutes(router *gin.Engine, ctrl *controllers.PaymentController) {
	paymentRoutes := router.Group("/payment")
	commonRoutes := paymentRoutes.Group("", middleware.AllowRoles(models.Photographer, models.Customer))
	{
		commonRoutes.GET("", ctrl.GetAllOwnedPayments)
		commonRoutes.GET("/:id", ctrl.GetPaymentById)
	}
	photographerRoutes := paymentRoutes.Group("", middleware.AllowRoles(models.Photographer))
	{
		photographerRoutes.GET("/onboardingURL", ctrl.GetOnBoardAccountURL)
		photographerRoutes.GET("/loginDashboard", ctrl.GetLoginLinkAccountURL)
	}
	paymentRoutes.POST("/charge/:appointmentId", ctrl.CreatePayment)
	paymentRoutes.POST("/webhook", ctrl.WebhookListener)
	paymentRoutes.GET("/test", ctrl.Test)
}
