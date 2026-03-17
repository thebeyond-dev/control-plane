package input

import "time"

type NewbieConfig struct {
	Devices       int
	Bandwidth     int
	TrialDuration time.Duration
	LanguageCode  string
	CurrencyCode  string
}

type Login struct {
	ReferrerID *string
}

type FeatureContext struct {
	UserID        string
	SessionID     string
	RemoteAddress string
	Environment   string
	AppName       string
	CurrentTime   string
	Properties    map[string]string
}
