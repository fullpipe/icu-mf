package message

import "testing"

func TestSelect_Eval(t *testing.T) {
	tests := []struct {
		name    string
		ArgName string
		Cases   map[string]Evalable
		ctx     Context
		want    string
		wantErr bool
	}{
		{
			"basic usage",
			"color",
			map[string]Evalable{
				"red":   Content("color is red"),
				"blue":  &Message{fragments: []Evalable{PlainArg("tone"), Content(" blue")}},
				"other": Content("color not exists"),
			},
			Context{"color": "red"},
			"color is red",
			false,
		},
		{
			"nested tree",
			"color",
			map[string]Evalable{
				"red":   Content("color is red"),
				"blue":  &Message{fragments: []Evalable{PlainArg("tone"), Content(" blue")}},
				"other": Content("color not exists"),
			},
			Context{"color": "blue", "tone": "deep"},
			"deep blue",
			false,
		},
		{
			"default case",
			"color",
			map[string]Evalable{
				"red":   Content("color is red"),
				"blue":  &Message{fragments: []Evalable{PlainArg("tone"), Content(" blue")}},
				"other": Content("color not exists"),
			},
			Context{"color": "nope", "tone": "deep"},
			"color not exists",
			false,
		},
		{
			"default case if no arg",
			"color",
			map[string]Evalable{
				"red":   Content("color is red"),
				"blue":  &Message{fragments: []Evalable{PlainArg("tone"), Content(" blue")}},
				"other": Content("color not exists"),
			},
			Context{},
			"color not exists",
			false,
		},
		{
			"error if default case not exists",
			"color",
			map[string]Evalable{
				"red":  Content("color is red"),
				"blue": &Message{fragments: []Evalable{PlainArg("tone"), Content(" blue")}},
			},
			Context{"color": "nope", "tone": "deep"},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Select{
				ArgName: tt.ArgName,
				Cases:   tt.Cases,
			}

			got, err := s.Eval(tt.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Select.Eval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("Select.Eval() = %v, want %v", got, tt.want)
			}
		})
	}
}
