package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	MigrationsPath             string `env:"MIGRATIONS_PATH" env-default:"./migrations"`
	DatabaseURL                string `env:"DATABASE_URL" env-required:"true"`
	TelegramBotToken           string `env:"TELEGRAM_BOT_TOKEN" env-required:"true"`
	TelegramWebhookSecretToken string `env:"TELEGRAM_WEBHOOK_SECRET_TOKEN" env-required:"true"`
	Port                       string `env:"PORT" env-default:"2000"`
	LogLevel                   string `env:"LOG_LEVEL" env-default:"info"`
	AuthSecret                 string `env:"AUTH_SECRET" env-required:"true"`
	YookassaShopID             string `env:"YOOKASSA_SHOP_ID" env-required:"true"`
	YookassaSecretKey          string `env:"YOOKASSA_SECRET_KEY" env-required:"true"`
	KassaReturnURL             string `env:"KASSA_RETURN_URL" env-required:"true"`

	BotUsername string `yaml:"bot_username"`

	DefaultBandwidth int         `yaml:"defaultBandwidth"`
	Newbie           NewbieDTO   `yaml:"newbie"`
	PriceListPath    string      `yaml:"price_list"`
	Bandwidths       []int       `yaml:"bandwidths"`
	Periods          []PeriodDTO `yaml:"periods"`
	Plans            []PlanDTO   `yaml:"plans"`

	PaymentMethods map[string]ItemDTO `yaml:"payment_methods"`
	AppURLs        map[string]string  `yaml:"app_urls"`
	Apps           map[string]ItemDTO `yaml:"apps"`
	Locations      map[string]ItemDTO `yaml:"locations"`
	Languages      map[string]ItemDTO `yaml:"languages"`
	Currencies     map[string]ItemDTO `yaml:"currencies"`
}

type NewbieDTO struct {
	Devices       int    `yaml:"devices"`
	Bandwidth     int    `yaml:"bandwidth"`
	TrialDuration string `yaml:"trial_duration"`
	LanguageCode  string `yaml:"language_code"`
	CurrencyCode  string `yaml:"currency_code"`
}

type PeriodDTO struct {
	Days     int    `yaml:"days"`
	Discount int    `yaml:"discount"`
	NameKey  string `yaml:"name_key"`
}

type PlanDTO struct {
	Devices  int                        `yaml:"devices"`
	Discount int                        `yaml:"discount"`
	Prices   map[int]map[string]float64 `yaml:"prices"`
}

type ItemDTO struct {
	Name   string `yaml:"name"`
	Icon   string `yaml:"icon"`
	Symbol string `yaml:"symbol"`
}

func Load() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig("config.yaml", &cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("env error: %w", err)
	}

	if cfg.PriceListPath != "" {
		if err := cleanenv.ReadConfig(cfg.PriceListPath, &cfg); err != nil {
			return nil, fmt.Errorf("error loading price list (%s): %w", cfg.PriceListPath, err)
		}
	}

	return &cfg, nil
}
