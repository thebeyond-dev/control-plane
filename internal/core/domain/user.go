package domain

import (
	"fmt"
	"time"
)

type User struct {
	ID string
	Identity
	Subscription
	Devices      int
	Bandwidth    int
	LanguageCode string
	CurrencyCode string
	CreatedAt    time.Time
}

func (user User) GetFormattedBandwidth() string {
	if user.Bandwidth >= 1000 {
		return "1 Gbps"
	}
	return fmt.Sprintf("%d Mbps", user.Bandwidth)
}
