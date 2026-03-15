package webhook

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/thebeyond-net/control-plane/internal/i18n"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) HandleUpdate(ctx context.Context, b *bot.Bot, update *models.Update) {
	languageCode := "en"
	if update.Message != nil && update.Message.From != nil {
		languageCode = update.Message.From.LanguageCode
	} else if update.CallbackQuery != nil {
		languageCode = update.CallbackQuery.From.LanguageCode
	}

	text := i18n.Get(languageCode, "CommandNotFound", nil, nil)
	switch {
	case update.Message != nil:
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   text,
		})
	case update.CallbackQuery != nil:
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			Text:            text,
		})
	}
}
