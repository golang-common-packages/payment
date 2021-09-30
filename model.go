package payment

import (
	"net/http"
	"time"
)

// Begin Payment Connection Models //

// Config model
type Config struct {
	PayPal PayPal `json:"paypal,omitempty"`
}

// Paypal model for Paypal connection config
type PayPal struct {
	ClientID string `json:"clientID"`
	SecretID string `json:"secretID"`
	APIBase  string `json:"apiBase"`
}

// End Payment Connection Models //

// -------------------------------------------------------------------------

// Begin PayPal Models //

// BillingPlanStatus has type is string. This type may change in the future
type BillingPlanStatus string

// JSONTime overrides MarshalJson method to format in ISO8601
type JSONTime time.Time

// TokenResponse is for API response for the /oauth2/token endpoint
type TokenResponse struct {
	RefreshToken string `json:"refresh_token"`
	Token        string `json:"access_token"`
	Type         string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

// ErrorResponse struct
// https://developer.paypal.com/docs/api/errors/
type ErrorResponse struct {
	Response        *http.Response        `json:"-"`
	Name            string                `json:"name"`
	DebugID         string                `json:"debug_id"`
	Message         string                `json:"message"`
	InformationLink string                `json:"information_link"`
	Details         []ErrorResponseDetail `json:"details"`
}

// ErrorResponseDetail struct
type ErrorResponseDetail struct {
	Field string `json:"field"`
	Issue string `json:"issue"`
}

// Payout struct
type Payout struct {
	SenderBatchHeader *SenderBatchHeader `json:"sender_batch_header"`
	Items             []PayoutItem       `json:"items"`
}

// SenderBatchHeader struct
type SenderBatchHeader struct {
	EmailSubject  string `json:"email_subject"`
	EmailMessage  string `json:"email_message"`
	SenderBatchID string `json:"sender_batch_id,omitempty"`
}

// PayoutItem struct
type PayoutItem struct {
	RecipientType   string        `json:"recipient_type"`
	RecipientWallet string        `json:"recipient_wallet"`
	Receiver        string        `json:"receiver"`
	Amount          *AmountPayout `json:"amount"`
	Note            string        `json:"note,omitempty"`
	SenderItemID    string        `json:"sender_item_id,omitempty"`
}

// AmountPayout struct
type AmountPayout struct {
	Currency string `json:"currency"`
	Value    string `json:"value"`
}

// PayoutResponse struct
type PayoutResponse struct {
	BatchHeader *BatchHeader         `json:"batch_header"`
	Items       []PayoutItemResponse `json:"items"`
	Links       []Link               `json:"links"`
}

// BatchHeader struct
type BatchHeader struct {
	Amount            *AmountPayout      `json:"amount,omitempty"`
	Fees              *AmountPayout      `json:"fees,omitempty"`
	PayoutBatchID     string             `json:"payout_batch_id,omitempty"`
	BatchStatus       string             `json:"batch_status,omitempty"`
	TimeCreated       *time.Time         `json:"time_created,omitempty"`
	TimeCompleted     *time.Time         `json:"time_completed,omitempty"`
	SenderBatchHeader *SenderBatchHeader `json:"sender_batch_header,omitempty"`
}

// PayoutItemResponse struct
type PayoutItemResponse struct {
	PayoutItemID      string        `json:"payout_item_id"`
	TransactionID     string        `json:"transaction_id"`
	TransactionStatus string        `json:"transaction_status"`
	PayoutBatchID     string        `json:"payout_batch_id,omitempty"`
	PayoutItemFee     *AmountPayout `json:"payout_item_fee,omitempty"`
	PayoutItem        *PayoutItem   `json:"payout_item"`
	TimeProcessed     *time.Time    `json:"time_processed,omitempty"`
	Links             []Link        `json:"links"`
	Error             ErrorResponse `json:"errors,omitempty"`
}

// Link struct
type Link struct {
	Href        string `json:"href"`
	Rel         string `json:"rel,omitempty"`
	Method      string `json:"method,omitempty"`
	Description string `json:"description,omitempty"`
	Enctype     string `json:"enctype,omitempty"`
}

// Sale struct
type Sale struct {
	ID                        string     `json:"id,omitempty"`
	Amount                    *Amount    `json:"amount,omitempty"`
	TransactionFee            *Currency  `json:"transaction_fee,omitempty"`
	Description               string     `json:"description,omitempty"`
	CreateTime                *time.Time `json:"create_time,omitempty"`
	State                     string     `json:"state,omitempty"`
	ParentPayment             string     `json:"parent_payment,omitempty"`
	UpdateTime                *time.Time `json:"update_time,omitempty"`
	PaymentMode               string     `json:"payment_mode,omitempty"`
	PendingReason             string     `json:"pending_reason,omitempty"`
	ReasonCode                string     `json:"reason_code,omitempty"`
	ClearingTime              string     `json:"clearing_time,omitempty"`
	ProtectionEligibility     string     `json:"protection_eligibility,omitempty"`
	ProtectionEligibilityType string     `json:"protection_eligibility_type,omitempty"`
	Links                     []Link     `json:"links,omitempty"`
}

// Amount struct
type Amount struct {
	Currency string  `json:"currency"`
	Total    string  `json:"total"`
	Details  Details `json:"details,omitempty"`
}

// Details structure used in Amount structures as optional value
type Details struct {
	Subtotal         string `json:"subtotal,omitempty"`
	Shipping         string `json:"shipping,omitempty"`
	Tax              string `json:"tax,omitempty"`
	HandlingFee      string `json:"handling_fee,omitempty"`
	ShippingDiscount string `json:"shipping_discount,omitempty"`
	Insurance        string `json:"insurance,omitempty"`
	GiftWrap         string `json:"gift_wrap,omitempty"`
}

// Currency struct
type Currency struct {
	Currency string `json:"currency,omitempty"`
	Value    string `json:"value,omitempty"`
}

// Refund struct
type Refund struct {
	ID            string     `json:"id,omitempty"`
	Amount        *Amount    `json:"amount,omitempty"`
	CreateTime    *time.Time `json:"create_time,omitempty"`
	State         string     `json:"state,omitempty"`
	CaptureID     string     `json:"capture_id,omitempty"`
	ParentPayment string     `json:"parent_payment,omitempty"`
	UpdateTime    *time.Time `json:"update_time,omitempty"`
}

// Authorization struct
type Authorization struct {
	ID               string                `json:"id,omitempty"`
	CustomID         string                `json:"custom_id,omitempty"`
	InvoiceID        string                `json:"invoice_id,omitempty"`
	Status           string                `json:"status,omitempty"`
	StatusDetails    *CaptureStatusDetails `json:"status_details,omitempty"`
	Amount           *PurchaseUnitAmount   `json:"amount,omitempty"`
	SellerProtection *SellerProtection     `json:"seller_protection,omitempty"`
	CreateTime       *time.Time            `json:"create_time,omitempty"`
	UpdateTime       *time.Time            `json:"update_time,omitempty"`
	ExpirationTime   *time.Time            `json:"expiration_time,omitempty"`
	Links            []Link                `json:"links,omitempty"`
}

// CaptureStatusDetails struct
// https://developer.paypal.com/docs/api/payments/v2/#definition-capture_status_details
type CaptureStatusDetails struct {
	Reason string `json:"reason,omitempty"`
}

// PurchaseUnitAmount struct
type PurchaseUnitAmount struct {
	Currency  string                       `json:"currency_code"`
	Value     string                       `json:"value"`
	Breakdown *PurchaseUnitAmountBreakdown `json:"breakdown,omitempty"`
}

// PurchaseUnitAmountBreakdown struct
type PurchaseUnitAmountBreakdown struct {
	ItemTotal        *Money `json:"item_total,omitempty"`
	Shipping         *Money `json:"shipping,omitempty"`
	Handling         *Money `json:"handling,omitempty"`
	TaxTotal         *Money `json:"tax_total,omitempty"`
	Insurance        *Money `json:"insurance,omitempty"`
	ShippingDiscount *Money `json:"shipping_discount,omitempty"`
	Discount         *Money `json:"discount,omitempty"`
}

// Money struct
// https://developer.paypal.com/docs/api/orders/v2/#definition-money
type Money struct {
	Currency string `json:"currency_code"`
	Value    string `json:"value"`
}

// SellerProtection struct
type SellerProtection struct {
	Status            string   `json:"status,omitempty"`
	DisputeCategories []string `json:"dispute_categories,omitempty"`
}

// PaymentCaptureRequest struct
// https://developer.paypal.com/docs/api/payments/v2/#authorizations_capture
type PaymentCaptureRequest struct {
	InvoiceID      string `json:"invoice_id,omitempty"`
	NoteToPayer    string `json:"note_to_payer,omitempty"`
	SoftDescriptor string `json:"soft_descriptor,omitempty"`
	Amount         *Money `json:"amount,omitempty"`
	FinalCapture   bool   `json:"final_capture,omitempty"`
}

// PaymentCaptureResponse struct
type PaymentCaptureResponse struct {
	Status           string                `json:"status,omitempty"`
	StatusDetails    *CaptureStatusDetails `json:"status_details,omitempty"`
	ID               string                `json:"id,omitempty"`
	Amount           *Money                `json:"amount,omitempty"`
	InvoiceID        string                `json:"invoice_id,omitempty"`
	FinalCapture     bool                  `json:"final_capture,omitempty"`
	DisbursementMode string                `json:"disbursement_mode,omitempty"`
	Links            []Link                `json:"links,omitempty"`
}

// Capture struct
type Capture struct {
	ID             string     `json:"id,omitempty"`
	Amount         *Amount    `json:"amount,omitempty"`
	State          string     `json:"state,omitempty"`
	ParentPayment  string     `json:"parent_payment,omitempty"`
	TransactionFee string     `json:"transaction_fee,omitempty"`
	IsFinalCapture bool       `json:"is_final_capture"`
	CreateTime     *time.Time `json:"create_time,omitempty"`
	UpdateTime     *time.Time `json:"update_time,omitempty"`
	Links          []Link     `json:"links,omitempty"`
}

// BillingPlanListParams struct
type BillingPlanListParams struct {
	ListParams
	Status string `json:"status,omitempty"` //Allowed values: CREATED, ACTIVE, INACTIVE, ALL.
}

// ListParams struct
type ListParams struct {
	Page          string `json:"page,omitempty"`           //Default: 0.
	PageSize      string `json:"page_size,omitempty"`      //Default: 10.
	TotalRequired string `json:"total_required,omitempty"` //Default: no.
}

// BillingPlanListResponse struct
type BillingPlanListResponse struct {
	SharedListResponse
	Plans []BillingPlan `json:"plans,omitempty"`
}

// SharedListResponse struct
type SharedListResponse struct {
	TotalItems int    `json:"total_items,omitempty"`
	TotalPages int    `json:"total_pages,omitempty"`
	Links      []Link `json:"links,omitempty"`
}

// BillingPlan struct
type BillingPlan struct {
	ID                  string               `json:"id,omitempty"`
	Name                string               `json:"name,omitempty"`
	Description         string               `json:"description,omitempty"`
	Type                string               `json:"type,omitempty"`
	PaymentDefinitions  []PaymentDefinition  `json:"payment_definitions,omitempty"`
	MerchantPreferences *MerchantPreferences `json:"merchant_preferences,omitempty"`
}

// PaymentDefinition struct
type PaymentDefinition struct {
	ID                string        `json:"id,omitempty"`
	Name              string        `json:"name,omitempty"`
	Type              string        `json:"type,omitempty"`
	Frequency         string        `json:"frequency,omitempty"`
	FrequencyInterval string        `json:"frequency_interval,omitempty"`
	Amount            AmountPayout  `json:"amount,omitempty"`
	Cycles            string        `json:"cycles,omitempty"`
	ChargeModels      []ChargeModel `json:"charge_models,omitempty"`
}

// MerchantPreferences struct
type MerchantPreferences struct {
	SetupFee                *AmountPayout `json:"setup_fee,omitempty"`
	ReturnURL               string        `json:"return_url,omitempty"`
	CancelURL               string        `json:"cancel_url,omitempty"`
	AutoBillAmount          string        `json:"auto_bill_amount,omitempty"`
	InitialFailAmountAction string        `json:"initial_fail_amount_action,omitempty"`
	MaxFailAttempts         string        `json:"max_fail_attempts,omitempty"`
}

// ChargeModel struct
type ChargeModel struct {
	Type   string       `json:"type,omitempty"`
	Amount AmountPayout `json:"amount,omitempty"`
}

// CreateBillingResponse struct
type CreateBillingResponse struct {
	ID                  string              `json:"id,omitempty"`
	State               string              `json:"state,omitempty"`
	PaymentDefinitions  []PaymentDefinition `json:"payment_definitions,omitempty"`
	MerchantPreferences MerchantPreferences `json:"merchant_preferences,omitempty"`
	CreateTime          time.Time           `json:"create_time,omitempty"`
	UpdateTime          time.Time           `json:"update_time,omitempty"`
	Links               []Link              `json:"links,omitempty"`
}

// Patch struct
type Patch struct {
	Operation string      `json:"op"`
	Path      string      `json:"path"`
	Value     interface{} `json:"value"`
}

// BillingAgreement struct
type BillingAgreement struct {
	Name                        string               `json:"name,omitempty"`
	Description                 string               `json:"description,omitempty"`
	StartDate                   JSONTime             `json:"start_date,omitempty"`
	Plan                        BillingPlan          `json:"plan,omitempty"`
	Payer                       Payer                `json:"payer,omitempty"`
	ShippingAddress             *ShippingAddress     `json:"shipping_address,omitempty"`
	OverrideMerchantPreferences *MerchantPreferences `json:"override_merchant_preferences,omitempty"`
}

// Payer struct
type Payer struct {
	PaymentMethod      string              `json:"payment_method"`
	FundingInstruments []FundingInstrument `json:"funding_instruments,omitempty"`
	PayerInfo          *PayerInfo          `json:"payer_info,omitempty"`
	Status             string              `json:"payer_status,omitempty"`
}

// FundingInstrument struct
type FundingInstrument struct {
	CreditCard      *CreditCard      `json:"credit_card,omitempty"`
	CreditCardToken *CreditCardToken `json:"credit_card_token,omitempty"`
}

// CreditCard struct
type CreditCard struct {
	ID                 string   `json:"id,omitempty"`
	PayerID            string   `json:"payer_id,omitempty"`
	ExternalCustomerID string   `json:"external_customer_id,omitempty"`
	Number             string   `json:"number"`
	Type               string   `json:"type"`
	ExpireMonth        string   `json:"expire_month"`
	ExpireYear         string   `json:"expire_year"`
	CVV2               string   `json:"cvv2,omitempty"`
	FirstName          string   `json:"first_name,omitempty"`
	LastName           string   `json:"last_name,omitempty"`
	BillingAddress     *Address `json:"billing_address,omitempty"`
	State              string   `json:"state,omitempty"`
	ValidUntil         string   `json:"valid_until,omitempty"`
}

// Address struct
type Address struct {
	Line1       string `json:"line1,omitempty"`
	Line2       string `json:"line2,omitempty"`
	City        string `json:"city,omitempty"`
	CountryCode string `json:"country_code,omitempty"`
	PostalCode  string `json:"postal_code,omitempty"`
	State       string `json:"state,omitempty"`
	Phone       string `json:"phone,omitempty"`
}

// CreditCardToken struct
type CreditCardToken struct {
	CreditCardID string `json:"credit_card_id"`
	PayerID      string `json:"payer_id,omitempty"`
	Last4        string `json:"last4,omitempty"`
	ExpireYear   string `json:"expire_year,omitempty"`
	ExpireMonth  string `json:"expire_month,omitempty"`
}

// PayerInfo struct
type PayerInfo struct {
	Email           string           `json:"email,omitempty"`
	FirstName       string           `json:"first_name,omitempty"`
	LastName        string           `json:"last_name,omitempty"`
	PayerID         string           `json:"payer_id,omitempty"`
	Phone           string           `json:"phone,omitempty"`
	ShippingAddress *ShippingAddress `json:"shipping_address,omitempty"`
	TaxIDType       string           `json:"tax_id_type,omitempty"`
	TaxID           string           `json:"tax_id,omitempty"`
	CountryCode     string           `json:"country_code"`
}

// ShippingAddress struct
type ShippingAddress struct {
	RecipientName string `json:"recipient_name,omitempty"`
	Type          string `json:"type,omitempty"`
	Line1         string `json:"line1"`
	Line2         string `json:"line2,omitempty"`
	City          string `json:"city"`
	CountryCode   string `json:"country_code"`
	PostalCode    string `json:"postal_code,omitempty"`
	State         string `json:"state,omitempty"`
	Phone         string `json:"phone,omitempty"`
}

// CreateAgreementResponse struct
type CreateAgreementResponse struct {
	Name        string      `json:"name,omitempty"`
	Description string      `json:"description,omitempty"`
	Plan        BillingPlan `json:"plan,omitempty"`
	Links       []Link      `json:"links,omitempty"`
	StartTime   time.Time   `json:"start_time,omitempty"`
}

// ExecuteAgreementResponse struct
type ExecuteAgreementResponse struct {
	ID               string           `json:"id"`
	State            string           `json:"state"`
	Description      string           `json:"description,omitempty"`
	Payer            Payer            `json:"payer"`
	Plan             BillingPlan      `json:"plan"`
	StartDate        time.Time        `json:"start_date"`
	ShippingAddress  ShippingAddress  `json:"shipping_address"`
	AgreementDetails AgreementDetails `json:"agreement_details"`
	Links            []Link           `json:"links"`
}

// AgreementDetails struct
type AgreementDetails struct {
	OutstandingBalance AmountPayout `json:"outstanding_balance"`
	CyclesRemaining    int          `json:"cycles_remaining,string"`
	CyclesCompleted    int          `json:"cycles_completed,string"`
	NextBillingDate    time.Time    `json:"next_billing_date"`
	LastPaymentDate    time.Time    `json:"last_payment_date"`
	LastPaymentAmount  AmountPayout `json:"last_payment_amount"`
	FinalPaymentDate   time.Time    `json:"final_payment_date"`
	FailedPaymentCount int          `json:"failed_payment_count,string"`
}

// UserInfo struct
type UserInfo struct {
	ID              string   `json:"user_id"`
	Name            string   `json:"name"`
	GivenName       string   `json:"given_name"`
	FamilyName      string   `json:"family_name"`
	Email           string   `json:"email"`
	Verified        bool     `json:"verified,omitempty,string"`
	Gender          string   `json:"gender,omitempty"`
	BirthDate       string   `json:"birthdate,omitempty"`
	ZoneInfo        string   `json:"zoneinfo,omitempty"`
	Locale          string   `json:"locale,omitempty"`
	Phone           string   `json:"phone_number,omitempty"`
	Address         *Address `json:"address,omitempty"`
	VerifiedAccount bool     `json:"verified_account,omitempty,string"`
	AccountType     string   `json:"account_type,omitempty"`
	AgeRange        string   `json:"age_range,omitempty"`
	PayerID         string   `json:"payer_id,omitempty"`
}

// WebProfile represents the configuration of the payment web payment experience.
// https://developer.paypal.com/docs/api/payment-experience/
type WebProfile struct {
	ID           string       `json:"id,omitempty"`
	Name         string       `json:"name"`
	Presentation Presentation `json:"presentation,omitempty"`
	InputFields  InputFields  `json:"input_fields,omitempty"`
	FlowConfig   FlowConfig   `json:"flow_config,omitempty"`
}

// Presentation represents the branding and locale that a customer sees on redirect payments.
// https://developer.paypal.com/docs/api/payment-experience/#definition-presentation
type Presentation struct {
	BrandName  string `json:"brand_name,omitempty"`
	LogoImage  string `json:"logo_image,omitempty"`
	LocaleCode string `json:"locale_code,omitempty"`
}

// InputFields represents the fields that are displayed to a customer on redirect payments.
// https://developer.paypal.com/docs/api/payment-experience/#definition-input_fields
type InputFields struct {
	AllowNote       bool `json:"allow_note,omitempty"`
	NoShipping      uint `json:"no_shipping,omitempty"`
	AddressOverride uint `json:"address_override,omitempty"`
}

// FlowConfig represents the general behaviour of redirect payment pages.
// https://developer.paypal.com/docs/api/payment-experience/#definition-flow_config
type FlowConfig struct {
	LandingPageType   string `json:"landing_page_type,omitempty"`
	BankTXNPendingURL string `json:"bank_txn_pending_url,omitempty"`
	UserAction        string `json:"user_action,omitempty"`
}

// End PayPal Models //
