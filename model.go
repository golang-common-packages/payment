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
