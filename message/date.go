package message

import (
	"golang.org/x/text/language"
)

type Datetime struct {
	argName string
	lang    language.Tag
	format  string
}

type DatetimeFormat int

const (
	NoneDatetimeFormat   DatetimeFormat = iota
	ShortDatetimeFormat                 // 1/2/06, 3:04 PM
	MediumDatetimeFormat                // Jan 2, 2006, 3:04:05 PM
	LongDatetimeFormat                  // January 2, 2006 at 3:04:05 PM UTC

	// TODO: go std time formater does not have full forms of local abbreviation
	FullDatetimeFormat // Monday, January 2, 2006 at 3:04:05 PM Coordinated Universal Time
)

var dateFormat = map[DatetimeFormat]string{
	NoneDatetimeFormat:   "",
	ShortDatetimeFormat:  "1/2/06",
	MediumDatetimeFormat: "Jan 2, 2006",
	LongDatetimeFormat:   "January 2, 2006",
	FullDatetimeFormat:   "Monday, January 2, 2006",
}

var timeFormat = map[DatetimeFormat]string{
	NoneDatetimeFormat:   "",
	ShortDatetimeFormat:  "3:04 PM",
	MediumDatetimeFormat: "3:04:05 PM",
	LongDatetimeFormat:   "3:04:05 PM MST",
	FullDatetimeFormat:   "3:04:05 PM MST",
}

var datetimeJoin = map[DatetimeFormat]string{
	NoneDatetimeFormat:   "",
	ShortDatetimeFormat:  ", ",
	MediumDatetimeFormat: ", ",
	LongDatetimeFormat:   " at ",
	FullDatetimeFormat:   " at ",
}

func NewDatetime(argName string, format DatetimeFormat, lang language.Tag) *Datetime {
	strFormat := dateFormat[format] + datetimeJoin[format] + timeFormat[format]

	return &Datetime{
		argName: argName,
		lang:    lang,
		format:  strFormat,
	}
}

func NewTime(argName string, format DatetimeFormat, lang language.Tag) *Datetime {
	return &Datetime{
		argName: argName,
		lang:    lang,
		format:  timeFormat[format],
	}
}

func NewDate(argName string, format DatetimeFormat, lang language.Tag) *Datetime {
	return &Datetime{
		argName: argName,
		lang:    lang,
		format:  dateFormat[format],
	}
}

func (dt Datetime) Eval(ctx Context) (string, error) {
	d, err := ctx.Time(dt.argName)
	if err != nil {
		return "", err
	}

	return d.Format(dt.format), nil
}

var _ Evalable = (*Datetime)(nil)
