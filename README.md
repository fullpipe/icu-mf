# ICU MessageFormat

Messages in your application are never static. They have variables, pluralization, and formatting.
To work with them easily, use [ICU MessageFormat](https://unicode-org.github.io/icu/userguide/format_parse/messages/).

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
    mf.WithDefaulLangFallback(language.English),

    // We could fine-tune fallbacks for some languages
    mf.WithLangFallback(language.BritishEnglish, language.English),
    mf.WithLangFallback(language.Portuguese, language.Spanish),


    // We assume that the translated messages are mostly correct.
    // However, if any errors occur during translation,
    // they will be directed to the error handler.
    mf.WithErrorHandler(func(err error, key string, ctx map[string]any) {
        slog.Error(err.Error(), slog.String("key", key), slog.Any("ctx", ctx))

        // or
        //panic(err)
    }),
)

if err != nil {
    log.Fatal(err)
}

// Load all yaml files in directory as messages
err = bundle.LoadDir(messagesDir)
if err != nil {
    log.Fatal(err)
}
```

Translate messages by their ID

```go
tr := bundle.Translator("en")

slog.Info(
    tr.Trans("say_hello", mf.Arg("name", "Bob"))
)

trEs := bundle.Translator("es")

slog.Info(
    trEs.Trans("say_hello", mf.Arg("name", "Aníbal"))
)
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

	tr := bundle.Translator("es")

	slog.Info(tr.Trans("say_hello", mf.Arg("name", "Bob")))
}
```
</details>

## MessageFormat overview

### Placeholders

MessageFormat allows to use placeholders in your messages.

```yaml
# translations/messages.en.yaml

say_hello: 'Hello {name}!'
```

Everything in `{...}` will be processed as an argument and will be replaced by the provided context arguments.

```go
log.Println(
    tr.Trans("say_hello", mf.Arg("name", "Bob"))
) // prints "Hello, Bob!"
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
log.Println(
    tr.Trans(
        "invitation.title",
        mf.Arg("organizer_name", "Ryan"),
        mf.Arg("organizer_gender", "male"),
    )
) // prints "Ryan has invited you to his party!"

log.Println(
    tr.Trans(
        "invitation.title",
        mf.Arg("organizer_name", "John & Jane"),
        mf.Arg("organizer_gender", "multiple"),
    )
) // prints "John & Jane have invited you to their party!"

log.Println(
    tr.Trans(
        "invitation.title",
        mf.Arg("organizer_name", "ACME Company"),
        mf.Arg("organizer_gender", "not_applicable"),
    )
) // prints "ACME Company has invited you to their party!"
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
log.Println(
    tr.Trans("num_of_apples", mf.Arg("apples", 5))
) // prints "I have 5 apples!"

// for RU
log.Println(
    tr.Trans("num_of_apples", mf.Arg("apples", 3))
) // prints "У меня 3 яблока"

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
log.Println(
    tr.Trans(
        "party_status",
        mf.Arg("num_guests", 1),
        mf.Arg("host", "Rogna"),
        mf.Arg("guest", "Azog"),
    )
) // prints "Rogna invites Azog to her party."

log.Println(
    tr.Trans(
        "party_status",
        mf.Arg("num_guests", 5),
        mf.Arg("host", "Rogna"),
        mf.Arg("guest", "Azog"),
    )
) // prints "Rogna invites Azog and 4 other people to her party."
```

First, we compare `num_guests` with the strict cases `=0`, `=1`, and `=2`.
If nothing matches, we subtract the `offset`, `num_guests = num_guests - offset`,
and then determine the plural case based on the result.

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
log.Println(
    tr.Trans("finish_place", mf.Arg("place", 1))
) // prints "You finished 1st!"
log.Println(
    tr.Trans("finish_place", mf.Arg("place", 9))
) // prints "You finished 9th!"
log.Println(
    tr.Trans("finish_place", mf.Arg("place", 43))
) // prints "You finished 43rd!"
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
log.Println(
    tr.Trans("big_num", mf.Arg("num", 123456789))
) // prints "big number 123,456,789!"

log.Println(
    bundle.Translator("es").Trans("big_num", mf.Arg("num", 123456789))
) // prints "gran numero 123.456.789!"
```

##### Percent

```yaml
# translations/messages.en.yaml

test_cover: we got {cover, number, percent} test coverage!
```

```go
log.Println(
    tr.Trans("test_cover", mf.Arg("cover", 0.42))
) // prints "we got 42% test coverage!"

log.Println(
    tr.Trans("test_cover", mf.Arg("cover", 1))
) // prints "we got 100% test coverage!"
```

#### Date and Time

TODO
