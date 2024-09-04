package message

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/fullpipe/icu-mf/parse"
	"golang.org/x/text/language"
)

func Build(in parse.Message, lang language.Tag) (Evalable, error) {
	if len(in.Fragments) == 1 {
		return buildFragment(*in.Fragments[0], lang)
	}

	root := Message{
		fragments: make([]Evalable, 0, len(in.Fragments)),
	}

	for _, f := range in.Fragments {
		eval, err := buildFragment(*f, lang)
		if err != nil {
			// TODO: collect errors
			return nil, err
		}

		root.fragments = append(root.fragments, eval)
	}

	return &root, nil
}

func buildFragment(f parse.Fragment, lang language.Tag) (Evalable, error) {
	if len(f.Text) > 0 {
		return Content(f.Text), nil
	}

	if f.Octothorpe {
		return PlainArg("#"), nil
	}

	if f.PlainArg != nil {
		return PlainArg(f.PlainArg.Name), nil
	}

	if f.Func != nil {
		switch f.Func.Func {
		case "number":
			return buildNumber(f.Func, lang)
		default:
			return nil, errors.New("empty fragment")
		}
	}

	if f.Expr == nil {
		return nil, errors.New("empty fragment")
	}

	e := f.Expr

	if e.Name != "" && e.Func == "select" {
		return buildSelect(e, lang)
	}

	if e.Name != "" && e.Func == "plural" {
		return buildPlural(e, lang)
	}

	if e.Name != "" && e.Func == "selectordinal" {
		return buildPlural(e, lang)
	}

	return nil, errors.New("empty fragment")
}

func buildSelect(e *parse.Expr, lang language.Tag) (Evalable, error) {
	if e == nil || e.Name == "" || e.Func != "select" {
		return nil, errors.New("no a select expresion")
	}

	if len(e.Cases) == 0 {
		return nil, errors.New("empty select cases")
	}

	eval := Select{
		ArgName: e.Name,
		Cases:   make(map[string]Evalable, len(e.Cases)),
	}
	hasDefaultCase := false

	for _, c := range e.Cases {
		if c.Name == DefaultCase {
			hasDefaultCase = true
		}

		caseEval, err := Build(*c.Message, lang)
		if err != nil {
			// TODO: collect errors
			return nil, err
		}

		eval.Cases[c.Name] = caseEval
	}

	if !hasDefaultCase {
		return nil, errors.New("no 'other' case in select")
	}

	return &eval, nil
}

func buildPlural(e *parse.Expr, lang language.Tag) (Evalable, error) {
	if e == nil || e.Name == "" || (e.Func != "plural" && e.Func != "selectordinal") {
		return nil, errors.New("no a select expresion")
	}

	if len(e.Cases) == 0 {
		return nil, fmt.Errorf("no cases for {%s, %s ...}", e.Name, e.Func)
	}

	if e.Offset < 0 {
		return nil, errors.New("offset should be positive")
	}

	var eval *Plural
	switch e.Func {
	case "plural":
		eval = NewPlural(e.Name, lang, e.Offset)
	case "selectordinal":
		eval = NewSelectOrdinal(e.Name, lang, e.Offset)
	default:
		return nil, fmt.Errorf("invalid plural func {%s, %s ...}", e.Name, e.Func)
	}

	hasDefaultCase := false

	for _, c := range e.Cases {
		if c.Name == DefaultCase {
			hasDefaultCase = true
		}

		caseEval, err := Build(*c.Message, lang)
		if err != nil {
			// TODO: collect errors
			return nil, err
		}

		form, ok := strToFormMap[c.Name]
		if ok {
			eval.Cases[form] = caseEval
			continue
		}

		if c.Name[0] != '=' {
			return nil, fmt.Errorf("invalid plural case %s", c.Name)
		}

		caseNum, err := strconv.ParseUint(c.Name[1:], 10, 64)
		if err != nil {
			// TODO: collect errors
			return nil, err
		}

		eval.EqCases[caseNum] = caseEval
	}

	if !hasDefaultCase {
		return nil, errors.New("no 'other' case in plural")
	}

	return eval, nil

}

func buildNumber(e *parse.Func, lang language.Tag) (Evalable, error) {
	if e == nil || e.Func != "number" {
		return nil, errors.New("no a number function")
	}

	// format, ok := strToNumberFormatMap[e.Param]
	format, ok := strToNumberFormatMap[e.Param]
	if !ok {
		return nil, fmt.Errorf("number format %s not supported", e.Param)
	}

	return NewNumber(e.ArgName, format, lang), nil
}
