package message

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"golang.org/x/text/feature/plural"
	"golang.org/x/text/language"
)

// { COUNT, plural,
//
//	    =0 {There are no results.}
//	    one {There is one result.}
//	    other {There are # results.}
//	}
type Plural struct {
	ArgName  string
	Lang     language.Tag
	Offset   int
	EqCases  map[uint64]Evalable
	Cases    map[plural.Form]Evalable
	formFunc func(lang language.Tag, i int, v int, w int, f int, t int) plural.Form
}

// NewPlural creates a new Plural for cardinal plurals.
func NewPlural(argName string, lang language.Tag, offset int) *Plural {
	return newPlural(argName, lang, offset, plural.Cardinal.MatchPlural)
}

// NewSelectOrdinal creates a new Plural for ordinal plurals.
func NewSelectOrdinal(argName string, lang language.Tag, offset int) *Plural {
	return newPlural(argName, lang, offset, plural.Ordinal.MatchPlural)
}

// newPlural is a helper function to create a new Plural.
func newPlural(
	argName string,
	lang language.Tag,
	offset int,
	formFunc func(lang language.Tag, i int, v int, w int, f int, t int) plural.Form,
) *Plural {
	return &Plural{
		ArgName:  argName,
		Lang:     lang,
		Cases:    map[plural.Form]Evalable{},
		EqCases:  map[uint64]Evalable{},
		Offset:   offset,
		formFunc: formFunc,
	}
}

var strToFormMap = map[string]plural.Form{
	DefaultCase: plural.Other,
	"zero":      plural.Zero,
	"one":       plural.One,
	"two":       plural.Two,
	"few":       plural.Few,
	"many":      plural.Many,
}

type PluralCase struct {
	Key  string
	Eval Evalable
}

type pm struct {
	i          uint64
	v, w, f, t int
}

func (p *Plural) Eval(ctx Context) (string, error) {
	num, err := ctx.Any(p.ArgName)
	if err != nil {
		return "", err
	}

	np, err := toPluralForm(num)
	if err != nil {
		return "", err
	}

	// FIX: plurals could be nested
	ctx.Set("#", num)

	if np.t == 0 {
		c, ok := p.EqCases[np.i]
		if ok {
			return c.Eval(ctx)
		}
	}

	if p.Offset > 0 {
		offset := uint64(p.Offset) //nolint: gosec
		if offset < np.i {
			np.i -= offset
		} else {
			np.i = 0
		}

		// FIX: num could be float
		ctx.Set("#", np.i)
	}

	pi := np.i
	if pi > math.MaxInt { // Use integer division for large numbers to avoid int overflow in plural.MatchPlural
		pi /= 10_000_000
	}
	form := p.formFunc(p.Lang, int(pi), np.v, np.w, np.f, np.t) //nolint: gosec

	c, ok := p.Cases[form]
	if ok {
		return c.Eval(ctx)
	}

	c, ok = p.Cases[plural.Other]
	if ok {
		return c.Eval(ctx)
	}

	return "", ErrNoDefaultCase
}

func toPluralForm(num any) (pm, error) {
	switch i := num.(type) {
	case int:
		if i < 0 {
			i = -i
		}
		return pm{i: uint64(i)}, nil //nolint: gosec
	case int8:
		if i < 0 {
			i = -i
		}
		return pm{i: uint64(i)}, nil //nolint: gosec
	case int16:
		if i < 0 {
			i = -i
		}
		return pm{i: uint64(i)}, nil //nolint: gosec
	case int32:
		if i < 0 {
			i = -i
		}
		return pm{i: uint64(i)}, nil //nolint: gosec
	case int64:
		if i < 0 {
			i = -i
		}
		return pm{i: uint64(i)}, nil //nolint: gosec
	case uint:
		return pm{i: uint64(i)}, nil
	case uint8:
		return pm{i: uint64(i)}, nil
	case uint16:
		return pm{i: uint64(i)}, nil
	case uint32:
		return pm{i: uint64(i)}, nil
	case uint64:
		return pm{i: i}, nil
	case float32:
		return parseString(strconv.FormatFloat(float64(i), 'f', -1, 32))
	case float64:
		return parseString(strconv.FormatFloat(i, 'f', -1, 64))
	case string:
		return parseString(i)
	default:
		return pm{}, fmt.Errorf("unable convert %v to plural form", num)
	}
}

func parseString(str string) (pm, error) {
	if strings.HasPrefix(str, "-") { // Remove negative if it is there
		str = str[1:]
	}

	parts := strings.SplitN(str, ".", 2)
	pmi, err := strconv.ParseUint(parts[0], 10, 32)

	if err != nil {
		return pm{}, fmt.Errorf("unable to parse uint part %s of %s", parts[0], str)
	}

	if len(parts) == 1 {
		return pm{i: uint64(pmi)}, nil
	}

	decimalPart := parts[1]
	decimalPartTrail := strings.TrimRight(decimalPart, "0")
	pmf, err := strconv.ParseUint(decimalPart, 10, 32)

	if err != nil {
		return pm{}, fmt.Errorf("unable to parse decimalPart part %s of %s", decimalPart, str)
	}

	pmt := uint64(0)
	if decimalPartTrail != "" {
		pmt, err = strconv.ParseUint(decimalPartTrail, 10, 32)
		if err != nil {
			return pm{}, fmt.Errorf("unable to parse decimalPartTrail part %s of %s", decimalPartTrail, str)
		}
	}

	return pm{
		i: uint64(pmi),
		v: len(decimalPart),
		w: len(decimalPartTrail),
		f: int(pmf), //nolint: gosec
		t: int(pmt), //nolint: gosec
	}, nil
}

var _ Evalable = (*Plural)(nil)
