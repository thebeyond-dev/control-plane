package plans

import (
	"context"
	"fmt"

	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/application/input"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/interaction"
	"github.com/thebeyond-net/control-plane/cmd/bot/pkg/keycap"
	"github.com/thebeyond-net/control-plane/internal/core/domain"
	"github.com/thebeyond-net/control-plane/internal/i18n"
)

var (
	IconBuy       = "5244828222935306011"
	IconBandwidth = "5242486782859185227"
)

func (uc *UseCase) renderPlanDetails(
	ctx context.Context,
	msg input.Message,
	user domain.User,
	plan domain.Plan,
	bandwidth int,
) error {
	buyBtn := i18n.Get(user.LanguageCode, "BuyButton", nil, nil)
	backBtn := i18n.Get(user.LanguageCode, "BackButton", nil, nil)

	prevID, nextID,
		hasPrev, hasNext := domain.GetNeighborIDs(plan.ID, len(uc.plans))
	prevBandwidth, nextBandwidth,
		hasPrevBandwidth, hasNextBandwidth := domain.GetNeighborBandwidths(
		bandwidth, uc.bandwidths,
	)

	text := uc.formatPlanDetails(plan, user.LanguageCode, user.CurrencyCode, bandwidth, 30)
	markup := interaction.NewReplyMarkup()

	markup.Next().AddButton(interaction.NewButton().
		Text(buyBtn).
		CallbackData(fmt.Sprintf("plan %d %d 0", plan.ID, bandwidth)).
		IconCustomEmojiID(IconBuy).
		Build())

	markup.Next()
	if hasPrev {
		payload := fmt.Sprintf("plan %d %d", prevID, bandwidth)
		markup.AddButton(interaction.NewButton().
			Text("▬").
			CallbackData(payload).
			Build())
	}

	markup.AddButton(interaction.NewButton().
		Text(keycap.Convert(fmt.Sprint(plan.ID + 1))).
		CallbackData("plan grid").
		Build())

	if hasNext {
		payload := fmt.Sprintf("plan %d %d", nextID, bandwidth)
		markup.AddButton(interaction.NewButton().
			Text("✚").
			CallbackData(payload).
			Build())
	}

	if hasNextBandwidth {
		btnText := formatBandwidth(nextBandwidth)
		payload := fmt.Sprintf("plan %d %d", plan.ID, nextBandwidth)
		markup.Next().AddButton(interaction.NewButton().
			Text(btnText).
			CallbackData(payload).
			IconCustomEmojiID(IconBandwidth).
			Build())
	}

	if hasPrevBandwidth {
		btnText := formatBandwidth(prevBandwidth)
		payload := fmt.Sprintf("plan %d %d", plan.ID, prevBandwidth)
		markup.Next().AddButton(interaction.NewButton().
			Text(btnText).
			CallbackData(payload).
			IconCustomEmojiID(IconBandwidth).
			Build())
	}

	markup.Next().AddButton(interaction.NewButton().
		Text(backBtn).
		CallbackData("menu").
		Build())

	return uc.bot.NewMessage(msg.ChatID, text).
		WithReplyMarkup(markup.Build()).
		Edit(ctx, msg.ID)
}
