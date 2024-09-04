package mf

import (
	"strings"

	"github.com/fullpipe/icu-mf/message"
	"github.com/fullpipe/icu-mf/parse"
	"golang.org/x/exp/constraints"
	"golang.org/x/text/language"
)

type Translator interface {
	Trans(id string, args ...TranslationArg) string
}

type translator struct {
	dictionary   Dictionary
	fallback     Translator
	errorHandler ErrorHandler
	lang         language.Tag
}

func (tr *translator) Trans(id string, args ...TranslationArg) string {
	yaml, err := tr.dictionary.Get(id)
	if err != nil {
		if tr.fallback != nil {
			return tr.fallback.Trans(id, args...)
		}

		tr.errorHandler(err, id, nil)

		return id
	}

	parser := parse.NewParser()

	msg, err := parser.Parse("", strings.NewReader(yaml))
	if err != nil {
		tr.errorHandler(err, id, nil)

		return id
	}

	eval, err := message.Build(*msg, tr.lang)
	if err != nil {
		tr.errorHandler(err, id, nil)

		return id
	}

	ctx := make(message.Context, len(args))
	for _, arg := range args {
		arg(&ctx)
	}

	translation, err := eval.Eval(ctx)
	if err != nil {
		tr.errorHandler(err, id, ctx)

		return id
	}

	return translation
}

type TranslationArg func(ctx *message.Context)

type Argument interface {
	constraints.Integer | constraints.Float | ~string
}

func Arg[T Argument](name string, value T) TranslationArg {
	return func(ctx *message.Context) {
		ctx.Set(name, value)
	}
}
