package message

import "testing"

func Test_context_AsString(t *testing.T) {
	tests := []struct {
		name    string
		c       Context
		argName string
		want    string
		wantErr bool
	}{
		{
			"string as string",
			Context{
				"foo":  "bar",
				"fizz": 42,
				"buzz": float64(3.14),
			},
			"foo",
			"bar",
			false,
		},
		{
			"int as string",
			Context{
				"foo":  "bar",
				"fizz": 42,
				"buzz": float64(3.14),
			},
			"fizz",
			"42",
			false,
		},
		{
			"float as string",
			Context{
				"foo":  "bar",
				"fizz": 42,
				"buzz": float64(3.14),
			},
			"buzz",
			"3.14",
			false,
		},
		{
			"error if key not exists",
			Context{
				"foo":  "bar",
				"fizz": 42,
				"buzz": float64(3.14),
			},
			"nope",
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.String(tt.argName)
			if (err != nil) != tt.wantErr {
				t.Errorf("context.AsString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("context.AsString() = %v, want %v", got, tt.want)
			}
		})
	}
}
