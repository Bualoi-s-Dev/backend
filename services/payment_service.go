package services

import (
	"context"
	"errors"
	"fmt"
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
	customer, err := service.StripeRepository.CreateCustomer(user.Email)
	if err != nil {
		return nil, err
	}

	user.StripeCustomerID = &customer.ID
	_, err = service.UserDatabaseRepository.ReplaceUser(ctx, user.ID, &user)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (service *PaymentService) RegisterConnectedAccount(ctx context.Context, user models.User) (*stripe.Account, error) {
	account, err := service.StripeRepository.CreateConnectedAccount(user.Email)
	if err != nil {
		return nil, err
	}

	err = service.StripeRepository.AttachBankAccount(account.ID, "TH", "thb", user.BankAccount)
	if err != nil {
		return nil, err
	}

	user.StripeAccountID = &account.ID
	_, err = service.UserDatabaseRepository.ReplaceUser(ctx, user.ID, &user)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (service *PaymentService) CreatePayment(ctx context.Context, appointmentId primitive.ObjectID, customer models.User) error {
	if customer.Role != models.Customer {
		return errors.New("user is not a customer")
	}
	var stripeCustomerId string
	if customer.StripeCustomerID == nil {
		stripeCustomer, err := service.RegisterCustomer(ctx, customer)
		if err != nil {
			return err
		}
		stripeCustomerId = stripeCustomer.ID
	} else {
		stripeCustomerId = *customer.StripeCustomerID
	}

	appointment, err := service.AppointmentDatabaseRepository.GetById(ctx, appointmentId)
	if err != nil {
		return err
	}
	subpackage, err := service.SubpackageDatabaseRepository.GetById(ctx, appointment.SubpackageID.Hex())
	if err != nil {
		return err
	}
	checkoutSession, err := service.CreateCheckoutSession(stripeCustomerId, subpackage.Title, int64(appointment.Price))
	if err != nil {
		return err
	}

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

	fmt.Println("Payment:", payment)
	return service.DatabaseRepository.Create(ctx, payment)
}

func (service *PaymentService) CreateCheckoutSession(customerId string, productName string, amount int64) (*stripe.CheckoutSession, error) {
	stripeCheckout, err := service.StripeRepository.CreateCheckoutSession(customerId, productName, amount*100, 1, "thb")
	if err != nil {
		return nil, err
	}
	return stripeCheckout, nil
}

func (service *PaymentService) CreatePayout(accountId string, amount int64) error {
	_, err := service.StripeRepository.CreatePayout(accountId, amount, "thb")
	if err != nil {
		return err
	}
	return nil
}

func (service *PaymentService) UpdateCustomerPaid(ctx context.Context, checkoutSession stripe.CheckoutSession) error {
	payment, err := service.DatabaseRepository.GetByCheckoutID(ctx, checkoutSession.ID)
	if err != nil {
		return err
	}

	payment.Customer.Status = models.Paid
	payment.Photographer.Status = models.InTransit

	err = service.DatabaseRepository.UpdateCustomerPayment(ctx, payment.ID.Hex(), &payment.Customer)
	if err != nil {
		return err
	}

	appointment, err := service.AppointmentDatabaseRepository.GetById(ctx, payment.AppointmentID)
	if err != nil {
		return err
	}
	photographer, err := service.UserDatabaseRepository.FindUserByID(ctx, appointment.PhotographerID)
	if err != nil {
		return err
	}

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

	err = service.CreatePayout(stripeAccountId, bt.Net)
	return err
}

func (service *PaymentService) UpdatePhotographerPayment(ctx context.Context, id primitive.ObjectID, photographer models.PhotographerPayment) error {
	return service.DatabaseRepository.UpdatePhotographerPayment(ctx, id.Hex(), &photographer)
}
