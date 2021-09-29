package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/golang-common-packages/hash"
)

// IPayPal interface for PayPal services
type IPayPal interface {
	GetAccessToken(ctx context.Context) (*TokenResponse, error)
}

// PayPalClient represents a Paypal REST API Client
type PayPalClient struct {
	sync.Mutex
	Client               *http.Client
	ClientID             string
	Secret               string
	APIBase              string
	Log                  io.Writer // If user set log file name all requests will be logged there
	Token                *TokenResponse
	tokenExpiresAt       time.Time
	returnRepresentation bool
}

// payPalClientSessionMapping singleton pattern
var payPalClientSessionMapping = make(map[string]*PayPalClient)

// newPayPal init new instance.
// APIBase is a base API URL, for testing you can use paypal.APIBaseSandBox
func newPayPal(config *PayPal) IPayPal {
	// Validate config file
	if config.ClientID == "" || config.SecretID == "" || config.APIBase == "" {
		log.Fatalln("ClientID, Secret and APIBase are required to create a Client")
	}

	// Init PayPal client with singleton pattern
	hasher := &hash.Client{}
	configAsJSON, err := json.Marshal(config)
	if err != nil {
		log.Fatalln("Unable to marshal PayPal configuration: ", err)
	}
	configAsString := hasher.SHA1(string(configAsJSON))

	currentPayPalSession := payPalClientSessionMapping[configAsString]
	if currentPayPalSession == nil {

		currentPayPalSession.Client = &http.Client{}
		currentPayPalSession.ClientID = config.ClientID
		currentPayPalSession.Secret = config.SecretID
		currentPayPalSession.APIBase = config.APIBase
		payPalClientSessionMapping[configAsString] = currentPayPalSession

		log.Println("Init PayPal client successfully")
	}

	return currentPayPalSession
}

// GetAccessToken returns struct of TokenResponse.
// No need to call SetAccessToken to apply new access token for current Client.
// Endpoint: POST /v1/oauth2/token
func (c *PayPalClient) GetAccessToken(ctx context.Context) (*TokenResponse, error) {
	buf := bytes.NewBuffer([]byte("grant_type=client_credentials"))
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/oauth2/token"), buf)
	if err != nil {
		return &TokenResponse{}, err
	}

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")

	response := &TokenResponse{}
	err = c.SendWithBasicAuth(req, response)

	// Set Token fur current Client
	if response.Token != "" {
		c.Token = response
		c.tokenExpiresAt = time.Now().Add(time.Duration(response.ExpiresIn) * time.Second)
	}

	return response, err
}
