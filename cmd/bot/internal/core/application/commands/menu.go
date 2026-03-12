package commands

import (
	"context"

	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/application/input"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/interaction"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/ports"
	"github.com/thebeyond-net/control-plane/internal/core/domain"
	"github.com/thebeyond-net/control-plane/internal/i18n"
)

var (
	IconConnection = "5300813830308276410"
	IconPlans      = "5244828222935306011"
	IconRef        = "5244653340456947899"
	IconNews       = "5269625714834977517"
	IconReviews    = "5244858867526965394"
	IconSupport    = "5242440092269712662"
	IconAbout      = "5242281857084594799"
	IconSettings   = "5301111368462669023"
	IconBack       = "5242383235492648578"
)

type MenuUseCase struct {
	bot        ports.Bot
	plans      []domain.Plan
	newsURL    string
	reviewsURL string
	supportURL string
}

func NewMenuUseCase(
	bot ports.Bot,
	plans []domain.Plan,
	newsURL,
	reviewsURL,
	supportURL string,
) ports.CommandHandler {
	return &MenuUseCase{
		bot:        bot,
		plans:      plans,
		newsURL:    newsURL,
		reviewsURL: reviewsURL,
		supportURL: supportURL,
	}
}

func (uc *MenuUseCase) Execute(ctx context.Context, msg input.Message, user domain.User) error {
	text := i18n.Get(user.LanguageCode, "Menu", map[string]any{
		"Remaining": i18n.FormatRemaining(user.LanguageCode, user.ExpiresAt),
		"Devices":   user.Devices,
		"Bandwidth": user.GetFormattedBandwidth(),
	}, nil)

	if err := uc.bot.DeleteMessage(ctx, msg.ChatID, msg.ID); err != nil {
		return err
	}

	return uc.bot.NewMessage(msg.ChatID, text).
		WithReplyMarkup(uc.buildReplyMarkup(user.LanguageCode)).
		Send(ctx)
}

func (uc *MenuUseCase) buildReplyMarkup(languageCode string) interaction.InlineKeyboardMarkup {
	connectionBtn := i18n.Get(languageCode, "ConnectionButton", nil, nil)
	plansBtn := i18n.Get(languageCode, "PlansButton", nil, nil)
	refBtn := i18n.Get(languageCode, "RefButton", nil, nil)
	newsBtn := i18n.Get(languageCode, "NewsButton", nil, nil)
	reviewsBtn := i18n.Get(languageCode, "ReviewsButton", nil, nil)
	supportBtn := i18n.Get(languageCode, "SupportButton", nil, nil)
	aboutBtn := i18n.Get(languageCode, "AboutButton", nil, nil)
	settingsBtn := i18n.Get(languageCode, "SettingsButton", nil, nil)

	markup := interaction.NewReplyMarkup()

	markup.Next().AddButton(interaction.NewButton().
		Text(connectionBtn).
		CallbackData("connection").
		IconCustomEmojiID(IconConnection).
		Build())

	markup.Next().AddButton(interaction.NewButton().
		Text(plansBtn).
		CallbackData("plan").
		IconCustomEmojiID(IconPlans).
		Build())
	markup.AddButton(interaction.NewButton().
		Text(refBtn).
		CallbackData("ref").
		IconCustomEmojiID(IconRef).
		Build())

	markup.Next().AddButton(interaction.NewButton().
		Text(newsBtn).
		URL(uc.newsURL).
		IconCustomEmojiID(IconNews).
		Build())
	markup.AddButton(interaction.NewButton().
		Text(reviewsBtn).
		URL(uc.reviewsURL).
		IconCustomEmojiID(IconReviews).
		Build())

	markup.Next().AddButton(interaction.NewButton().
		Text(supportBtn).
		URL(uc.supportURL).
		IconCustomEmojiID(IconSupport).
		Build())

	markup.Next().AddButton(interaction.NewButton().
		Text(aboutBtn).
		CallbackData("about").
		IconCustomEmojiID(IconAbout).
		Build())

	markup.Next().AddButton(interaction.NewButton().
		Text(settingsBtn).
		CallbackData("settings").
		IconCustomEmojiID(IconSettings).
		Build())

	return markup.Build()
}
