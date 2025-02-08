package message

import "testing"

func TestMessage_Eval(t *testing.T) {
	tests := []struct {
		name      string
		fragments []Evalable
		ctx       Context
		want      string
		wantErr   bool
	}{
		{
			"basic usage",
			[]Evalable{Content("foo "), PlainArg("bar")},
			Context{"bar": 42},
			"foo 42",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Message{
				fragments: tt.fragments,
			}

			got, err := m.Eval(tt.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Message.Eval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("Message.Eval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContent_Eval(t *testing.T) {
	tests := []struct {
		name    string
		c       Content
		ctx     Context
		want    string
		wantErr bool
	}{
		{
			"basic usage",
			Content("foo"),
			Context{},
			"foo",
			false,
		},
		{
			"fine with empty string",
			Content(""),
			Context{},
			"",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.Eval(tt.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Content.Eval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("Content.Eval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlainArg_Eval(t *testing.T) {
	tests := []struct {
		name    string
		pa      PlainArg
		ctx     Context
		want    string
		wantErr bool
	}{
		{
			"string as string",
			PlainArg("foo"),
			Context{"foo": "bar"},
			"bar",
			false,
		},
		{
			"int as string",
			PlainArg("foo"),
			Context{"foo": 42},
			"42",
			false,
		},
		{
			"error if no key in context",
			PlainArg("foo"),
			Context{"bar": 42},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.pa.Eval(tt.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("PlainArg.Eval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("PlainArg.Eval() = %v, want %v", got, tt.want)
			}
		})
	}
}
