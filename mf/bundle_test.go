package mf

import (
	"reflect"
	"runtime"
	"testing"
	"testing/fstest"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

func TestNewBundle(t *testing.T) {
	b, err := NewBundle(
		WithDefaulLangFallback(language.English),
		WithLangFallback(language.Portuguese, language.Spanish),
		WithProvider(new(MockedProvider)),
	)

	require.NoError(t, err)
	assert.NotNil(t, b)
}

func TestWithDefaulLangFallback(t *testing.T) {
	b := &bundle{
		fallbacks:   map[language.Tag]language.Tag{},
		translators: map[language.Tag]Translator{},

		defaultLang:         language.Und,
		defaultErrorHandler: func(_ error, _ string, _ map[string]any) {},
	}

	assert.Equal(t, language.Und, b.defaultLang)
	require.NoError(t, WithDefaulLangFallback(language.Afrikaans)(b))
	assert.Equal(t, language.Afrikaans, b.defaultLang)
}

func TestWithLangFallback(t *testing.T) {
	b := &bundle{
		fallbacks:   map[language.Tag]language.Tag{},
		translators: map[language.Tag]Translator{},

		defaultLang:         language.Und,
		defaultErrorHandler: func(_ error, _ string, _ map[string]any) {},
	}

	assert.Equal(t, language.Und, b.fallbacks[language.AmericanEnglish])
	require.NoError(t, WithLangFallback(language.AmericanEnglish, language.English)(b))
	assert.Equal(t, language.English, b.fallbacks[language.AmericanEnglish])
}

func TestWithErrorHandler(t *testing.T) {
	b := &bundle{
		fallbacks:   map[language.Tag]language.Tag{},
		translators: map[language.Tag]Translator{},

		defaultLang:         language.Und,
		defaultErrorHandler: func(_ error, _ string, _ map[string]any) {},
	}

	assert.NotNil(t, b.defaultErrorHandler)
	errHandler := func(_ error, _ string, _ map[string]any) {}
	require.NoError(t, WithErrorHandler(errHandler)(b))

	funcName1 := runtime.FuncForPC(reflect.ValueOf(errHandler).Pointer()).Name()
	funcName2 := runtime.FuncForPC(reflect.ValueOf(b.defaultErrorHandler).Pointer()).Name()
	assert.Equal(t, funcName1, funcName2)
}

func TestBundle_Translator(t *testing.T) {
	provider := new(MockedProvider)
	b, err := NewBundle(WithProvider(provider))
	require.NoError(t, err)
	assert.NotNil(t, b.Translator("ru"), "even empty bundle returns some translator")

	provider.On("Get", language.Russian, "msg_id").Return("", errors.New("no message")) // call from actual translator
	provider.On("Get", language.Und, "msg_id").Return("", errors.New("no message"))     // call from fallback translator
	assert.Equal(t, "msg_id", b.Translator("ru").Trans("msg_id"), "even empty bundle returns some translator")

	b, err = NewBundle(
		WithDefaulLangFallback(language.English),
		WithLangFallback(language.Portuguese, language.Spanish),
		WithYamlProvider(fstest.MapFS{
			"messages.en.yaml": {Data: []byte("foo: en\nbar_id: enbar")},
			"messages.es.yaml": {Data: []byte("foo: es")},
		}),
	)
	require.NoError(t, err)

	assert.Equal(t, "en", b.Translator("en").Trans("foo"), "en loaded from file")
	assert.Equal(t, "en", b.Translator("pl").Trans("foo"), "fallback to defalt lang")
	assert.Equal(t, "en", b.Translator("FOO").Trans("foo"), "invalid lang fallback to defalt lang")
	assert.Equal(t, "enbar", b.Translator("es").Trans("bar_id"), "lang with dictionary fallbacks to default lang")
	assert.Equal(t, "es", b.Translator("pt").Trans("es"), "lang fallback")
	assert.Equal(t, "none_id", b.Translator("pl").Trans("none_id"), "dummy translator if nothing works")
	assert.Equal(t, "none_id", b.Translator("en").Trans("none_id"), "dummy translator if nothing works")
}
