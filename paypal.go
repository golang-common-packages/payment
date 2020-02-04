package payment

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/plutov/paypal/v3"
)

// PaypalClient ...
type PaypalClient struct {
	client *paypal.Client
	token  *paypal.TokenResponse
}

// NewPaypal ...
func NewPaypal(clientID, secretID string) *PaypalClient {
	currentSesstion := &PaypalClient{nil, nil}

	client, err := paypal.NewClient("clientID", "secretID", paypal.APIBaseLive)
	if err != nil {
		log.Println("Error when init paypal client: ", err)
		panic(err)
	}
	client.SetLog(os.Stdout) // Set log to terminal stdout

	tokenResult, err := client.GetAccessToken()
	if err != nil {
		log.Println("Error when get paypal token: ", err)
		panic(err)
	}

	currentSesstion.client = client
	currentSesstion.token = tokenResult

	return currentSesstion
}

// GetAccessTokenFromAuthCode ...
func (pp *PaypalClient) GetAccessTokenFromAuthCode(code, redirectURL string) (result *paypal.TokenResponse, err error) {
	return pp.client.GrantNewAccessTokenFromAuthCode(code, redirectURL)
}

// GetAccessTokenFromRefreshToken ...
func (pp *PaypalClient) GetAccessTokenFromRefreshToken(refreshToken string) (result *paypal.TokenResponse, err error) {
	return pp.client.GrantNewAccessTokenFromRefreshToken(refreshToken)
}

// SetAccessToken ...
func (pp *PaypalClient) SetAccessToken(accessToken string) {
	pp.client.SetAccessToken(accessToken)
}

// SetHTTPClient ...
func (pp *PaypalClient) SetHTTPClient(httpClient *http.Client) {
	pp.client.SetHTTPClient(httpClient)
}

// SetLog ...
func (pp *PaypalClient) SetLog(log io.Writer) {
	pp.client.SetLog(log)
}

// SubmitPayment ...
func (pp *PaypalClient) SubmitPayment(emailSubject, recipientType, receiver, amount, currencyType, sendingNote string) (result *paypal.PayoutResponse, err error) {
	payout := paypal.Payout{
		SenderBatchHeader: &paypal.SenderBatchHeader{
			EmailSubject: emailSubject,
		},
		Items: []paypal.PayoutItem{
			paypal.PayoutItem{
				RecipientType: recipientType,
				Receiver:      receiver,
				Amount: &paypal.AmountPayout{
					Value:    amount,
					Currency: currencyType,
				},
				Note: sendingNote,
			},
		},
	}

	payoutResp, err := pp.client.CreateSinglePayout(payout)
	if err != nil {
		return nil, err
	}

	return payoutResp, nil
}

// ConvertRequestToJSON ...
func (pp *PaypalClient) ConvertRequestToJSON(method, url string, payload interface{}) (httpClient *http.Request, err error) {
	return pp.client.NewRequest(method, url, payload)
}

// SendRequest ...
func (pp *PaypalClient) SendRequest(req *http.Request, value interface{}) error {
	return pp.client.Send(req, value)
}

// SendRequestWithAuth ...
func (pp *PaypalClient) SendRequestWithAuth(req *http.Request, value interface{}) error {
	return pp.client.SendWithAuth(req, value)
}

// SendRequestWithBaseAuth ...
func (pp *PaypalClient) SendRequestWithBaseAuth(req *http.Request, value interface{}) error {
	return pp.client.SendWithBasicAuth(req, value)
}

// CreateBillingAgreement ...
func (pp *PaypalClient) CreateBillingAgreement(bill paypal.BillingAgreement) (result *paypal.CreateAgreementResp, err error) {
	return pp.client.CreateBillingAgreement(bill)
}

// CreateBillingPlan ...
func (pp *PaypalClient) CreateBillingPlan(bill paypal.BillingPlan) (result *paypal.CreateBillingResp, err error) {
	return pp.client.CreateBillingPlan(bill)
}

// GetPayment ...
func (pp *PaypalClient) GetPayment(payoutBatchID string) (result *paypal.PayoutResponse, err error) {
	return pp.client.GetPayout(payoutBatchID)
}

// GetPaymentItem ...
func (pp *PaypalClient) GetPaymentItem(payoutBatchID string) (result *paypal.PayoutItemResponse, err error) {
	return pp.client.GetPayoutItem(payoutBatchID)
}

// ListCreditCards ...
func (pp *PaypalClient) ListCreditCards(page, pageSize int) (result *paypal.CreditCards, err error) {
	creditCardFilter := &paypal.CreditCardsFilter{
		PageSize: pageSize,
		Page:     page,
	}

	return pp.client.GetCreditCards(creditCardFilter)
}

// GetCreditCardDetail ...
func (pp *PaypalClient) GetCreditCardDetail(creditCardID string) (result *paypal.CreditCard, err error) {
	return pp.client.GetCreditCard(creditCardID)
}

// StoreCreditCardDetail ...
func (pp *PaypalClient) StoreCreditCardDetail(line1, line2, city, countryCode, postalCode, state, phone, id, payerID, externalCustomerID, number, typeCard, expireMonth, expireYear, cvv2, firstName, lastName, State, ValidUntil string) (result *paypal.CreditCard, err error) {
	billingAddress := generateAddress(line1, line2, city, countryCode, postalCode, state, phone)
	creditCard := generateCreditDetail(id, payerID, externalCustomerID, number, typeCard, expireMonth, expireYear, cvv2, firstName, lastName, State, ValidUntil, billingAddress)

	return pp.client.StoreCreditCard(creditCard)
}

// RemoveCreditCardDetail ...
func (pp *PaypalClient) RemoveCreditCardDetail(creditCardID string) error {
	err := pp.client.DeleteCreditCard(creditCardID)
	return err
}

// Paypal util

// generateAddress ...
func generateAddress(line1, line2, city, countryCode, postalCode, state, phone string) *paypal.Address {
	address := &paypal.Address{
		Line1:       line1,
		Line2:       line2,
		City:        city,
		CountryCode: countryCode,
		PostalCode:  postalCode,
		State:       state,
		Phone:       phone,
	}

	return address
}

// generateCreditDetail ...
func generateCreditDetail(id, payerID, externalCustomerID, number, typeCard, expireMonth, expireYear, cvv2, firstName, lastName, State, ValidUntil string, billingAddress *paypal.Address) paypal.CreditCard {
	creditCardDetail := paypal.CreditCard{
		ID:                 id,
		PayerID:            payerID,
		ExternalCustomerID: externalCustomerID,
		Number:             number,
		Type:               typeCard,
		ExpireMonth:        expireMonth,
		ExpireYear:         expireYear,
		CVV2:               cvv2,
		FirstName:          firstName,
		LastName:           lastName,
		BillingAddress:     billingAddress,
		State:              State,
		ValidUntil:         ValidUntil,
	}

	return creditCardDetail
}
