package interaction

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton
}

type InlineKeyboardButton struct {
	Text              string
	CallbackData      string
	URL               string
	Pay               bool
	IconCustomEmojiID string
	Style             string
}

type ReplyMarkupBuilder struct {
	markup          InlineKeyboardMarkup
	currentRowIndex int
}

func NewReplyMarkup() *ReplyMarkupBuilder {
	return &ReplyMarkupBuilder{
		markup: InlineKeyboardMarkup{
			InlineKeyboard: [][]InlineKeyboardButton{},
		},
		currentRowIndex: -1,
	}
}

func (b *ReplyMarkupBuilder) Next() *ReplyMarkupBuilder {
	b.currentRowIndex++
	b.markup.InlineKeyboard = append(b.markup.InlineKeyboard, []InlineKeyboardButton{})
	return b
}

func (b *ReplyMarkupBuilder) AddButton(button InlineKeyboardButton) *ReplyMarkupBuilder {
	b.markup.InlineKeyboard[b.currentRowIndex] = append(b.markup.InlineKeyboard[b.currentRowIndex], button)
	return b
}

func (b *ReplyMarkupBuilder) Build() InlineKeyboardMarkup {
	return b.markup
}

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
