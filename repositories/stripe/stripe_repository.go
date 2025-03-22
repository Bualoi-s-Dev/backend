package stripe

import (
	stripe "github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/account"
	"github.com/stripe/stripe-go/v81/bankaccount"
	"github.com/stripe/stripe-go/v81/checkout/session"
	"github.com/stripe/stripe-go/v81/customer"
	"github.com/stripe/stripe-go/v81/payout"
)

type StripeRepository struct {
}

func NewStripeRepository() *StripeRepository {
	return &StripeRepository{}
}

func (s *StripeRepository) CreateConnectedAccount(email string) (*stripe.Account, error) {
	params := &stripe.AccountParams{
		Type:  stripe.String("custom"), // Can be "standard", "express" or "custom"
		Email: stripe.String(email),
		Capabilities: &stripe.AccountCapabilitiesParams{
			CardPayments: &stripe.AccountCapabilitiesCardPaymentsParams{
				Requested: stripe.Bool(true),
			},
			Transfers: &stripe.AccountCapabilitiesTransfersParams{
				Requested: stripe.Bool(true),
			},
		},
	}
	acc, err := account.New(params)
	if err != nil {
		return nil, err
	}
	return acc, nil
}

func (s *StripeRepository) CreateCustomer(email string) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Email: stripe.String(email),
	}
	cust, err := customer.New(params)
	if err != nil {
		return nil, err
	}
	return cust, nil
}

func (s *StripeRepository) GetCustomerByEmail(email string) *stripe.Customer {
	params := &stripe.CustomerListParams{}
	params.Filters.AddFilter("email", "", email) // Filter by email

	iter := customer.List(params)
	for iter.Next() {
		c := iter.Customer()
		if c.Email == email {
			return c
		}
	}

	return nil
}

func (s *StripeRepository) GetAccountByEmail(email string) *stripe.Account {
	params := &stripe.AccountListParams{}
	params.Filters.AddFilter("email", "", email) // Filter by email

	iter := account.List(params)
	for iter.Next() {
		a := iter.Account()
		if a.Email == email {
			return a
		}
	}

	return nil
}

func (s *StripeRepository) AttachBankAccount(accountID, country, currency, accountNumber string) error {
	params := &stripe.BankAccountParams{
		Country:       stripe.String(country),
		Currency:      stripe.String(currency),
		AccountNumber: stripe.String(accountNumber),
	}
	params.SetStripeAccount(accountID)
	_, err := bankaccount.New(params)
	return err
}

// func CreateInvoice(customerID string, amount int64, currency string) (*stripe.Invoice, error) {
// 	// Step 1: Add an invoice item (the charge)
// 	_, err := invoiceitem.New(&stripe.InvoiceItemParams{
// 		Customer: stripe.String(customerID),
// 		Amount:   stripe.Int64(amount),
// 		Currency: stripe.String(currency),
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Step 2: Create the invoice
// 	inv, err := invoice.New(&stripe.InvoiceParams{
// 		Customer:    stripe.String(customerID),
// 		AutoAdvance: stripe.Bool(true), // Automatically finalize invoice
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return inv, nil
// }

func (s *StripeRepository) CreateCheckoutSession(customerId string, productName string, amount int64, quantity int64, currency string) (*stripe.CheckoutSession, error) {
	params := &stripe.CheckoutSessionParams{
		Customer:           stripe.String(customerId),
		Mode:               stripe.String(string(stripe.CheckoutSessionModePayment)),
		PaymentMethodTypes: stripe.StringSlice([]string{"card", "promptpay"}),
		SuccessURL:         stripe.String("https://your-website.com/success"),
		CancelURL:          stripe.String("https://your-website.com/cancel"),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(currency),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(productName),
					},
					UnitAmount: stripe.Int64(amount),
				},
				Quantity: stripe.Int64(quantity),
			},
		},
	}

	return session.New(params)
}

func (s *StripeRepository) CreatePayout(accountID string, amount int64, currency string) (*stripe.Payout, error) {
	params := &stripe.PayoutParams{
		Amount:   stripe.Int64(amount),
		Currency: stripe.String(currency),
	}
	params.SetStripeAccount(accountID)
	p, err := payout.New(params)
	if err != nil {
		return nil, err
	}
	return p, nil
}
