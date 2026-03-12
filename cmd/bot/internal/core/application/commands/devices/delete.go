package devices

import (
	"context"

	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/application/input"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/interaction"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/ports"
	"github.com/thebeyond-net/control-plane/internal/core/domain"
	sharedPorts "github.com/thebeyond-net/control-plane/internal/core/ports"
	"github.com/thebeyond-net/control-plane/internal/i18n"
)

type Destroyer struct {
	bot           ports.Bot
	deviceUseCase sharedPorts.DeviceUseCase
}

func NewDestroyer(
	bot ports.Bot,
	deviceUseCase sharedPorts.DeviceUseCase,
) ports.CommandHandler {
	return &Destroyer{bot, deviceUseCase}
}

func (uc *Destroyer) Execute(ctx context.Context, msg input.Message, user domain.User) error {
	if err := uc.deviceUseCase.Delete(ctx, user.ID, msg.Args[0]); err != nil {
		return err
	}

	backBtn := i18n.Get(user.LanguageCode, "BackButton", nil, nil)

	text := i18n.Get(user.LanguageCode, "DeviceDeleted", nil, nil)
	markup := interaction.NewReplyMarkup()

	markup.Next().AddButton(interaction.NewButton().
		Text(backBtn).
		CallbackData("device").
		Build())

	return uc.bot.NewMessage(msg.ChatID, text).
		WithReplyMarkup(markup.Build()).
		Edit(ctx, msg.ID)
}
