package mf

import (
	"io/fs"

	"github.com/pkg/errors"
	"golang.org/x/text/language"
)

type Bundle interface {
	Translator(lang string) Translator
}

type bundle struct {
	fallbacks   map[language.Tag]language.Tag
	translators map[language.Tag]Translator
	provider    MessageProvider

	defaultLang         language.Tag
	defaultErrorHandler ErrorHandler
}

type ErrorHandler func(err error, id string, ctx map[string]any)

type BundleOption func(b *bundle) error

func NewBundle(options ...BundleOption) (Bundle, error) {
	bundle := &bundle{
		fallbacks:   map[language.Tag]language.Tag{},
		translators: map[language.Tag]Translator{},

		defaultLang:         language.Und,
		defaultErrorHandler: func(_ error, _ string, _ map[string]any) {},
	}

	for _, option := range options {
		err := option(bundle)
		if err != nil {
			return nil, err
		}
	}

	if bundle.provider == nil {
		return nil, errors.New("you have add message provider with WithFSProvider or WithProvider")
	}

	// TODO: check fallbacks for cicles en -> es -> en -> ...

	return bundle, nil
}

func (b *bundle) Translator(lang string) Translator {
	tag, err := language.Parse(lang)
	if err != nil {
		tag = b.defaultLang
	}

	tr, ok := b.translators[tag]
	if ok {
		return tr
	}

	b.translators[tag] = b.getTranlator(tag)

	return b.translators[tag]
}

func (b *bundle) getTranlator(tag language.Tag) Translator {
	tr, ok := b.translators[tag]
	if ok {
		return tr
	}

	var fallback Translator
	fallbackTag, hasFallback := b.fallbacks[tag]
	if hasFallback {
		fallback = b.getTranlator(fallbackTag)
	} else if tag != b.defaultLang {
		fallback = b.getTranlator(b.defaultLang)
	}

	return &translator{
		provider:     b.provider,
		fallback:     fallback,
		errorHandler: b.defaultErrorHandler,
		lang:         tag,
	}
}

func WithDefaulLangFallback(l language.Tag) BundleOption {
	return func(b *bundle) error {
		b.defaultLang = l

		return nil
	}
}

func WithYamlProvider(dir fs.FS) BundleOption {
	return func(b *bundle) error {
		provider, err := NewYamlMessageProvider(dir)
		b.provider = provider

		return err
	}
}

func WithProvider(provider MessageProvider) BundleOption {
	return func(b *bundle) error {
		b.provider = provider

		return nil
	}
}

func WithLangFallback(from language.Tag, to language.Tag) BundleOption {
	return func(b *bundle) error {
		b.fallbacks[from] = to

		return nil
	}
}

func WithErrorHandler(handler ErrorHandler) BundleOption {
	return func(b *bundle) error {
		b.defaultErrorHandler = handler

		return nil
	}
}
