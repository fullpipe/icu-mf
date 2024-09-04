package message

import (
	"reflect"
	"testing"

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
