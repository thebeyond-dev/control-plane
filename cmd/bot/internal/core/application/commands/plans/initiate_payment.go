package plans

import (
	"context"
	"fmt"
	"strings"

	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/application/input"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/interaction"
	"github.com/thebeyond-net/control-plane/internal/core/domain"
	"github.com/thebeyond-net/control-plane/internal/i18n"
)

func (uc *UseCase) initiatePayment(
	ctx context.Context,
	msg input.Message,
	user domain.User,
	plan domain.Plan,
	state requestState,
) error {
	currencyCode := ""
	switch state.PaymentMethod {
	case "yookassa":
		currencyCode = "rub"
	case "stars":
		currencyCode = "xtr"
	}

	periodDiscount := int((float64(uc.getPeriodDiscountFraction(state.Period)) / 100.0) * 100)
	price, err := plan.GetPrice(currencyCode, state.Bandwidth, state.Period, periodDiscount)
	if err != nil {
		return err
	}

	switch state.PaymentMethod {
	case "yookassa":
		return uc.handleYookassaPayment(ctx, msg, user, plan, state, price, currencyCode)
	case "stars":
		return uc.handleStarsPayment(ctx, user, plan, state, price, currencyCode)
	}

	return fmt.Errorf("unknown payment method: %s", state.PaymentMethod)
}

func (uc *UseCase) handleYookassaPayment(
	ctx context.Context,
	msg input.Message,
	user domain.User,
	plan domain.Plan,
	state requestState,
	price float64,
	currencyCode string,
) error {
	payBtn := i18n.Get(user.LanguageCode, "PayButton", nil, nil)
	backBtn := i18n.Get(user.LanguageCode, "BackButton", nil, nil)

	url, err := uc.yookassa.NewPayment(
		ctx, user,
		currencyCode,
		plan.Devices,
		state.Bandwidth,
		state.Period,
		price,
	)
	if err != nil {
		return fmt.Errorf("yookassa payment creation failed: %w", err)
	}

	text := i18n.Get(user.LanguageCode, "AgreementPayment", nil, nil)
	markup := interaction.NewReplyMarkup()

	markup.Next().AddButton(interaction.NewButton().
		Text(payBtn).
		IconCustomEmojiID(IconBuy).
		URL(url).
		Build())
	markup.Next().AddButton(interaction.NewButton().
		Text(backBtn).
		CallbackData(fmt.Sprintf("plan %d %d %d 0", state.PlanID, state.Bandwidth, state.Period)).
		Build())

	return uc.bot.NewMessage(msg.ChatID, text).
		WithReplyMarkup(markup.Build()).
		Edit(ctx, msg.ID)
}

func (uc *UseCase) handleStarsPayment(
	ctx context.Context,
	user domain.User,
	plan domain.Plan,
	state requestState,
	price float64,
	currencyCode string,
) error {
	_, err := uc.telegramStars.NewPayment(
		ctx, user,
		strings.ToUpper(currencyCode),
		plan.Devices,
		state.Bandwidth,
		state.Period,
		price,
	)
	if err != nil {
		return fmt.Errorf("telegram stars payment creation failed: %w", err)
	}
	return nil
}
