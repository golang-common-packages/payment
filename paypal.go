package payment

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/plutov/paypal/v3"
)

// PaypalClient model for PayPal instance
type PaypalClient struct {
	client *paypal.Client
	auth   *paypal.TokenResponse
}

// NewPaypal return new PayPal instance
func NewPaypal(clientID, secretID string) *PaypalClient {
	currentClient := &PaypalClient{nil, nil}

	client, err := paypal.NewClient("clientID", "secretID", paypal.APIBaseLive)
	if err != nil {
		log.Println("Can not init PayPal client: ", err)
		panic(err)
	}
	client.SetLog(os.Stdout) // Set log to terminal stdout

	result, err := client.GetAccessToken()
	if err != nil {
		log.Println("Can not get access token: ", err)
		panic(err)
	}

	currentClient.client = client
	currentClient.auth = result
	log.Println("Init PayPal service success")

	return currentClient
}

// SetAccessToken set new access token to current PayPal client
func (pp *PaypalClient) SetAccessToken(accessToken string) {
	pp.client.SetAccessToken(accessToken)
}

// SetHTTPClient set new HTTP client to current PayPal client
func (pp *PaypalClient) SetHTTPClient(httpClient *http.Client) {
	pp.client.SetHTTPClient(httpClient)
}

// SetLog  set new log service to current PayPal client
func (pp *PaypalClient) SetLog(log io.Writer) {
	pp.client.SetLog(log)
}

// https://developer.paypal.com/docs/payouts/
// TransferMoney to send money to multiple people at the same time
func (pp *PaypalClient) TransferMoney(transferInfo *MoneyTransfer) (result interface{}, err error) {
	payout := paypal.Payout{
		SenderBatchHeader: &paypal.SenderBatchHeader{
			EmailSubject: transferInfo.EmailSubject,
		},
		Items: []paypal.PayoutItem{
			{
				RecipientType: transferInfo.ReceiverType,
				Receiver:      transferInfo.Receiver,
				Amount: &paypal.AmountPayout{
					Value:    transferInfo.Value,
					Currency: transferInfo.Currency,
				},
				Note: transferInfo.Comment,
			},
		},
	}

	payoutResp, err := pp.client.CreateSinglePayout(payout)
	if err != nil {
		return nil, err
	}

	return payoutResp, nil
}

// https://developer.paypal.com/docs/archive/adaptive-accounts/api/add-bank-account/
// LinkBankAccount to current user
func (pp *PaypalClient) LinkBankAccount(linkToPayPal PayPalLinkBank) (interface{}, error) {
	httpRequest, err := convertPayPalRequestToJSON(pp, "POST", "", linkToPayPal)
	if err != nil {
		return nil, err
	}

	var result interface{}
	err = pp.client.SendWithBasicAuth(httpRequest, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// StoreCreditCard to PalPay system based on current user
func (pp *PaypalClient) StoreCreditCard(creditCardInfo paypal.CreditCard) (result *paypal.CreditCard, err error) {
	return pp.client.StoreCreditCard(creditCardInfo)
}

// ListCreditCards belong to current user
func (pp *PaypalClient) ListCreditCards(page, pageSize int) (result *paypal.CreditCards, err error) {
	creditCardFilter := &paypal.CreditCardsFilter{
		PageSize: pageSize,
		Page:     page,
	}

	return pp.client.GetCreditCards(creditCardFilter)
}

// GetCreditCardDetail belong to current user by creditCardID
func (pp *PaypalClient) GetCreditCardDetail(creditCardID string) (result *paypal.CreditCard, err error) {
	return pp.client.GetCreditCard(creditCardID)
}

// RemoveCreditCard out of current user
func (pp *PaypalClient) RemoveCreditCard(creditCardID string) (err error) {
	return pp.client.DeleteCreditCard(creditCardID)
}

///// PayPal Util /////

// convertPayPalRequestToJSON ...
func convertPayPalRequestToJSON(pp *PaypalClient, method, url string, payload interface{}) (httpClient *http.Request, err error) {
	return pp.client.NewRequest(method, url, payload)
}
