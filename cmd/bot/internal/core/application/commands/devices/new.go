package devices

import (
	"context"
	"fmt"
	"strings"

	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/application/commands/helpers"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/application/input"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/interaction"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/ports"
	"github.com/thebeyond-net/control-plane/internal/core/domain"
	sharedPorts "github.com/thebeyond-net/control-plane/internal/core/ports"
	"github.com/thebeyond-net/control-plane/internal/i18n"
)

type Creator struct {
	bot           ports.Bot
	deviceUseCase sharedPorts.DeviceUseCase
	nodeUseCase   sharedPorts.NodeUseCase
	apps          domain.Items
	locations     domain.Items
}

func NewCreator(
	bot ports.Bot,
	deviceUseCase sharedPorts.DeviceUseCase,
	nodeUseCase sharedPorts.NodeUseCase,
	apps domain.Items,
	locations domain.Items,
) ports.CommandHandler {
	return &Creator{bot, deviceUseCase, nodeUseCase, apps, locations}
}

func (uc *Creator) Execute(ctx context.Context, msg input.Message, user domain.User) error {
	switch len(msg.Args) {
	case 1:
		return uc.selectLocation(ctx, msg, user)
	case 2:
		return uc.createDevice(ctx, msg, user)
	default:
		return uc.selectDeviceType(ctx, msg, user)
	}
}

func (uc *Creator) selectDeviceType(ctx context.Context, msg input.Message, user domain.User) error {
	const rowWidth = 2
	text := i18n.Get(user.LanguageCode, "SelectDevice", nil, nil)
	backBtn := i18n.Get(user.LanguageCode, "BackButton", nil, nil)

	markup := helpers.BuildSelectionMarkup(uc.apps, "newdevice", rowWidth)
	markup.Next().AddButton(interaction.NewButton().
		Text(backBtn).
		CallbackData("device").
		Build())

	return uc.bot.NewMessage(msg.ChatID, text).
		WithReplyMarkup(markup.Build()).
		Edit(ctx, msg.ID)
}

func (uc *Creator) selectLocation(ctx context.Context, msg input.Message, user domain.User) error {
	nodes, err := uc.nodeUseCase.ListLocations(ctx)
	if err != nil {
		return err
	}

	markup := interaction.NewReplyMarkup()
	deviceType := msg.Args[0]

	for i, node := range nodes {
		if i%2 == 0 {
			markup.Next()
		}

		locationCode := node.ID[:2]
		location, ok := uc.locations.Get(locationCode)
		if !ok {
			continue
		}

		locationName := i18n.Get(user.LanguageCode, location.Name, nil, nil)
		payload := fmt.Sprintf("newdevice %s %s", deviceType, node.ID)

		markup.AddButton(interaction.NewButton().
			Text(locationName).
			CallbackData(payload).
			IconCustomEmojiID(location.Icon).
			Build())
	}

	text := i18n.Get(user.LanguageCode, "SelectLocation", nil, nil)
	backBtn := i18n.Get(user.LanguageCode, "BackButton", nil, nil)

	markup.Next().AddButton(interaction.NewButton().
		Text(backBtn).
		CallbackData("newdevice").
		Build())

	return uc.bot.NewMessage(msg.ChatID, text).
		WithReplyMarkup(markup.Build()).
		Edit(ctx, msg.ID)
}

func (uc *Creator) createDevice(ctx context.Context, msg input.Message, user domain.User) error {
	devices, err := uc.deviceUseCase.List(ctx, user.ID)
	if err != nil {
		return err
	}

	if len(devices) >= user.Devices {
		text := i18n.Get(user.LanguageCode, "DeviceLimitReached", nil, nil)
		return uc.bot.ShowNotification(ctx, msg.InteractionID, text, false)
	}

	os := msg.Args[0]
	nodeID := msg.Args[1]

	config, err := uc.deviceUseCase.Create(ctx, user.ID, nodeID, os, user.Bandwidth)
	if err != nil {
		return err
	}

	cfgName := fmt.Sprintf("%s-%s.conf", nodeID[:2], os)
	return uc.bot.NewMessage(msg.ChatID, "").
		WithFile(cfgName, strings.NewReader(config)).
		Send(ctx)
}
