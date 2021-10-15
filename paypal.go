package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang-common-packages/hash"
)

// IPayPal interface for PayPal services
type IPayPal interface {
	GetAccessToken(ctx context.Context) (*TokenResponse, error)
	CreatePayout(ctx context.Context, p Payout) (*PayoutResponse, error)
	GetPayout(ctx context.Context, payoutBatchID string) (*PayoutResponse, error)
	GetPayoutItem(ctx context.Context, payoutItemID string) (*PayoutItemResponse, error)
	CancelPayoutItem(ctx context.Context, payoutItemID string) (*PayoutItemResponse, error)
	GetSale(ctx context.Context, saleID string) (*Sale, error)
	RefundSale(ctx context.Context, saleID string, a *Amount) (*Refund, error)
	ListBillingPlans(ctx context.Context, bplp BillingPlanListParams) (*BillingPlanListResponse, error)
	CreateBillingPlan(ctx context.Context, plan BillingPlan) (*CreateBillingResponse, error)
	UpdateBillingPlan(ctx context.Context, planId string, pathValues map[string]map[string]interface{}) error
	ActivatePlan(ctx context.Context, planID string) error
	CreateBillingAgreement(ctx context.Context, a BillingAgreement) (*CreateAgreementResponse, error)
	ExecuteApprovedAgreement(ctx context.Context, token string) (*ExecuteAgreementResponse, error)
	GetAuthorization(ctx context.Context, authID string) (*Authorization, error)
	CaptureAuthorization(ctx context.Context, authID string, paymentCaptureRequest *PaymentCaptureRequest) (*PaymentCaptureResponse, error)
	CaptureAuthorizationWithPaypalRequestId(ctx context.Context, authID string, paymentCaptureRequest *PaymentCaptureRequest, requestID string) (*PaymentCaptureResponse, error)
	VoidAuthorization(ctx context.Context, authID string) (*Authorization, error)
	ReauthorizeAuthorization(ctx context.Context, authID string, a *Amount) (*Authorization, error)
	GetCapturedPaymentDetails(ctx context.Context, id string) (*Capture, error)
	GetRefund(ctx context.Context, refundID string) (*Refund, error)
	GetUserInfo(ctx context.Context, schema string) (*UserInfo, error)
	GrantNewAccessTokenFromAuthCode(ctx context.Context, code, redirectURI string) (*TokenResponse, error)
	GrantNewAccessTokenFromRefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error)
	CreateWebProfile(ctx context.Context, wp WebProfile) (*WebProfile, error)
	GetWebProfile(ctx context.Context, profileID string) (*WebProfile, error)
	GetWebProfiles(ctx context.Context) ([]WebProfile, error)
	SetWebProfile(ctx context.Context, wp WebProfile) error
	DeleteWebProfile(ctx context.Context, profileID string) error
	ListTransactions(ctx context.Context, req *TransactionSearchRequest) (*TransactionSearchResponse, error)
	StoreCreditCard(ctx context.Context, cc CreditCard) (*CreditCard, error)
	DeleteCreditCard(ctx context.Context, id string) error
	GetCreditCard(ctx context.Context, id string) (*CreditCard, error)
	GetCreditCards(ctx context.Context, ccf *CreditCardsFilter) (*CreditCards, error)
	PatchCreditCard(ctx context.Context, id string, ccf []CreditCardField) (*CreditCard, error)
	GetOrder(ctx context.Context, orderID string) (*Order, error)
	CreateOrder(ctx context.Context, intent string, purchaseUnits []PurchaseUnitRequest, payer *CreateOrderPayer, appContext *ApplicationContext) (*Order, error)
	UpdateOrder(ctx context.Context, orderID string, purchaseUnits []PurchaseUnitRequest) (*Order, error)
	AuthorizeOrder(ctx context.Context, orderID string, authorizeOrderRequest AuthorizeOrderRequest) (*Authorization, error)
	CaptureOrder(ctx context.Context, orderID string, captureOrderRequest CaptureOrderRequest) (*CaptureOrderResponse, error)
	CaptureOrderWithPaypalRequestId(ctx context.Context, orderID string, captureOrderRequest CaptureOrderRequest, requestID string) (*CaptureOrderResponse, error)
	CreateWebhook(ctx context.Context, createWebhookRequest *CreateWebhookRequest) (*Webhook, error)
	GetWebhook(ctx context.Context, webhookID string) (*Webhook, error)
	UpdateWebhook(ctx context.Context, webhookID string, fields []WebhookField) (*Webhook, error)
	ListWebhooks(ctx context.Context, anchorType string) (*ListWebhookResponse, error)
	DeleteWebhook(ctx context.Context, webhookID string) error
	VerifyWebhookSignature(ctx context.Context, httpReq *http.Request, webhookID string) (*VerifyWebhookResponse, error)
	GetWebhookEventTypes(ctx context.Context) (*WebhookEventTypesResponse, error)
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

	// BillingPlanStatusActive is used by BillingPlan and a few others
	BillingPlanStatusActive BillingPlanStatus = "ACTIVE"

	AncorTypeApplication string = "APPLICATION"
	AncorTypeAccount     string = "ACCOUNT"
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

// ListBillingPlans lists billing-plans
// Endpoint: GET /v1/payments/billing-plans
func (c *PayPalClient) ListBillingPlans(ctx context.Context, bplp BillingPlanListParams) (*BillingPlanListResponse, error) {
	req, err := c.NewRequest(ctx, "GET", fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/billing-plans"), nil)
	response := &BillingPlanListResponse{}
	if err != nil {
		return response, err
	}

	q := req.URL.Query()
	q.Add("page", bplp.Page)
	q.Add("page_size", bplp.PageSize)
	q.Add("status", bplp.Status)
	q.Add("total_required", bplp.TotalRequired)
	req.URL.RawQuery = q.Encode()

	err = c.SendWithAuth(req, response)
	return response, err
}

// CreateBillingPlan creates a billing plan in Paypal
// Endpoint: POST /v1/payments/billing-plans
func (c *PayPalClient) CreateBillingPlan(ctx context.Context, plan BillingPlan) (*CreateBillingResponse, error) {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/billing-plans"), plan)
	response := &CreateBillingResponse{}
	if err != nil {
		return response, err
	}
	err = c.SendWithAuth(req, response)
	return response, err
}

// UpdateBillingPlan updates values inside a billing plan
// Endpoint: PATCH /v1/payments/billing-plans
func (c *PayPalClient) UpdateBillingPlan(ctx context.Context, planId string, pathValues map[string]map[string]interface{}) error {
	patchData := []Patch{}
	for path, data := range pathValues {
		patchData = append(patchData, Patch{
			Operation: "replace",
			Path:      path,
			Value:     data,
		})
	}

	req, err := c.NewRequest(ctx, http.MethodPatch, fmt.Sprintf("%s%s%s", c.APIBase, "/v1/payments/billing-plans/", planId), patchData)
	if err != nil {
		return err
	}
	err = c.SendWithAuth(req, nil)
	return err
}

// ActivatePlan activates a billing plan
// By default, a new plan is not activated
// Endpoint: PATCH /v1/payments/billing-plans/
func (c *PayPalClient) ActivatePlan(ctx context.Context, planID string) error {
	return c.UpdateBillingPlan(ctx, planID, map[string]map[string]interface{}{
		"/": {"state": BillingPlanStatusActive},
	})
}

// CreateBillingAgreement creates an agreement for specified plan
// Endpoint: POST /v1/payments/billing-agreements
// Deprecated: Use POST /v1/billing-agreements/agreements
func (c *PayPalClient) CreateBillingAgreement(ctx context.Context, a BillingAgreement) (*CreateAgreementResponse, error) {
	// PayPal needs only ID, so we will remove all fields except Plan ID
	a.Plan = BillingPlan{
		ID: a.Plan.ID,
	}

	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s%s", c.APIBase, "/v1/payments/billing-agreements"), a)
	response := &CreateAgreementResponse{}
	if err != nil {
		return response, err
	}
	err = c.SendWithAuth(req, response)
	return response, err
}

// ExecuteApprovedAgreement - Use this call to execute (complete) a PayPal agreement that has been approved by the payer.
// Endpoint: POST /v1/payments/billing-agreements/token/agreement-execute
func (c *PayPalClient) ExecuteApprovedAgreement(ctx context.Context, token string) (*ExecuteAgreementResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/v1/payments/billing-agreements/%s/agreement-execute", c.APIBase, token), nil)
	response := &ExecuteAgreementResponse{}

	if err != nil {
		return response, err
	}

	req.SetBasicAuth(c.ClientID, c.Secret)
	req.Header.Set("Authorization", "Bearer "+c.Token.Token)

	if err = c.SendWithAuth(req, response); err != nil {
		return response, err
	}

	if response.ID == "" {
		return response, errors.New("Unable to execute agreement with token=" + token)
	}

	return response, err
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
func (c *PayPalClient) CaptureAuthorizationWithPaypalRequestId(ctx context.Context, authID string, paymentCaptureRequest *PaymentCaptureRequest, requestID string) (*PaymentCaptureResponse, error) {
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

// GetUserInfo for retrieve user profile attributes.
// Pass the schema that is used to return as per openidconnect protocol. The only supported schema value is openid.
// Endpoint: GET /v1/identity/openidconnect/userinfo/?schema=<Schema>
func (c *PayPalClient) GetUserInfo(ctx context.Context, schema string) (*UserInfo, error) {
	u := &UserInfo{}

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s%s%s", c.APIBase, "/v1/identity/openidconnect/userinfo/?schema=", schema), nil)
	if err != nil {
		return u, err
	}

	if err = c.SendWithAuth(req, u); err != nil {
		return u, err
	}

	return u, nil
}

// GrantNewAccessTokenFromAuthCode - Use this call to grant a new access token, using the previously obtained authorization code.
// Endpoint: POST /v1/identity/openidconnect/tokenservice
func (c *PayPalClient) GrantNewAccessTokenFromAuthCode(ctx context.Context, code, redirectURI string) (*TokenResponse, error) {
	token := &TokenResponse{}

	q := url.Values{}
	q.Set("grant_type", "authorization_code")
	q.Set("code", code)
	q.Set("redirect_uri", redirectURI)

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/identity/openidconnect/tokenservice"), strings.NewReader(q.Encode()))
	if err != nil {
		return token, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if err = c.SendWithBasicAuth(req, token); err != nil {
		return token, err
	}

	return token, nil
}

// GrantNewAccessTokenFromRefreshToken - Use this call to grant a new access token, using a refresh token.
// Endpoint: POST /v1/identity/openidconnect/tokenservice
func (c *PayPalClient) GrantNewAccessTokenFromRefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	type request struct {
		GrantType    string `json:"grant_type"`
		RefreshToken string `json:"refresh_token"`
	}

	token := &TokenResponse{}

	req, err := c.NewRequest(ctx, "POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/identity/openidconnect/tokenservice"), request{GrantType: "refresh_token", RefreshToken: refreshToken})
	if err != nil {
		return token, err
	}

	if err = c.SendWithAuth(req, token); err != nil {
		return token, err
	}

	return token, nil
}

// CreateWebProfile creates a new web experience profile in Paypal.
// Allows for the customisation of the payment experience.
// Endpoint: POST /v1/payment-experience/web-profiles
func (c *PayPalClient) CreateWebProfile(ctx context.Context, wp WebProfile) (*WebProfile, error) {
	url := fmt.Sprintf("%s%s", c.APIBase, "/v1/payment-experience/web-profiles")
	req, err := c.NewRequest(ctx, "POST", url, wp)
	response := &WebProfile{}

	if err != nil {
		return response, err
	}

	if err = c.SendWithAuth(req, response); err != nil {
		return response, err
	}

	return response, nil
}

// GetWebProfile gets an exists payment experience from Paypal.
// Endpoint: GET /v1/payment-experience/web-profiles/<profile-id>
func (c *PayPalClient) GetWebProfile(ctx context.Context, profileID string) (*WebProfile, error) {
	var wp WebProfile

	url := fmt.Sprintf("%s%s%s", c.APIBase, "/v1/payment-experience/web-profiles/", profileID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	if err != nil {
		return &wp, err
	}

	if err = c.SendWithAuth(req, &wp); err != nil {
		return &wp, err
	}

	if wp.ID == "" {
		return &wp, fmt.Errorf("paypal: unable to get web profile with ID = %s", profileID)
	}

	return &wp, nil
}

// GetWebProfiles retrieves web experience profiles from Paypal.
// Endpoint: GET /v1/payment-experience/web-profiles
func (c *PayPalClient) GetWebProfiles(ctx context.Context) ([]WebProfile, error) {
	var wps []WebProfile

	url := fmt.Sprintf("%s%s", c.APIBase, "/v1/payment-experience/web-profiles")
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	if err != nil {
		return wps, err
	}

	if err = c.SendWithAuth(req, &wps); err != nil {
		return wps, err
	}

	return wps, nil
}

// SetWebProfile sets a web experience profile in Paypal with given id.
// Endpoint: PUT /v1/payment-experience/web-profiles
func (c *PayPalClient) SetWebProfile(ctx context.Context, wp WebProfile) error {

	if wp.ID == "" {
		return fmt.Errorf("paypal: no ID specified for WebProfile")
	}

	url := fmt.Sprintf("%s%s%s", c.APIBase, "/v1/payment-experience/web-profiles/", wp.ID)

	req, err := c.NewRequest(ctx, "PUT", url, wp)

	if err != nil {
		return err
	}

	if err = c.SendWithAuth(req, nil); err != nil {
		return err
	}

	return nil
}

// DeleteWebProfile deletes a web experience profile from Paypal with given id.
// Endpoint: DELETE /v1/payment-experience/web-profiles
func (c *PayPalClient) DeleteWebProfile(ctx context.Context, profileID string) error {

	url := fmt.Sprintf("%s%s%s", c.APIBase, "/v1/payment-experience/web-profiles/", profileID)

	req, err := c.NewRequest(ctx, "DELETE", url, nil)

	if err != nil {
		return err
	}

	if err = c.SendWithAuth(req, nil); err != nil {
		return err
	}

	return nil
}

// ListTransactions for search transactions from the last 31 days.
// Endpoint: GET /v1/reporting/transactions
func (c *PayPalClient) ListTransactions(ctx context.Context, req *TransactionSearchRequest) (*TransactionSearchResponse, error) {
	response := &TransactionSearchResponse{}

	r, err := c.NewRequest(ctx, "GET", fmt.Sprintf("%s%s", c.APIBase, "/v1/reporting/transactions"), nil)
	if err != nil {
		return nil, err
	}

	q := r.URL.Query()

	q.Add("start_date", req.StartDate.Format(time.RFC3339))
	q.Add("end_date", req.EndDate.Format(time.RFC3339))

	if req.TransactionID != nil {
		q.Add("transaction_id", *req.TransactionID)
	}
	if req.TransactionType != nil {
		q.Add("transaction_type", *req.TransactionType)
	}
	if req.TransactionStatus != nil {
		q.Add("transaction_status", *req.TransactionStatus)
	}
	if req.TransactionAmount != nil {
		q.Add("transaction_amount", *req.TransactionAmount)
	}
	if req.TransactionCurrency != nil {
		q.Add("transaction_currency", *req.TransactionCurrency)
	}
	if req.PaymentInstrumentType != nil {
		q.Add("payment_instrument_type", *req.PaymentInstrumentType)
	}
	if req.StoreID != nil {
		q.Add("store_id", *req.StoreID)
	}
	if req.TerminalID != nil {
		q.Add("terminal_id", *req.TerminalID)
	}
	if req.Fields != nil {
		q.Add("fields", *req.Fields)
	}
	if req.BalanceAffectingRecordsOnly != nil {
		q.Add("balance_affecting_records_only", *req.BalanceAffectingRecordsOnly)
	}
	if req.PageSize != nil {
		q.Add("page_size", strconv.Itoa(*req.PageSize))
	}
	if req.Page != nil {
		q.Add("page", strconv.Itoa(*req.Page))
	}

	r.URL.RawQuery = q.Encode()

	if err = c.SendWithAuth(r, response); err != nil {
		return nil, err
	}

	return response, nil
}

// StoreCreditCard function.
// Endpoint: POST /v1/vault/credit-cards
func (c *PayPalClient) StoreCreditCard(ctx context.Context, cc CreditCard) (*CreditCard, error) {
	req, err := c.NewRequest(ctx, "POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/vault/credit-cards"), cc)
	if err != nil {
		return nil, err
	}

	response := &CreditCard{}

	if err = c.SendWithAuth(req, response); err != nil {
		return nil, err
	}

	return response, nil
}

// DeleteCreditCard function.
// Endpoint: DELETE /v1/vault/credit-cards/credit_card_id
func (c *PayPalClient) DeleteCreditCard(ctx context.Context, id string) error {
	req, err := c.NewRequest(ctx, "DELETE", fmt.Sprintf("%s/v1/vault/credit-cards/%s", c.APIBase, id), nil)
	if err != nil {
		return err
	}

	if err = c.SendWithAuth(req, nil); err != nil {
		return err
	}

	return nil
}

// GetCreditCard function.
// Endpoint: GET /v1/vault/credit-cards/credit_card_id
func (c *PayPalClient) GetCreditCard(ctx context.Context, id string) (*CreditCard, error) {
	req, err := c.NewRequest(ctx, "GET", fmt.Sprintf("%s/v1/vault/credit-cards/%s", c.APIBase, id), nil)
	if err != nil {
		return nil, err
	}

	response := &CreditCard{}

	if err = c.SendWithAuth(req, response); err != nil {
		return nil, err
	}

	return response, nil
}

// GetCreditCards function.
// Endpoint: GET /v1/vault/credit-cards
func (c *PayPalClient) GetCreditCards(ctx context.Context, ccf *CreditCardsFilter) (*CreditCards, error) {
	page := 1
	if ccf != nil && ccf.Page > 0 {
		page = ccf.Page
	}
	pageSize := 10
	if ccf != nil && ccf.PageSize > 0 {
		pageSize = ccf.PageSize
	}

	req, err := c.NewRequest(ctx, "GET", fmt.Sprintf("%s/v1/vault/credit-cards?page=%d&page_size=%d", c.APIBase, page, pageSize), nil)
	if err != nil {
		return nil, err
	}

	response := &CreditCards{}

	if err = c.SendWithAuth(req, response); err != nil {
		return nil, err
	}

	return response, nil
}

// PatchCreditCard function.
// Endpoint: PATCH /v1/vault/credit-cards/credit_card_id
func (c *PayPalClient) PatchCreditCard(ctx context.Context, id string, ccf []CreditCardField) (*CreditCard, error) {
	req, err := c.NewRequest(ctx, "PATCH", fmt.Sprintf("%s/v1/vault/credit-cards/%s", c.APIBase, id), ccf)
	if err != nil {
		return nil, err
	}

	response := &CreditCard{}

	if err = c.SendWithAuth(req, response); err != nil {
		return nil, err
	}

	return response, nil
}

// GetOrder retrieves order by ID
// Endpoint: GET /v2/checkout/orders/ID
func (c *PayPalClient) GetOrder(ctx context.Context, orderID string) (*Order, error) {
	order := &Order{}

	req, err := c.NewRequest(ctx, "GET", fmt.Sprintf("%s%s%s", c.APIBase, "/v2/checkout/orders/", orderID), nil)
	if err != nil {
		return order, err
	}

	if err = c.SendWithAuth(req, order); err != nil {
		return order, err
	}

	return order, nil
}

// CreateOrder - Use this call to create an order
// Endpoint: POST /v2/checkout/orders
func (c *PayPalClient) CreateOrder(ctx context.Context, intent string, purchaseUnits []PurchaseUnitRequest, payer *CreateOrderPayer, appContext *ApplicationContext) (*Order, error) {
	return c.CreateOrderWithPaypalRequestID(ctx, intent, purchaseUnits, payer, appContext, "")
}

// CreateOrderWithPaypalRequestID - Use this call to create an order with idempotency
// Endpoint: POST /v2/checkout/orders
func (c *PayPalClient) CreateOrderWithPaypalRequestID(ctx context.Context, intent string, purchaseUnits []PurchaseUnitRequest, payer *CreateOrderPayer, appContext *ApplicationContext, requestID string) (*Order, error) {
	type createOrderRequest struct {
		Intent             string                `json:"intent"`
		Payer              *CreateOrderPayer     `json:"payer,omitempty"`
		PurchaseUnits      []PurchaseUnitRequest `json:"purchase_units"`
		ApplicationContext *ApplicationContext   `json:"application_context,omitempty"`
	}

	order := &Order{}

	req, err := c.NewRequest(ctx, "POST", fmt.Sprintf("%s%s", c.APIBase, "/v2/checkout/orders"), createOrderRequest{Intent: intent, PurchaseUnits: purchaseUnits, Payer: payer, ApplicationContext: appContext})
	if err != nil {
		return order, err
	}

	if requestID != "" {
		req.Header.Set("PayPal-Request-Id", requestID)
	}

	if err = c.SendWithAuth(req, order); err != nil {
		return order, err
	}

	return order, nil
}

// UpdateOrder updates the order by ID
// Endpoint: PATCH /v2/checkout/orders/ID
func (c *PayPalClient) UpdateOrder(ctx context.Context, orderID string, purchaseUnits []PurchaseUnitRequest) (*Order, error) {
	order := &Order{}

	req, err := c.NewRequest(ctx, "PATCH", fmt.Sprintf("%s%s%s", c.APIBase, "/v2/checkout/orders/", orderID), purchaseUnits)
	if err != nil {
		return order, err
	}

	if err = c.SendWithAuth(req, order); err != nil {
		return order, err
	}

	return order, nil
}

// AuthorizeOrder - https://developer.paypal.com/docs/api/orders/v2/#orders_authorize
// Endpoint: POST /v2/checkout/orders/ID/authorize
func (c *PayPalClient) AuthorizeOrder(ctx context.Context, orderID string, authorizeOrderRequest AuthorizeOrderRequest) (*Authorization, error) {
	auth := &Authorization{}

	req, err := c.NewRequest(ctx, "POST", fmt.Sprintf("%s%s", c.APIBase, "/v2/checkout/orders/"+orderID+"/authorize"), authorizeOrderRequest)
	if err != nil {
		return auth, err
	}

	if err = c.SendWithAuth(req, auth); err != nil {
		return auth, err
	}

	return auth, nil
}

// CaptureOrder - https://developer.paypal.com/docs/api/orders/v2/#orders_capture
// Endpoint: POST /v2/checkout/orders/ID/capture
func (c *PayPalClient) CaptureOrder(ctx context.Context, orderID string, captureOrderRequest CaptureOrderRequest) (*CaptureOrderResponse, error) {
	return c.CaptureOrderWithPaypalRequestId(ctx, orderID, captureOrderRequest, "")
}

// CaptureOrder with idempotency - https://developer.paypal.com/docs/api/orders/v2/#orders_capture
// Endpoint: POST /v2/checkout/orders/ID/capture
// https://developer.paypal.com/docs/api/reference/api-requests/#http-request-headers
func (c *PayPalClient) CaptureOrderWithPaypalRequestId(ctx context.Context, orderID string, captureOrderRequest CaptureOrderRequest, requestID string) (*CaptureOrderResponse, error) {
	capture := &CaptureOrderResponse{}

	c.SetReturnRepresentation()
	req, err := c.NewRequest(ctx, "POST", fmt.Sprintf("%s%s", c.APIBase, "/v2/checkout/orders/"+orderID+"/capture"), captureOrderRequest)
	if err != nil {
		return capture, err
	}

	if requestID != "" {
		req.Header.Set("PayPal-Request-Id", requestID)
	}

	if err = c.SendWithAuth(req, capture); err != nil {
		return capture, err
	}

	return capture, nil
}

// CreateWebhook - Subscribes your webhook listener to events.
// Endpoint: POST /v1/notifications/webhooks
func (c *PayPalClient) CreateWebhook(ctx context.Context, createWebhookRequest *CreateWebhookRequest) (*Webhook, error) {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s%s", c.APIBase, "/v1/notifications/webhooks"), createWebhookRequest)
	webhook := &Webhook{}
	if err != nil {
		return webhook, err
	}

	err = c.SendWithAuth(req, webhook)
	return webhook, err
}

// GetWebhook - Shows details for a webhook, by ID.
// Endpoint: GET /v1/notifications/webhooks/ID
func (c *PayPalClient) GetWebhook(ctx context.Context, webhookID string) (*Webhook, error) {
	req, err := c.NewRequest(ctx, http.MethodGet, fmt.Sprintf("%s%s%s", c.APIBase, "/v1/notifications/webhooks/", webhookID), nil)
	webhook := &Webhook{}
	if err != nil {
		return webhook, err
	}

	err = c.SendWithAuth(req, webhook)
	return webhook, err
}

// UpdateWebhook - Updates a webhook to replace webhook fields with new values.
// Endpoint: PATCH /v1/notifications/webhooks/ID
func (c *PayPalClient) UpdateWebhook(ctx context.Context, webhookID string, fields []WebhookField) (*Webhook, error) {
	req, err := c.NewRequest(ctx, http.MethodPatch, fmt.Sprintf("%s/v1/notifications/webhooks/%s", c.APIBase, webhookID), fields)
	webhook := &Webhook{}
	if err != nil {
		return webhook, err
	}

	err = c.SendWithAuth(req, webhook)
	return webhook, err
}

// ListWebhooks - Lists webhooks for an app.
// Endpoint: GET /v1/notifications/webhooks
func (c *PayPalClient) ListWebhooks(ctx context.Context, anchorType string) (*ListWebhookResponse, error) {
	if len(anchorType) == 0 {
		anchorType = AncorTypeApplication
	}
	req, err := c.NewRequest(ctx, http.MethodGet, fmt.Sprintf("%s%s", c.APIBase, "/v1/notifications/webhooks"), nil)
	q := req.URL.Query()
	q.Add("anchor_type", anchorType)
	req.URL.RawQuery = q.Encode()
	resp := &ListWebhookResponse{}
	if err != nil {
		return nil, err
	}

	err = c.SendWithAuth(req, resp)
	return resp, err
}

// DeleteWebhook - Deletes a webhook, by ID.
// Endpoint: DELETE /v1/notifications/webhooks/ID
func (c *PayPalClient) DeleteWebhook(ctx context.Context, webhookID string) error {
	req, err := c.NewRequest(ctx, http.MethodDelete, fmt.Sprintf("%s/v1/notifications/webhooks/%s", c.APIBase, webhookID), nil)
	if err != nil {
		return err
	}

	err = c.SendWithAuth(req, nil)
	return err
}

// VerifyWebhookSignature - Use this to verify the signature of a webhook recieved from paypal.
// Endpoint: POST /v1/notifications/verify-webhook-signature
func (c *PayPalClient) VerifyWebhookSignature(ctx context.Context, httpReq *http.Request, webhookID string) (*VerifyWebhookResponse, error) {
	type verifyWebhookSignatureRequest struct {
		AuthAlgo         string          `json:"auth_algo,omitempty"`
		CertURL          string          `json:"cert_url,omitempty"`
		TransmissionID   string          `json:"transmission_id,omitempty"`
		TransmissionSig  string          `json:"transmission_sig,omitempty"`
		TransmissionTime string          `json:"transmission_time,omitempty"`
		WebhookID        string          `json:"webhook_id,omitempty"`
		Event            json.RawMessage `json:"webhook_event"`
	}

	// Read the content
	var bodyBytes []byte
	if httpReq.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(httpReq.Body)
	}
	// Restore the io.ReadCloser to its original state
	httpReq.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	verifyRequest := verifyWebhookSignatureRequest{
		AuthAlgo:         httpReq.Header.Get("PAYPAL-AUTH-ALGO"),
		CertURL:          httpReq.Header.Get("PAYPAL-CERT-URL"),
		TransmissionID:   httpReq.Header.Get("PAYPAL-TRANSMISSION-ID"),
		TransmissionSig:  httpReq.Header.Get("PAYPAL-TRANSMISSION-SIG"),
		TransmissionTime: httpReq.Header.Get("PAYPAL-TRANSMISSION-TIME"),
		WebhookID:        webhookID,
		Event:            json.RawMessage(bodyBytes),
	}

	response := &VerifyWebhookResponse{}

	req, err := c.NewRequest(ctx, "POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/notifications/verify-webhook-signature"), verifyRequest)
	if err != nil {
		return nil, err
	}

	if err = c.SendWithAuth(req, response); err != nil {
		return nil, err
	}

	return response, nil
}

// GetWebhookEventTypes - Lists all webhook event types.
// Endpoint: GET /v1/notifications/webhooks-event-types
func (c *PayPalClient) GetWebhookEventTypes(ctx context.Context) (*WebhookEventTypesResponse, error) {
	req, err := c.NewRequest(ctx, http.MethodGet, fmt.Sprintf("%s%s", c.APIBase, "/v1/notifications/webhooks-event-types"), nil)
	q := req.URL.Query()

	req.URL.RawQuery = q.Encode()
	resp := &WebhookEventTypesResponse{}
	if err != nil {
		return nil, err
	}

	err = c.SendWithAuth(req, resp)
	return resp, err
}

// CreateProduct creates a product
// Doc: https://developer.paypal.com/docs/api/catalog-products/v1/#products_create
// Endpoint: POST /v1/catalogs/products
func (c *PayPalClient) CreateProduct(ctx context.Context, product Product) (*CreateProductResponse, error) {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s%s", c.APIBase, "/v1/catalogs/products"), product)
	response := &CreateProductResponse{}
	if err != nil {
		return response, err
	}
	err = c.SendWithAuth(req, response)
	return response, err
}

// UpdateProduct. updates a product information
// Doc: https://developer.paypal.com/docs/api/catalog-products/v1/#products_patch
// Endpoint: PATCH /v1/catalogs/products/:product_id
func (c *PayPalClient) UpdateProduct(ctx context.Context, product Product) error {
	req, err := c.NewRequest(ctx, http.MethodPatch, fmt.Sprintf("%s%s%s", c.APIBase, "/v1/catalogs/products/", product.ID), product.GetUpdatePatch())
	if err != nil {
		return err
	}
	err = c.SendWithAuth(req, nil)
	return err
}

// Get product details
// Doc: https://developer.paypal.com/docs/api/catalog-products/v1/#products_get
// Endpoint: GET /v1/catalogs/products/:product_id
func (c *PayPalClient) GetProduct(ctx context.Context, productId string) (*Product, error) {
	req, err := c.NewRequest(ctx, http.MethodGet, fmt.Sprintf("%s%s%s", c.APIBase, "/v1/catalogs/products/", productId), nil)
	response := &Product{}
	if err != nil {
		return response, err
	}
	err = c.SendWithAuth(req, response)
	return response, err
}

// List all products
// Doc: https://developer.paypal.com/docs/api/catalog-products/v1/#products_list
// Endpoint: GET /v1/catalogs/products
func (c *PayPalClient) ListProducts(ctx context.Context, params *ProductListParameters) (*ListProductsResponse, error) {
	req, err := c.NewRequest(ctx, http.MethodGet, fmt.Sprintf("%s%s", c.APIBase, "/v1/catalogs/products"), nil)
	response := &ListProductsResponse{}
	if err != nil {
		return response, err
	}

	if params != nil {
		q := req.URL.Query()
		q.Add("page", params.Page)
		q.Add("page_size", params.PageSize)
		q.Add("total_required", params.TotalRequired)
		req.URL.RawQuery = q.Encode()
	}

	err = c.SendWithAuth(req, response)
	return response, err
}

// CreateSubscriptionPlan creates a subscriptionPlan
// Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#plans_create
// Endpoint: POST /v1/billing/plans
func (c *PayPalClient) CreateSubscriptionPlan(ctx context.Context, newPlan SubscriptionPlan) (*CreateSubscriptionPlanResponse, error) {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s%s", c.APIBase, "/v1/billing/plans"), newPlan)
	response := &CreateSubscriptionPlanResponse{}
	if err != nil {
		return response, err
	}
	err = c.SendWithAuth(req, response)
	return response, err
}

// UpdateSubscriptionPlan. updates a plan
// Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#plans_patch
// Endpoint: PATCH /v1/billing/plans/:plan_id
func (c *PayPalClient) UpdateSubscriptionPlan(ctx context.Context, updatedPlan SubscriptionPlan) error {
	req, err := c.NewRequest(ctx, http.MethodPatch, fmt.Sprintf("%s%s%s", c.APIBase, "/v1/billing/plans/", updatedPlan.ID), updatedPlan.GetUpdatePatch())
	if err != nil {
		return err
	}
	err = c.SendWithAuth(req, nil)
	return err
}

// UpdateSubscriptionPlan. updates a plan
// Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#plans_get
// Endpoint: GET /v1/billing/plans/:plan_id
func (c *PayPalClient) GetSubscriptionPlan(ctx context.Context, planId string) (*SubscriptionPlan, error) {
	req, err := c.NewRequest(ctx, http.MethodGet, fmt.Sprintf("%s%s%s", c.APIBase, "/v1/billing/plans/", planId), nil)
	response := &SubscriptionPlan{}
	if err != nil {
		return response, err
	}
	err = c.SendWithAuth(req, response)
	return response, err
}

// List all plans
// Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#plans_list
// Endpoint: GET /v1/billing/plans
func (c *PayPalClient) ListSubscriptionPlans(ctx context.Context, params *SubscriptionPlanListParameters) (*ListSubscriptionPlansResponse, error) {
	req, err := c.NewRequest(ctx, http.MethodGet, fmt.Sprintf("%s%s", c.APIBase, "/v1/billing/plans"), nil)
	response := &ListSubscriptionPlansResponse{}
	if err != nil {
		return response, err
	}

	if params != nil {
		q := req.URL.Query()
		q.Add("page", params.Page)
		q.Add("page_size", params.PageSize)
		q.Add("total_required", params.TotalRequired)
		q.Add("product_id", params.ProductId)
		q.Add("plan_ids", params.PlanIds)
		req.URL.RawQuery = q.Encode()
	}

	err = c.SendWithAuth(req, response)
	return response, err
}

// Activates a plan
// Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#plans_activate
// Endpoint: POST /v1/billing/plans/{id}/activate
func (c *PayPalClient) ActivateSubscriptionPlan(ctx context.Context, planId string) error {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s/v1/billing/plans/%s/activate", c.APIBase, planId), nil)
	if err != nil {
		return err
	}

	err = c.SendWithAuth(req, nil)
	return err
}

// Deactivates a plan
// Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#plans_deactivate
// Endpoint: POST /v1/billing/plans/{id}/deactivate
func (c *PayPalClient) DeactivateSubscriptionPlans(ctx context.Context, planId string) error {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s/v1/billing/plans/%s/deactivate", c.APIBase, planId), nil)
	if err != nil {
		return err
	}

	err = c.SendWithAuth(req, nil)
	return err
}

// Updates pricing for a plan. For example, you can update a regular billing cycle from $5 per month to $7 per month.
// Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#plans_update-pricing-schemes
// Endpoint: POST /v1/billing/plans/{id}/update-pricing-schemes
func (c *PayPalClient) UpdateSubscriptionPlanPricing(ctx context.Context, planId string, pricingSchemes []PricingSchemeUpdate) error {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s/v1/billing/plans/%s/update-pricing-schemes", c.APIBase, planId), PricingSchemeUpdateRequest{
		Schemes: pricingSchemes,
	})
	if err != nil {
		return err
	}

	err = c.SendWithAuth(req, nil)
	return err
}

// CreateSubscriptionPlan creates a subscriptionPlan
// Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_create
// Endpoint: POST /v1/billing/subscriptions
func (c *PayPalClient) CreateSubscription(ctx context.Context, newSubscription SubscriptionBase) (*SubscriptionDetailResp, error) {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s%s", c.APIBase, "/v1/billing/subscriptions"), newSubscription)
	req.Header.Add("Prefer", "return=representation")
	response := &SubscriptionDetailResp{}
	if err != nil {
		return response, err
	}
	err = c.SendWithAuth(req, response)
	return response, err
}

// UpdateSubscriptionPlan. updates a plan
// Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_patch
// Endpoint: PATCH /v1/billing/subscriptions/:subscription_id
func (c *PayPalClient) UpdateSubscription(ctx context.Context, updatedSubscription Subscription) error {
	req, err := c.NewRequest(ctx, http.MethodPatch, fmt.Sprintf("%s%s%s", c.APIBase, "/v1/billing/subscriptions/", updatedSubscription.ID), updatedSubscription.GetUpdatePatch())
	if err != nil {
		return err
	}
	err = c.SendWithAuth(req, nil)
	return err
}

// GetSubscriptionDetails shows details for a subscription, by ID.
// Endpoint: GET /v1/billing/subscriptions/
func (c *PayPalClient) GetSubscriptionDetails(ctx context.Context, subscriptionID string) (*SubscriptionDetailResp, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/v1/billing/subscriptions/%s", c.APIBase, subscriptionID), nil)
	response := &SubscriptionDetailResp{}
	if err != nil {
		return response, err
	}
	err = c.SendWithAuth(req, response)
	return response, err
}

// Activates the subscription.
// Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_activate
// Endpoint: POST /v1/billing/subscriptions/{id}/activate
func (c *PayPalClient) ActivateSubscription(ctx context.Context, subscriptionId, activateReason string) error {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s/v1/billing/subscriptions/%s/activate", c.APIBase, subscriptionId), map[string]string{"reason": activateReason})
	if err != nil {
		return err
	}
	err = c.SendWithAuth(req, nil)
	return err
}

// Cancels the subscription.
// Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_cancel
// Endpoint: POST /v1/billing/subscriptions/{id}/cancel
func (c *PayPalClient) CancelSubscription(ctx context.Context, subscriptionId, cancelReason string) error {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s/v1/billing/subscriptions/%s/cancel", c.APIBase, subscriptionId), map[string]string{"reason": cancelReason})
	if err != nil {
		return err
	}
	err = c.SendWithAuth(req, nil)
	return err
}

// Captures an authorized payment from the subscriber on the subscription.
// Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_capture
// Endpoint: POST /v1/billing/subscriptions/{id}/capture
func (c *PayPalClient) CaptureSubscription(ctx context.Context, subscriptionId string, request CaptureReqeust) (*SubscriptionCaptureResponse, error) {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s/v1/billing/subscriptions/%s/capture", c.APIBase, subscriptionId), request)
	response := &SubscriptionCaptureResponse{}
	if err != nil {
		return response, err
	}
	err = c.SendWithAuth(req, response)
	return response, err
}

// Suspends the subscription.
// Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_suspend
// Endpoint: POST /v1/billing/subscriptions/{id}/suspend
func (c *PayPalClient) SuspendSubscription(ctx context.Context, subscriptionId, reason string) error {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s/v1/billing/subscriptions/%s/suspend", c.APIBase, subscriptionId), map[string]string{"reason": reason})
	if err != nil {
		return err
	}
	err = c.SendWithAuth(req, nil)
	return err
}

// Lists transactions for a subscription.
// Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_transactions
// Endpoint: GET /v1/billing/subscriptions/{id}/transactions
func (c *PayPalClient) GetSubscriptionTransactions(ctx context.Context, requestParams SubscriptionTransactionsParams) (*SubscriptionTransactionsResponse, error) {
	startTime := requestParams.StartTime.Format("2006-01-02T15:04:05Z")
	endTime := requestParams.EndTime.Format("2006-01-02T15:04:05Z")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/v1/billing/subscriptions/%s/transactions?start_time=%s&end_time=%s", c.APIBase, requestParams.SubscriptionId, startTime, endTime), nil)
	response := &SubscriptionTransactionsResponse{}
	if err != nil {
		return response, err
	}

	err = c.SendWithAuth(req, response)
	return response, err
}

// Revise plan or quantity of subscription
// Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#subscriptions_revise
// Endpoint: POST /v1/billing/subscriptions/{id}/revise
func (c *PayPalClient) ReviseSubscription(ctx context.Context, subscriptionId string, reviseSubscription SubscriptionBase) (*SubscriptionDetailResp, error) {
	req, err := c.NewRequest(ctx, http.MethodPost, fmt.Sprintf("%s/v1/billing/subscriptions/%s/revise", c.APIBase, subscriptionId), reviseSubscription)
	response := &SubscriptionDetailResp{}
	if err != nil {
		return response, err
	}

	req.Header.Add("Content-Type", "application/json")
	err = c.SendWithAuth(req, response)

	return response, err
}
