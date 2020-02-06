package payment

import "context"

type IPayment interface {
	TransferMoney(transferInfo *MoneyTransfer) (result interface{}, err error)
}

var ctx = context.Background()

const (
	PAYPAL = iota
	PLAID
	STRIPE
)

func NewPayment(paymentProvider int, config *Config) IPayment {
	switch paymentProvider {
	case PAYPAL:
		return NewPaypal(config.ClientID, config.SecretID)
	case PLAID:
		return NewPlaid(config.ClientID, config.SecretID, config.PublicKey)
	case STRIPE:
		return NewStripe(config.SecretID)
	}
}
