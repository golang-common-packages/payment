package payment

import (
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/balance"
	"github.com/stripe/stripe-go/bankaccount"
	"github.com/stripe/stripe-go/paymentmethod"
	"github.com/stripe/stripe-go/paymentsource"
	"github.com/stripe/stripe-go/topup"
	"github.com/stripe/stripe-go/transfer"
)

type StripeClient struct{}

func NewStripeClient(apiKey string) *StripeClient {
	currentSesstion := &StripeClient{}
	stripe.Key = apiKey
	return currentSesstion
}

func (s *StripeClient) RetrieveBalance() (*stripe.Balance, error) {
	accountBalance, err := balance.Get(nil)
	return accountBalance, err
}

func (s *StripeClient) TopUpStripeBalance(amount int64, typeCurrentcy stripe.Currency, description string) (*stripe.Topup, error) {
	params := &stripe.TopupParams{
		Amount:              stripe.Int64(amount),
		Currency:            stripe.String(string(typeCurrentcy)),
		Description:         stripe.String(description),
		StatementDescriptor: stripe.String("Top-up"),
	}
	result, err := topup.New(params)

	return result, err
}

func (s *StripeClient) GetTopUpDetail(topUpID string) (*stripe.Topup, error) {
	detail, err := topup.Get(topUpID, nil)

	return detail, err
}

func (s *StripeClient) AddTopUpMetadata(topUpID, key, value string) (*stripe.Topup, error) {
	params := &stripe.TopupParams{}
	params.AddMetadata(key, value)
	result, err := topup.Update(topUpID, params)

	return result, err
}

func (s *StripeClient) ListTopUps(searchType, option, value string) *topup.Iter {
	params := &stripe.TopupListParams{}
	params.Filters.AddFilter(searchType, option, value)
	result := topup.List(params)

	return result
}

func (s *StripeClient) CancelPendingTopUp(topUpID string) (*stripe.Topup, error) {
	result, err := topup.Cancel("tu_123456789", nil)

	return result, err
}

func (s *StripeClient) Transfer(amount int64, typeCurrentcy stripe.Currency, method, description string) (*stripe.Transfer, error) {
	params := &stripe.TransferParams{
		Amount:      stripe.Int64(amount),
		Currency:    stripe.String(string(typeCurrentcy)),
		Destination: stripe.String(description),
		SourceType:  &method,
	}
	detail, err := transfer.New(params)

	return detail, err
}

func (s *StripeClient) GetTransferDetail(transferID string) (*stripe.Transfer, error) {
	detail, err := transfer.Get(transferID, nil)

	return detail, err
}

func (s *StripeClient) AddTransferMetadata(transferID, key, value string) (*stripe.Transfer, error) {
	params := &stripe.TransferParams{}
	params.AddMetadata(key, value)
	result, err := transfer.Update(transferID, params)

	return result, err
}

func (s *StripeClient) ListTransfers(searchType, option, value string) *transfer.Iter {
	params := &stripe.TransferListParams{}
	params.Filters.AddFilter(searchType, option, value)
	result := transfer.List(params)

	return result
}

func (s *StripeClient) addBankAccount(customerID, token, accountHolderName, accountHolderType, accountNumber, country, currency string) (*stripe.BankAccount, error) {
	params := &stripe.BankAccountParams{
		AccountHolderName: stripe.String(accountHolderName),
		AccountHolderType: stripe.String(accountHolderType),
		AccountNumber:     stripe.String(accountNumber),
		Country:           stripe.String(country),
		Currency:          stripe.String(currency),
		Customer:          stripe.String(customerID),
		Token:             stripe.String(token),
	}
	result, err := bankaccount.New(params)

	return result, err
}

func (s *StripeClient) RetrieveBankAccount(customerID, bankID string) (*stripe.BankAccount, error) {
	params := &stripe.BankAccountParams{
		Customer: stripe.String(customerID),
	}
	result, err := bankaccount.Get(
		bankID,
		params,
	)

	return result, err
}

func (s *StripeClient) AddBankAccountMetadata(customerID, bankID, key, value string) (*stripe.BankAccount, error) {
	params := &stripe.BankAccountParams{
		Customer: stripe.String(customerID),
	}
	params.AddMetadata(key, value)
	result, err := bankaccount.Update(
		bankID,
		params,
	)

	return result, err
}

func (s *StripeClient) VerifyBankAccount(customerID, bankID string, amounts [2]int64) (*stripe.PaymentSource, error) {
	params := &stripe.SourceVerifyParams{
		Amounts:  amounts,
		Customer: stripe.String(customerID),
	}
	result, err := paymentsource.Verify(bankID, params)

	return result, err
}

func (s *StripeClient) RemoveBankAccount(customerID, bankID string) (*stripe.BankAccount, error) {
	params := &stripe.BankAccountParams{
		Customer: stripe.String(customerID),
	}
	result, err := bankaccount.Del(
		bankID,
		params,
	)

	return result, err
}

func (s *StripeClient) ListBankAccounts(customerID, searchType, option, value string) *bankaccount.Iter {
	params := &stripe.BankAccountListParams{
		Customer: stripe.String(customerID),
	}
	params.Filters.AddFilter(searchType, option, value)
	result := bankaccount.List(params)

	return result
}

func (s *StripeClient) CreatePayment(cardNumber, expMonth, expYear, cvc string) (*stripe.PaymentMethod, error) {
	params := &stripe.PaymentMethodParams{
		Type: stripe.String("card"),
		Card: &stripe.PaymentMethodCardParams{
			Number:   stripe.String(cardNumber),
			ExpMonth: stripe.String(expMonth),
			ExpYear:  stripe.String(expYear),
			CVC:      stripe.String(cvc),
		},
	}
	result, err := paymentmethod.New(params)

	return result, err
}

func (s *StripeClient) RetrievePayment(paymentID string) (*stripe.PaymentMethod, error) {
	result, err := paymentmethod.Get(
		paymentID,
		nil,
	)

	return result, err
}

func (s *StripeClient) AddPaymentMetadata(paymentID, key, value string) (*stripe.PaymentMethod, error) {
	params := &stripe.PaymentMethodParams{}
	params.AddMetadata(key, value)
	result, err := paymentmethod.Update(
		paymentID,
		params,
	)

	return result, err
}

func (s *StripeClient) ListPaymentByCustermerID(customerID, paymentType string) *paymentmethod.Iter {
	params := &stripe.PaymentMethodListParams{
		Customer: stripe.String(customerID),
		Type:     stripe.String(paymentType),
	}
	detail := paymentmethod.List(params)

	return detail
}

func (s *StripeClient) AttachPaymentToCustomer(customerID, paymentID string) (*stripe.PaymentMethod, error) {
	params := &stripe.PaymentMethodAttachParams{
		Customer: stripe.String(customerID),
	}
	result, err := paymentmethod.Attach(
		paymentID,
		params,
	)

	return result, err
}

func (s *StripeClient) DetachPaymentFromCustomer(customerID, paymentID string) (*stripe.PaymentMethod, error) {
	params := &stripe.PaymentMethodAttachParams{
		Customer: stripe.String(customerID),
	}
	result, err := paymentmethod.Attach(
		paymentID,
		params,
	)

	return result, err
}
