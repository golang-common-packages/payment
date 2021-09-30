package payment

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
