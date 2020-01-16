package payment

// IMailClient store function in email package
type IMailClient interface {
	SubmitPayment(emailSubject, recipientType, receiver, amount, currencyType, sendingNote string) (interface{}, error)
	GetPayment(payoutBatchID string) (interface{}, error)
	GetPaymentItem(payoutBatchID string) (interface{}, error)
	ListCreditCards(page, pageSize int) (interface{}, error)
	GetCreditCardDetail(creditCardID string) (interface{}, error)
	StoreCreditCardDetail(line1, line2, city, countryCode, postalCode, state, phone, id, payerID, externalCustomerID, number, typeCard, expireMonth, expireYear, cvv2, firstName, lastName, State, ValidUntil string) (interface{}, error)
}

/*
	@PAYPAL: Paypal service
*/
const (
	PAYPAL = iota
)

// NewPaymentClient function for Factory Pattern
func NewPaymentClient(paymentType int, config *Config) IMailClient {
	switch paymentType {
	case PAYPAL:
		return NewPaypalClient(config.ClientID, config.SecretID)
	}

	return nil
}
