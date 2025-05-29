# ICU MessageFormat for Golang

[![test](https://github.com/fullpipe/icu-mf/actions/workflows/test.yml/badge.svg)](https://github.com/fullpipe/icu-mf/actions/workflows/test.yml)
[![codecov](https://codecov.io/github/fullpipe/icu-mf/graph/badge.svg?token=W6C02M3BFQ)](https://codecov.io/github/fullpipe/icu-mf)
[![lint](https://github.com/fullpipe/icu-mf/actions/workflows/lint.yml/badge.svg)](https://github.com/fullpipe/icu-mf/actions/workflows/lint.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/fullpipe/icu-mf.svg)](https://pkg.go.dev/github.com/fullpipe/icu-mf)

Messages in your application are never static. They have variables, pluralization, and formatting.
To translate them easily, use [ICU MessageFormat](https://unicode-org.github.io/icu/userguide/format_parse/messages/).

## Why?

There is a great package for translations called [nicksnyder/go-i18n](https://github.com/nicksnyder/go-i18n).
However, once I had a lot of translations in chatbots, it started to feel cumbersome.

So, I tried to make translations simpler. Now, instead:

```go
localizer.Localize(&i18n.LocalizeConfig{
    DefaultMessage: &i18n.Message{
        ID: "PersonCats",
        One: "{{.Name}} has {{.Count}} cat.",
        Other: "{{.Name}} has {{.Count}} cats.",
    },
    TemplateData: map[string]interface{}{
        "Name": "Nick",
        "Count": 2,
    },
    PluralCount: 2,
}) // Nick has 2 cats.
```

I got:

```go
tr.Trans("person.cats", mf.Arg("name", "Nick"), mf.Arg("cats_num", 2))
// Nick has 2 cats.
```

## Usage

Import package

```go
import "github.com/fullpipe/icu-mf/mf"
```

Locate messages with go:embed

```go
//go:embed var/messages.*.yaml
var messagesDir embed.FS

// or you could load messages dynamically
messagesDir := os.DirFS("var")
```

Create translations bundle

```go
bundle, err := mf.NewBundle(
    // If not possible to find a message for the specific language, fallback to English (EN)
    mf.WithDefaultLangFallback(language.English),

    // We could fine-tune fallbacks for some languages
    mf.WithLangFallback(language.BritishEnglish, language.English),
    mf.WithLangFallback(language.Portuguese, language.Spanish),

    // Load all yaml files in directory as messages
    mf.WithYamlProvider(messagesDir),

    // or you could use your own custom message provider
    // mf.WithProvider(sqlMessageProvider),

    // We assume that the translated messages are mostly correct.
    // However, if any errors occur during translation,
    // they will be directed to the error handler.
    mf.WithErrorHandler(func(err error, id string, ctx map[string]any) {
        slog.Error(err.Error(), slog.String("id", id), slog.Any("ctx", ctx))

        // or
        //panic(err)
    }),
)

if err != nil {
    log.Fatal(err)
}
```

Translate messages by their ID

```go
tr := bundle.Translator("en")

tr.Trans("invitation.status",
    mf.Arg("gender_of_host", "female"),
    mf.Arg("num_guests", 5),
    mf.Arg("guest", "Sionia"),
    mf.Arg("host", "Rina"),
) // Rina invites Sionia and 4 other people to her party.


trEs := bundle.Translator("es")

trEs.Trans("say_hello", mf.Arg("name", "Aníbal"))
```

Find the best translation from user preferences

```go

handler( w http.ResponseWriter, r *http.Request) {
    langCookie, _ := r.Cookie("lang")
	// langCookie = "es-AR"
	acceptHeader := r.Header.Get("Accept-Language")
	// acceptHeader = "da, en-gb;q=0.8, en;q=0.7"
	
	tr := bundle.Translator(langCookie, acceptHeader)
	// tr will translate to the closest matching language, e.g. Spanish, "es"
}
```

<details>
  <summary>Full example</summary>

```go
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
		mf.WithDefaultLangFallback(language.English),

		mf.WithLangFallback(language.BritishEnglish, language.English),
		mf.WithLangFallback(language.Portuguese, language.Spanish),

		mf.WithYamlProvider(messagesDir),

		mf.WithErrorHandler(func(err error, id string, ctx map[string]any) {
			slog.Error(err.Error(), slog.String("id", id), slog.Any("ctx", ctx))
		}),
	)

	if err != nil {
		log.Fatal(err)
	}

	tr := bundle.Translator("es")

	slog.Info(tr.Trans("say_hello", mf.Arg("name", "Bob")))
}
```

</details>

### YAML

YAML allows you to organize your translations in a tree-like structure.

```yaml
user:
  profile:
    name: My name is {name}
    age: I'm {age, plural, one {# year} other {# years}} old
  account_form:
    username_field: 'Enter your username:'
    error: >-
      {name, select
          required {specify {field}}
          min {{field} requires at least 10 chars}
          other {some unknown error with {field}}
      }

payments: ...

server:
  http:
    404: Page not found
    503: Oops!
```

And you get messages by their "path"

```go
tr.Trans("user.profile.age", mf.Arg("age", 42))
tr.Trans(
    "user.account_form.error",
    mf.Arg("name", "min"), mf.Arg("field", "description"),
)
```

### Escaping

Sometimes you need to print `{`, `'`, or `#`. You could escape them with `'` char.

```yaml
# translations/messages.en.yaml
escape: "'{foo} is ''{foo}''"
```

```go
tr.Trans("escape", mf.Arg("foo", "bar"))
// {foo} is 'bar'
```

## MessageFormat overview

### Placeholders

MessageFormat allows to use placeholders in your messages.

```yaml
# translations/messages.en.yaml

say_hello: 'Hello, {name}!'
```

Everything in `{...}` will be processed as an argument and will be replaced by the provided context arguments.

```go
tr.Trans("say_hello", mf.Arg("name", "Bob"))
// Hello, Bob!
```

### Simple select

```yaml
# translations/messages.en.yaml

# the 'other' key is required, and is selected if no other case matches
invitation:
  title: >-
    {organizer_gender, select,
        female   {{organizer_name} has invited you to her party!}
        male     {{organizer_name} has invited you to his party!}
        multiple {{organizer_name} have invited you to their party!}
        other    {{organizer_name} has invited you to their party!}
    }
  body: ...
```

```go
tr.Trans(
    "invitation.title",
    mf.Arg("organizer_name", "Ryan"),
    mf.Arg("organizer_gender", "male"),
) // Ryan has invited you to his party!

tr.Trans(
    "invitation.title",
    mf.Arg("organizer_name", "John & Jane"),
    mf.Arg("organizer_gender", "multiple"),
) // John & Jane have invited you to their party!

tr.Trans(
    "invitation.title",
    mf.Arg("organizer_name", "ACME Company"),
    mf.Arg("organizer_gender", "not_applicable"),
) // ACME Company has invited you to their party!
```

As you can see, the `{...}` syntax behaves differently here:

1. The first `{organizer_gender, select, ...}` block starts "code" mode, meaning `organizer_gender` is processed as a variable.
2. The inner `{... has invited you to her party!}` block switches to "literal" mode, meaning the text inside is processed as sub-message.
3. Inside this block, `{organizer_name}` starts "code" mode again, allowing `organizer_name` to be processed as a variable.

### Pluralization

There is another function, `plural`, similar to `select`. It allows you to handle pluralization in your messages (e.g., `There are 3 apples` vs. `There is one apple`).

```yaml
# translations/messages.en.yaml
num_of_apples: >-
  {apples, plural,
      =0    {I don't have an apple}
      one   {I have one apple}
      other {I have # apples!}
  }
```

Pluralization rules are actually quite complex and differ for each language.
For instance, Russian uses different plural forms for numbers ending with 1;
numbers ending with 2, 3, or 4; numbers ending with 5, 6, 7, 8, or 9;
and even some exceptions to this!

To properly translate plural forms, the possible cases in the `plural` function
are also different for each language. For instance, Russian has `one`, `few`, `many`,
and `other`, while English has only `one` and `other`.
The full list of possible cases can be found
in Unicode's [Language Plural Rules](https://www.unicode.org/cldr/charts/43/supplemental/language_plural_rules.html) document.

By prefixing with `=`, you can match exact values (like 0 in the above example).

```yaml
# translations/messages.ru.yaml

num_of_apples: >-
  {apples, plural,
      =0    {У меня нет яблок}
      =1    {У меня одно яблоко}
      one   {У меня # яблоко}
      few   {У меня # яблока}
      many  {У меня # яблок}
      other {У меня # яблок}
  }
```

The usage of this string is the same as with `select`:

```go
// for EN

tr.Trans("num_of_apples", mf.Arg("apples", 5))
// I have 5 apples!

// for RU
trRU.Trans("num_of_apples", mf.Arg("apples", 3))
// У меня 3 яблока

```

You can use the `#` placeholder to display the pluralized number.

#### Offset

You can also set an `offset` variable to determine whether the pluralization should be adjusted. For example, in sentences like `You and # other people` / `You and # other person`.

```yaml
# translations/messages.en.yaml
party_status: >-
  {num_guests, plural, offset:1
      =0    {{host} does not give a party.}
      =1    {{host} invites {guest} to her party.}
      =2    {{host} invites {guest} and one other person to her party.}
      other {{host} invites {guest} and # other people to her party.}
  }
```

```go
tr.Trans(
    "party_status",
    mf.Arg("num_guests", 1),
    mf.Arg("host", "Rogna"),
    mf.Arg("guest", "Azog"),
) // Rogna invites Azog to her party.

tr.Trans(
    "party_status",
    mf.Arg("num_guests", 5),
    mf.Arg("host", "Rogna"),
    mf.Arg("guest", "Azog"),
) // Rogna invites Azog and 4 other people to her party.
```

First, we compare `num_guests` with the strict cases `=0`, `=1`, and `=2`.
If nothing matches, we subtract the `offset`, `num_guests = num_guests - offset`,
and then determine the plural case based on the result.

#### Nesting

You could make pretty complex nested messages if needed.

```yaml
# translations/messages.en.yaml

invitation_status: >-
  {gender_of_host, select,
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
  }
```

#### Inline

Cases in `plural`, `select` or `selectordinal` could be inlined

```yaml
# translations/messages.en.yaml
num_of_apples: 'There {apples, plural, =0 {are no} one {is one} other {are # apples}} apples'
```

### Additional Functions

#### Ordinal

Similar to `plural`, `selectordinal` allows you to use numbers as ordinal scale:

```yaml
# translations/messages.en.yaml
finish_place: >-
  You finished {place, selectordinal,
      one   {#st}
      two   {#nd}
      few   {#rd}
      other {#th}
  }!
```

```go
tr.Trans("finish_place", mf.Arg("place", 1))
// You finished 1st!

tr.Trans("finish_place", mf.Arg("place", 9))
// You finished 9th!

tr.Trans("finish_place", mf.Arg("place", 43))
// You finished 43rd!
```

The possible cases for this are also shown in Unicode's [Language Plural Rules](https://www.unicode.org/cldr/charts/43/supplemental/language_plural_rules.html) document.

#### Numbers

There are some minor functions to work with numbers and dates.

##### Integer

```yaml
# translations/messages.en.yaml

big_num: big number {num, number, integer}!


# translations/messages.es.yaml

big_num: gran numero {num, number, integer}!
```

```go
tr.Trans("big_num", mf.Arg("num", 123456789))
// big number 123,456,789!

bundle.Translator("es").Trans("big_num", mf.Arg("num", 123456789))
// gran numero 123.456.789!
```

##### Percent

```yaml
# translations/messages.en.yaml

test_cover: we got {cover, number, percent} test coverage!
```

```go
tr.Trans("test_cover", mf.Arg("cover", 0.42))
// we got 42% test coverage!

tr.Trans("test_cover", mf.Arg("cover", 1))
// we got 100% test coverage!
```

#### Date and Time

There are `date`, `time`, and `datetime` functions to format `time.Time` arguments.
Additionally, there are four different formats: `short`, `medium`, `long`, and `full`.

```yaml
# translations/messages.en.yaml

vostok:
  start: Vostok-1 start {start_date, datetime, long}.
  landing: Vostok-1 landing time {land_time, time, medium}.
apollo:
  step: First step on the Moon on {step_date, date, long}.
```

```go
start := time.Date(1961, 4, 12, 6, 7, 3, 0, time.UTC)
land := time.Date(1961, 4, 12, 7, 55, 0, 0, time.UTC)
step := time.Date(1969, 7, 21, 2, 56, 0, 0, time.UTC)

tr.Trans("vostok.start", mf.Time("start_date", start))
// Vostok-1 start April 12, 1961 at 6:07:03 AM UTC.

tr.Trans("vostok.landing", mf.Time("land_time", land))
// Vostok-1 landing time 7:55:00 AM.

tr.Trans("apollo.step", mf.Time("step_date", step))
// First step on the Moon on July 21, 1969.
```
