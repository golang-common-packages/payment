package payment

import (
	"log"
	"os"

	"github.com/plutov/paypal"
)

type PaypalClient struct {
	client *paypal.Client
}

func NewPaypalClient(clientID, secretID string) IMailClient {
	currentSesstion := &PaypalClient{nil}

	client, err := paypal.NewClient("clientID", "secretID", paypal.APIBaseSandBox)
	if err != nil {
		log.Println("Error when try to make strconv port from config: ", err)
		panic(err)
	}
	client.SetLog(os.Stdout) // Set log to terminal stdout

	currentSesstion.client = client

	return currentSesstion
}

func (pp *PaypalClient) SubmitPayment(emailSubject, recipientType, receiver, amount, currencyType, sendingNote string) (interface{}, error) {
	payout := paypal.Payout{
		SenderBatchHeader: &paypal.SenderBatchHeader{
			EmailSubject: emailSubject,
		},
		Items: []paypal.PayoutItem{
			paypal.PayoutItem{
				RecipientType: recipientType,
				Receiver:      receiver,
				Amount: &paypal.AmountPayout{
					Value:    amount,
					Currency: currencyType,
				},
				Note: sendingNote,
			},
		},
	}

	payoutResp, err := pp.client.CreateSinglePayout(payout)
	if err != nil {
		return nil, err
	}

	return payoutResp, nil
}

func (pp *PaypalClient) GetPayment(payoutBatchID string) (interface{}, error) {
	payloads, err := pp.client.GetPayout(payoutBatchID)
	if err != nil {
		return nil, err
	}

	return payloads, nil
}

func (pp *PaypalClient) GetPaymentItem(payoutBatchID string) (interface{}, error) {
	payload, err := pp.client.GetPayoutItem(payoutBatchID)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (pp *PaypalClient) ListCreditCards(page, pageSize int) (interface{}, error) {
	creditCardFilter := &paypal.CreditCardsFilter{
		PageSize: pageSize,
		Page:     page,
	}

	creaditCards, err := pp.client.GetCreditCards(creditCardFilter)
	if err != nil {
		return nil, err
	}

	return creaditCards, nil
}

func (pp *PaypalClient) GetCreditCardDetail(creditCardID string) (interface{}, error) {
	detail, err := pp.client.GetCreditCard(creditCardID)
	if err != nil {
		return nil, err
	}

	return detail, nil
}

func (pp *PaypalClient) StoreCreditCardDetail(line1, line2, city, countryCode, postalCode, state, phone, id, payerID, externalCustomerID, number, typeCard, expireMonth, expireYear, cvv2, firstName, lastName, State, ValidUntil string) (interface{}, error) {
	billingAddress := generateAddress(line1, line2, city, countryCode, postalCode, state, phone)
	creditCard := generateCreditDetail(id, payerID, externalCustomerID, number, typeCard, expireMonth, expireYear, cvv2, firstName, lastName, State, ValidUntil, billingAddress)

	result, err := pp.client.StoreCreditCard(creditCard)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (pp *PaypalClient) RemoveCreditCardDetail(creditCardID string) error {
	err := pp.client.DeleteCreditCard(creditCardID)
	return err
}

// Paypal util
func generateAddress(line1, line2, city, countryCode, postalCode, state, phone string) *paypal.Address {
	address := &paypal.Address{
		Line1:       line1,
		Line2:       line2,
		City:        city,
		CountryCode: countryCode,
		PostalCode:  postalCode,
		State:       state,
		Phone:       phone,
	}

	return address
}

func generateCreditDetail(id, payerID, externalCustomerID, number, typeCard, expireMonth, expireYear, cvv2, firstName, lastName, State, ValidUntil string, billingAddress *paypal.Address) paypal.CreditCard {
	creditCardDetail := paypal.CreditCard{
		ID:                 id,
		PayerID:            payerID,
		ExternalCustomerID: externalCustomerID,
		Number:             number,
		Type:               typeCard,
		ExpireMonth:        expireMonth,
		ExpireYear:         expireYear,
		CVV2:               cvv2,
		FirstName:          firstName,
		LastName:           lastName,
		BillingAddress:     billingAddress,
		State:              State,
		ValidUntil:         ValidUntil,
	}

	return creditCardDetail
}
