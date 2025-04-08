package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Bualoi-s-Dev/backend/dto"
	"github.com/Bualoi-s-Dev/backend/middleware"
	"github.com/Bualoi-s-Dev/backend/models"
	"github.com/Bualoi-s-Dev/backend/services"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/webhook"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentController struct {
	Service            *services.PaymentService
	AppointmentService *services.AppointmentService
	PackageService     *services.PackageService
}

func NewPaymentController(service *services.PaymentService, appointmentService *services.AppointmentService, packageService *services.PackageService) *PaymentController {
	return &PaymentController{Service: service, AppointmentService: appointmentService, PackageService: packageService}
}

// GetAllOwnedPayments godoc
// @Tags Payment
// @Summary Get a list of payment owned by the user
// @Description Retrieve all payments owned by the user in the jwt
// @Success 200 {object} []dto.PaymentResponse
// @Failure 400 {object} string "Bad Request"
// @Router /payment [get]
func (ctrl *PaymentController) GetAllOwnedPayments(c *gin.Context) {
	user := middleware.GetUserFromContext(c)

	payments, err := ctrl.Service.GetAllOwnedPayments(c.Request.Context(), *user)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var response []dto.PaymentResponse
	for _, payment := range payments {
		dto, err := ctrl.mapToPaymentResponse(c, payment)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		response = append(response, *dto)
	}
	if response == nil {
		response = []dto.PaymentResponse{}
	}
	c.JSON(200, response)
}

// GetPaymentById godoc
// @Tags Payment
// @Summary Get a payment by id
// @Description Retrieve a payment from given id which is owned by the user in the jwt
// @Param id path string true "Payment ID"
// @Success 200 {object} dto.PaymentResponse
// @Failure 400 {object} string "Bad Request"
// @Router /payment/{id} [get]
func (ctrl *PaymentController) GetPaymentById(c *gin.Context) {
	id := c.Param("id")
	user := middleware.GetUserFromContext(c)

	oid, _ := primitive.ObjectIDFromHex(id)
	payment, err := ctrl.Service.GetPaymentById(c.Request.Context(), oid)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			c.JSON(404, gin.H{"error": "Payment not found"})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	response, err := ctrl.mapToPaymentResponse(c, *payment)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if response.Appointment.CustomerID != user.ID && response.Appointment.PhotographerID != user.ID {
		c.JSON(403, gin.H{"error": "You are not authorized to view this payment"})
		return
	}
	c.JSON(200, response)
}

// CreatePayment godoc
// @Tags Payment
// @Summary Create a payment for the appointment, this usually called after the appointment is completed
// @Description Create a payment for the appointment
// @Param id path string true "Appointment ID"
// @Success 200 {object} dto.PaymentResponse
// @Failure 400 {object} string "Bad Request"
// @Router /payment/charge/{appointmentId} [post]
func (ctrl *PaymentController) CreatePayment(c *gin.Context) {
	id := c.Param("appointmentId")

	appointmentId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid appointment ID format"})
		return
	}

	payment, err := ctrl.Service.CreatePayment(c.Request.Context(), appointmentId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	response, err := ctrl.mapToPaymentResponse(c, *payment)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, response)
}

// GetOnBoardAccountURL godoc
// @Tags Payment
// @Summary Create stripe onboarding account URL for photographer
// @Description Create stripe onboarding account URL for photographer
// @Success 200 {object} dto.PaymentURL
// @Failure 400 {object} string "Bad Request"
// @Router /payment/onboardingURL [get]
func (ctrl *PaymentController) GetOnBoardAccountURL(c *gin.Context) {
	user := middleware.GetUserFromContext(c)

	var accountId string
	if user.StripeAccountID == nil {
		account, err := ctrl.Service.RegisterConnectedAccount(c.Request.Context(), *user)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		accountId = account.ID
	} else {
		accountId = *user.StripeAccountID
	}

	accountLink, err := ctrl.Service.CreateAccountLink(c.Request.Context(), accountId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	res := dto.PaymentURL{URL: accountLink.URL}
	c.JSON(200, res)
}

// // GetLoginLinkAccountURL godoc
// // @Tags Payment
// // @Summary Create stripe login account URL for photographer
// // @Description Create stripe login account URL for photographer
// // @Success 200 {object} dto.PaymentURL
// // @Failure 400 {object} string "Bad Request"
// // @Router /payment/loginDashboardURL [get]
// func (ctrl *PaymentController) GetLoginLinkAccountURL(c *gin.Context) {
// 	user := middleware.GetUserFromContext(c)

// 	if user.StripeAccountID == nil {
// 		c.JSON(400, gin.H{"error": "User does not have stripe account yet"})
// 		return
// 	}

// 	loginLink, err := ctrl.Service.CreateLoginLink(c.Request.Context(), *user.StripeAccountID)
// 	if err != nil {
// 		c.JSON(500, gin.H{"error": err.Error()})
// 		return
// 	}

// 	res := dto.PaymentURL{URL: loginLink.URL}
// 	c.JSON(200, res)
// }

func (ctrl *PaymentController) WebhookListener(c *gin.Context) {
	const MaxBodyBytes = int64(65536)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Error reading request body"})
		return
	}

	myAccountSecret := os.Getenv("STRIPE_WEBHOOK_MY_ACCOUNT_SECRET")
	connectedAccountSecret := os.Getenv("STRIPE_WEBHOOK_CONNECTED_ACCOUNT_SECRET")
	localSecret := os.Getenv("STRIPE_WEBHOOK_LOCAL_SECRET")
	signatureHeader := c.Request.Header.Get("Stripe-Signature")

	var event stripe.Event
	// Try to verify with "My Account" secret
	event, err = webhook.ConstructEvent(payload, signatureHeader, myAccountSecret)
	if err != nil {
		// If "My Account" verification fails, try with "Connected Account" secret
		event, err = webhook.ConstructEvent(payload, signatureHeader, connectedAccountSecret)
		if err != nil {

			event, err = webhook.ConstructEvent(payload, signatureHeader, localSecret)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Error verifying webhook signature"})
				return
			}
		}
	}

	fmt.Println("Received event: ", event.Type)
	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing checkout session JSON"})
			return
		}
		fmt.Printf("Checkout session %s completed\n", session.ID)
		// Handle the successful session completed
		ctrl.Service.UpdateCheckoutCompleted(c.Request.Context(), session)
	case "charge.updated":
		var charge stripe.Charge
		if err := json.Unmarshal(event.Data.Raw, &charge); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing charge JSON"})
			return
		}
		fmt.Printf("Charge %s updated\n", charge.ID)
		// Handle the successful charge updated
		err := ctrl.Service.PaidPhotographer(c.Request.Context(), charge)
		if err != nil {
			fmt.Println("Error updating photographer payment status: ", err)
		}
	case "payout.paid":
		var payout stripe.Payout
		if err := json.Unmarshal(event.Data.Raw, &payout); err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing payout JSON"})
			return
		}
		fmt.Printf("Payout %s paid\n", payout.ID)
		// Handle the successful payout paid
		ctrl.Service.UpdateSuccessPayoutPhotographer(c.Request.Context(), payout)
	default:
		fmt.Printf("Unhandled event type: %s\n", event.Type)
	}

	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

func (ctrl *PaymentController) Test(c *gin.Context) {
	// user := middleware.GetUserFromContext(c)
	// id := "67d2497084553958bcfc0f4b"
	// objectID, err := primitive.ObjectIDFromHex(id)
	// if err != nil {
	// 	c.JSON(400, gin.H{"error": "Invalid ID format"})
	// 	return
	// }
	// ctrl.Service.CreatePayment(c.Request.Context(), objectID, *user)
}

// it need to be here because of the circular dependency
func (ctrl *PaymentController) mapToPaymentResponse(ctx context.Context, payment models.Payment) (*dto.PaymentResponse, error) {
	appointment, err := ctrl.AppointmentService.AppointmentRepo.GetById(ctx, payment.AppointmentID)
	if err != nil {
		return nil, err
	}
	appointmentDetail, err := ctrl.AppointmentService.GetAppointmentDetailById(ctx, nil, appointment)
	if err != nil {
		return nil, err
	}
	pkg, err := ctrl.PackageService.GetById(ctx, appointment.Package.ID.Hex())
	if err != nil {
		return nil, err
	}
	pkgResponse, err := ctrl.PackageService.MappedToPackageResponse(ctx, pkg)
	if err != nil {
		return nil, err
	}
	return &dto.PaymentResponse{
		Payment:     payment,
		Appointment: *appointmentDetail,
		Package:     *pkgResponse,
	}, nil
}
