package message

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/feature/plural"
	"golang.org/x/text/language"
)

func TestPlural_Eval(t *testing.T) {
	type fields struct {
		ArgName string
		Lang    language.Tag
		Offset  int
		EqCases map[uint64]Evalable
		Cases   map[plural.Form]Evalable
	}
	tests := []struct {
		name    string
		fields  fields
		ctx     Context
		want    string
		wantErr bool
	}{
		{
			"basic usage",
			fields{
				"count",
				language.English,
				0,
				nil,
				map[plural.Form]Evalable{
					plural.One: Content("one"),
				},
			},
			Context{"count": 1},
			"one",
			false,
		},
		{
			"error if no argument in context",
			fields{
				"count",
				language.English,
				0,
				nil,
				map[plural.Form]Evalable{
					plural.One: Content("one"),
				},
			},
			Context{"none": 1},
			"",
			true,
		},
		{
			"error if argument is not plural",
			fields{
				"count",
				language.English,
				0,
				nil,
				map[plural.Form]Evalable{
					plural.One:   Content("one"),
					plural.Other: Content("other"),
				},
			},
			Context{"count": "foo"},
			"",
			true,
		},
		{
			"cases like =1 are in priority",
			fields{
				"count",
				language.English,
				0,
				map[uint64]Evalable{
					1: Content("eq one"),
				},
				map[plural.Form]Evalable{
					plural.One:   Content("one"),
					plural.Other: Content("other"),
				},
			},
			Context{"count": 1},
			"eq one",
			false,
		},
		{
			"different forms of plural [string]",
			fields{
				"count",
				language.Afrikaans,
				0,
				nil,
				map[plural.Form]Evalable{
					plural.One:   Content("one"),
					plural.Other: Content("other"),
				},
			},
			Context{"count": "1.0"},
			"one",
			false,
		},
		{
			"different forms of plural [float]",
			fields{
				"count",
				language.English,
				0,
				nil,
				map[plural.Form]Evalable{
					plural.One:   Content("one"),
					plural.Other: Content("other"),
				},
			},
			Context{"count": 3.14159},
			"other",
			false,
		},
		{
			"other case if no other case",
			fields{
				"count",
				language.English,
				0,
				nil,
				map[plural.Form]Evalable{
					plural.Few:   Content("few"),
					plural.Other: Content("other"),
				},
			},
			Context{"count": 1},
			"other",
			false,
		},

		{
			"offset",
			fields{
				"count",
				language.English,
				2,
				nil,
				map[plural.Form]Evalable{
					plural.One:   Content("one"),
					plural.Other: PlainArg("#"),
				},
			},
			Context{"count": 4},
			"2",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plural := NewPlural(tt.fields.ArgName, tt.fields.Lang, tt.fields.Offset)
			plural.EqCases = tt.fields.EqCases
			plural.Cases = tt.fields.Cases

			got, err := plural.Eval(tt.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Plural.Eval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("Plural.Eval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toPluralForm(t *testing.T) {
	tests := []struct {
		name    string
		num     any
		want    pm
		wantErr bool
	}{
		{
			"int",
			int(42),
			pm{i: 42},
			false,
		},
		{
			"-int",
			int(-42),
			pm{i: 42},
			false,
		},
		{
			"int8",
			int8(42),
			pm{i: 42},
			false,
		},
		{
			"-int8",
			int8(-42),
			pm{i: 42},
			false,
		},
		{
			"int16",
			int16(42),
			pm{i: 42},
			false,
		},
		{
			"-int16",
			int16(-42),
			pm{i: 42},
			false,
		},
		{
			"int32",
			int32(42),
			pm{i: 42},
			false,
		},
		{
			"-int32",
			int32(-42),
			pm{i: 42},
			false,
		},
		{
			"int64",
			int64(42),
			pm{i: 42},
			false,
		},
		{
			"-int64",
			int64(-42),
			pm{i: 42},
			false,
		},
		{
			"uint",
			uint(42),
			pm{i: 42},
			false,
		},
		{
			"uint8",
			uint8(42),
			pm{i: 42},
			false,
		},
		{
			"uint16",
			uint16(42),
			pm{i: 42},
			false,
		},
		{
			"uint32",
			uint32(42),
			pm{i: 42},
			false,
		},
		{
			"uint64",
			uint64(42),
			pm{i: 42},
			false,
		},
		{
			"error on unknown type",
			[]byte("foo"),
			pm{},
			true,
		},
		{
			"float string",
			"42.4200",
			pm{i: 42, v: 4, w: 2, f: 4200, t: 42},
			false,
		},
		{
			"float string",
			"1200.50",
			pm{i: 1200, v: 2, w: 1, f: 50, t: 5},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toPluralForm(tt.num)
			if (err != nil) != tt.wantErr {
				t.Errorf("toPluralForm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toPluralForm() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseString(t *testing.T) {
	tests := []struct {
		input    string
		expected pm
		wantErr  bool
	}{
		{"1", pm{i: 1}, false},
		{"-1", pm{i: 1}, false}, // Test negative number handling
		{"1.0", pm{i: 1, v: 1, w: 0, f: 0, t: 0}, false},
		{"1.23", pm{i: 1, v: 2, w: 2, f: 23, t: 23}, false},
		{"1.230", pm{i: 1, v: 3, w: 2, f: 230, t: 23}, false},
		{"1.000", pm{i: 1, v: 3, w: 0, f: 0, t: 0}, false},
		{"-1.23", pm{i: 1, v: 2, w: 2, f: 23, t: 23}, false},                 // Test negative number with decimals
		{"-12345.67890", pm{i: 12345, v: 5, w: 4, f: 67890, t: 6789}, false}, // Test longer numbers
		{"12345", pm{i: 12345}, false},                                       // Test larger integer
		{"abc", pm{}, true},
		{"1.abc", pm{}, true},
		{"1.2abc", pm{}, true},
		{"1.230abc", pm{}, true},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			actual, err := parseString(test.input)
			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expected, actual)
			}
		})
	}
}
