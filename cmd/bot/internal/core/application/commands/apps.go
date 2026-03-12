package commands

import (
	"context"

	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/application/input"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/interaction"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/ports"
	"github.com/thebeyond-net/control-plane/internal/core/domain"
	"github.com/thebeyond-net/control-plane/internal/i18n"
)

type AppsUseCase struct {
	bot     ports.Bot
	apps    domain.Items
	appURLs map[string]string
}

func NewAppsUseCase(
	bot ports.Bot,
	apps domain.Items,
	appURLs map[string]string,
) ports.CommandHandler {
	return &AppsUseCase{bot, apps, appURLs}
}

func (uc *AppsUseCase) Execute(ctx context.Context, msg input.Message, user domain.User) error {
	text := i18n.Get(user.LanguageCode, "SelectDevice", nil, nil)
	markup := uc.buildReplyMarkup(user.LanguageCode)

	return uc.bot.NewMessage(msg.ChatID, text).
		WithReplyMarkup(markup).
		Edit(ctx, msg.ID)
}

func (uc *AppsUseCase) buildReplyMarkup(languageCode string) interaction.InlineKeyboardMarkup {
	backBtn := i18n.Get(languageCode, "BackButton", nil, nil)

	markup := interaction.NewReplyMarkup()
	for i, device := range uc.apps.All() {
		if i%2 == 0 {
			markup.Next()
		}

		markup.AddButton(interaction.NewButton().
			Text(device.Name).
			URL(uc.appURLs[device.Code]).
			IconCustomEmojiID(device.Icon).
			Build(),
		)
	}

	markup.Next().AddButton(interaction.NewButton().
		Text(backBtn).
		CallbackData("connection").
		IconCustomEmojiID(IconBack).
		Build())

	return markup.Build()
}
