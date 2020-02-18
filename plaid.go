package payment

import (
	"log"
	"net/http"
	"time"

	"github.com/plaid/plaid-go/plaid"
)

// PlaidClient model for Plaid instance
type PlaidClient struct {
	client                           *plaid.Client
	publicToken, accessToken, itemID string
}

// NewPlaid return new Plaid instance
func NewPlaid(clientID, secretKey, publicKey string) *PlaidClient {
	currentClient := &PlaidClient{nil, "", "", ""}

	plaidClientOptions := plaid.ClientOptions{
		ClientID:    clientID,
		Secret:      secretKey,
		PublicKey:   publicKey,
		Environment: plaid.Production, // Available environments are Sandbox, Development, and Production
		HTTPClient:  &http.Client{},   // This parameter is optional
	}

	client, err := plaid.NewClient(plaidClientOptions)
	if err != nil {
		log.Println("Error when try to init Plaid client: ", err.Error())
		panic(err)
	}

	currentClient.client = client

	return currentClient
}

// GenerateAccessToken generate 'publicToken', 'accessToken', 'itemID' based on 'publicToken'
// and set them to Plaid instance
// 'publicToken' return from Plaid link bank WebUI
func (pc *PlaidClient) GenerateAccessToken(publicToken string) error {
	response, err := pc.client.ExchangePublicToken(publicToken)
	if err != nil {
		return err
	}

	pc.publicToken = publicToken
	pc.accessToken = response.AccessToken
	pc.itemID = response.ItemID

	return nil
}

// GetAccounts retrieves high-level information about all accounts associated with an bank
func (pc *PlaidClient) GetAccounts() (interface{}, error) {
	response, err := pc.client.GetAccounts(pc.accessToken)
	if err != nil {
		return nil, err
	}

	return response.Accounts, nil
}

// GetBalances return all balance for each account
func (pc *PlaidClient) GetBalances() (interface{}, error) {
	response, err := pc.client.GetBalances(pc.accessToken)
	if err != nil {
		return nil, err
	}

	return response.Accounts, nil
}

// CreatePayment for goods and return 'recipientID', 'paymentID' and 'paymentToken'
func (pc *PlaidClient) CreatePayment(plaidPayment PlaidPayment) (interface{}, error) {
	recipientCreateResp, err := pc.client.CreatePaymentRecipient(plaidPayment.ProductName, plaidPayment.IBAN, plaid.PaymentRecipientAddress{
		Street:     plaidPayment.Street,
		City:       plaidPayment.City,
		PostalCode: plaidPayment.PostalCode,
		Country:    plaidPayment.Country,
	})
	if err != nil {
		return nil, err
	}
	recipientID := recipientCreateResp.RecipientID

	paymentCreateResp, err := pc.client.CreatePayment(recipientID, "payment-ref", plaid.PaymentAmount{
		Currency: plaidPayment.Currency,
		Value:    plaidPayment.Amount,
	})
	if err != nil {
		return nil, err
	}
	paymentID := paymentCreateResp.PaymentID

	paymentTokenCreateResp, err := pc.client.CreatePaymentToken(paymentID)
	if err != nil {
		return nil, err
	}
	paymentToken := paymentTokenCreateResp.PaymentToken

	plaidPaymentResult := PlaidPaymentResult{
		RecipientID:  recipientID,
		PaymentID:    paymentID,
		PaymentToken: paymentToken,
	}

	return plaidPaymentResult, nil
}

// GetPaymentsHistory return Transactions history
func (pc *PlaidClient) GetPaymentsHistory(startDate, endDate string) (interface{}, error) {
	// By default, pull Transactions for the past 30 days
	if startDate == "" || endDate == "" {
		endDate = time.Now().Local().Format("2020-01-01")
		startDate = time.Now().Local().Add(-30 * 24 * time.Hour).Format("2020-01-01")
	}

	response, err := pc.client.GetTransactions(pc.accessToken, startDate, endDate)
	if err != nil {
		return nil, err
	}

	transactions := PlaidTransactionsHistory{
		Accounts:     response.Accounts,
		Transactions: response.Transactions,
	}

	return transactions, nil
}
