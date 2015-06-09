package main

import (
	"log"
	"os"

	"github.com/kyokomi/paypal"
	"github.com/k0kubun/pp"
)

func main() {
	log.SetFlags(log.Llongfile)

	clientID := os.Getenv("PAYPAL_CLIENTID")
	if clientID == "" {
		log.Fatalln("get env")
	}
	secret := os.Getenv("PAYPAL_SECRET")
	if secret == "" {
		log.Fatalln("get env")
	}

	opts := paypal.NewOptions(clientID, secret)
	opts.Sandbox = true
	client := paypal.NewClient(opts)

	req := paypal.PaymentCreateRequest{}
	req.Intent = "sale"
	req.Payer.PaymentMethod = "paypal"
	req.RedirectURLs.CancelURL = "http://localhost:8000//paypal/payment/cancel"
	req.RedirectURLs.ReturnURL = "http://localhost:8000/paypal/payment/execute"
	req.Transactions = []paypal.Transaction{
		{
			Amount: paypal.Amount{
				Total:    "9.99",
				Currency: "USD",
			},
			Description: "example paypal",
		},
	}

	adminToken, err := client.OAuth2.GetToken()
	if err != nil {
		log.Println(err)
	}
	pp.Println(adminToken)

	client.Admin = adminToken
	if response, err := client.Payment.Create(req); err != nil {
		log.Println(err)
	} else {
		pp.Println(response)
	}
}
