package payment

import "net/http"

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

// End PayPal Models //
