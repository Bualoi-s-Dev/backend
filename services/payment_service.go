package services

import (
	"context"
	"errors"
	"log"

	"github.com/Bualoi-s-Dev/backend/models"
	databaseRepo "github.com/Bualoi-s-Dev/backend/repositories/database"
	stripeRepo "github.com/Bualoi-s-Dev/backend/repositories/stripe"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/balancetransaction"
	"github.com/stripe/stripe-go/v81/charge"
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
	return service.DatabaseRepository.GetByUserIDAndRole(ctx, user.Role, user.ID.Hex())
}

func (service *PaymentService) GetPaymentById(ctx context.Context, id string) (*models.Payment, error) {
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
	if err != nil {
		return nil, err
	}

	// Attach bank account
	err = service.StripeRepository.AttachBankAccount(account.ID, "TH", "thb", user.BankAccount)
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

func (service *PaymentService) CreatePayment(ctx context.Context, appointmentId primitive.ObjectID) error {
	// Get customer from appointment
	appointment, err := service.AppointmentDatabaseRepository.GetById(ctx, appointmentId)
	if err != nil {
		return err
	}
	customer, err := service.UserDatabaseRepository.FindUserByID(ctx, appointment.CustomerID)
	if err != nil {
		return err
	}

	// Get customer stripe id, if not exist create new stripe customer
	if customer.Role != models.Customer {
		return errors.New("user is not a customer")
	}
	var stripeCustomerId string
	if customer.StripeCustomerID == nil {
		stripeCustomer, err := service.RegisterCustomer(ctx, *customer)
		if err != nil {
			return err
		}
		stripeCustomerId = stripeCustomer.ID
	} else {
		stripeCustomerId = *customer.StripeCustomerID
	}

	// Create checkout session
	subpackage, err := service.SubpackageDatabaseRepository.GetById(ctx, appointment.SubpackageID.Hex())
	if err != nil {
		return err
	}
	checkoutSession, err := service.CreateCheckoutSession(stripeCustomerId, subpackage.Title, int64(appointment.Price))
	if err != nil {
		return err
	}

	// Create payment
	payment := &models.Payment{
		AppointmentID: appointmentId,
		Customer: models.CustomerPayment{
			Status:     models.Unpaid,
			CheckoutID: &checkoutSession.ID,
		},
		Photographer: models.PhotographerPayment{
			Status: models.Unpaid,
		},
	}
	return service.DatabaseRepository.Create(ctx, payment)
}

func (service *PaymentService) CreateCheckoutSession(customerId string, productName string, amount int64) (*stripe.CheckoutSession, error) {
	stripeCheckout, err := service.StripeRepository.CreateCheckoutSession(customerId, productName, amount*100, 1, "thb")
	if err != nil {
		return nil, err
	}
	return stripeCheckout, nil
}

func (service *PaymentService) CreatePayout(accountId string, amount int64) (*stripe.Payout, error) {
	payout, err := service.StripeRepository.CreatePayout(accountId, amount, "thb")
	if err != nil {
		return nil, err
	}
	return payout, nil
}

func (service *PaymentService) UpdateCustomerPaid(ctx context.Context, checkoutSession stripe.CheckoutSession) error {
	// Get checkout session from payment
	payment, err := service.DatabaseRepository.GetByCheckoutID(ctx, checkoutSession.ID)
	if err != nil {
		return err
	}

	// Update customer payment status
	payment.Customer.Status = models.Paid

	appointment, err := service.AppointmentDatabaseRepository.GetById(ctx, payment.AppointmentID)
	if err != nil {
		return err
	}
	photographer, err := service.UserDatabaseRepository.FindUserByID(ctx, appointment.PhotographerID)
	if err != nil {
		return err
	}

	// Get net charge from payment intent in checkout session
	params := &stripe.ChargeListParams{}
	params.Filters.AddFilter("payment_intent", "", checkoutSession.PaymentIntent.ID) // Filter by PaymentIntent ID
	chargeIter := charge.List(params)

	var balanceTransactionID string

	if chargeIter.Next() {
		ch := chargeIter.Charge()
		balanceTransactionID = ch.BalanceTransaction.ID
	} else {
		log.Fatalf("No charges found for PaymentIntent")
	}

	bt, err := balancetransaction.Get(balanceTransactionID, nil)
	if err != nil {
		log.Fatalf("Error retrieving balance transaction: %v", err)
	}

	// Get stripe account id, if not exist create new stripe account
	var stripeAccountId string
	if photographer.StripeAccountID == nil {
		stripeAccount, err := service.RegisterConnectedAccount(ctx, *photographer)
		if err != nil {
			return err
		}
		stripeAccountId = stripeAccount.ID
	} else {
		stripeAccountId = *photographer.StripeCustomerID
	}

	// Create payout
	payout, err := service.CreatePayout(stripeAccountId, bt.Net)
	if err != nil {
		return err
	}

	// Update photographer payout and status in payment
	payment.Photographer.PayoutID = &payout.ID
	payment.Photographer.Status = models.InTransit

	// Update payment in database
	err = service.DatabaseRepository.Replace(ctx, payment.ID.Hex(), payment)
	return err
}

func (service *PaymentService) UpdatePaidPhotographer(ctx context.Context, payout stripe.Payout) error {
	// Get payment by payout id
	payment, err := service.DatabaseRepository.GetByPayoutID(ctx, payout.ID)
	if err != nil {
		return err
	}

	// Update photographer payment status
	payment.Photographer.Status = models.Paid

	// Update payment in database
	return service.DatabaseRepository.Replace(ctx, payment.ID.Hex(), payment)
}
