package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Bualoi-s-Dev/backend/middleware"
	"github.com/Bualoi-s-Dev/backend/services"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v81"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentController struct {
	Service *services.PaymentService
}

func NewPaymentController(service *services.PaymentService) *PaymentController {
	return &PaymentController{Service: service}
}

func (ctrl *PaymentController) GetAllOwnedPayments(c *gin.Context) {
	user := middleware.GetUserFromContext(c)

	payments, err := ctrl.Service.GetAllOwnedPayments(c.Request.Context(), *user)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, payments)
}

func (ctrl *PaymentController) GetPaymentById(c *gin.Context) {
	id := c.Param("id")

	payment, err := ctrl.Service.GetPaymentById(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, payment)
}

func (ctrl *PaymentController) WebhookListener(c *gin.Context) {
	const MaxBodyBytes = int64(65536)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)
	payload, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Error reading request body"})
		return
	}

	event := stripe.Event{}

	if err := json.Unmarshal(payload, &event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse webhook JSON"})
		return
	}

	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing checkout session JSON"})
			return
		}
		fmt.Printf("Checkout session %s completed", session.ID)
		// Handle the successful session completed
		ctrl.Service.UpdateCustomerPaid(c.Request.Context(), session)
	case "payout.paid":
		var payout stripe.Payout
		if err := json.Unmarshal(event.Data.Raw, &payout); err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing payout JSON"})
			return
		}
		fmt.Printf("Payout %s paid", payout.ID)
		// Handle the successful payout paid
		ctrl.Service.UpdatePaidPhotographer(c.Request.Context(), payout)
	default:
		fmt.Printf("Unhandled event type: %s\n", event.Type)
	}

	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

func (ctrl *PaymentController) Test(c *gin.Context) {
	user := middleware.GetUserFromContext(c)
	id := "67d2497084553958bcfc0f4b"
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}
	ctrl.Service.CreatePayment(c.Request.Context(), objectID, *user)
}
