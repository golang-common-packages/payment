package payment

type Config struct {
	ClientID  string `json:"clientID,omitempty"`
	SecretID  string `json:"secretID,omitempty"`
	PublicKey string `json:"publicKey,omitempty"`
}

type MoneyTransfer struct {
	RecipientIBAN  string `json:"recipientIBAN,omitempty"`
	Recipient      string `json:"recipient,omitempty"`
	TransferMethod string `json:"transferMethod,omitempty"`
	Amount         string `json:"amount,omitempty"`
	CurrencyType   string `json:"currencyType,omitempty"`
	Comment        string `json:"comment,omitempty"`
	EmailSubject   string `json:"emailSubject,omitempty"`
	Address        `json:"address,omitempty"`
}

type Address struct {
	Street     string `json:"street,omitempty"`
	City       string `json:"city,omitempty"`
	PostalCode string `json:"postalCode,omitempty"`
	Country    string `json:"country,omitempty"`
}

type BankAccount struct {
	LinkToPayPal LinkToPayPal `json:"linkToPayPal,omitempty"`
	LinkToPlaid  LinkToPlaid  `json:"linkToPlaid,omitempty"`
	LinkToStripe LinkToStripe `json:"linkToStripe,omitempty"`
}

type LinkToPayPal struct {
	BankAccountNumber              string `json:"bankAccountNumber,omitempty"`
	InternationalBankAccountNumber string `json:"iban,omitempty"`
	BankAccountType                string `json:"bankAccountType,omitempty"`
	BankCountryCode                string `json:"bankCountryCode,omitempty"`
	BankName                       string `json:"bankName,omitempty"`
	CLABE                          string `json:"clabe,omitempty"`
	ConfirmationType               string `json:"confirmationType,omitempty"`
}

type LinkToPlaid struct {
	RecipientName                  string `json:"recipientName,omitempty"`
	InternationalBankAccountNumber string `json:"iban,omitempty"`
	Address                        `json:"address,omitempty"`
}

type LinkToStripe struct {
	CustomerID        string `json:"customerID,omitempty"`
	Token             string `json:"token,omitempty"`
	AccountHolderName string `json:"accountHolderName,omitempty"`
	AccountHolderType string `json:"accountHolderType,omitempty"`
	AccountNumber     string `json:"accountNumber,omitempty"`
	Country           string `json:"country,omitempty"`
	Currency          string `json:"currency,omitempty"`
}
