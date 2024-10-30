package main

import (
	"embed"
	"log"
	"log/slog"
	"math/rand/v2"

	"github.com/fullpipe/icu-mf/mf"
	"golang.org/x/text/language"
)

//go:embed var/messages.*.yaml
var messagesDir embed.FS

func main() {
	bundle, err := mf.NewBundle(
		mf.WithDefaulLangFallback(language.English),

		mf.WithLangFallback(language.BritishEnglish, language.English),
		mf.WithLangFallback(language.Portuguese, language.Spanish),

		mf.WithYamlProvider(messagesDir),

		mf.WithErrorHandler(func(err error, id string, ctx map[string]any) {
			slog.Error(err.Error(), slog.String("id", id), slog.Any("ctx", ctx))

			// or
			// panic(err)
		}),
	)

	if err != nil {
		log.Fatal(err)
	}

	tr := bundle.Translator("en")

	slog.Info(tr.Trans("title", mf.Arg("lang", "en")))
	slog.Info(tr.Trans("subtitle", mf.Arg("lang", "en")))
	slog.Info(tr.Trans("say.hello", mf.Arg("name", "Bob")))
	slog.Info(tr.Trans("say.goodbye", mf.Arg("name", "Bob")))
	slog.Info(tr.Trans("user.profile.cats", mf.Arg("num", rand.IntN(10)))) //nolint
	slog.Info(tr.Trans("user.profile.dogs", mf.Arg("num", rand.IntN(10)))) //nolint

	tr = bundle.Translator("es")

	slog.Info(tr.Trans("title", mf.Arg("lang", "es")))
	slog.Info(tr.Trans("subtitle", mf.Arg("lang", "es")))
	slog.Info(tr.Trans("say.hello", mf.Arg("name", "Bob")))
	slog.Info(tr.Trans("say.goodbye", mf.Arg("name", "Bob")))
	slog.Info(tr.Trans("user.profile.cats", mf.Arg("num", rand.IntN(10)))) //nolint
	slog.Info(tr.Trans("user.profile.dogs", mf.Arg("num", rand.IntN(10)))) //nolint

	tr = bundle.Translator("pt")

	slog.Info(tr.Trans("title", mf.Arg("lang", "pt")))
	slog.Info(tr.Trans("subtitle", mf.Arg("lang", "pt")))
	slog.Info(tr.Trans("say.hello", mf.Arg("name", "Bob")))
	slog.Info(tr.Trans("say.goodbye", mf.Arg("name", "Bob")))
	slog.Info(tr.Trans("user.profile.cats", mf.Arg("num", rand.IntN(10)))) //nolint
	slog.Info(tr.Trans("user.profile.dogs", mf.Arg("num", rand.IntN(10)))) //nolint

	tr = bundle.Translator("ru")

	slog.Info(tr.Trans("title", mf.Arg("lang", "ru")))
	slog.Info(tr.Trans("subtitle", mf.Arg("lang", "ru")))
	slog.Info(tr.Trans("say.hello", mf.Arg("name", "Bob")))
	slog.Info(tr.Trans("say.goodbye", mf.Arg("name", "Bob")))
	slog.Info(tr.Trans("user.profile.cats", mf.Arg("num", rand.IntN(10)))) //nolint
	slog.Info(tr.Trans("user.profile.dogs", mf.Arg("num", rand.IntN(10)))) //nolint
}
