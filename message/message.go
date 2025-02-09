package message

import (
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

	switch v := v.(type) {
	case string:
		return v, nil
	case int, int32, int64, float32, float64:
		return fmt.Sprintf("%v", v), nil
	default:
		return fmt.Sprint(v), nil
	}
}

// {age, number, integer}
type FunctionArg struct {
	Name    string
	ArgName string
	Param   string
}

var (
	_ Evalable = (*Content)(nil)
	_ Evalable = (*Message)(nil)
	_ Evalable = (*PlainArg)(nil)
)
