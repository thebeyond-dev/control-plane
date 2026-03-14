package dto

type PaymentRequest struct {
	Amount            *AmountRequest            `json:"amount,omitempty"`
	PaymentMethodData *PaymentMethodDataRequest `json:"payment_method_data,omitempty"`
	Confirmation      *ConfirmationRequest      `json:"confirmation,omitempty"`
	Capture           bool                      `json:"capture,omitempty"`
	Description       string                    `json:"description,omitempty"`
	Metadata          map[string]any            `json:"metadata,omitempty"`
	SavePaymentMethod string                    `json:"save_payment_method,omitempty"`
	Test              bool                      `json:"test,omitempty"`
}

type AmountRequest struct {
	Value    string `json:"value,omitempty"`
	Currency string `json:"currency,omitempty"`
}

type PaymentMethodDataRequest struct {
	Type string `json:"type,omitempty"`
}

type ConfirmationRequest struct {
	Type      string `json:"type,omitempty"`
	ReturnURL string `json:"return_url,omitempty"`
}
