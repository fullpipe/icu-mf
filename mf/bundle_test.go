package mf

import (
	"io/fs"
	"reflect"
	"runtime"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

func TestNewBundle(t *testing.T) {
	b, err := NewBundle(
		WithDefaulLangFallback(language.English),
		WithLangFallback(language.Portuguese, language.Spanish),
	)

	require.NoError(t, err)
	assert.NotNil(t, b)
}

func TestWithDefaulLangFallback(t *testing.T) {
	b := &bundle{
		fallbacks:    map[language.Tag]language.Tag{},
		translators:  map[language.Tag]Translator{},
		dictionaries: map[language.Tag]Dictionary{},

		defaultLang:         language.Und,
		defaultErrorHandler: func(_ error, _ string, _ map[string]any) {},
	}

	assert.Equal(t, language.Und, b.defaultLang)
	WithDefaulLangFallback(language.Afrikaans)(b)
	assert.Equal(t, language.Afrikaans, b.defaultLang)
}

func TestWithLangFallback(t *testing.T) {
	b := &bundle{
		fallbacks:    map[language.Tag]language.Tag{},
		translators:  map[language.Tag]Translator{},
		dictionaries: map[language.Tag]Dictionary{},

		defaultLang:         language.Und,
		defaultErrorHandler: func(_ error, _ string, _ map[string]any) {},
	}

	assert.Equal(t, language.Und, b.fallbacks[language.AmericanEnglish])
	WithLangFallback(language.AmericanEnglish, language.English)(b)
	assert.Equal(t, language.English, b.fallbacks[language.AmericanEnglish])
}

func TestWithErrorHandler(t *testing.T) {
	b := &bundle{
		fallbacks:    map[language.Tag]language.Tag{},
		translators:  map[language.Tag]Translator{},
		dictionaries: map[language.Tag]Dictionary{},

		defaultLang:         language.Und,
		defaultErrorHandler: func(_ error, _ string, _ map[string]any) {},
	}

	assert.NotNil(t, b.defaultErrorHandler)
	errHandler := func(_ error, _ string, _ map[string]any) {}
	WithErrorHandler(errHandler)(b)

	funcName1 := runtime.FuncForPC(reflect.ValueOf(errHandler).Pointer()).Name()
	funcName2 := runtime.FuncForPC(reflect.ValueOf(b.defaultErrorHandler).Pointer()).Name()
	assert.Equal(t, funcName1, funcName2)
}

func TestBundle_LoadMessages(t *testing.T) {
	b := &bundle{
		fallbacks:    map[language.Tag]language.Tag{},
		translators:  map[language.Tag]Translator{},
		dictionaries: map[language.Tag]Dictionary{},

		defaultLang:         language.Und,
		defaultErrorHandler: func(_ error, _ string, _ map[string]any) {},
	}

	fs := fstest.MapFS{
		"non_readable": {
			Mode: fs.ModeDir,
		},
		"foo.en.yaml": {
			Data: []byte("foo: bar"),
		},
	}

	require.Error(t, b.LoadMessages(fs, "file_not_exists.yaml", language.English))
	assert.Nil(t, b.dictionaries[language.English])

	require.Error(t, b.LoadMessages(fs, "non_readable", language.English))
	assert.Nil(t, b.dictionaries[language.English])

	require.NoError(t, b.LoadMessages(fs, "foo.en.yaml", language.English))
	assert.NotNil(t, b.dictionaries[language.English])

	require.Error(t, b.LoadMessages(fs, "foo.en.yaml", language.English), "error when lang already loaded")
	assert.NotNil(t, b.dictionaries[language.English])
}

func TestBundle_LoadDir(t *testing.T) {
	b := &bundle{
		fallbacks:    map[language.Tag]language.Tag{},
		translators:  map[language.Tag]Translator{},
		dictionaries: map[language.Tag]Dictionary{},

		defaultLang:         language.Und,
		defaultErrorHandler: func(_ error, _ string, _ map[string]any) {},
	}

	require.NoError(t, b.LoadDir(fstest.MapFS{"dir.en": {Mode: fs.ModeDir}}), "no error on empty fs")
	assert.Nil(t, b.dictionaries[language.English], "does not loads dirs")

	require.NoError(t, b.LoadDir(fstest.MapFS{"messages.en.toml": {Data: []byte("foo=bar")}}), "no error on non yaml files")
	assert.Nil(t, b.dictionaries[language.English], "but does not loads them")

	require.Error(t, b.LoadDir(fstest.MapFS{".yaml": {Data: []byte("foo: bar")}}), "error on invalid filename")
	require.Error(t, b.LoadDir(fstest.MapFS{"messages.FOO.yaml": {Data: []byte("foo: bar")}}), "error on invalid lang in filename")

	require.NoError(t, b.LoadDir(fstest.MapFS{
		"messages.en.yaml": {Data: []byte("foo: bar")},
		"messages.es.yaml": {Data: []byte("foo: bar")},
	}), "no error on normal yaml files")
	assert.NotNil(t, b.dictionaries[language.English])
	assert.NotNil(t, b.dictionaries[language.Spanish])
}

func TestBundle_Translator(t *testing.T) {
	b, _ := NewBundle()
	assert.NotNil(t, b.Translator("ru"), "even empty bundle returns some translator")
	assert.Equal(t, "msg_id", b.Translator("ru").Trans("msg_id"), "even empty bundle returns some translator")

	b, _ = NewBundle(
		WithDefaulLangFallback(language.English),
		WithLangFallback(language.Portuguese, language.Spanish),
	)

	require.NoError(t, b.LoadDir(fstest.MapFS{
		"messages.en.yaml": {Data: []byte("foo: en\nbar_id: enbar")},
		"messages.es.yaml": {Data: []byte("foo: es")},
	}))

	assert.Equal(t, "en", b.Translator("en").Trans("foo"), "en loaded from file")
	assert.Equal(t, "en", b.Translator("pl").Trans("foo"), "fallback to defalt lang")
	assert.Equal(t, "en", b.Translator("FOO").Trans("foo"), "invalid lang fallback to defalt lang")
	assert.Equal(t, "enbar", b.Translator("es").Trans("bar_id"), "lang with dictionary fallbacks to default lang")
	assert.Equal(t, "es", b.Translator("pt").Trans("es"), "lang fallback")
	assert.Equal(t, "none_id", b.Translator("pl").Trans("none_id"), "dummy translator if nothing works")
	assert.Equal(t, "none_id", b.Translator("en").Trans("none_id"), "dummy translator if nothing works")
}
