package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/golang-common-packages/hash"
)

// IPayPal interface for PayPal services
type IPayPal interface {
	GetAccessToken(ctx context.Context) (*TokenResponse, error)
}

// PayPalClient represents a Paypal REST API Client
type PayPalClient struct {
	sync.Mutex
	Client               *http.Client
	ClientID             string
	Secret               string
	APIBase              string
	Log                  io.Writer // If user set log file name all requests will be logged there
	Token                *TokenResponse
	tokenExpiresAt       time.Time
	returnRepresentation bool
}

const (
	// APIBaseSandBox points to the sandbox (for testing) version of the API
	APIBaseSandBox = "https://api.sandbox.paypal.com"

	// APIBaseLive points to the live version of the API
	APIBaseLive = "https://api.paypal.com"

	// RequestNewTokenBeforeExpiresIn is used by SendWithAuth and try to get new Token when it's about to expire
	RequestNewTokenBeforeExpiresIn = time.Duration(60) * time.Second
)

// payPalClientSessionMapping singleton pattern
var payPalClientSessionMapping = make(map[string]*PayPalClient)

// newPayPal init new instance.
// APIBase is a base API URL, for testing you can use paypal.APIBaseSandBox
func newPayPal(config *PayPal) IPayPal {
	// Validate config file
	if config.ClientID == "" || config.SecretID == "" || config.APIBase == "" {
		log.Fatalln("ClientID, Secret and APIBase are required to create a Client")
	}

	// Init PayPal client with singleton pattern
	hasher := &hash.Client{}
	configAsJSON, err := json.Marshal(config)
	if err != nil {
		log.Fatalln("Unable to marshal PayPal configuration: ", err)
	}
	configAsString := hasher.SHA1(string(configAsJSON))

	currentPayPalSession := payPalClientSessionMapping[configAsString]
	if currentPayPalSession == nil {

		currentPayPalSession.Client = &http.Client{}
		currentPayPalSession.ClientID = config.ClientID
		currentPayPalSession.Secret = config.SecretID
		currentPayPalSession.APIBase = config.APIBase
		payPalClientSessionMapping[configAsString] = currentPayPalSession

		log.Println("Init PayPal client successfully")
	}

	return currentPayPalSession
}

// GetAccessToken returns struct of TokenResponse.
// No need to call SetAccessToken to apply new access token for current Client.
// Endpoint: POST /v1/oauth2/token
func (c *PayPalClient) GetAccessToken(ctx context.Context) (*TokenResponse, error) {
	buf := bytes.NewBuffer([]byte("grant_type=client_credentials"))
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/oauth2/token"), buf)
	if err != nil {
		return &TokenResponse{}, err
	}

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")

	response := &TokenResponse{}
	err = c.SendWithBasicAuth(req, response)

	// Set Token for current Client
	if response.Token != "" {
		c.Token = response
		c.tokenExpiresAt = time.Now().Add(time.Duration(response.ExpiresIn) * time.Second)
	}

	return response, err
}

// CreatePayout submits a payout with an asynchronous API call, which immediately returns the results of a PayPal payment.
// For email payout set RecipientType: "EMAIL" and receiver email into Receiver
// Endpoint: POST /v1/payments/payouts
func (c *PayPalClient) CreatePayout(ctx context.Context, p Payout) (*PayoutResponse, error) {
	req, err := c.NewRequest(ctx, "POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/payouts"), p)
	response := &PayoutResponse{}

	if err != nil {
		return response, err
	}

	if err = c.SendWithAuth(req, response); err != nil {
		return response, err
	}

	return response, nil
}

// GetPayout shows the latest status of a batch payout along with the transaction status and other data for individual items.
// Also, returns IDs for the individual payout items. You can use these item IDs in other calls.
// Endpoint: GET /v1/payments/payouts/ID
func (c *PayPalClient) GetPayout(ctx context.Context, payoutBatchID string) (*PayoutResponse, error) {
	req, err := c.NewRequest(ctx, "GET", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/payouts/"+payoutBatchID), nil)
	response := &PayoutResponse{}

	if err != nil {
		return response, err
	}

	if err = c.SendWithAuth(req, response); err != nil {
		return response, err
	}

	return response, nil
}

// GetPayoutItem shows the details for a payout item.
// Use this call to review the current status of a previously unclaimed, or pending, payout item.
// Endpoint: GET /v1/payments/payouts-item/ID
func (c *PayPalClient) GetPayoutItem(ctx context.Context, payoutItemID string) (*PayoutItemResponse, error) {
	req, err := c.NewRequest(ctx, "GET", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/payouts-item/"+payoutItemID), nil)
	response := &PayoutItemResponse{}

	if err != nil {
		return response, err
	}

	if err = c.SendWithAuth(req, response); err != nil {
		return response, err
	}

	return response, nil
}

// CancelPayoutItem cancels an unclaimed Payout Item. If no one claims the unclaimed item within 30 days,
// the funds are automatically returned to the sender. Use this call to cancel the unclaimed item before the automatic 30-day refund.
// Endpoint: POST /v1/payments/payouts-item/ID/cancel
func (c *PayPalClient) CancelPayoutItem(ctx context.Context, payoutItemID string) (*PayoutItemResponse, error) {
	req, err := c.NewRequest(ctx, "POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/payouts-item/"+payoutItemID+"/cancel"), nil)
	response := &PayoutItemResponse{}

	if err != nil {
		return response, err
	}

	if err = c.SendWithAuth(req, response); err != nil {
		return response, err
	}

	return response, nil
}

// GetSale returns a sale by ID
// Use this call to get details about a sale transaction.
// Note: This call returns only the sales that were created via the REST API.
// Endpoint: GET /v1/payments/sale/ID
func (c *PayPalClient) GetSale(ctx context.Context, saleID string) (*Sale, error) {
	sale := &Sale{}

	req, err := c.NewRequest(ctx, "GET", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/sale/"+saleID), nil)
	if err != nil {
		return sale, err
	}

	if err = c.SendWithAuth(req, sale); err != nil {
		return sale, err
	}

	return sale, nil
}

// RefundSale refunds a completed payment.
// Use this call to refund a completed payment. Provide the sale_id in the URI and an empty JSON payload for a full refund. For partial refunds, you can include an amount.
// Endpoint: POST /v1/payments/sale/ID/refund
func (c *PayPalClient) RefundSale(ctx context.Context, saleID string, a *Amount) (*Refund, error) {
	type refundRequest struct {
		Amount *Amount `json:"amount"`
	}

	refund := &Refund{}

	req, err := c.NewRequest(ctx, "POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/sale/"+saleID+"/refund"), &refundRequest{Amount: a})
	if err != nil {
		return refund, err
	}

	if err = c.SendWithAuth(req, refund); err != nil {
		return refund, err
	}

	return refund, nil
}

// GetRefund by ID
// Use it to look up details of a specific refund on direct and captured payments.
// Endpoint: GET /v2/payments/refund/ID
func (c *PayPalClient) GetRefund(ctx context.Context, refundID string) (*Refund, error) {
	refund := &Refund{}

	req, err := c.NewRequest(ctx, "GET", fmt.Sprintf("%s%s", c.APIBase, "/v2/payments/refund/"+refundID), nil)
	if err != nil {
		return refund, err
	}

	if err = c.SendWithAuth(req, refund); err != nil {
		return refund, err
	}

	return refund, nil
}

// GetAuthorization returns an authorization by ID
// Endpoint: GET /v2/payments/authorizations/ID
func (c *PayPalClient) GetAuthorization(ctx context.Context, authID string) (*Authorization, error) {
	buf := bytes.NewBuffer([]byte(""))
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s%s%s", c.APIBase, "/v2/payments/authorizations/", authID), buf)
	auth := &Authorization{}

	if err != nil {
		return auth, err
	}

	err = c.SendWithAuth(req, auth)
	return auth, err
}

// CaptureAuthorization captures and process an existing authorization.
// To use this method, the original payment must have Intent set to "authorize"
// Endpoint: POST /v2/payments/authorizations/ID/capture
func (c *PayPalClient) CaptureAuthorization(ctx context.Context, authID string, paymentCaptureRequest *PaymentCaptureRequest) (*PaymentCaptureResponse, error) {
	return c.CaptureAuthorizationWithPaypalRequestId(ctx, authID, paymentCaptureRequest, "")
}

// CaptureAuthorization captures and process an existing authorization with idempotency.
// To use this method, the original payment must have Intent set to "authorize"
// Endpoint: POST /v2/payments/authorizations/ID/capture
func (c *PayPalClient) CaptureAuthorizationWithPaypalRequestId(ctx context.Context,
	authID string,
	paymentCaptureRequest *PaymentCaptureRequest,
	requestID string,
) (*PaymentCaptureResponse, error) {
	req, err := c.NewRequest(ctx, "POST", fmt.Sprintf("%s%s", c.APIBase, "/v2/payments/authorizations/"+authID+"/capture"), paymentCaptureRequest)
	paymentCaptureResponse := &PaymentCaptureResponse{}

	if err != nil {
		return paymentCaptureResponse, err
	}

	if requestID != "" {
		req.Header.Set("PayPal-Request-Id", requestID)
	}

	err = c.SendWithAuth(req, paymentCaptureResponse)
	return paymentCaptureResponse, err
}

// VoidAuthorization voids a previously authorized payment
// Endpoint: POST /v2/payments/authorizations/ID/void
func (c *PayPalClient) VoidAuthorization(ctx context.Context, authID string) (*Authorization, error) {
	buf := bytes.NewBuffer([]byte(""))
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s%s", c.APIBase, "/v2/payments/authorizations/"+authID+"/void"), buf)
	auth := &Authorization{}

	if err != nil {
		return auth, err
	}

	err = c.SendWithAuth(req, auth)
	return auth, err
}

// ReauthorizeAuthorization reauthorize a Paypal account payment.
// PayPal recommends reauthorizing payment after ~3 days
// Endpoint: POST /v2/payments/authorizations/ID/reauthorize
func (c *PayPalClient) ReauthorizeAuthorization(ctx context.Context, authID string, a *Amount) (*Authorization, error) {
	buf := bytes.NewBuffer([]byte(`{"amount":{"currency_code":"` + a.Currency + `","value":"` + a.Total + `"}}`))
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s%s", c.APIBase, "/v2/payments/authorizations/"+authID+"/reauthorize"), buf)
	auth := &Authorization{}

	if err != nil {
		return auth, err
	}

	err = c.SendWithAuth(req, auth)
	return auth, err
}

// GetCapturedPaymentDetails.
// Endpoint: GET /v1/payments/capture/:id
func (c *PayPalClient) GetCapturedPaymentDetails(ctx context.Context, id string) (*Capture, error) {
	res := &Capture{}

	req, err := c.NewRequest(ctx, "GET", fmt.Sprintf("%s%s%s", c.APIBase, "/v1/payments/capture/", id), nil)
	if err != nil {
		return res, err
	}

	if err = c.SendWithAuth(req, res); err != nil {
		return res, err
	}

	return res, nil
}