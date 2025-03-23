package stripe

import (
	stripe "github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/account"
	"github.com/stripe/stripe-go/v81/accountlink"
	"github.com/stripe/stripe-go/v81/accountsession"
	"github.com/stripe/stripe-go/v81/bankaccount"
	"github.com/stripe/stripe-go/v81/checkout/session"
	"github.com/stripe/stripe-go/v81/customer"
	"github.com/stripe/stripe-go/v81/loginlink"
	"github.com/stripe/stripe-go/v81/payout"
)

type StripeRepository struct {
}

func NewStripeRepository() *StripeRepository {
	return &StripeRepository{}
}

// TODO: {"status":400,"message":"Platforms in TH cannot create accounts where the platform is loss-liable, due to risk control measures. Please refer to our guide (https://support.stripe.com/questions/stripe-thailand-support-for-marketplaces) for more details.","request_id":"req_YqPlgTCvJiVNzw","request_log_url":"https://dashboard.stripe.com/test/logs/req_YqPlgTCvJiVNzw?t=1742682891","type":"invalid_request_error"}
func (s *StripeRepository) CreateConnectedAccount(email string) (*stripe.Account, error) {
	params := &stripe.AccountParams{
		Type:  stripe.String("standard"), // Can be "standard", "express" or "custom"
		Email: stripe.String(email),
		// Capabilities: &stripe.AccountCapabilitiesParams{
		// 	CardPayments: &stripe.AccountCapabilitiesCardPaymentsParams{
		// 		Requested: stripe.Bool(true),
		// 	},
		// 	Transfers: &stripe.AccountCapabilitiesTransfersParams{
		// 		Requested: stripe.Bool(true),
		// 	},
		// },
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
		Account:       stripe.String(accountID),
	}
	_, err := bankaccount.New(params)
	return err
}

func (s *StripeRepository) UpdateBankAccount(accountID, accountNumber string) error {
	// Delete all bank accounts in the account
	params := &stripe.BankAccountListParams{}
	params.SetStripeAccount(accountID)

	iter := bankaccount.List(params)
	for iter.Next() {
		ba := iter.BankAccount()

		_, err := bankaccount.Del(ba.ID, &stripe.BankAccountParams{Account: stripe.String(accountID)})
		if err != nil {
			return err
		}
	}

	// Attach new bank account
	err := s.AttachBankAccount(accountID, "TH", "thb", accountNumber)
	return err
}

func (s *StripeRepository) CreateAccountLink(accountID string) (*stripe.AccountLink, error) {
	params := &stripe.AccountLinkParams{
		Account: stripe.String(accountID),
		// TODO: Change the URL
		RefreshURL: stripe.String("https://example.com/reauth"),
		ReturnURL:  stripe.String("https://example.com/return"),
		Type:       stripe.String("account_onboarding"),
		Collect:    stripe.String("eventually_due"),
	}
	return accountlink.New(params)
}

func (s *StripeRepository) CreateAccountSession(accountID string) (*stripe.AccountSession, error) {
	params := &stripe.AccountSessionParams{
		Account: stripe.String(accountID),
		Components: &stripe.AccountSessionComponentsParams{
			AccountOnboarding: &stripe.AccountSessionComponentsAccountOnboardingParams{
				Enabled: stripe.Bool(true),
			},
			Payments: &stripe.AccountSessionComponentsPaymentsParams{
				Enabled: stripe.Bool(true),
			},
			Payouts: &stripe.AccountSessionComponentsPayoutsParams{
				Enabled: stripe.Bool(true),
			},
			Balances: &stripe.AccountSessionComponentsBalancesParams{
				Enabled: stripe.Bool(true),
			},
		},
	}
	return accountsession.New(params)
}

func (s *StripeRepository) AttachAccountSetting(accountID string) error {
	params := &stripe.AccountParams{
		Settings: &stripe.AccountSettingsParams{
			Payouts: &stripe.AccountSettingsPayoutsParams{
				Schedule: &stripe.AccountSettingsPayoutsScheduleParams{
					Interval: stripe.String("daily"), // Options: "daily", "weekly", "monthly", "manual"
				},
			},
		},
	}
	_, err := account.Update(accountID, params)
	return err
}

func (s *StripeRepository) CreateCheckoutSession(customerId string, sellerAccountId string, productName string, amount int64, quantity int64, currency string) (*stripe.CheckoutSession, error) {
	params := &stripe.CheckoutSessionParams{
		Customer:           stripe.String(customerId),
		Mode:               stripe.String(string(stripe.CheckoutSessionModePayment)),
		PaymentMethodTypes: stripe.StringSlice([]string{"card", "promptpay"}),
		// TODO: Change the URLs
		SuccessURL: stripe.String("https://your-website.com/success"),
		CancelURL:  stripe.String("https://your-website.com/cancel"),
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
		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			OnBehalfOf: stripe.String(sellerAccountId), // Seller's connected account ID
			// TODO: Change the application fee amount
			ApplicationFeeAmount: stripe.Int64(500), // Platform fee (5 THB)
			TransferData: &stripe.CheckoutSessionPaymentIntentDataTransferDataParams{
				Destination: stripe.String(sellerAccountId), // Seller's connected account ID
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
