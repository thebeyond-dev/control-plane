package settings

import (
	"context"
	"errors"

	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/application/commands/helpers"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/application/input"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/interaction"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/ports"
	"github.com/thebeyond-net/control-plane/internal/core/domain"
	sharedPorts "github.com/thebeyond-net/control-plane/internal/core/ports"
	"github.com/thebeyond-net/control-plane/internal/i18n"
)

type LanguageUseCase struct {
	bot          ports.Bot
	userSettings sharedPorts.UserSettingsUseCase
	languages    domain.Items
}

func NewLanguageUseCase(
	bot ports.Bot,
	userSettings sharedPorts.UserSettingsUseCase,
	languages domain.Items,
) ports.CommandHandler {
	return &LanguageUseCase{bot, userSettings, languages}
}

func (uc *LanguageUseCase) Execute(ctx context.Context, msg input.Message, user domain.User) error {
	if len(msg.Args) > 0 {
		return uc.applyChoice(ctx, msg, user)
	}
	return uc.presentOptions(ctx, msg, user)
}

func (uc *LanguageUseCase) applyChoice(ctx context.Context, msg input.Message, user domain.User) error {
	languageCode := msg.Args[0]
	language, ok := uc.languages.Get(languageCode)
	if !ok {
		return errors.New("language not found")
	}

	if err := uc.userSettings.SetLanguage(ctx, user.ID, languageCode); err != nil {
		return err
	}

	user.LanguageCode = languageCode
	text := i18n.Get(user.LanguageCode, "LanguageSelected", map[string]any{
		"Name": language.Name,
	}, nil)

	return uc.bot.ShowNotification(ctx, msg.InteractionID, text, false)
}

func (uc *LanguageUseCase) presentOptions(ctx context.Context, msg input.Message, user domain.User) error {
	const rowWidth = 3

	settingsBtn := i18n.Get(user.LanguageCode, "BackButton", nil, nil)

	text := i18n.Get(user.LanguageCode, "SelectLanguage", nil, nil)
	markup := helpers.BuildSelectionMarkup(
		uc.languages,
		"language",
		rowWidth,
	)

	markup.Next().AddButton(interaction.NewButton().
		Text(settingsBtn).
		CallbackData("settings").
		Build())

	return uc.bot.NewMessage(msg.ChatID, text).
		WithReplyMarkup(markup.Build()).
		Edit(ctx, msg.ID)
}
