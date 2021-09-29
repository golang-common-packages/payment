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

// End PayPal Models //
