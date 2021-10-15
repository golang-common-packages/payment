package payment

import "context"

const (
	// Paypal services
	PAYPAL = iota
)

var (
	// Init context with default value
	ctx = context.Background()
)

// New payment by abstract factory pattern
func New(context context.Context, paymentCompany int, config *Config) interface{} {
	SetContext(context)

	switch paymentCompany {
	case PAYPAL:
		return newPayPal(&config.PayPal)
	default:
		return nil
	}
}
