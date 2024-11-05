package tr

import (
	"context"
	"embed"
	"log/slog"

	"github.com/fullpipe/icu-mf/mf"
	"golang.org/x/text/language"
)

type contextKey string

//go:embed messages/messages.*.yaml
var messagesDir embed.FS
var bundle mf.Bundle
var langKey contextKey = "lang"

func init() {
	var err error

	bundle, err = mf.NewBundle(
		mf.WithDefaulLangFallback(language.English),

		mf.WithLangFallback(language.BritishEnglish, language.English),
		mf.WithLangFallback(language.Portuguese, language.Spanish),

		mf.WithYamlProvider(messagesDir),

		mf.WithErrorHandler(func(err error, id string, ctx map[string]any) {
			slog.Error(err.Error(), slog.String("id", id), slog.Any("ctx", ctx))
		}),
	)

	if err != nil {
		panic(err)
	}
}

// Tr translates message
func Tr(ctx context.Context, path string, args ...mf.TranslationArg) string {
	lang := ctx.Value(langKey).(string)
	if lang == "" {
		lang = "en"
	}

	return bundle.Translator(lang).Trans(path, args...)
}
