package payment

import (
	"github.com/plaid/plaid-go/plaid"
)

type Config struct {
	ClientID  string `json:"clientID,omitempty"`
	SecretID  string `json:"secretID,omitempty"`
	PublicKey string `json:"publicKey,omitempty"`
}

type MoneyTransfer struct {
	ReceiverIBAN string                    `json:"recipientIBAN,omitempty"`
	Receiver     string                    `json:"recipient,omitempty"`      /// PayPal
	ReceiverType string                    `json:"transferMethod,omitempty"` /// PayPal
	Amount       `json:"amount,omitempty"` /// PayPal
	Comment      string                    `json:"comment,omitempty"`      /// PayPal
	EmailSubject string                    `json:"emailSubject,omitempty"` /// PayPal
	Address      `json:"address,omitempty"`
}

type Amount struct {
	Value    string `json:"amount,omitempty"`
	Currency string `json:"currencyType,omitempty"`
}

type Address struct {
	Street     string `json:"street,omitempty"`
	City       string `json:"city,omitempty"`
	PostalCode string `json:"postalCode,omitempty"`
	Country    string `json:"country,omitempty"`
}

type PayPalLinkBank struct {
	BankAccountNumber              string `json:"bankAccountNumber,omitempty"`
	InternationalBankAccountNumber string `json:"iban,omitempty"`
	BankAccountType                string `json:"bankAccountType,omitempty"`
	BankCountryCode                string `json:"bankCountryCode,omitempty"`
	BankName                       string `json:"bankName,omitempty"`
	CLABE                          string `json:"clabe,omitempty"`
	ConfirmationType               string `json:"confirmationType,omitempty"`
}

type StripeLinkBank struct {
	CustomerID        string `json:"customerID,omitempty"`
	Token             string `json:"auth,omitempty"`
	AccountHolderName string `json:"accountHolderName,omitempty"`
	AccountHolderType string `json:"accountHolderType,omitempty"`
	AccountNumber     string `json:"accountNumber,omitempty"`
	Country           string `json:"country,omitempty"`
	Currency          string `json:"currency,omitempty"`
}

///// Plaid models /////

type PlaidTransactionsHistory struct {
	Accounts     []plaid.Account     `json:"Accounts,omitempty"`
	Transactions []plaid.Transaction `json:"Transactions,omitempty"`
}

type PlaidPayment struct {
	ProductName         string `json:"productName,omitempty"`
	IBAN                string `json:"iban,omitempty"`
	PlaidAmount         `json:"plaidAmount,omitempty"`
	PlaidPaymentAddress `json:"plaidPaymentAddress,omitempty"`
}

type PlaidAmount struct {
	Currency string  `json:"currency,omitempty"`
	Amount   float64 `json:"value,omitempty"`
}

type PlaidPaymentAddress struct {
	Street     []string `json:"street,omitempty"`
	City       string   `json:"city,omitempty"`
	PostalCode string   `json:"postalCode,omitempty"`
	Country    string   `json:"country,omitempty"`
}

type PlaidPaymentResult struct {
	RecipientID  string `json:"recipientID,omitempty"`
	PaymentID    string `json:"paymentID,omitempty"`
	PaymentToken string `json:"paymentToken,omitempty"`
}
