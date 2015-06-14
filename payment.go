package paypal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	paymentListURL    = "/v1/payments/payment"
	paymentCreateURL  = "/v1/payments/payment"
	paymentExecuteURL = "/v1/payments/payment/%s/execute/"
	paymentPayoutURL  = "/v1/payments/payouts?sync_mode=%v"
)

// PaymentService payment api service
type PaymentService struct {
	client *PayPalClient
}

type PaymentListResponse struct {
	Count    int       `json:"count"`
	Payments []Payment `json:"payments"`
}

/*
List For the /v1/payments/payment resource, the following input parameters can be used.

TODO: not support
	count	Number of items to return. Default is 10 with a maximum value of 20.
	start_id	Resource ID that indicates the starting resource to return. When results are paged, you can use the next_id response value as the start_id to continue with the next set of results.
	start_index	Start index of the resources to be returned. Typically used to jump to a specific position in the resource history based on its order. Example for starting at the second item in a list of results: ?start_index=2
	start_time	Resource creation time as defined in RFC 3339 Section 5.6 that indicates the start of a range of results. Example: start_time=2013-03-06T11:00:00Z
	end_time	Resource creation time that indicates the end of a range of results.
	sort_by	Sort based on create_time or update_time.
	sort_order	Sort based on order of results. Options include asc for ascending order or desc for descending order (default).
*/
func (s PaymentService) List() (PaymentListResponse, error) {
	req, err := http.NewRequest("GET", s.client.URL(paymentListURL), nil)
	if err != nil {
		return PaymentListResponse{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", s.client.Authorization())

	res, err := s.client.Do(req)
	if err != nil {
		return PaymentListResponse{}, err
	}
	defer res.Body.Close()

	outData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return PaymentListResponse{}, err
	}

	if res.StatusCode >= 400 {
		return PaymentListResponse{}, fmt.Errorf("response error %d %s", res.StatusCode, string(outData))
	}

	var r PaymentListResponse
	return r, json.Unmarshal(outData, &r)
}

type PaymentCreateRequest struct {
	Intent       string `json:"intent"` // TODO: enumにする
	Payer        Payer  `json:"payer"`
	RedirectURLs struct {
		CancelURL string `json:"cancel_url"`
		ReturnURL string `json:"return_url"`
	} `json:"redirect_urls"`
	Transactions []Transaction `json:"transactions"`
}

type PaymentCreateResponse struct {
	ID           string        `json:"id"`
	Intent       string        `json:"intent"` // TODO: enumにする
	State        string        `json:"state"`  // TODO: enum?
	Payer        Payer         `json:"payer"`
	Links        []Link        `json:"links"`
	Transactions []Transaction `json:"transactions"`
	CreateTime   string        `json:"create_time"`
	UpdateTime   string        `json:"update_time"`
}

func (r PaymentCreateResponse) LinkByRel(rel Rel) Link {
	for _, l := range r.Links {
		if l.Rel == rel {
			return l
		}
	}
	return Link{}
}

func (s PaymentService) Create(request PaymentCreateRequest) (PaymentCreateResponse, error) {
	inData, err := json.Marshal(request)
	if err != nil {
		return PaymentCreateResponse{}, err
	}

	req, err := http.NewRequest("POST", s.client.URL(paymentCreateURL), bytes.NewBuffer(inData))
	if err != nil {
		return PaymentCreateResponse{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", s.client.Authorization())

	res, err := s.client.Do(req)
	if err != nil {
		return PaymentCreateResponse{}, err
	}
	defer res.Body.Close()

	outData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return PaymentCreateResponse{}, err
	}

	if res.StatusCode >= 400 {
		return PaymentCreateResponse{}, fmt.Errorf("response error %d %s", res.StatusCode, string(outData))
	}

	var result PaymentCreateResponse
	return result, json.Unmarshal(outData, &result)
}

type PaymentExecuteRequest struct {
	PayerID string `json:"payer_id"`
}

func (s PaymentService) Execute(paymentID string, executeReq PaymentExecuteRequest) error {
	inData, err := json.Marshal(executeReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", s.client.URL(fmt.Sprintf(paymentExecuteURL, paymentID)), bytes.NewBuffer(inData))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", s.client.Authorization())

	res, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	outData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	// TODO: あとで
	fmt.Println("payment/execute", string(outData))

	return nil
}

type PaymentPayoutRequest struct {
	Items             []PayoutItem `json:"items"`
	SenderBatchHeader struct {
		EmailSubject  string        `json:"email_subject"`
		RecipientType RecipientType `json:"recipient_type"`
		SenderBatchID string        `json:"sender_batch_id"`
	} `json:"sender_batch_header"`
}

func (s PaymentService) Payout(syncMode bool, payoutReq PaymentPayoutRequest) error {
	inData, err := json.Marshal(payoutReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", s.client.URL(fmt.Sprintf(paymentPayoutURL, syncMode)), bytes.NewBuffer(inData))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", s.client.Authorization())

	res, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	outData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		return fmt.Errorf("response error %d %s", res.StatusCode, string(outData))
	}

	// TODO: あとで
	fmt.Println("payment/execute", string(outData))

	return nil
}
