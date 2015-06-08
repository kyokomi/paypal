package paypal

type PayPalOptions struct {
	ClientID       string
	Secret         string
	AcceptLanguage string
	Sandbox        bool
}

func NewOptions(clientID, secret string) PayPalOptions {
	c := PayPalOptions{}
	c.ClientID = clientID
	c.Secret = secret
	c.AcceptLanguage = "en_US"
	c.Sandbox = false
	return c
}
