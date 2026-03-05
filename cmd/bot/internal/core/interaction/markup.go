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
