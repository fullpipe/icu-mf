package mf

import (
	"fmt"
	"io"
	"io/fs"
	"path"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/text/language"
)

type Bundle struct {
	fallbacks    map[language.Tag]language.Tag
	translators  map[language.Tag]Translator
	dictionaries map[language.Tag]Dictionary

	defaultLang         language.Tag
	defaultErrorHandler ErrorHandler
}

type ErrorHandler func(err error, id string, ctx map[string]any)

type BundleOption func(b *Bundle)

func NewBundle(options ...BundleOption) (*Bundle, error) {
	bundle := &Bundle{
		fallbacks:    map[language.Tag]language.Tag{},
		translators:  map[language.Tag]Translator{},
		dictionaries: map[language.Tag]Dictionary{},

		defaultLang:         language.Und,
		defaultErrorHandler: func(_ error, _ string, _ map[string]any) {},
	}

	for _, option := range options {
		option(bundle)
	}

	// TODO: check fallbacks for cicles en -> es -> en -> ...

	return bundle, nil
}

func (b *Bundle) LoadMessages(rd fs.FS, path string, lang language.Tag) error {
	yamlFile, err := rd.Open(path)
	if err != nil {
		return errors.Wrap(err, "unable to open file")
	}

	yamlData, err := io.ReadAll(yamlFile)
	if err != nil {
		return errors.Wrap(err, "unable to read file")
	}

	_, hasDictionary := b.dictionaries[lang]
	if hasDictionary {
		return fmt.Errorf("unable to load %s: language %s already has messages loaded", path, lang)
	}

	b.dictionaries[lang], err = NewDictionary(yamlData)
	if err != nil {
		return errors.Wrap(err, "unable to create dictionary")
	}

	return nil
}

func (b *Bundle) LoadDir(dir fs.FS) error {
	return fs.WalkDir(dir, ".", func(p string, f fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if f.IsDir() {
			return nil
		}

		if path.Ext(f.Name()) != ".yaml" {
			return nil
		}

		nameParts := strings.Split(f.Name(), ".")
		if len(nameParts) < 2 {
			return fmt.Errorf("no lang in file %s", f.Name())
		}

		tag, err := language.Parse(nameParts[len(nameParts)-2])
		if err != nil {
			return errors.Wrap(err, "unable to parse language from filename")
		}

		return b.LoadMessages(dir, p, tag)
	})
}

func (b *Bundle) Translator(lang string) Translator {
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

func (b *Bundle) getTranlator(tag language.Tag) Translator {
	tr, ok := b.translators[tag]
	if ok {
		return tr
	}

	dictionary, hasDictionary := b.dictionaries[tag]
	fallbackTag, hasFallback := b.fallbacks[tag]

	if hasDictionary {
		var fallback Translator

		if hasFallback {
			fallback = b.getTranlator(fallbackTag)
		} else if tag != b.defaultLang {
			fallback = b.getTranlator(b.defaultLang)
		}

		return &translator{
			dictionary:   dictionary,
			fallback:     fallback,
			errorHandler: b.defaultErrorHandler,
			lang:         tag,
		}
	}

	if hasFallback {
		return b.getTranlator(fallbackTag)
	}

	tr, ok = b.translators[b.defaultLang]
	if ok {
		return tr
	}

	dictionary, hasDictionary = b.dictionaries[b.defaultLang]
	if !hasDictionary {
		dictionary = &dummyDictionary{}
	}

	return &translator{
		dictionary:   dictionary,
		errorHandler: b.defaultErrorHandler,
		lang:         tag,
	}
}

func WithDefaulLangFallback(l language.Tag) BundleOption {
	return func(b *Bundle) {
		b.defaultLang = l
	}
}

func WithLangFallback(from language.Tag, to language.Tag) BundleOption {
	return func(b *Bundle) {
		b.fallbacks[from] = to
	}
}

func WithErrorHandler(handler ErrorHandler) BundleOption {
	return func(b *Bundle) {
		b.defaultErrorHandler = handler
	}
}
