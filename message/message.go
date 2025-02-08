package message

import (
	"errors"
	"fmt"
	"strings"
)

type Evalable interface {
	Eval(ctx Context) (string, error)
}

type Message struct {
	fragments []Evalable
}

func (m *Message) Eval(ctx Context) (string, error) {
	var builder strings.Builder

	for _, child := range m.fragments {
		childRes, err := child.Eval(ctx)
		if err != nil {
			return "", err
		}

		builder.WriteString(childRes)
	}

	return builder.String(), nil
}

type Content string

func (c Content) Eval(_ Context) (string, error) {
	return string(c), nil
}

type PlainArg string

func (pa PlainArg) Eval(ctx Context) (string, error) {
	v, err := ctx.Any(string(pa))
	if err != nil {
		return "", err
	}

	return fmt.Sprint(v), nil
}

// {age, number, integer}
type FunctionArg struct {
	Name    string
	ArgName string
	Param   string
}

// { GENDER, select,
//
//	    male {He}
//	    female {She}
//	    other {They}
//	} liked this.
const DefaultCase = "other"

type Select struct {
	ArgName string
	Cases   map[string]Evalable
}

var ErrNoDefaultCase = errors.New("no default case")

func (s *Select) Eval(ctx Context) (string, error) {
	v, err := ctx.String(s.ArgName)
	if err == nil {
		c, ok := s.Cases[v]
		if ok {
			return c.Eval(ctx)
		}
	}

	c, ok := s.Cases[DefaultCase]
	if ok {
		return c.Eval(ctx)
	}

	return "", ErrNoDefaultCase
}

var (
	_ Evalable = (*Content)(nil)
	_ Evalable = (*Message)(nil)
	_ Evalable = (*PlainArg)(nil)
	_ Evalable = (*Select)(nil)
)
