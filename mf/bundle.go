package mf

import (
	"io/fs"

	"github.com/pkg/errors"
	"golang.org/x/text/language"
)

type Bundle interface {
	Translator(lang ...string) Translator
}

type bundle struct {
	fallbacks           map[language.Tag]language.Tag
	translators         map[language.Tag]Translator
	provider            MessageProvider
	matcher             language.Matcher
	languages           []language.Tag
	defaultLang         language.Tag
	defaultErrorHandler ErrorHandler
}

type ErrorHandler func(err error, id string, ctx map[string]any)

type BundleOption func(b *bundle) error

func NewBundle(options ...BundleOption) (Bundle, error) {
	bundle := &bundle{
		fallbacks:   make(map[language.Tag]language.Tag),
		translators: make(map[language.Tag]Translator),
		defaultLang: language.Und,
		defaultErrorHandler: func(_ error, _ string, _ map[string]any) {
			// Provide meaningful logging or handling here
		},
	}

	for _, option := range options {
		if err := option(bundle); err != nil {
			return nil, err
		}
	}

	if bundle.provider == nil {
		return nil, errors.New("you must add a message provider with WithFSProvider or WithProvider")
	}
	languagesProvider, ok := bundle.provider.(LanguagesProvider)
	if ok {
		languages := languagesProvider.Languages()
		for j, v := range languages {
			if v == bundle.defaultLang {
				for i := j; i > 0; i-- {
					languages[i] = languages[i-1]
				}
				languages[0] = bundle.defaultLang
				break
			}
		}
		bundle.languages = languages
		bundle.matcher = language.NewMatcher(languages)
	}

	// Check for cyclic fallbacks
	if err := checkCyclicFallbacks(bundle.fallbacks); err != nil {
		return nil, err
	}

	return bundle, nil
}

func (b *bundle) Translator(lang ...string) Translator {
	var tag language.Tag
	if b.matcher != nil {
		var idx int
		_, idx = language.MatchStrings(b.matcher, lang...)
		tag = b.languages[idx]
	} else if len(lang) == 0 {
		tag = b.defaultLang
	} else {
		var err error
		tag, err = language.Parse(lang[0])
		if err != nil {
			tag = b.defaultLang
		}
	}

	if tr, ok := b.translators[tag]; ok {
		return tr
	}

	tr := b.getTranslator(tag)
	b.translators[tag] = tr

	return tr
}

func (b *bundle) getTranslator(tag language.Tag) Translator {
	if tr, ok := b.translators[tag]; ok {
		return tr
	}

	var fallback Translator
	if fallbackTag, hasFallback := b.fallbacks[tag]; hasFallback {
		fallback = b.getTranslator(fallbackTag)
	} else if tag != b.defaultLang {
		fallback = b.getTranslator(b.defaultLang)
	}

	return &translator{
		provider:     b.provider,
		fallback:     fallback,
		errorHandler: b.defaultErrorHandler,
		lang:         tag,
	}
}

func WithDefaultLangFallback(l language.Tag) BundleOption {
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

// checkCyclicFallbacks checks for cyclic fallbacks to prevent infinite loops
func checkCyclicFallbacks(fallbacks map[language.Tag]language.Tag) error {
	visited := make(map[language.Tag]bool)
	for tag := range fallbacks {
		if hasCycle(tag, fallbacks, visited) {
			return errors.New("cyclic fallback detected")
		}
	}
	return nil
}

func hasCycle(tag language.Tag, fallbacks map[language.Tag]language.Tag, visited map[language.Tag]bool) bool {
	if visited[tag] {
		return true
	}

	visited[tag] = true
	defer delete(visited, tag)

	next, ok := fallbacks[tag]
	if !ok {
		return false
	}

	return hasCycle(next, fallbacks, visited)
}
