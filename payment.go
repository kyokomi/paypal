package paypal

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	paymentListURL = "/v1/payments/payment"
)

type PaymentListResponse struct {
	Count    int       `json:"count"`
	Payments []Payment `json:"payments"`
}

type Payment struct {
	CreateTime   string        `json:"create_time"`
	ID           string        `json:"id"`
	Intent       string        `json:"intent"`
	Links        []Link        `json:"links"`
	Payer        Payer         `json:"payer"`
	State        string        `json:"state"`
	Transactions []Transaction `json:"transactions"`
	UpdateTime   string        `json:"update_time"`
}

type Transaction struct {
	Amount           Amount            `json:"amount"`
	Description      string            `json:"description"`
	RelatedResources []RelatedResource `json:"related_resources"`
}

type RelatedResource struct {
	Sale Sale `json:"sale"`
}

type Amount struct {
	Currency string `json:"currency"`
	Details  struct {
		Subtotal string `json:"subtotal"`
	} `json:"details"`
	Total string `json:"total"`
}

type Sale struct {
	Amount                    Amount `json:"amount"`
	CreateTime                string `json:"create_time"`
	ID                        string `json:"id"`
	Links                     []Link `json:"links"`
	ParentPayment             string `json:"parent_payment"`
	PaymentMode               string `json:"payment_mode"`
	ProtectionEligibility     string `json:"protection_eligibility"`
	ProtectionEligibilityType string `json:"protection_eligibility_type"`
	State                     string `json:"state"`
	TransactionFee            struct {
		Currency string `json:"currency"`
		Value    string `json:"value"`
	} `json:"transaction_fee"`
	UpdateTime string `json:"update_time"`
}

type Payer struct {
	PayerInfo     PayerInfo `json:"payer_info"`
	PaymentMethod string    `json:"payment_method"`
	Status        string    `json:"status"`
}

type PayerInfo struct {
	Email           string          `json:"email"`
	FirstName       string          `json:"first_name"`
	LastName        string          `json:"last_name"`
	PayerID         string          `json:"payer_id"`
	ShippingAddress ShippingAddress `json:"shipping_address"`
}

type ShippingAddress struct {
	City          string `json:"city"`
	CountryCode   string `json:"country_code"`
	Line1         string `json:"line1"`
	PostalCode    string `json:"postal_code"`
	RecipientName string `json:"recipient_name"`
	State         string `json:"state"`
}

type Link struct {
	URL    string `json:"href"`
	Rel    Rel    `json:"rel"`
	Method Method `json:"method"`
}

// Rel component
// The rel component provides the relation type for the URL in question.
type Rel string

// Here are the possible relation types:
const (
	RelSelf          = "self"           // Link to get information about the call itself. For example, the self link in response to a PayPal account payment provides you with more information about the payment resource itself. Similarly, the self link in the response to a refund will provide you with information about the refund that just completed.
	RelParentPayment = "parent_payment" // Link to get information about the originally created payment resource. All payment related calls (/payments/) through the PayPal REST payment API, including refunds, authorized payments, and captured payments involve a parent payment resource.
	RelSale          = "sale"           // Link to get information about a completed sale.
	RelUpdate        = "update"         // Link to execute and complete user-approved PayPal payments.
	RelAuthorization = "authorization"  // Link to look up the original authorized payment for a captured payment.
	RelReauthorize   = "reauthorize"    // Link to reauthorize a previously authorized PayPal payment.
	RelCapture       = "capture"        // Link to capture authorized but uncaptured payments.
	RelVoid          = "void"           // Link to void an authorized payment.
	RelRefund        = "refund"         // Link to refund a completed sale.
	RelDelete        = "delete"         // Link to delete a credit card from the vault.
)

// Method component
// The method component provides the HTTP methods required to interact with the provided HATEOAS URL.
type Method string

// Here are the possible methods:
const (
	MethodPOST     = "POST"     // Use this method to create or act upon resources, including:
	MethodGET      = "GET"      // Use this method to get information about existing resources, including:
	MethodDELETE   = "DELETE"   // Use this method to remove a resource. Currently, you can use this method to delete stored credit cards.
	MethodREDIRECT = "REDIRECT" // This method is actually not an HTTP method. It instead provides a redirect URL where payers are redirected to approve a PayPal account payment.
)

// PaymentService payment api service
type PaymentService struct {
	client *PayPalClient
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

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return PaymentListResponse{}, err
	}

	var r PaymentListResponse
	return r, json.Unmarshal(data, &r)
}
