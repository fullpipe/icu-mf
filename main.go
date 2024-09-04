package main

import (
	"embed"
	"log"
	"log/slog"

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

		mf.WithErrorHandler(func(err error, key string, ctx map[string]any) {
			slog.Error(err.Error(), slog.String("key", key), slog.Any("ctx", ctx))
		}),
	)

	if err != nil {
		log.Fatal(err)
	}

	err = bundle.LoadDir(messagesDir)
	if err != nil {
		log.Fatal(err)
	}

	tr := bundle.Translator("ru")

	// log.Println(bundle)
	slog.Info(tr.Trans("one.two", mf.Arg("foo", "bar")))
	slog.Info(tr.Trans("user.gender", mf.Arg("gender", "male")))
	slog.Info(tr.Trans("user.gender"))
	slog.Info(tr.Trans("user.gender_adsasdw", mf.Arg("gender", "male")))
	log.Println(tr.Trans("apples", mf.Arg("apples", 0)))
	log.Println(
		tr.Trans("say_hello", mf.Arg("name", "Bob")),
	) // prints "Hello Bob!"

	slog.Info(tr.Trans("apples", mf.Arg("apples", 1)))
	slog.Info(tr.Trans("apples", mf.Arg("apples", 2)))
	slog.Info(tr.Trans("apples", mf.Arg("apples", 2)))
	slog.Info(tr.Trans("apples", mf.Arg("apples", 5)))
	slog.Info(tr.Trans("apples", mf.Arg("apples", 7)))
	slog.Info(tr.Trans("apples", mf.Arg("apples", "23")))
	slog.Info(tr.Trans("apples", mf.Arg("apples", "51")))
}
