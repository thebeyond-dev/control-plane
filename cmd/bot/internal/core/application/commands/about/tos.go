package about

import (
	"context"

	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/application/input"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/interaction"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/ports"
	"github.com/thebeyond-net/control-plane/internal/core/domain"
	"github.com/thebeyond-net/control-plane/internal/i18n"
)

type TermsOfService struct {
	bot ports.Bot
}

func NewTermsOfService(bot ports.Bot) ports.CommandHandler {
	return &TermsOfService{bot}
}

func (uc *TermsOfService) Execute(ctx context.Context, msg input.Message, user domain.User) error {
	text := i18n.Get(user.LanguageCode, "TermsOfService", nil, nil)
	markup := uc.buildReplyMarkup(user.LanguageCode)

	return uc.bot.NewMessage(msg.ChatID, text).
		WithReplyMarkup(markup).
		Edit(ctx, msg.ID)
}

func (uc *TermsOfService) buildReplyMarkup(languageCode string) interaction.InlineKeyboardMarkup {
	backBtn := i18n.Get(languageCode, "BackButton", nil, nil)

	markup := interaction.NewReplyMarkup()

	markup.Next()
	markup.AddButton(interaction.NewButton().
		Text(backBtn).
		CallbackData("about").
		Build())

	return markup.Build()
}
