package about

import (
	"context"

	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/application/input"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/interaction"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/ports"
	"github.com/thebeyond-net/control-plane/internal/core/domain"
	"github.com/thebeyond-net/control-plane/internal/i18n"
)

var (
	IconTermsOfService = "5269557261646206128"
	IconPrivacyPolicy  = "5244460818547904859"
	IconRefundPolicy   = "5244606924745380898"
	IconBack           = "5242383235492648578"
)

type UseCase struct {
	bot ports.Bot
}

func NewUseCase(bot ports.Bot) ports.CommandHandler {
	return &UseCase{bot}
}

func (uc *UseCase) Execute(ctx context.Context, msg input.Message, user domain.User) error {
	text := i18n.Get(user.LanguageCode, "SelectAction", nil, nil)
	markup := uc.buildReplyMarkup(user.LanguageCode)

	return uc.bot.NewMessage(msg.ChatID, text).
		WithReplyMarkup(markup).
		Edit(ctx, msg.ID)
}

func (uc *UseCase) buildReplyMarkup(languageCode string) interaction.InlineKeyboardMarkup {
	tosBtn := i18n.Get(languageCode, "TermsOfServiceButton", nil, nil)
	privacyPolicyBtn := i18n.Get(languageCode, "PrivacyPolicyButton", nil, nil)
	refundPolicyBtn := i18n.Get(languageCode, "RefundPolicyButton", nil, nil)
	backBtn := i18n.Get(languageCode, "BackButton", nil, nil)

	markup := interaction.NewReplyMarkup()

	markup.Next().AddButton(interaction.NewButton().
		Text(tosBtn).
		CallbackData("tos").
		IconCustomEmojiID(IconTermsOfService).
		Build())

	markup.Next().AddButton(interaction.NewButton().
		Text(privacyPolicyBtn).
		CallbackData("privacypolicy").
		IconCustomEmojiID(IconPrivacyPolicy).
		Build())

	markup.Next().AddButton(interaction.NewButton().
		Text(refundPolicyBtn).
		CallbackData("refundpolicy").
		IconCustomEmojiID(IconRefundPolicy).
		Build())

	markup.Next().AddButton(interaction.NewButton().
		Text(backBtn).
		CallbackData("menu").
		IconCustomEmojiID(IconBack).
		Build())

	return markup.Build()
}
