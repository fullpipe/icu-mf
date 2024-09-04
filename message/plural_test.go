package message

import (
	"testing"

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
