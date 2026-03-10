package helpers

import (
	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/interaction"
	"github.com/thebeyond-net/control-plane/internal/core/ports"
)

func BuildSelectionMarkup[T ports.Selectable](
	layout [][]string,
	items map[string]T,
	callbackPrefix string,
) *interaction.ReplyMarkupBuilder {
	markup := interaction.NewReplyMarkup()
	for _, row := range layout {
		markup.Next()
		for _, code := range row {
			item, ok := items[code]
			if !ok {
				continue
			}

			markup.AddButton(interaction.NewButton().
				Text(item.GetName()).
				CallbackData(callbackPrefix + " " + item.GetCode()).
				IconCustomEmojiID(item.GetIcon()).
				Build())
		}
	}
	return markup
}
