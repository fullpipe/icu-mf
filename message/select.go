package message

import "errors"

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

var _ Evalable = (*Select)(nil)
