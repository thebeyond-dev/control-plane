package middlewares

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/thebeyond-net/control-plane/internal/core/application/input"
	"github.com/thebeyond-net/control-plane/internal/core/ports"
	"go.uber.org/zap"
)

type PaymentHandler struct {
	authUseCase          ports.AuthUseCase
	billingUseCase       ports.BillingUseCase
	telegramAlertUseCase ports.AlertUseCase
	appLogger            *zap.Logger
}

func NewPaymentHandler(
	authUseCase ports.AuthUseCase,
	billingUseCase ports.BillingUseCase,
	telegramAlertUseCase ports.AlertUseCase,
	appLogger *zap.Logger,
) *PaymentHandler {
	return &PaymentHandler{
		authUseCase,
		billingUseCase,
		telegramAlertUseCase,
		appLogger,
	}
}

func (h *PaymentHandler) Handle(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.PreCheckoutQuery != nil {
			h.HandlePreCheckout(ctx, b, update)
			return
		}

		message := update.Message
		if message != nil && message.SuccessfulPayment != nil {
			h.HandleSuccessfulPayment(ctx, b, update)
			return
		}

		next(ctx, b, update)
	}
}

func (h *PaymentHandler) HandlePreCheckout(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.PreCheckoutQuery == nil {
		return
	}

	query := update.PreCheckoutQuery

	if query.Currency != "XTR" {
		h.rejectCheckout(ctx, b, query.ID, "Unsupported currency")
		return
	}

	var devices, bandwidth, days int
	if _, err := fmt.Sscanf(query.InvoicePayload, "devices=%d;bandwidth=%d;days=%d",
		&devices, &bandwidth, &days); err != nil {
		h.appLogger.Warn("invalid invoice payload",
			zap.String("payload", query.InvoicePayload),
			zap.Error(err))
		h.rejectCheckout(ctx, b, query.ID, "Invalid payload format")
		return
	}

	_, err := b.AnswerPreCheckoutQuery(ctx, &bot.AnswerPreCheckoutQueryParams{
		PreCheckoutQueryID: query.ID,
		OK:                 true,
	})
	if err != nil {
		h.appLogger.Error("failed to answer pre-checkout", zap.Error(err))
	}
}

func (h *PaymentHandler) HandleSuccessfulPayment(ctx context.Context, b *bot.Bot, update *models.Update) {
	message := update.Message
	if message == nil || message.SuccessfulPayment == nil {
		return
	}

	payment := message.SuccessfulPayment

	var devices, bandwidth, days int
	if _, err := fmt.Sscanf(payment.InvoicePayload, "devices=%d;bandwidth=%d;days=%d",
		&devices, &bandwidth, &days); err != nil {
		h.appLogger.Error("failed to parse payload on success",
			zap.String("payload", payment.InvoicePayload),
			zap.Error(err))
		return
	}

	payerID := message.From.ID
	providerID := strconv.FormatInt(payerID, 10)

	user, err := h.authUseCase.Login(ctx, "telegram", providerID, input.Login{})
	if err != nil {
		h.appLogger.Error("user login failed after payment",
			zap.String("payer_id", providerID),
			zap.Error(err))
		return
	}

	if err := h.billingUseCase.RenewSubscription(ctx, user.ID, devices, bandwidth, days); err != nil {
		h.appLogger.Error("renew subscription failed",
			zap.String("user_id", user.ID),
			zap.Error(err))
		return
	}

	user.Devices = devices
	user.Bandwidth = bandwidth
	user.Subscription.ExpiresAt = user.Subscription.ExpiresAt.
		Add(time.Hour * 24 * time.Duration(days))

	if err := h.telegramAlertUseCase.NotifyNewPayment(ctx, user, payment.TotalAmount, days); err != nil {
		h.appLogger.Error("notify new payment failed", zap.Error(err))
	}
}

func (h *PaymentHandler) rejectCheckout(ctx context.Context, b *bot.Bot, queryID string, reason string) {
	_, _ = b.AnswerPreCheckoutQuery(ctx, &bot.AnswerPreCheckoutQueryParams{
		PreCheckoutQueryID: queryID,
		OK:                 false,
		ErrorMessage:       reason,
	})
}
