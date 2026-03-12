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

var IconDelete = "5269688764954875572"

type UseCase struct {
	bot           ports.Bot
	deviceUseCase sharedPorts.DeviceUseCase
	apps          domain.Items
	locations     domain.Items
}

func NewUseCase(
	bot ports.Bot,
	deviceUseCase sharedPorts.DeviceUseCase,
	apps domain.Items,
	locations domain.Items,
) ports.CommandHandler {
	return &UseCase{bot, deviceUseCase, apps, locations}
}

func (uc *UseCase) Execute(ctx context.Context, msg input.Message, user domain.User) error {
	if len(msg.Args) > 0 {
		return uc.renderDeviceActions(ctx, msg, user)
	}
	return uc.renderDeviceList(ctx, msg, user)
}

func (uc *UseCase) renderDeviceList(ctx context.Context, msg input.Message, user domain.User) error {
	createBtn := i18n.Get(user.LanguageCode, "CreateButton", nil, nil)
	backBtn := i18n.Get(user.LanguageCode, "BackButton", nil, nil)

	text := i18n.Get(user.LanguageCode, "SelectDevice", nil, nil)
	markup := interaction.NewReplyMarkup()

	devices, err := uc.deviceUseCase.List(ctx, user.ID)
	if err != nil {
		return err
	}

	if len(devices) < user.Devices {
		markup.Next().AddButton(interaction.NewButton().
			Text(createBtn).
			CallbackData("newdevice").
			Build())
	}

	for i, device := range devices {
		if i%2 == 0 {
			markup.Next()
		}

		app, ok := uc.apps.Get(device.Name)
		if !ok {
			continue
		}

		location, ok := uc.locations.Get(device.NodeID[:2])
		if !ok {
			continue
		}

		markup.AddButton(interaction.NewButton().
			Text(app.Name).
			CallbackData("device " + device.PublicKey).
			IconCustomEmojiID(location.Icon).
			Build())
	}

	markup.Next().AddButton(interaction.NewButton().
		Text(backBtn).
		CallbackData("connection").
		Build())

	return uc.bot.NewMessage(msg.ChatID, text).
		WithReplyMarkup(markup.Build()).
		Edit(ctx, msg.ID)
}

func (uc *UseCase) renderDeviceActions(ctx context.Context, msg input.Message, user domain.User) error {
	deleteBtn := i18n.Get(user.LanguageCode, "DeleteButton", nil, nil)
	backBtn := i18n.Get(user.LanguageCode, "BackButton", nil, nil)

	text := i18n.Get(user.LanguageCode, "SelectAction", nil, nil)
	markup := interaction.NewReplyMarkup()

	markup.Next().AddButton(interaction.NewButton().
		Text(deleteBtn).
		CallbackData("deletedevice " + msg.Args[0]).
		IconCustomEmojiID(IconDelete).
		Build())

	markup.Next().AddButton(interaction.NewButton().
		Text(backBtn).
		CallbackData("device").
		Build())

	return uc.bot.NewMessage(msg.ChatID, text).
		WithReplyMarkup(markup.Build()).
		Edit(ctx, msg.ID)
}
