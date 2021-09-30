package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
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
