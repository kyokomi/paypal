package paypal

import (
	"net/http"
)

const (
	sandboxURL = "https://api.sandbox.paypal.com" // for testing
	liveURL    = "https://api.paypal.com"         // production
)

type PayPalClient struct {
	*http.Client

	Options PayPalOptions
	Admin   AdminAuthToken

	// service
	OAuth2  OAuth2Service
	Payment PaymentService
}

func NewClient(options PayPalOptions) *PayPalClient {
	c := &PayPalClient{}
	c.Client = &http.Client{}
	c.Options = options

	// service
	c.OAuth2 = OAuth2Service{client: c}
	c.Payment = PaymentService{client: c}

	return c
}

func (c PayPalClient) Authorization() string {
	return c.Admin.Authorization()
}

func (c PayPalClient) URL(url string) string {
	if c.Options.Sandbox {
		return sandboxURL + url
	} else {
		return liveURL + url
	}
}
