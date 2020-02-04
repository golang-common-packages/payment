package payment

import (
	"log"
	"net/http"

	"github.com/plaid/plaid-go/plaid"
)

// PlaidClient ...
type PlaidClient struct {
	client *plaid.Client
}

// NewPlaid ...
func NewPlaid(clientID, secretKey, publicKey string) *PlaidClient {
	currentSesstion := &PlaidClient{nil}

	clientOptions := plaid.ClientOptions{
		clientID,
		secretKey,
		publicKey,
		plaid.Production, // Available environments are Sandbox, Development, and Production
		&http.Client{},   // This parameter is optional
	}

	client, err := plaid.NewClient(clientOptions)
	if err != nil {
		log.Println("Error when try to init Plaid client: ", err.Error())
		panic(err)
	}

	currentSesstion.client = client

	return currentSesstion
}

// Auth ...
func (p *PlaidClient) Auth(accessToken string) (result plaid.GetAuthResponse, err error) {
	return p.client.GetAuth(accessToken)
}

// RotateAccessToken ...
func (p *PlaidClient) RotateAccessToken(accessToken string) (result plaid.InvalidateAccessTokenResponse, err error) {
	return p.client.InvalidateAccessToken(accessToken)
}

// OnetimeToken ...
func (p *PlaidClient) OnetimeToken(accessToken string) (result plaid.CreatePublicTokenResponse, err error) {
	return p.client.CreatePublicToken(accessToken)
}

// ConvertToAccessToken ...
func (p *PlaidClient) ConvertToAccessToken(publicToken string) (result plaid.ExchangePublicTokenResponse, err error) {
	return p.client.ExchangePublicToken(publicToken)
}

// GetItem ...
func (p *PlaidClient) GetItem(accessToken string) (result plaid.GetItemResponse, err error) {
	return p.client.GetItem(accessToken)
}

// RemoveItem ...
func (p *PlaidClient) RemoveItem(accessToken string) (result plaid.RemoveItemResponse, err error) {
	return p.client.RemoveItem(accessToken)
}

// GetBankByID ...
func (p *PlaidClient) GetBankByID(ID string) (result plaid.GetInstitutionByIDResponse, err error) {
	return p.client.GetInstitutionByID(ID)
}

// GetBanks ...
func (p *PlaidClient) GetBanks(count, offset int) (result plaid.GetInstitutionsResponse, err error) {
	return p.client.GetInstitutions(count, offset)
}

// GetBalances ...
func (p *PlaidClient) GetBalances(accessToken string) (result plaid.GetBalancesResponse, err error) {
	return p.client.GetBalances(accessToken)
}

// SendToPlaidAccount ...
func (p *PlaidClient) SendToPlaidAccount(recipientID, reference, moneyType string, amount float64) (result plaid.CreatePaymentResponse, err error) {
	return p.client.CreatePayment(recipientID, reference, plaid.PaymentAmount{
		Currency: moneyType,
		Value:    amount,
	})
}

// SendToAddress ...
func (p *PlaidClient) SendToAddress(street, city, postalCode, country, recipientName, iban, reference, moneyType string, amount float64) (result plaid.CreatePaymentResponse, err error) {
	paymentRecipientResponse, err := p.client.CreatePaymentRecipient(recipientName, iban, plaid.PaymentRecipientAddress{
		Street:     []string{street},
		City:       city,
		PostalCode: postalCode,
		Country:    country,
	})
	if err != nil {
		return plaid.CreatePaymentResponse{
			APIResponse: plaid.APIResponse{},
			PaymentID:   "",
			Status:      "",
		}, err
	}

	return p.SendToPlaidAccount(paymentRecipientResponse.RecipientID, reference, moneyType, amount)
}

// RegistryRecipientFromAddress ...
func (p *PlaidClient) RegistryRecipientFromAddress(street, city, postalCode, country, recipientName, iban string) (result plaid.CreatePaymentRecipientResponse, err error) {
	return p.client.CreatePaymentRecipient(recipientName, iban, plaid.PaymentRecipientAddress{
		Street:     []string{street},
		City:       city,
		PostalCode: postalCode,
		Country:    country,
	})
}

// GetAccountsInfo ...
func (p *PlaidClient) GetAccountsInfo(accessToken string) (result plaid.GetAccountsResponse, err error) {
	return p.client.GetAccounts(accessToken)
}

// GetReport ...
func (p *PlaidClient) GetReport(reportToken string) (result plaid.GetAssetReportResponse, err error) {
	return p.client.GetAssetReport(reportToken)
}

// For Sandbox //

// OnetimeSandboxToken ...
func (p *PlaidClient) OnetimeSandboxToken(institutionID string, initialProducts []string) (result plaid.CreateSandboxPublicTokenResponse, err error) {
	return p.client.CreateSandboxPublicToken(institutionID, initialProducts)
}

// ResetSandboxItem ...
func (p *PlaidClient) ResetSandboxItem(accessToken string) (result plaid.ResetSandboxItemResponse, err error) {
	return p.client.ResetSandboxItem(accessToken)
}