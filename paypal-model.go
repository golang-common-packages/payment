package payment

import (
	"net/http"
	"time"
)

// BillingPlanStatus has type is string. This type may change in the future
type BillingPlanStatus string

// ShippingPreference has type is string. This type may change in the future
type ShippingPreference string

// UserAction has type is string. This type may change in the future
type UserAction string

// JSONTime overrides MarshalJson method to format in ISO8601
type JSONTime time.Time

//Doc: https://developer.paypal.com/docs/api/catalog-products/v1/#definition-product_category
type ProductCategory string

type ProductType string

type SubscriptionPlanStatus string

type IntervalUnit string

type TenureType string

type SetupFeeFailureAction string

type CaptureType string

type SubscriptionTransactionStatus string

type SubscriptionStatus string

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

// TransactionSearchRequest struct
type TransactionSearchRequest struct {
	TransactionID               *string
	TransactionType             *string
	TransactionStatus           *string
	TransactionAmount           *string
	TransactionCurrency         *string
	StartDate                   time.Time
	EndDate                     time.Time
	PaymentInstrumentType       *string
	StoreID                     *string
	TerminalID                  *string
	Fields                      *string
	BalanceAffectingRecordsOnly *string
	PageSize                    *int
	Page                        *int
}

// TransactionSearchResponse struct
type TransactionSearchResponse struct {
	TransactionDetails  []SearchTransactionDetails `json:"transaction_details"`
	AccountNumber       string                     `json:"account_number"`
	StartDate           JSONTime                   `json:"start_date"`
	EndDate             JSONTime                   `json:"end_date"`
	LastRefreshDatetime JSONTime                   `json:"last_refreshed_datetime"`
	Page                int                        `json:"page"`
	SharedListResponse
}

// SearchTransactionDetails struct
type SearchTransactionDetails struct {
	TransactionInfo SearchTransactionInfo `json:"transaction_info"`
	PayerInfo       *SearchPayerInfo      `json:"payer_info"`
	ShippingInfo    *SearchShippingInfo   `json:"shipping_info"`
	CartInfo        *SearchCartInfo       `json:"cart_info"`
}

// SearchTransactionInfo struct
type SearchTransactionInfo struct {
	PayPalAccountID           string   `json:"paypal_account_id"`
	TransactionID             string   `json:"transaction_id"`
	PayPalReferenceID         string   `json:"paypal_reference_id"`
	PayPalReferenceIDType     string   `json:"paypal_reference_id_type"`
	TransactionEventCode      string   `json:"transaction_event_code"`
	TransactionInitiationDate JSONTime `json:"transaction_initiation_date"`
	TransactionUpdatedDate    JSONTime `json:"transaction_updated_date"`
	TransactionAmount         Money    `json:"transaction_amount"`
	FeeAmount                 *Money   `json:"fee_amount"`
	InsuranceAmount           *Money   `json:"insurance_amount"`
	ShippingAmount            *Money   `json:"shipping_amount"`
	ShippingDiscountAmount    *Money   `json:"shipping_discount_amount"`
	ShippingTaxAmount         *Money   `json:"shipping_tax_amount"`
	OtherAmount               *Money   `json:"other_amount"`
	TipAmount                 *Money   `json:"tip_amount"`
	TransactionStatus         string   `json:"transaction_status"`
	TransactionSubject        string   `json:"transaction_subject"`
	PaymentTrackingID         string   `json:"payment_tracking_id"`
	BankReferenceID           string   `json:"bank_reference_id"`
	TransactionNote           string   `json:"transaction_note"`
	EndingBalance             *Money   `json:"ending_balance"`
	AvailableBalance          *Money   `json:"available_balance"`
	InvoiceID                 string   `json:"invoice_id"`
	CustomField               string   `json:"custom_field"`
	ProtectionEligibility     string   `json:"protection_eligibility"`
	CreditTerm                string   `json:"credit_term"`
	CreditTransactionalFee    *Money   `json:"credit_transactional_fee"`
	CreditPromotionalFee      *Money   `json:"credit_promotional_fee"`
	AnnualPercentageRate      string   `json:"annual_percentage_rate"`
	PaymentMethodType         string   `json:"payment_method_type"`
}

// SearchPayerInfo struct
type SearchPayerInfo struct {
	AccountID     string               `json:"account_id"`
	EmailAddress  string               `json:"email_address"`
	PhoneNumber   *PhoneWithTypeNumber `json:"phone_number"`
	AddressStatus string               `json:"address_status"`
	PayerStatus   string               `json:"payer_status"`
	PayerName     SearchPayerName      `json:"payer_name"`
	CountryCode   string               `json:"country_code"`
	Address       *Address             `json:"address"`
}

// PhoneWithTypeNumber struct for PhoneWithType
type PhoneWithTypeNumber struct {
	NationalNumber string `json:"national_number,omitempty"`
}

// SearchPayerName struct
type SearchPayerName struct {
	GivenName string `json:"given_name"`
	Surname   string `json:"surname"`
}

// SearchShippingInfo struct
type SearchShippingInfo struct {
	Name                     string   `json:"name"`
	Method                   string   `json:"method"`
	Address                  Address  `json:"address"`
	SecondaryShippingAddress *Address `json:"secondary_shipping_address"`
}

// SearchCartInfo struct
type SearchCartInfo struct {
	ItemDetails     []SearchItemDetails `json:"item_details"`
	TaxInclusive    *bool               `json:"tax_inclusive"`
	PayPalInvoiceID string              `json:"paypal_invoice_id"`
}

// SearchItemDetails struct
type SearchItemDetails struct {
	ItemCode            string                 `json:"item_code"`
	ItemName            string                 `json:"item_name"`
	ItemDescription     string                 `json:"item_description"`
	ItemOptions         string                 `json:"item_options"`
	ItemQuantity        string                 `json:"item_quantity"`
	ItemUnitPrice       Money                  `json:"item_unit_price"`
	ItemAmount          Money                  `json:"item_amount"`
	DiscountAmount      *Money                 `json:"discount_amount"`
	AdjustmentAmount    *Money                 `json:"adjustment_amount"`
	GiftWrapAmount      *Money                 `json:"gift_wrap_amount"`
	TaxPercentage       string                 `json:"tax_percentage"`
	TaxAmounts          []SearchTaxAmount      `json:"tax_amounts"`
	BasicShippingAmount *Money                 `json:"basic_shipping_amount"`
	ExtraShippingAmount *Money                 `json:"extra_shipping_amount"`
	HandlingAmount      *Money                 `json:"handling_amount"`
	InsuranceAmount     *Money                 `json:"insurance_amount"`
	TotalItemAmount     Money                  `json:"total_item_amount"`
	InvoiceNumber       string                 `json:"invoice_number"`
	CheckoutOptions     []SearchCheckoutOption `json:"checkout_options"`
}

// SearchTaxAmount struct
type SearchTaxAmount struct {
	TaxAmount Money `json:"tax_amount"`
}

// SearchCheckoutOption struct
type SearchCheckoutOption struct {
	CheckoutOptionName  string `json:"checkout_option_name"`
	CheckoutOptionValue string `json:"checkout_option_value"`
}

// CreditCardsFilter struct
type CreditCardsFilter struct {
	PageSize int
	Page     int
}

// CreditCards struct
type CreditCards struct {
	Items []CreditCard `json:"items"`
	SharedListResponse
}

// CreditCardField struct
type CreditCardField struct {
	Operation string `json:"op"`
	Path      string `json:"path"`
	Value     string `json:"value"`
}

// Order struct
type Order struct {
	ID            string                 `json:"id,omitempty"`
	Status        string                 `json:"status,omitempty"`
	Intent        string                 `json:"intent,omitempty"`
	Payer         *PayerWithNameAndPhone `json:"payer,omitempty"`
	PurchaseUnits []PurchaseUnit         `json:"purchase_units,omitempty"`
	Links         []Link                 `json:"links,omitempty"`
	CreateTime    *time.Time             `json:"create_time,omitempty"`
	UpdateTime    *time.Time             `json:"update_time,omitempty"`
}

// PayerWithNameAndPhone struct
type PayerWithNameAndPhone struct {
	Name         *CreateOrderPayerName          `json:"name,omitempty"`
	EmailAddress string                         `json:"email_address,omitempty"`
	Phone        *PhoneWithType                 `json:"phone,omitempty"`
	PayerID      string                         `json:"payer_id,omitempty"`
	BirthDate    string                         `json:"birth_date,omitempty"`
	TaxInfo      *TaxInfo                       `json:"tax_info,omitempty"`
	Address      *ShippingDetailAddressPortable `json:"address,omitempty"`
}

// CreateOrderPayerName create order payer name
type CreateOrderPayerName struct {
	GivenName string `json:"given_name,omitempty"`
	Surname   string `json:"surname,omitempty"`
}

// PhoneWithType struct used for orders
type PhoneWithType struct {
	PhoneType   string               `json:"phone_type,omitempty"`
	PhoneNumber *PhoneWithTypeNumber `json:"phone_number,omitempty"`
}

// TaxInfo used for orders.
type TaxInfo struct {
	TaxID     string `json:"tax_id,omitempty"`
	TaxIDType string `json:"tax_id_type,omitempty"`
}

// ShippingDetailAddressPortable used with create orders
type ShippingDetailAddressPortable struct {
	AddressLine1 string `json:"address_line_1,omitempty"`
	AddressLine2 string `json:"address_line_2,omitempty"`
	AdminArea1   string `json:"admin_area_1,omitempty"`
	AdminArea2   string `json:"admin_area_2,omitempty"`
	PostalCode   string `json:"postal_code,omitempty"`
	CountryCode  string `json:"country_code,omitempty"`
}

// PurchaseUnit struct
type PurchaseUnit struct {
	ReferenceID        string              `json:"reference_id"`
	Amount             *PurchaseUnitAmount `json:"amount,omitempty"`
	Payee              *PayeeForOrders     `json:"payee,omitempty"`
	Payments           *CapturedPayments   `json:"payments,omitempty"`
	PaymentInstruction *PaymentInstruction `json:"payment_instruction,omitempty"`
	Description        string              `json:"description,omitempty"`
	CustomID           string              `json:"custom_id,omitempty"`
	InvoiceID          string              `json:"invoice_id,omitempty"`
	ID                 string              `json:"id,omitempty"`
	SoftDescriptor     string              `json:"soft_descriptor,omitempty"`
	Shipping           *ShippingDetail     `json:"shipping,omitempty"`
	Items              []Item              `json:"items,omitempty"`
}

// PayeeForOrders struct
type PayeeForOrders struct {
	EmailAddress string `json:"email_address,omitempty"`
	MerchantID   string `json:"merchant_id,omitempty"`
}

// CapturedPayments has the amounts for a captured order
type CapturedPayments struct {
	Captures []CaptureAmount `json:"captures,omitempty"`
}

// https://developer.paypal.com/docs/api/payments/v2/#definition-payment_instruction
type PaymentInstruction struct {
	PlatformFees     []PlatformFee `json:"platform_fees,omitempty"`
	DisbursementMode string        `json:"disbursement_mode,omitempty"`
}

// ShippingDetail struct
type ShippingDetail struct {
	Name    *Name                          `json:"name,omitempty"`
	Address *ShippingDetailAddressPortable `json:"address,omitempty"`
}

// Name struct.
// https://developer.paypal.com/docs/api/subscriptions/v1/#definition-name
type Name struct {
	FullName   string `json:"full_name,omitempty"`
	Suffix     string `json:"suffix,omitempty"`
	Prefix     string `json:"prefix,omitempty"`
	GivenName  string `json:"given_name,omitempty"`
	Surname    string `json:"surname,omitempty"`
	MiddleName string `json:"middle_name,omitempty"`
}

// CaptureAmount struct
type CaptureAmount struct {
	ID                        string                     `json:"id,omitempty"`
	CustomID                  string                     `json:"custom_id,omitempty"`
	Amount                    *PurchaseUnitAmount        `json:"amount,omitempty"`
	SellerProtection          *SellerProtection          `json:"seller_protection,omitempty"`
	SellerReceivableBreakdown *SellerReceivableBreakdown `json:"seller_receivable_breakdown,omitempty"`
}

// SellerReceivableBreakdown has the detailed breakdown of the capture activity.
type SellerReceivableBreakdown struct {
	GrossAmount                   *Money        `json:"gross_amount,omitempty"`
	PaypalFee                     *Money        `json:"paypal_fee,omitempty"`
	PaypalFeeInReceivableCurrency *Money        `json:"paypal_fee_in_receivable_currency,omitempty"`
	NetAmount                     *Money        `json:"net_amount,omitempty"`
	ReceivableAmount              *Money        `json:"receivable_amount,omitempty"`
	ExchangeRate                  *ExchangeRate `json:"exchange_rate,omitempty"`
	PlatformFees                  []PlatformFee `json:"platform_fees,omitempty"`
}

// ExchangeRate struct.
// https://developer.paypal.com/docs/api/orders/v2/#definition-exchange_rate
type ExchangeRate struct {
	SourceCurrency string `json:"source_currency"`
	TargetCurrency string `json:"target_currency"`
	Value          string `json:"value"`
}

// PlatformFee struct.
// https://developer.paypal.com/docs/api/payments/v2/#definition-platform_fee
type PlatformFee struct {
	Amount *Money          `json:"amount,omitempty"`
	Payee  *PayeeForOrders `json:"payee,omitempty"`
}

// Item struct
type Item struct {
	Name        string `json:"name"`
	UnitAmount  *Money `json:"unit_amount,omitempty"`
	Tax         *Money `json:"tax,omitempty"`
	Quantity    string `json:"quantity"`
	Description string `json:"description,omitempty"`
	SKU         string `json:"sku,omitempty"`
	Category    string `json:"category,omitempty"`
}

// PurchaseUnitRequest struct
type PurchaseUnitRequest struct {
	ReferenceID        string              `json:"reference_id,omitempty"`
	Amount             *PurchaseUnitAmount `json:"amount"`
	Payee              *PayeeForOrders     `json:"payee,omitempty"`
	Description        string              `json:"description,omitempty"`
	CustomID           string              `json:"custom_id,omitempty"`
	InvoiceID          string              `json:"invoice_id,omitempty"`
	SoftDescriptor     string              `json:"soft_descriptor,omitempty"`
	Items              []Item              `json:"items,omitempty"`
	Shipping           *ShippingDetail     `json:"shipping,omitempty"`
	PaymentInstruction *PaymentInstruction `json:"payment_instruction,omitempty"`
}

// CreateOrderPayer used with create order requests
type CreateOrderPayer struct {
	Name         *CreateOrderPayerName          `json:"name,omitempty"`
	EmailAddress string                         `json:"email_address,omitempty"`
	PayerID      string                         `json:"payer_id,omitempty"`
	Phone        *PhoneWithType                 `json:"phone,omitempty"`
	BirthDate    string                         `json:"birth_date,omitempty"`
	TaxInfo      *TaxInfo                       `json:"tax_info,omitempty"`
	Address      *ShippingDetailAddressPortable `json:"address,omitempty"`
}

// ApplicationContext struct
type ApplicationContext struct {
	BrandName          string             `json:"brand_name,omitempty"`
	Locale             string             `json:"locale,omitempty"`
	ShippingPreference ShippingPreference `json:"shipping_preference,omitempty"`
	UserAction         UserAction         `json:"user_action,omitempty"`
	ReturnURL          string             `json:"return_url,omitempty"`
	CancelURL          string             `json:"cancel_url,omitempty"`
}

// AuthorizeOrderRequest struct.
// https://developer.paypal.com/docs/api/orders/v2/#orders_authorize
type AuthorizeOrderRequest struct {
	PaymentSource      *PaymentSource     `json:"payment_source,omitempty"`
	ApplicationContext ApplicationContext `json:"application_context,omitempty"`
}

// PaymentSource structure
type PaymentSource struct {
	Card  *PaymentSourceCard  `json:"card,omitempty"`
	Token *PaymentSourceToken `json:"token,omitempty"`
}

// PaymentSourceCard struct
type PaymentSourceCard struct {
	ID             string              `json:"id"`
	Name           string              `json:"name"`
	Number         string              `json:"number"`
	Expiry         string              `json:"expiry"`
	SecurityCode   string              `json:"security_code"`
	LastDigits     string              `json:"last_digits"`
	CardType       string              `json:"card_type"`
	BillingAddress *CardBillingAddress `json:"billing_address"`
}

// CardBillingAddress struct
type CardBillingAddress struct {
	AddressLine1 string `json:"address_line_1"`
	AddressLine2 string `json:"address_line_2"`
	AdminArea2   string `json:"admin_area_2"`
	AdminArea1   string `json:"admin_area_1"`
	PostalCode   string `json:"postal_code"`
	CountryCode  string `json:"country_code"`
}

// PaymentSourceToken struct
type PaymentSourceToken struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// CaptureOrderRequest.
// https://developer.paypal.com/docs/api/orders/v2/#orders_capture
type CaptureOrderRequest struct {
	PaymentSource *PaymentSource `json:"payment_source"`
}

// CaptureOrderResponse is the response for capture order
type CaptureOrderResponse struct {
	ID            string                 `json:"id,omitempty"`
	Status        string                 `json:"status,omitempty"`
	Payer         *PayerWithNameAndPhone `json:"payer,omitempty"`
	Address       *Address               `json:"address,omitempty"`
	PurchaseUnits []CapturedPurchaseUnit `json:"purchase_units,omitempty"`
}

// CapturedPurchaseUnit are purchase units for a captured order
type CapturedPurchaseUnit struct {
	Items       []CapturedPurchaseItem       `json:"items,omitempty"`
	ReferenceID string                       `json:"reference_id"`
	Shipping    CapturedPurchaseUnitShipping `json:"shipping,omitempty"`
	Payments    *CapturedPayments            `json:"payments,omitempty"`
}

// CapturedPurchaseItem are items for a captured order
type CapturedPurchaseItem struct {
	Quantity    string `json:"quantity"`
	Name        string `json:"name"`
	SKU         string `json:"sku,omitempty"`
	Description string `json:"description,omitempty"`
}

// CapturedPurchaseUnitShipping struct
type CapturedPurchaseUnitShipping struct {
	Address ShippingDetailAddressPortable `json:"address,omitempty"`
}

// CreateWebhookRequest struct
type CreateWebhookRequest struct {
	URL        string             `json:"url"`
	EventTypes []WebhookEventType `json:"event_types"`
}

// WebhookEventType struct
type WebhookEventType struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status,omitempty"`
}

// ListWebhookResponse struct
type ListWebhookResponse struct {
	Webhooks []Webhook `json:"webhooks"`
}

// Webhook struct
type Webhook struct {
	ID         string             `json:"id"`
	URL        string             `json:"url"`
	EventTypes []WebhookEventType `json:"event_types"`
	Links      []Link             `json:"links"`
}

// WebhookField struct
type WebhookField struct {
	Operation string      `json:"op"`
	Path      string      `json:"path"`
	Value     interface{} `json:"value"`
}

// VerifyWebhookResponse struct
type VerifyWebhookResponse struct {
	VerificationStatus string `json:"verification_status,omitempty"`
}

// WebhookEventTypesResponse struct
type WebhookEventTypesResponse struct {
	EventTypes []WebhookEventType `json:"event_types"`
}

type Product struct {
	ID          string          `json:"id,omitempty"`
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Category    ProductCategory `json:"category,omitempty"`
	Type        ProductType     `json:"type"`
	ImageUrl    string          `json:"image_url,omitempty"`
	HomeUrl     string          `json:"home_url,omitempty"`
}

type CreateProductResponse struct {
	Product
	SharedResponse
}

type SharedResponse struct {
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
	Links      []Link `json:"links"`
}

type ListProductsResponse struct {
	Products []Product `json:"products"`
	SharedListResponse
}

type ProductListParameters struct {
	ListParams
}

type SubscriptionPlan struct {
	ID                 string                 `json:"id,omitempty"`
	ProductId          string                 `json:"product_id"`
	Name               string                 `json:"name"`
	Status             SubscriptionPlanStatus `json:"status"`
	Description        string                 `json:"description,omitempty"`
	BillingCycles      []BillingCycle         `json:"billing_cycles"`
	PaymentPreferences *PaymentPreferences    `json:"payment_preferences"`
	Taxes              *Taxes                 `json:"taxes"`
	QuantitySupported  bool                   `json:"quantity_supported"` //Indicates whether you can subscribe to this plan by providing a quantity for the goods or service.
}

// Doc https://developer.paypal.com/docs/api/subscriptions/v1/#definition-billing_cycle
type BillingCycle struct {
	PricingScheme PricingScheme `json:"pricing_scheme"` // The active pricing scheme for this billing cycle. A free trial billing cycle does not require a pricing scheme.
	Frequency     Frequency     `json:"frequency"`      // The frequency details for this billing cycle.
	TenureType    TenureType    `json:"tenure_type"`    // The tenure type of the billing cycle. In case of a plan having trial cycle, only 2 trial cycles are allowed per plan. The possible values are:
	Sequence      int           `json:"sequence"`       // The order in which this cycle is to run among other billing cycles. For example, a trial billing cycle has a sequence of 1 while a regular billing cycle has a sequence of 2, so that trial cycle runs before the regular cycle.
	TotalCycles   int           `json:"total_cycles"`   // The number of times this billing cycle gets executed. Trial billing cycles can only be executed a finite number of times (value between 1 and 999 for total_cycles). Regular billing cycles can be executed infinite times (value of 0 for total_cycles) or a finite number of times (value between 1 and 999 for total_cycles).
}

type PricingScheme struct {
	Version    int       `json:"version"`
	FixedPrice Money     `json:"fixed_price"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

//doc: https://developer.paypal.com/docs/api/subscriptions/v1/#definition-frequency
type Frequency struct {
	IntervalUnit  IntervalUnit `json:"interval_unit"`
	IntervalCount int          `json:"interval_count"` //different per unit. check documentation
}

type PaymentPreferences struct {
	AutoBillOutstanding     bool                  `json:"auto_bill_outstanding"`
	SetupFee                *Money                `json:"setup_fee"`
	SetupFeeFailureAction   SetupFeeFailureAction `json:"setup_fee_failure_action"`
	PaymentFailureThreshold int                   `json:"payment_failure_threshold"`
}

type Taxes struct {
	Percentage string `json:"percentage"`
	Inclusive  bool   `json:"inclusive"`
}

type CreateSubscriptionPlanResponse struct {
	SubscriptionPlan
	SharedResponse
}

type SubscriptionPlanListParameters struct {
	ProductId string `json:"product_id"`
	PlanIds   string `json:"plan_ids"` // Filters the response by list of plan IDs. Filter supports upto 10 plan IDs.
	ListParams
}

type ListSubscriptionPlansResponse struct {
	Plans []SubscriptionPlan `json:"plans"`
	SharedListResponse
}

type PricingSchemeUpdate struct {
	BillingCycleSequence int           `json:"billing_cycle_sequence"`
	PricingScheme        PricingScheme `json:"pricing_scheme"`
}

type PricingSchemeUpdateRequest struct {
	Schemes []PricingSchemeUpdate `json:"pricing_schemes"`
}

type SubscriptionBase struct {
	PlanID             string              `json:"plan_id"`
	StartTime          *JSONTime           `json:"start_time,omitempty"`
	EffectiveTime      *JSONTime           `json:"effective_time,omitempty"`
	Quantity           string              `json:"quantity,omitempty"`
	ShippingAmount     *Money              `json:"shipping_amount,omitempty"`
	Subscriber         *Subscriber         `json:"subscriber,omitempty"`
	AutoRenewal        bool                `json:"auto_renewal,omitempty"`
	ApplicationContext *ApplicationContext `json:"application_context,omitempty"`
	CustomID           string              `json:"custom_id,omitempty"`
}

type Subscriber struct {
	ShippingAddress ShippingDetail       `json:"shipping_address,omitempty"`
	Name            CreateOrderPayerName `json:"name,omitempty"`
	EmailAddress    string               `json:"email_address,omitempty"`
}

type Subscription struct {
	SubscriptionDetailResp
}

type SubscriptionDetailResp struct {
	SubscriptionBase
	SubscriptionDetails
	BillingInfo BillingInfo `json:"billing_info,omitempty"` // not found in documentation
	SharedResponse
}

type BillingInfo struct {
	OutstandingBalance  AmountPayout      `json:"outstanding_balance,omitempty"`
	CycleExecutions     []CycleExecutions `json:"cycle_executions,omitempty"`
	LastPayment         LastPayment       `json:"last_payment,omitempty"`
	NextBillingTime     time.Time         `json:"next_billing_time,omitempty"`
	FailedPaymentsCount int               `json:"failed_payments_count,omitempty"`
}

type CycleExecutions struct {
	TenureType      string `json:"tenure_type,omitempty"`
	Sequence        int    `json:"sequence,omitempty"`
	CyclesCompleted int    `json:"cycles_completed,omitempty"`
	CyclesRemaining int    `json:"cycles_remaining,omitempty"`
	TotalCycles     int    `json:"total_cycles,omitempty"`
}

type LastPayment struct {
	Amount Money     `json:"amount,omitempty"`
	Time   time.Time `json:"time,omitempty"`
}

type CaptureReqeust struct {
	Note        string      `json:"note"`
	CaptureType CaptureType `json:"capture_type"`
	Amount      Money       `json:"amount"`
}

type SubscriptionCaptureResponse struct {
	Status              SubscriptionTransactionStatus `json:"status"`
	Id                  string                        `json:"id"`
	AmountWithBreakdown AmountWithBreakdown           `json:"amount_with_breakdown"`
	PayerName           Name                          `json:"payer_name"`
	PayerEmail          string                        `json:"payer_email"`
	Time                time.Time                     `json:"time"`
}

//Doc: https://developer.paypal.com/docs/api/subscriptions/v1/#definition-amount_with_breakdown
type AmountWithBreakdown struct {
	GrossAmount    Money `json:"gross_amount"`
	FeeAmount      Money `json:"fee_amount"`
	ShippingAmount Money `json:"shipping_amount"`
	TaxAmount      Money `json:"tax_amount"`
	NetAmount      Money `json:"net_amount"`
}

type SubscriptionTransactionsParams struct {
	SubscriptionId string
	StartTime      time.Time
	EndTime        time.Time
}

type SubscriptionTransactionsResponse struct {
	Transactions []SubscriptionCaptureResponse `json:"transactions"`
	SharedListResponse
}

type SubscriptionDetails struct {
	ID                           string             `json:"id,omitempty"`
	SubscriptionStatus           SubscriptionStatus `json:"status,omitempty"`
	SubscriptionStatusChangeNote string             `json:"status_change_note,omitempty"`
	StatusUpdateTime             time.Time          `json:"status_update_time,omitempty"`
}
