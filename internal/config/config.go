package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	LogLevel      string `env:"LOG_LEVEL" env-default:"info"`
	DatabaseURI   string `env:"DB_URI" env-required:"true"`
	TelegramToken string `env:"TG_TOKEN" env-required:"true"`
	AuthToken     string `env:"AUTH_TOKEN"`
}

func New() (cfg Config, err error) {
	return cfg, cleanenv.ReadConfig(".env", &cfg)
}
