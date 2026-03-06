package bot

import (
	"context"
	"io"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/interaction"
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/ports"
)

type MessageBuilder struct {
	client      *bot.Bot
	ChatID      int
	Text        string
	ReplyMarkup models.ReplyMarkup
	File        File
}

type File struct {
	Name    string
	Content io.Reader
}

func (b *Bot) NewMessage(chatID int, text string) ports.MessageBuilder {
	return &MessageBuilder{
		client: b.client,
		ChatID: chatID,
		Text:   text,
	}
}

func (b *MessageBuilder) WithReplyMarkup(
	markup interaction.InlineKeyboardMarkup,
) ports.MessageBuilder {
	inlineKeyboard := [][]models.InlineKeyboardButton{}
	for _, rows := range markup.InlineKeyboard {
		buttons := []models.InlineKeyboardButton{}
		for _, row := range rows {
			button := models.InlineKeyboardButton{Text: row.Text}

			if row.URL != "" {
				button.URL = row.URL
			}

			if row.Pay {
				button.Pay = true
			}

			if row.CallbackData != "" {
				button.CallbackData = row.CallbackData
			}

			if row.IconCustomEmojiID != "" {
				button.IconCustomEmojiID = row.IconCustomEmojiID
			}

			if row.Style != "" {
				button.Style = row.Style
			}

			buttons = append(buttons, button)
		}
		inlineKeyboard = append(inlineKeyboard, buttons)
	}

	b.ReplyMarkup = &models.InlineKeyboardMarkup{
		InlineKeyboard: inlineKeyboard,
	}

	return b
}

func (b *MessageBuilder) WithFile(name string, content io.Reader) ports.MessageBuilder {
	b.File = File{name, content}
	return b
}

func (b *MessageBuilder) Send(ctx context.Context) error {
	if b.File.Content != nil {
		params := &bot.SendDocumentParams{
			ChatID:      b.ChatID,
			Caption:     b.Text,
			ReplyMarkup: b.ReplyMarkup,
			Document: &models.InputFileUpload{
				Filename: b.File.Name,
				Data:     b.File.Content,
			},
		}
		_, err := b.client.SendDocument(ctx, params)
		return err
	}

	params := &bot.SendMessageParams{
		ChatID:      b.ChatID,
		Text:        b.Text,
		ReplyMarkup: b.ReplyMarkup,
		ParseMode:   models.ParseModeMarkdown,
	}

	_, err := b.client.SendMessage(ctx, params)
	return err
}

func (b *MessageBuilder) Edit(ctx context.Context, messageID int) error {
	params := &bot.EditMessageTextParams{
		ChatID:      b.ChatID,
		MessageID:   messageID,
		Text:        b.Text,
		ReplyMarkup: b.ReplyMarkup,
		ParseMode:   models.ParseModeMarkdown,
	}

	_, err := b.client.EditMessageText(ctx, params)
	return err
}
