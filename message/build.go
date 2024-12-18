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

	root := &Message{
		fragments: make([]Evalable, 0, len(in.Fragments)),
	}

	for _, f := range in.Fragments {
		eval, err := buildFragment(*f, lang)
		if err != nil {
			return nil, err
		}
		root.fragments = append(root.fragments, eval)
	}

	return root, nil
}

func buildFragment(f parse.Fragment, lang language.Tag) (Evalable, error) {
	switch {
	case len(f.Escaped) > 0:
		return Content(f.Escaped[1:]), nil
	case len(f.Text) > 0:
		return Content(f.Text), nil
	case f.Octothorpe:
		return PlainArg("#"), nil
	case f.PlainArg != nil:
		return PlainArg(f.PlainArg.Name), nil
	case f.Func != nil:
		return buildFunc(f.Func, lang)
	case f.Expr != nil:
		return buildExpr(f.Expr, lang)
	default:
		return nil, errors.New("empty fragment")
	}
}

func buildFunc(f *parse.Func, lang language.Tag) (Evalable, error) {
	switch f.Func {
	case "number":
		return buildNumber(f, lang)
	case "date", "time", "datetime":
		return buildDatetime(f, lang)
	default:
		return nil, fmt.Errorf("unsupported function: %s", f.Func)
	}
}

func buildExpr(e *parse.Expr, lang language.Tag) (Evalable, error) {
	switch e.Func {
	case "select":
		return buildSelect(e, lang)
	case "plural", "selectordinal":
		return buildPlural(e, lang)
	default:
		return nil, fmt.Errorf("unsupported expression: %s", e.Func)
	}
}

func buildSelect(e *parse.Expr, lang language.Tag) (Evalable, error) {
	if e == nil || e.Name == "" || e.Func != "select" {
		return nil, errors.New("invalid select expression")
	}

	if len(e.Cases) == 0 {
		return nil, errors.New("empty select cases")
	}

	eval := &Select{
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
			return nil, err
		}

		eval.Cases[c.Name] = caseEval
	}

	if !hasDefaultCase {
		return nil, errors.New("no 'other' case in select")
	}

	return eval, nil
}

func buildPlural(e *parse.Expr, lang language.Tag) (Evalable, error) {
	if e == nil || e.Name == "" || (e.Func != "plural" && e.Func != "selectordinal") {
		return nil, errors.New("invalid plural expression")
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
			return nil, err
		}

		if form, ok := strToFormMap[c.Name]; ok {
			eval.Cases[form] = caseEval
		} else if c.Name[0] == '=' {
			caseNum, err := strconv.ParseUint(c.Name[1:], 10, 64)
			if err != nil {
				return nil, err
			}
			eval.EqCases[caseNum] = caseEval
		} else {
			return nil, fmt.Errorf("invalid plural case %s", c.Name)
		}
	}

	if !hasDefaultCase {
		return nil, errors.New("no 'other' case in plural")
	}

	return eval, nil
}

func buildNumber(f *parse.Func, lang language.Tag) (Evalable, error) {
	format, ok := strToNumberFormatMap[f.Param]
	if !ok {
		return nil, fmt.Errorf("number format %s not supported", f.Param)
	}

	return NewNumber(f.ArgName, format, lang), nil
}

func buildDatetime(f *parse.Func, lang language.Tag) (Evalable, error) {
	format, ok := strToDatetimeFormatMap[f.Param]
	if !ok {
		return nil, fmt.Errorf("date format %s not supported", f.Param)
	}

	switch f.Func {
	case "date":
		return NewDate(f.ArgName, format, lang), nil
	case "time":
		return NewTime(f.ArgName, format, lang), nil
	case "datetime":
		return NewDatetime(f.ArgName, format, lang), nil
	default:
		return nil, fmt.Errorf("unsupported datetime function: %s", f.Func)
	}
}
