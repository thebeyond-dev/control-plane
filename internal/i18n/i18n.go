package i18n

import (
	"embed"
	"fmt"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/zap"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle

func Init(fs embed.FS, dir string, logger *zap.Logger) error {
	logger = logger.Named("i18n")

	logger.Debug("initializing i18n bundle",
		zap.String("default_language", language.English.String()),
		zap.String("dir", dir),
	)

	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	entries, err := fs.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("i18n: failed to read locales directory %q: %w", dir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".toml") {
			continue
		}

		filePath := path.Join(dir, entry.Name())
		if _, err := bundle.LoadMessageFileFS(fs, filePath); err != nil {
			return fmt.Errorf("i18n: failed to load message file %q: %w", filePath, err)
		}

		logger.Debug("loaded locale file",
			zap.String("file", entry.Name()),
		)
	}

	return nil
}

func Get(lang, messageID string, templateData map[string]any, pluralCount any) string {
	if bundle == nil {
		return messageID
	}

	localizer := i18n.NewLocalizer(bundle, lang)

	config := &i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
		PluralCount:  pluralCount,
	}

	translated, err := localizer.Localize(config)
	if err != nil {
		return messageID
	}

	return translated
}
