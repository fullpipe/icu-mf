package message

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
)

type Number struct {
	ArgName string
	Format  NumberFormat
	Lang    language.Tag
	printer *message.Printer
}

var strToNumberFormatMap = map[string]NumberFormat{
	"":        NoneNumberFormat,
	"integer": IntegerNumberFormat,
	"percent": PercentNumberFormat,
}

func NewNumber(argName string, format NumberFormat, lang language.Tag) *Number {
	return &Number{
		ArgName: argName,
		Format:  format,
		Lang:    lang,
		printer: message.NewPrinter(lang),
	}
}

type NumberFormat int

const (
	NoneNumberFormat NumberFormat = iota
	IntegerNumberFormat
	PercentNumberFormat
)

func (n Number) Eval(ctx Context) (string, error) {
	switch n.Format {
	case NoneNumberFormat:
		v, err := ctx.Float64(n.ArgName)
		if err != nil {
			return "", err
		}

		return n.printer.Sprint(number.Decimal(v)), nil
	case IntegerNumberFormat:
		v, err := ctx.Int64(n.ArgName)
		if err != nil {
			return "", err
		}

		return n.printer.Sprint(number.Decimal(v)), nil
	case PercentNumberFormat:
		v, err := ctx.Float64(n.ArgName)
		if err != nil {
			return "", err
		}

		return n.printer.Sprint(number.Percent(v, number.MaxFractionDigits(2))), nil
	}

	v, err := ctx.Float64(n.ArgName)
	if err != nil {
		return "", err
	}

	return n.printer.Sprint(number.Decimal(v)), nil
}

var _ Evalable = (*Number)(nil)
