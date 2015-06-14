package paypal

type PaymentMethod string

const (
	PaymentMethodPayPal = "paypal"
)

type PayIntent string

const (
	IntentSale = "sale"
)

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
	RelApprovalURL   = "approval_url"   // Link to approval_url.
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

type RecipientType string

const (
	RECIPIENT_EMAIL     = "EMAIL"
	RECIPIENT_PHONE     = "PHONE"
	RECIPIENT_PAYPAL_ID = "PAYPAL_ID"
)

type Payer struct {
	PayerInfo struct {
		PayerID         string `json:"payer_id"`
		FirstName       string `json:"first_name"`
		LastName        string `json:"last_name"`
		Email           string `json:"email"`
		ShippingAddress struct {
			City          string `json:"city"`
			CountryCode   string `json:"country_code"`
			Line1         string `json:"line1"`
			PostalCode    string `json:"postal_code"`
			RecipientName string `json:"recipient_name"`
			State         string `json:"state"`
		} `json:"shipping_address"`
	} `json:"payer_info"`
	PaymentMethod PaymentMethod `json:"payment_method"`
	Status        string        `json:"status"` // TODO: enum?
}

type Link struct {
	URL    string `json:"href"`
	Rel    Rel    `json:"rel"`
	Method Method `json:"method"`
}

type Payment struct {
	ID           string        `json:"id"`
	Intent       PayIntent     `json:"intent"`
	Links        []Link        `json:"links"`
	Payer        Payer         `json:"payer"`
	State        string        `json:"state"` // TODO: enum
	Transactions []Transaction `json:"transactions"`
	CreateTime   string        `json:"create_time"`
	UpdateTime   string        `json:"update_time"`
}

type Transaction struct {
	Amount           Amount `json:"amount"`
	Description      string `json:"description"`
	RelatedResources []struct {
		Sale Sale `json:"sale"`
	} `json:"related_resources"`
}

type Amount struct {
	Currency string `json:"currency"`
	Details  struct {
		Subtotal string `json:"subtotal,omitempty"`
	} `json:"details,omitempty"`
	Total string `json:"total"`
}

// Sale This object defines a sale object.
type Sale struct {
	ID                        string `json:"id"`
	Amount                    Amount `json:"amount"`
	State                     string `json:"state"` // TODO: enum : (pending completed refunded) or partially_refunded
	Links                     []Link `json:"links"`
	ParentPayment             string `json:"parent_payment"`
	PaymentMode               string `json:"payment_mode"` // TODO: enum? INSTANT_TRANSFER, MANUAL_BANK_TRANSFER, DELAYED_TRANSFER, or ECHECK.
	ProtectionEligibility     string `json:"protection_eligibility"`
	ProtectionEligibilityType string `json:"protection_eligibility_type"`
	TransactionFee            struct {
		Currency string `json:"currency"`
		Value    string `json:"value"`
	} `json:"transaction_fee"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}

// PayoutItem Sender-created description of a payout to a single recipient.
type PayoutItem struct {
	RecipientType RecipientType `json:"recipient_type"`
	Amount        struct {
		Currency string `json:"currency"`
		Value    string `json:"value"`
	} `json:"amount"`
	Note         string `json:"note,omitempty"`
	Receiver     string `json:"receiver"`
	SenderItemID string `json:"sender_item_id,omitempty"`
}
