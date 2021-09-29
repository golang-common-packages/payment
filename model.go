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

// ErrorResponse based on https://developer.paypal.com/docs/api/errors/
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

// End PayPal Models //
