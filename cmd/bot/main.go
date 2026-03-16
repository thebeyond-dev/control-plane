package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/go-telegram/bot"
	"github.com/shizumico/arcane/pkg/logger"
	telegramBot "github.com/thebeyond-net/control-plane/cmd/bot/internal/adapters/telegram/bot"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/adapters/telegram/launcher"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/adapters/telegram/middlewares"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/adapters/telegram/router"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/adapters/telegram/webhook"
	"github.com/thebeyond-net/control-plane/config"
	"github.com/thebeyond-net/control-plane/internal/adapters/nodepool"
	"github.com/thebeyond-net/control-plane/internal/adapters/repositories/postgres/devices"
	"github.com/thebeyond-net/control-plane/internal/adapters/repositories/postgres/nodes"
	"github.com/thebeyond-net/control-plane/internal/adapters/repositories/postgres/subscription"
	"github.com/thebeyond-net/control-plane/internal/adapters/repositories/postgres/users"
	"github.com/thebeyond-net/control-plane/internal/adapters/tgnotifier"
	telegramStars "github.com/thebeyond-net/control-plane/internal/adapters/tgstars"
	"github.com/thebeyond-net/control-plane/internal/adapters/yookassa"
	alertApplication "github.com/thebeyond-net/control-plane/internal/core/application/alert"
	authApplication "github.com/thebeyond-net/control-plane/internal/core/application/auth"
	billingApplication "github.com/thebeyond-net/control-plane/internal/core/application/billing"
	deviceApplication "github.com/thebeyond-net/control-plane/internal/core/application/devices"
	"github.com/thebeyond-net/control-plane/internal/core/application/input"
	nodeApplication "github.com/thebeyond-net/control-plane/internal/core/application/nodes"
	userSettingsApplication "github.com/thebeyond-net/control-plane/internal/core/application/user_settings"
	"github.com/thebeyond-net/control-plane/internal/i18n"
	"github.com/thebeyond-net/control-plane/pkg/postgres"
	"go.uber.org/zap"
)

const (
	defaultNewsURL    = "https://t.me/thebeyondnews"
	defaultReviewsURL = "https://t.me/thebeyondreviews"
	defaultSupportURL = "https://t.me/thebeyondnews?direct"
)

func main() {
	newsURL := os.Getenv("NEWS_URL")
	if newsURL == "" {
		newsURL = defaultNewsURL
	}

	reviewsURL := os.Getenv("REVIEWS_URL")
	if reviewsURL == "" {
		reviewsURL = defaultReviewsURL
	}

	supportURL := os.Getenv("SUPPORT_URL")
	if supportURL == "" {
		supportURL = defaultSupportURL
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v\n", err)
	}

	appLogger, err := logger.Init(cfg.LogLevel)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v\n", err)
	}

	defer appLogger.Sync()

	if err := i18n.Init(LocaleFS, "assets/locales", appLogger); err != nil {
		appLogger.Fatal("Failed to init i18n", zap.Error(err))
	}

	dbCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	pool, err := postgres.New(dbCtx, cfg.DatabaseURL)
	if err != nil {
		appLogger.Fatal("Failed to connect to database", zap.Error(err))
	}

	botContainer := new(BotContainer)

	userRepo := users.NewRepository(pool, appLogger)
	deviceRepo := devices.NewRepository(pool)
	nodeRepo := nodes.NewRepository(pool)
	nodePool := nodepool.New(nodeRepo, cfg.AuthSecret)
	subscriptionRepo := subscription.NewRepository(pool)
	telegramNotifier := tgnotifier.New(botContainer)

	trialDuration, err := time.ParseDuration(cfg.Newbie.TrialDuration)
	if err != nil {
		appLogger.Fatal("Failed to parse trial duration", zap.Error(err))
	}

	authUseCase := authApplication.NewInteractor(userRepo, input.NewbieConfig{
		Devices:       cfg.Newbie.Devices,
		Bandwidth:     cfg.Newbie.Bandwidth,
		TrialDuration: trialDuration,
		LanguageCode:  cfg.Newbie.LanguageCode,
		CurrencyCode:  cfg.Newbie.CurrencyCode,
	})
	userSettingsUseCase := userSettingsApplication.NewInteractor(userRepo)
	deviceUseCase := deviceApplication.NewInteractor(deviceRepo, nodePool)
	nodeUseCase := nodeApplication.NewInteractor(nodeRepo, nodePool)
	billingUseCase := billingApplication.NewInteractor(subscriptionRepo)
	telegramAlertUseCase := alertApplication.New(telegramNotifier, config.ToDomain(cfg.Currencies))

	handler := webhook.NewHandler()
	paymentMiddleware := middlewares.NewPaymentHandler(
		authUseCase,
		billingUseCase,
		telegramAlertUseCase,
		appLogger,
	)

	botInstance, err := bot.New(
		cfg.TelegramBotToken,
		bot.WithDefaultHandler(handler.HandleUpdate),
		bot.WithWebhookSecretToken(cfg.TelegramWebhookSecretToken),
		bot.WithMiddlewares(paymentMiddleware.Handle),
	)
	if err != nil {
		appLogger.Fatal("Failed to create bot", zap.Error(err))
	}

	botContainer.Instance = botInstance

	yookassa := yookassa.New(
		cfg.YookassaShopID,
		cfg.YookassaSecretKey,
		cfg.KassaReturnURL,
	)

	telegramBot := telegramBot.New(botContainer.Instance, appLogger)
	telegramStars := telegramStars.New(botContainer.Instance)
	router := router.NewRouter(authUseCase, appLogger)

	router.RegisterAllHandlers(
		botContainer.Instance,
		telegramBot,
		userSettingsUseCase,
		deviceUseCase,
		nodeUseCase,
		yookassa,
		telegramStars,
		cfg,
		newsURL,
		reviewsURL,
		supportURL,
		cfg.DefaultBandwidth,
	)

	var ln launcher.Launcher

	switch os.Getenv("BOT_MODE") {
	case "webhook":
		ln = launcher.NewWebHook(botContainer.Instance, appLogger, cfg.Port)
	default:
		ln = launcher.NewLongPoll(botContainer.Instance)
	}

	ln.Launch(ctx)
}
