package interaction

type ReplyMarkupButtonBuilder struct {
	button InlineKeyboardButton
}

func NewButton() *ReplyMarkupButtonBuilder {
	return &ReplyMarkupButtonBuilder{}
}

func (b *ReplyMarkupButtonBuilder) Text(text string) *ReplyMarkupButtonBuilder {
	b.button.Text = text
	return b
}

func (b *ReplyMarkupButtonBuilder) CallbackData(callbackData string) *ReplyMarkupButtonBuilder {
	b.button.CallbackData = callbackData
	return b
}

func (b *ReplyMarkupButtonBuilder) URL(url string) *ReplyMarkupButtonBuilder {
	b.button.URL = url
	return b
}

func (b *ReplyMarkupButtonBuilder) Pay() *ReplyMarkupButtonBuilder {
	b.button.Pay = true
	return b
}

func (b *ReplyMarkupButtonBuilder) IconCustomEmojiID(iconCustomEmojiID string) *ReplyMarkupButtonBuilder {
	b.button.IconCustomEmojiID = iconCustomEmojiID
	return b
}

func (b *ReplyMarkupButtonBuilder) Style(style string) *ReplyMarkupButtonBuilder {
	b.button.Style = style
	return b
}

func (b *ReplyMarkupButtonBuilder) Build() InlineKeyboardButton {
	return b.button
}
