package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/Bualoi-s-Dev/backend/models"
	databaseRepo "github.com/Bualoi-s-Dev/backend/repositories/database"
	stripeRepo "github.com/Bualoi-s-Dev/backend/repositories/stripe"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/balancetransaction"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentService struct {
	DatabaseRepository            *databaseRepo.PaymentRepository
	UserDatabaseRepository        *databaseRepo.UserRepository
	AppointmentDatabaseRepository *databaseRepo.AppointmentRepository
	AppointmentService            *AppointmentService
	SubpackageDatabaseRepository  *databaseRepo.SubpackageRepository
	PackageDatabaseRepository     *databaseRepo.PackageRepository
	StripeRepository              *stripeRepo.StripeRepository
}

func NewPaymentService(databaseRepository *databaseRepo.PaymentRepository, userRepo *databaseRepo.UserRepository, appointmentRepo *databaseRepo.AppointmentRepository,
	subpackageRepo *databaseRepo.SubpackageRepository, packageRepo *databaseRepo.PackageRepository, stripeRepository *stripeRepo.StripeRepository) *PaymentService {
	return &PaymentService{
		DatabaseRepository:            databaseRepository,
		UserDatabaseRepository:        userRepo,
		AppointmentDatabaseRepository: appointmentRepo,
		SubpackageDatabaseRepository:  subpackageRepo,
		PackageDatabaseRepository:     packageRepo,
		StripeRepository:              stripeRepository,
	}
}

func (service *PaymentService) GetAllOwnedPayments(ctx context.Context, user models.User) ([]models.Payment, error) {
	return service.DatabaseRepository.GetByUserIDAndRole(ctx, user.Role, user.ID)
}

func (service *PaymentService) GetPaymentById(ctx context.Context, id primitive.ObjectID) (*models.Payment, error) {
	return service.DatabaseRepository.GetById(ctx, id)
}

func (service *PaymentService) RegisterCustomer(ctx context.Context, user models.User) (*stripe.Customer, error) {
	// Create stripe customer
	customer, err := service.StripeRepository.CreateCustomer(user.Email)
	if err != nil {
		return nil, err
	}

	// Update stripe customer id in user
	user.StripeCustomerID = &customer.ID
	_, err = service.UserDatabaseRepository.ReplaceUser(ctx, user.ID, &user)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (service *PaymentService) RegisterConnectedAccount(ctx context.Context, user models.User) (*stripe.Account, error) {
	// Create stripe connected account
	account, err := service.StripeRepository.CreateConnectedAccount(user.Email)
	fmt.Println("Create account", account, err)
	if err != nil {
		return nil, err
	}

	// Attach bank account
	// err = service.StripeRepository.AttachBankAccount(account.ID, "TH", "thb", user.BankAccount)
	// fmt.Println("Attach bank account", err)
	// if err != nil {
	// 	return nil, err
	// }

	// Attach account setting
	err = service.StripeRepository.AttachAccountSetting(account.ID)
	fmt.Println("Attach account setting", err)
	if err != nil {
		return nil, err
	}

	// Update stripe account id in user
	user.StripeAccountID = &account.ID
	_, err = service.UserDatabaseRepository.ReplaceUser(ctx, user.ID, &user)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (service *PaymentService) CreatePayment(ctx context.Context, appointmentId primitive.ObjectID) (*models.Payment, error) {
	// Get customer and photographer from appointment
	appointment, err := service.AppointmentDatabaseRepository.GetById(ctx, appointmentId)
	if err != nil {
		return nil, err
	}
	customer, err := service.UserDatabaseRepository.FindUserByID(ctx, appointment.CustomerID)
	if err != nil {
		return nil, err
	}
	photographer, err := service.UserDatabaseRepository.FindUserByID(ctx, appointment.PhotographerID)
	if err != nil {
		return nil, err
	}

	// Get customer stripe id, if not exist create new stripe customer
	if customer.Role != models.Customer {
		return nil, errors.New("user is not a customer")
	}
	var stripeCustomerId string
	if customer.StripeCustomerID == nil {
		stripeCustomer, err := service.RegisterCustomer(ctx, *customer)
		if err != nil {
			return nil, err
		}
		stripeCustomerId = stripeCustomer.ID
	} else {
		stripeCustomerId = *customer.StripeCustomerID
	}

	// Create photographer stripe account, if not exist create new stripe account
	if photographer.Role != models.Photographer {
		return nil, errors.New("user is not a photographer")
	}
	var stripeAccountId string
	if photographer.StripeAccountID == nil {
		stripeAccount, err := service.RegisterConnectedAccount(ctx, *photographer)
		if err != nil {
			return nil, err
		}
		stripeAccountId = stripeAccount.ID
	} else {
		stripeAccountId = *photographer.StripeAccountID
	}

	// Create checkout session for customer into photographer account
	subpackage := appointment.Subpackage
	checkoutSession, err := service.CreateCheckoutSession(stripeCustomerId, stripeAccountId, subpackage.Title, int64(appointment.Price))
	if err != nil {
		return nil, err
	}

	// Create payment
	payment := &models.Payment{
		ID:            primitive.NewObjectID(),
		AppointmentID: appointmentId,
		Customer: models.CustomerPayment{
			Status:     models.Unpaid,
			CheckoutID: &checkoutSession.ID,
		},
		Photographer: models.PhotographerPayment{
			Status: models.Wait,
		},
	}
	return payment, service.DatabaseRepository.Create(ctx, payment)
}

func (service *PaymentService) CreateAccountLink(ctx context.Context, accountId string) (*stripe.AccountLink, error) {
	return service.StripeRepository.CreateAccountLink(accountId)
}

func (service *PaymentService) CreateLoginLink(ctx context.Context, accountId string) (*stripe.LoginLink, error) {
	return service.StripeRepository.CreateLoginLink(accountId)
}

func (service *PaymentService) CreateCheckoutSession(customerId string, sellerAccountId string, productName string, amount int64) (*stripe.CheckoutSession, error) {
	stripeCheckout, err := service.StripeRepository.CreateCheckoutSession(customerId, sellerAccountId, productName, amount*100, 1, "thb")
	if err != nil {
		return nil, err
	}
	return stripeCheckout, nil
}

func (service *PaymentService) UpdateAccount(ctx context.Context, user models.User) error {
	// Re-Attach bank account
	err := service.StripeRepository.UpdateBankAccount(*user.StripeAccountID, user.BankAccount)
	if err != nil {
		return err
	}

	// Re-Attach account setting
	err = service.StripeRepository.AttachAccountSetting(*user.StripeAccountID)
	return err
}

func (service *PaymentService) UpdateCheckoutCompleted(ctx context.Context, checkoutSession stripe.CheckoutSession) error {
	// Get checkout session from payment
	payment, err := service.DatabaseRepository.GetByCheckoutID(ctx, checkoutSession.ID)
	if err != nil {
		return err
	}

	// Update customer payment status
	payment.Customer.Status = models.Paid
	payment.Customer.PaymentIntentID = &checkoutSession.PaymentIntent.ID

	// Update payment in database
	err = service.DatabaseRepository.Replace(ctx, payment.ID, payment)
	return err
}

func (service *PaymentService) PaidPhotographer(ctx context.Context, charge stripe.Charge) error {
	// // Get payment by payment intent id
	payment, err := service.DatabaseRepository.GetByPaymentIntentID(ctx, charge.PaymentIntent.ID)
	if err != nil {
		return err
	}
	if payment.Photographer.Status != models.Wait {
		return nil // Already paid
	}

	// Update photographer payout and status in payment
	payment.Photographer.BalanceTransactionID = &charge.BalanceTransaction.ID
	payment.Photographer.Status = models.InProcess

	// Update payment in database
	err = service.DatabaseRepository.Replace(ctx, payment.ID, payment)
	return err
}

func (service *PaymentService) UpdateSuccessPayoutPhotographer(ctx context.Context, payout stripe.Payout) error {
	// Get transaction list from payout
	params := &stripe.BalanceTransactionListParams{
		Payout: stripe.String(payout.ID),
	}

	i := balancetransaction.List(params)
	for i.Next() {
		tx := i.BalanceTransaction()

		// Get payment by balance transaction id
		payment, err := service.DatabaseRepository.GetByBalanceTransactionID(ctx, tx.ID)
		if err != nil {
			return err
		}

		// Update photographer payment status
		payment.Photographer.Status = models.Completed
		err = service.DatabaseRepository.Replace(ctx, payment.ID, payment)
		if err != nil {
			return err
		}
	}
	return nil
}
