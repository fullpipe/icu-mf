package message

import (
	"reflect"
	"testing"
	"time"

	"github.com/fullpipe/icu-mf/parse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

func TestBuild(t *testing.T) {
	tests := []struct {
		name    string
		in      parse.Message
		ctx     Context
		want    string
		wantErr bool
	}{
		{
			"basic usage",
			parse.Message{Fragments: []*parse.Fragment{
				{Text: "foo "},
				{PlainArg: &parse.PlainArg{Name: "foo"}},
			}},
			Context{"foo": "bar"},
			"foo bar",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eval, err := Build(tt.in, language.English)
			if (err != nil) != tt.wantErr {
				t.Errorf("Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got, err := eval.Eval(tt.ctx)
			require.NoError(t, err)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Build() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildFragment(t *testing.T) {
	tests := []struct {
		name    string
		f       parse.Fragment
		ctx     Context
		want    string
		wantErr bool
	}{
		{
			"text solo",
			parse.Fragment{Text: "foo"},
			Context{},
			"foo",
			false,
		},
		{
			"octothorpe",
			parse.Fragment{Octothorpe: true},
			Context{"#": "octo"},
			"octo",
			false,
		},
		{
			"error if no text, octothorpe or expretion",
			parse.Fragment{},
			Context{"#": "octo"},
			"",
			true,
		},
		{
			"simple expretion with name",
			parse.Fragment{PlainArg: &parse.PlainArg{Name: "foo"}},
			Context{"foo": "bar"},
			"bar",
			false,
		},
		{
			"builds valid date function",
			parse.Fragment{Func: &parse.Func{ArgName: "foo", Func: "date", Param: "short"}},
			Context{"foo": time.Date(1961, 4, 12, 6, 7, 3, 0, time.UTC)},
			"4/12/61",
			false,
		},
		{
			"builds valid time function",
			parse.Fragment{Func: &parse.Func{ArgName: "foo", Func: "time", Param: "short"}},
			Context{"foo": time.Date(1961, 4, 12, 6, 7, 3, 0, time.UTC)},
			"6:07 AM",
			false,
		},
		{
			"builds valid datetime function",
			parse.Fragment{Func: &parse.Func{ArgName: "foo", Func: "datetime", Param: "short"}},
			Context{"foo": time.Date(1961, 4, 12, 6, 7, 3, 0, time.UTC)},
			"4/12/61, 6:07 AM",
			false,
		},
		{
			"error on invalid datetime format",
			parse.Fragment{Func: &parse.Func{ArgName: "foo", Func: "datetime", Param: "invalid_format"}},
			Context{"foo": time.Date(1961, 4, 12, 6, 7, 3, 0, time.UTC)},
			"",
			true,
		},
		{
			"error on missed arg",
			parse.Fragment{Func: &parse.Func{ArgName: "foo", Func: "datetime", Param: "invalid_format"}},
			Context{},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eval, err := buildFragment(tt.f, language.English)
			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, eval)
				return
			}

			got, err := eval.Eval(tt.ctx)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
