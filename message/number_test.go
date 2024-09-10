package message

import (
	"testing"

	"golang.org/x/text/language"
)

func TestNumber_Eval(t *testing.T) {
	tests := []struct {
		name    string
		argName string
		format  NumberFormat
		ctx     Context
		want    string
		wantErr bool
	}{
		{
			"float to integer format",
			"n",
			IntegerNumberFormat,
			Context{"n": 3.14},
			"3",
			false,
		},
		{
			"string float to integer format",
			"n",
			IntegerNumberFormat,
			Context{"n": "3.14"},
			"3",
			false,
		},
		{
			"string int to integer format",
			"n",
			IntegerNumberFormat,
			Context{"n": "3"},
			"3",
			false,
		},
		{
			"float to percent format",
			"n",
			PercentNumberFormat,
			Context{"n": 0.0314},
			"3.14%",
			false,
		},
		{
			"no fraction if not required",
			"n",
			PercentNumberFormat,
			Context{"n": 0.5},
			"50%",
			false,
		},
		{
			"percent greater then 1",
			"n",
			PercentNumberFormat,
			Context{"n": "1.314"},
			"131.4%",
			false,
		},
		{
			"big int",
			"n",
			IntegerNumberFormat,
			Context{"n": 123456789},
			"123,456,789",
			false,
		},

		{
			"to decimal if no format",
			"n",
			NoneNumberFormat,
			Context{"n": 3.14},
			"3.14",
			false,
		},
		{
			"to decimal if invalid format",
			"n",
			NumberFormat(42),
			Context{"n": 3.14},
			"3.14",
			false,
		},
		{
			"error on invalid arg",
			"n",
			NoneNumberFormat,
			Context{"n": "foo"},
			"",
			true,
		},
		{
			"error on invalid arg",
			"n",
			IntegerNumberFormat,
			Context{"n": "foo"},
			"",
			true,
		},
		{
			"error on invalid arg",
			"n",
			PercentNumberFormat,
			Context{"n": "foo"},
			"",
			true,
		},
		{
			"error on invalid arg",
			"n",
			NumberFormat(42),
			Context{"n": "foo"},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NewNumber(tt.argName, tt.format, language.English)
			got, err := n.Eval(tt.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Number.Eval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Number.Eval() = %v, want %v", got, tt.want)
			}
		})
	}
}
