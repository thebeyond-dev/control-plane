package dto

import "time"

type APIError struct {
	Type        string `json:"type,omitempty"`
	ID          string `json:"id,omitempty"`
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
	Parameter   string `json:"parameter,omitempty"`
}

type PaymentResponse struct {
	Error         *APIError
	ID            string                `json:"id,omitempty"`
	Status        string                `json:"status,omitempty"`
	Paid          bool                  `json:"paid,omitempty"`
	Amount        AmountResponse        `json:"amount,omitempty"`
	Confirmation  ConfirmationResponse  `json:"confirmation,omitempty"`
	CreatedAt     time.Time             `json:"created_at,omitempty"`
	Description   string                `json:"description,omitempty"`
	Metadata      map[string]any        `json:"metadata,omitempty"`
	PaymentMethod PaymentMethodResponse `json:"payment_method,omitempty"`
	Recipient     RecipientResponse     `json:"recipient,omitempty"`
	Refundable    bool                  `json:"refundable,omitempty"`
	Test          bool                  `json:"test,omitempty"`
}

type AmountResponse struct {
	Value    string `json:"value,omitempty"`
	Currency string `json:"currency,omitempty"`
}

type ConfirmationResponse struct {
	Type            string `json:"type,omitempty"`
	ReturnURL       string `json:"return_url,omitempty"`
	ConfirmationURL string `json:"confirmation_url,omitempty"`
}

type PaymentMethodResponse struct {
	Type  string `json:"type,omitempty"`
	Id    string `json:"id,omitempty"`
	Saved bool   `json:"saved,omitempty"`
}

type RecipientResponse struct {
	AccountID string `json:"account_id,omitempty"`
	GatewayID string `json:"gateway_id,omitempty"`
}
