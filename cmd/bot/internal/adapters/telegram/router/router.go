package router

import (
	"github.com/go-telegram/bot"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/application/commands"
	aboutCommands "github.com/thebeyond-net/control-plane/cmd/bot/internal/core/application/commands/about"
	devicesCommands "github.com/thebeyond-net/control-plane/cmd/bot/internal/core/application/commands/devices"
	planCommands "github.com/thebeyond-net/control-plane/cmd/bot/internal/core/application/commands/plans"
	settingsCommands "github.com/thebeyond-net/control-plane/cmd/bot/internal/core/application/commands/settings"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/ports"
	"github.com/thebeyond-net/control-plane/config"
	sharedPeers "github.com/thebeyond-net/control-plane/internal/core/ports"
	"go.uber.org/zap"
)

type Router struct {
	authUseCase sharedPeers.AuthUseCase
	appLogger   *zap.Logger
}

func NewRouter(
	authUseCase sharedPeers.AuthUseCase,
	appLogger *zap.Logger,
) *Router {
	return &Router{authUseCase, appLogger}
}

func (r *Router) RegisterAllHandlers(
	b *bot.Bot,
	cfg *config.Config,
	telegramBotAdapter ports.Bot,
	userSettingsUseCase sharedPeers.UserSettingsUseCase,
	deviceUseCase sharedPeers.DeviceUseCase,
	nodeUseCase sharedPeers.NodeUseCase,
	yookassa sharedPeers.Invoice,
	telegramStars sharedPeers.Invoice,
	featureFlags sharedPeers.FeatureFlags,
	newsURL, reviewsURL, supportURL string,
	defaultBandwidth int,
) {
	periods := config.ToDomainPeriods(cfg.Periods)
	plans := config.ToDomainPlans(cfg.Plans)
	paymentMethods := config.ToDomain(cfg.PaymentMethods)
	apps := config.ToDomain(cfg.Apps)
	locations := config.ToDomain(cfg.Locations)
	languages := config.ToDomain(cfg.Languages)
	currencies := config.ToDomain(cfg.Currencies)

	menuCommand := commands.NewMenuUseCase(telegramBotAdapter, plans, newsURL, reviewsURL, supportURL)
	connectionCommand := commands.NewConnectionUseCase(telegramBotAdapter)
	appsCommand := commands.NewAppsUseCase(telegramBotAdapter, apps, cfg.AppURLs)

	plansCommand := planCommands.NewUseCase(
		telegramBotAdapter,
		cfg.Bandwidths,
		periods, plans,
		paymentMethods,
		currencies,
		yookassa,
		telegramStars,
		featureFlags,
		defaultBandwidth,
	)
	refCommand := commands.NewRefUseCase(telegramBotAdapter, cfg.BotUsername, supportURL)

	aboutCommand := aboutCommands.NewUseCase(telegramBotAdapter)
	tosCommand := aboutCommands.NewTermsOfService(telegramBotAdapter)
	privacyPolicyCommand := aboutCommands.NewPrivacyPolicy(telegramBotAdapter)
	refundPolicyCommand := aboutCommands.NewRefundPolicy(telegramBotAdapter)

	deviceCommand := devicesCommands.NewUseCase(telegramBotAdapter, deviceUseCase, apps, locations)
	mkDeviceCommand := devicesCommands.NewCreator(telegramBotAdapter, deviceUseCase, nodeUseCase, apps, locations)
	rmDeviceCommand := devicesCommands.NewDestroyer(telegramBotAdapter, deviceUseCase)

	settingsCommand := settingsCommands.NewSettingsUseCase(telegramBotAdapter)
	languageCommand := settingsCommands.NewLanguageUseCase(telegramBotAdapter, userSettingsUseCase, languages)
	currencyCommand := settingsCommands.NewCurrencyUseCase(telegramBotAdapter, userSettingsUseCase, currencies)

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypePrefix, r.CommandWrapper(menuCommand, WithName("start")))

	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "menu", bot.MatchTypeExact, r.CommandWrapper(menuCommand))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "connection", bot.MatchTypeExact, r.CommandWrapper(connectionCommand))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "app", bot.MatchTypePrefix, r.CommandWrapper(appsCommand))

	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "device", bot.MatchTypePrefix, r.CommandWrapper(deviceCommand))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "newdevice", bot.MatchTypePrefix, r.CommandWrapper(mkDeviceCommand))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "deletedevice", bot.MatchTypePrefix, r.CommandWrapper(rmDeviceCommand))

	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "plan", bot.MatchTypePrefix, r.CommandWrapper(plansCommand))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "ref", bot.MatchTypeExact, r.CommandWrapper(refCommand))

	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "about", bot.MatchTypeExact, r.CommandWrapper(aboutCommand))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "tos", bot.MatchTypePrefix, r.CommandWrapper(tosCommand))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "privacypolicy", bot.MatchTypePrefix, r.CommandWrapper(privacyPolicyCommand))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "refundpolicy", bot.MatchTypePrefix, r.CommandWrapper(refundPolicyCommand))

	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "settings", bot.MatchTypeExact, r.CommandWrapper(settingsCommand))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "language", bot.MatchTypePrefix, r.CommandWrapper(languageCommand))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "currency", bot.MatchTypePrefix, r.CommandWrapper(currencyCommand))
}
