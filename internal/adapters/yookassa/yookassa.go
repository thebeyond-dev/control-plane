package yookassa

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/thebeyond-net/control-plane/internal/adapters/yookassa/dto"
	"github.com/thebeyond-net/control-plane/internal/core/domain"
	"github.com/thebeyond-net/control-plane/internal/core/ports"
	"github.com/thebeyond-net/control-plane/internal/i18n"
)

const defaultEndpoint = "https://api.yookassa.ru/v3"

type Adapter struct {
	httpClient *http.Client
	endpoint   string
	shopID     string
	secretKey  string
	returnURL  string
}

func New(shopID, secretKey, returnURL string) ports.Invoice {
	return &Adapter{
		httpClient: http.DefaultClient,
		endpoint:   defaultEndpoint,
		shopID:     shopID,
		secretKey:  secretKey,
		returnURL:  returnURL,
	}
}

func (a *Adapter) NewPayment(
	ctx context.Context,
	user domain.User,
	currency string,
	devices, bandwidth, days int,
	price float64,
) (string, error) {
	devicesCount := i18n.Get(user.LanguageCode, "DevicesCount", map[string]any{
		"Count": devices,
	}, devices)

	daysCount := i18n.Get(user.LanguageCode, "DaysCount", map[string]any{
		"Count": days,
	}, days)

	description := i18n.Get(user.LanguageCode, "InvoiceDescription", map[string]any{
		"Devices":   devicesCount,
		"Bandwidth": bandwidth,
		"Days":      daysCount,
	}, nil)

	jsonPayload, err := json.Marshal(dto.PaymentRequest{
		Amount: &dto.AmountRequest{
			Value:    fmt.Sprintf("%.2f", price),
			Currency: strings.ToUpper(currency),
		},
		Confirmation: &dto.ConfirmationRequest{
			Type:      "redirect",
			ReturnURL: a.returnURL,
		},
		Capture:     true,
		Description: description,
		Metadata: map[string]any{
			"user_id":   user.ID,
			"devices":   devices,
			"bandwidth": bandwidth,
			"days":      days,
		},
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		a.endpoint+"/payments",
		bytes.NewReader(jsonPayload),
	)
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(a.shopID, a.secretKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotence-Key", uuid.NewString())

	res, err := a.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	data := dto.PaymentResponse{}
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return "", err
	}

	if res.StatusCode != 200 {
		return "", fmt.Errorf("status code %d", res.StatusCode)
	}

	if data.Error != nil {
		return "", fmt.Errorf("%s: %s", data.Error.Code, data.Error.Description)
	}

	return data.Confirmation.ConfirmationURL, nil
}
