package mf

import (
	"fmt"
	"testing"
	"testing/fstest"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

func Test_translator_Trans(t *testing.T) {
	tests := []struct {
		name    string
		msg     string
		lang    language.Tag
		args    []TranslationArg
		want    string
		wantErr bool
	}{
		{
			"simple text just works",
			"foo bar!",
			language.English,
			[]TranslationArg{Arg("foo", "wow")},
			"foo bar!",
			false,
		},
		{
			"plain args works fine with different spacing",
			"I {verb} { ArT } {\ntarGet3\n\t}.",
			language.English,
			[]TranslationArg{Arg("verb", "have"), Arg("ArT", "an"), Arg("tarGet3", "apple")},
			"I have an apple.",
			false,
		},
		{
			"error on missing plain arg",
			"so {bar}!",
			language.English,
			[]TranslationArg{Arg("foo", "wow")},
			"msg_id",
			true,
		},
		{
			"error on invalid message",
			"so {foo!",
			language.English,
			[]TranslationArg{Arg("foo", "wow")},
			"msg_id",
			true,
		},
		{
			"error on empty expression",
			"so {}!",
			language.English,
			[]TranslationArg{Arg("foo", "wow")},
			"msg_id",
			true,
		},
		{
			"error on invalid argument name",
			"so {#$%}!",
			language.English,
			[]TranslationArg{Arg("#$%", "wow")},
			"msg_id",
			true,
		},

		// select
		{
			"error on empty select",
			"so {foo, select}!",
			language.English,
			[]TranslationArg{Arg("foo", "wow")},
			"msg_id",
			true,
		},
		{
			"error on select without 'other' case",
			"so {foo, select, wow {good}}!",
			language.English,
			[]TranslationArg{Arg("foo", "wow")},
			"msg_id",
			true,
		},
		{
			"no error on select with 'other' case",
			"so {foo, select, wow {good} other {better}}!",
			language.English,
			[]TranslationArg{Arg("foo", "wow")},
			"so good!",
			false,
		},
		{
			"no error on select with 'other' case and new lines",
			`so {foo, select,
                wow {good}
                other {better}
            }!`,
			language.English,
			[]TranslationArg{Arg("foo", "wow")},
			"so good!",
			false,
		},
		{
			"no error on select with 'other' case and new lines",
			`{lang, select,
                en {no fallback}
                other {fallback to EN}
            }`,
			language.English,
			[]TranslationArg{Arg("lang", "en")},
			"no fallback",
			false,
		},

		// plural
		{
			"error on empty plural",
			"so {foo, plural}!",
			language.English,
			[]TranslationArg{Arg("foo", 1)},
			"msg_id",
			true,
		},
		{
			"error on plural without 'other' case",
			"so {foo, plural, one {good}}!",
			language.English,
			[]TranslationArg{Arg("foo", 1)},
			"msg_id",
			true,
		},
		{
			"no error on plural with 'other' case",
			"so {foo, plural, one {good} other {better}}!",
			language.English,
			[]TranslationArg{Arg("foo", 1)},
			"so good!",
			false,
		},
		{
			"strict cases first",
			"so {foo, plural, =1 {good} one {bad} other {better}}!",
			language.English,
			[]TranslationArg{Arg("foo", 1)},
			"so good!",
			false,
		},
		{
			"Welsh plural zero",
			`{foo, plural,
                zero {# cŵn, # cathod}
                one {# ci, # gath}
                two {# gi, # gath}
                few {# chi, # cath}
                many {# chi, # chath}
                other {# ci, # cath}
            }`,
			language.MustParse("cy"),
			[]TranslationArg{Arg("foo", 0)},
			"0 cŵn, 0 cathod",
			false,
		},
		{
			"Welsh plural one",
			`{foo, plural,
                zero {# cŵn, # cathod}
                one {# ci, # gath}
                two {# gi, # gath}
                few {# chi, # cath}
                many {# chi, # chath}
                other {# ci, # cath}
            }`,
			language.MustParse("cy"),
			[]TranslationArg{Arg("foo", 1)},
			"1 ci, 1 gath",
			false,
		},
		{
			"Welsh plural two",
			`{foo, plural,
                zero {# cŵn, # cathod}
                one {# ci, # gath}
                two {# gi, # gath}
                few {# chi, # cath}
                many {# chi, # chath}
                other {# ci, # cath}
            }`,
			language.MustParse("cy"),
			[]TranslationArg{Arg("foo", 2)},
			"2 gi, 2 gath",
			false,
		},
		{
			"Welsh plural few",
			`{foo, plural,
                zero {# cŵn, # cathod}
                one {# ci, # gath}
                two {# gi, # gath}
                few {# chi, # cath}
                many {# chi, # chath}
                other {# ci, # cath}
            }`,
			language.MustParse("cy"),
			[]TranslationArg{Arg("foo", 3)},
			"3 chi, 3 cath",
			false,
		},
		{
			"Welsh plural many",
			`{foo, plural,
                zero {# cŵn, # cathod}
                one {# ci, # gath}
                two {# gi, # gath}
                few {# chi, # cath}
                many {# chi, # chath}
                other {# ci, # cath}
            }`,
			language.MustParse("cy"),
			[]TranslationArg{Arg("foo", 6)},
			"6 chi, 6 chath",
			false,
		},
		{
			"Welsh plural other",
			`{foo, plural,
                zero {# cŵn, # cathod}
                one {# ci, # gath}
                two {# gi, # gath}
                few {# chi, # cath}
                many {# chi, # chath}
                other {# ci, # cath}
            }`,
			language.MustParse("cy"),
			[]TranslationArg{Arg("foo", 4)},
			"4 ci, 4 cath",
			false,
		},
		{
			"Offset used for plural cases #1",
			`{foo, plural, offset: 1
                =1 {one before offset}
                one {one after offset}
                other {other}
            }`,
			language.English,
			[]TranslationArg{Arg("foo", 1)},
			"one before offset",
			false,
		},
		{
			"Offset used for plural cases #2",
			`{foo, plural, offset: 1
                =1 {one before offset}
                one {one after offset}
                other {other}
            }`,
			language.English,
			[]TranslationArg{Arg("foo", 2)},
			"one after offset",
			false,
		},
		{
			"Offset used for plural cases #3",
			`{foo, plural, offset: 1
                =1 {one before offset}
                one {one after offset}
                other {# other}
            }`,
			language.English,
			[]TranslationArg{Arg("foo", 3)},
			"2 other",
			false,
		},

		// selectordinal
		{
			"Select ordinal basic usage #1",
			`You finished {place, selectordinal,
                one   {#st}
                two   {#nd}
                few   {#rd}
                other {#th}
            }!`,
			language.English,
			[]TranslationArg{Arg("place", 1)},
			"You finished 1st!",
			false,
		},
		{
			"Select ordinal basic usage #2",
			`You finished {place, selectordinal,
                one   {#st}
                two   {#nd}
                few   {#rd}
                other {#th}
            }!`,
			language.English,
			[]TranslationArg{Arg("place", 2)},
			"You finished 2nd!",
			false,
		},
		{
			"Select ordinal basic usage #3",
			`You finished {place, selectordinal,
                one   {#st}
                two   {#nd}
                few   {#rd}
                other {#th}
            }!`,
			language.English,
			[]TranslationArg{Arg("place", 9)},
			"You finished 9th!",
			false,
		},
		{
			"Select ordinal basic usage #4",
			`You finished {place, selectordinal,
                one   {#st}
                two   {#nd}
                few   {#rd}
                other {#th}
            }!`,
			language.English,
			[]TranslationArg{Arg("place", 23)},
			"You finished 23rd!",
			false,
		},
		{
			"percent number format",
			`this lib is {progress, number, percent} ready!`,
			language.English,
			[]TranslationArg{Arg("progress", 0.423)},
			"this lib is 42.3% ready!",
			false,
		},
		{
			"percent number format greater 100",
			`this lib is {progress, number, percent} ready!`,
			language.English,
			[]TranslationArg{Arg("progress", "1.3")},
			"this lib is 130% ready!",
			false,
		},
		{
			"integer number format",
			`big number {num, number, integer}!`,
			language.English,
			[]TranslationArg{Arg("num", 123456789)},
			"big number 123,456,789!",
			false,
		},
		{
			"integer number format",
			`big number {num, number, integer}!`,
			language.Spanish,
			[]TranslationArg{Arg("num", 123456789)},
			"big number 123.456.789!",
			false,
		},

		{
			"nested example",
			`{gender_of_host, select,
    female {{num_guests, plural, offset:1
        =0    {{host} does not give a party.}
        =1    {{host} invites {guest} to her party.}
        =2    {{host} invites {guest} and one other person to her party.}
        other {{host} invites {guest} and # other people to her party.}
    }}
    male {{num_guests, plural, offset:1
        =0    {{host} does not give a party.}
        =1    {{host} invites {guest} to his party.}
        =2    {{host} invites {guest} and one other person to his party.}
        other {{host} invites {guest} and # other people to his party.}
    }}
    other {{num_guests, plural, offset:1
        =0    {{host} does not give a party.}
        =1    {{host} invites {guest} to their party.}
        =2    {{host} invites {guest} and one other person to their party.}
        other {{host} invites {guest} and # other people to their party.}
    }}
}`,
			language.English,
			[]TranslationArg{
				Arg("gender_of_host", "female"),
				Arg("num_guests", 2),
				Arg("guest", "Sionia"),
				Arg("host", "Rina"),
			},
			"Rina invites Sionia and one other person to her party.",
			false,
		},

		{
			"datetime long",
			"Vostok-1 start {st, datetime, long}.",
			language.English,
			[]TranslationArg{Time("st", time.Date(1961, 4, 12, 6, 7, 3, 0, time.UTC))},
			"Vostok-1 start April 12, 1961 at 6:07:03 AM UTC.",
			false,
		},
		{
			"time long",
			"Vostok-1 landing time {st, time, medium}.",
			language.English,
			[]TranslationArg{Time("st", time.Date(1961, 4, 12, 7, 55, 0, 0, time.UTC))},
			"Vostok-1 landing time 7:55:00 AM.",
			false,
		},
		{
			"date long",
			"First step on the Moon on {st, date, long}.",
			language.English,
			[]TranslationArg{Time("st", time.Date(1969, 7, 21, 2, 56, 0, 0, time.UTC))},
			"First step on the Moon on July 21, 1969.",
			false,
		},
		{
			"escaping quotes not required",
			"foo ' {foo}.",
			language.English,
			[]TranslationArg{Arg("foo", "bar")},
			"foo ' bar.",
			false,
		},
		{
			"escaping single curly brace",
			"foo '{ {foo}.",
			language.English,
			[]TranslationArg{Arg("foo", "bar")},
			"foo { bar.",
			false,
		},
		{
			"escaping single quote",
			"foo '' {foo}.",
			language.English,
			[]TranslationArg{Arg("foo", "bar")},
			"foo ' bar.",
			false,
		},
		{
			"escaping single quote before curly brace",
			"foo ''{foo}''.",
			language.English,
			[]TranslationArg{Arg("foo", "bar")},
			"foo 'bar'.",
			false,
		},
		{
			"escaping",
			"foo '{ ''{foo} {num, plural, one {''#'' '# ' '{ one}, other {other}}.",
			language.English,
			[]TranslationArg{Arg("foo", "bar"), Arg("num", 1)},
			"foo { 'bar '1' # ' { one.",
			false,
		},
		{
			"escaping",
			"'{foo} is ''{foo}''",
			language.English,
			[]TranslationArg{Arg("foo", "bar")},
			"{foo} is 'bar'",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := new(MockedProvider)
			provider.On("Get", tt.lang, "msg_id").Return(tt.msg, nil)

			t.Log(len(provider.ExpectedCalls))

			var testErr error

			tr := &translator{
				provider: provider,
				errorHandler: func(err error, _ string, _ map[string]any) {
					t.Log(err.Error())
					testErr = err
				},
				lang: tt.lang,
			}

			got := tr.Trans("msg_id", tt.args...)
			assert.Equal(t, tt.want, got)

			if tt.wantErr {
				assert.Error(t, testErr)
			} else {
				assert.NoError(t, testErr)
			}
		})
	}
}

func Test_translator_SimilarLanguages(t *testing.T) {
	tests := []struct {
		name           string
		translatorLang language.Tag
		wantedLang     language.Tag
		want           string
	}{
		{
			"identical languages",
			language.English,
			language.English,
			"englishTranslation",
		},
		{
			"british to english",
			language.BritishEnglish,
			language.English,
			"englishTranslation",
		},
		{
			"english to british",
			language.English,
			language.BritishEnglish,
			"englishTranslation",
		},
		{
			"british to american",
			language.BritishEnglish,
			language.AmericanEnglish,
			"englishTranslation",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := NewBundle(
				WithDefaultLangFallback(language.Spanish),
				WithYamlProvider(fstest.MapFS{
					"messages.es.yaml": {Data: []byte("msg_id: spanishTranslation")},
					fmt.Sprintf("messages.%s.yaml", tt.translatorLang): {Data: []byte("msg_id: englishTranslation")},
				}),
			)
			require.NoError(t, err)

			var testErr error

			tr := b.Translator(tt.wantedLang.String())

			got := tr.Trans("msg_id")
			assert.Equal(t, tt.want, got)

			assert.NoError(t, testErr)
		})
	}
}

func Test_translator_Web(t *testing.T) {
	tests := []struct {
		name           string
		cookie         string
		acceptLanguage string
		want           string
	}{
		{
			"defaults",
			"",
			"",
			"spanishTranslation",
		},
		{
			"cookie",
			"en-GB",
			"",
			"englishTranslation",
		},
		{
			"Accept-Language",
			"",
			"da, en-gb;q=0.7, de;q=0.8",
			"germanTranslation",
		},
		{
			"Cookie wins",
			"fr-CA",
			"en",
			"frenchTranslation",
		},
		{
			"Fallback",
			"da",
			"it",
			"spanishTranslation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := NewBundle(
				WithDefaultLangFallback(language.Spanish),
				WithYamlProvider(fstest.MapFS{
					"messages.es.yaml": {Data: []byte("msg_id: spanishTranslation")},
					"messages.en.yaml": {Data: []byte("msg_id: englishTranslation")},
					"messages.de.yaml": {Data: []byte("msg_id: germanTranslation")},
					"messages.fr.yaml": {Data: []byte("msg_id: frenchTranslation")},
				}),
			)
			require.NoError(t, err)

			var testErr error

			tr := b.Translator(tt.cookie, tt.acceptLanguage)

			got := tr.Trans("msg_id")
			assert.Equal(t, tt.want, got)

			assert.NoError(t, testErr)
		})
	}
}

type MockedProvider struct {
	mock.Mock
}

func (m *MockedProvider) Get(lang language.Tag, path string) (string, error) {
	args := m.Called(lang, path)

	return args.String(0), args.Error(1)
}

type MockedDictionary struct {
	mock.Mock
}

func (m *MockedDictionary) Get(path string) (string, error) {
	args := m.Called(path)

	return args.String(0), args.Error(1)
}
