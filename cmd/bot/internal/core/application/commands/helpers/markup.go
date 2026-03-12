package helpers

import (
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/interaction"
	"github.com/thebeyond-net/control-plane/internal/core/domain"
)

func BuildSelectionMarkup(
	items domain.Items,
	callbackPrefix string,
	rowWidth int,
) *interaction.ReplyMarkupBuilder {
	markup := interaction.NewReplyMarkup()
	for i, item := range items.All() {
		if i%rowWidth == 0 {
			markup.Next()
		}

		markup.AddButton(interaction.NewButton().
			Text(item.Name).
			CallbackData(callbackPrefix + " " + item.Code).
			IconCustomEmojiID(item.Icon).
			Build())
	}
	return markup
}
