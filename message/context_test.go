package message

import (
	"math"
	"testing"
)

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

func TestContext_Int64(t *testing.T) {
	tests := []struct {
		name    string
		c       Context
		key     string
		want    int64
		wantErr bool
	}{
		{
			"int",
			Context{"foo": int(42)},
			"foo",
			42,
			false,
		},
		{
			"int8",
			Context{"foo": int8(42)},
			"foo",
			42,
			false,
		},
		{
			"int16",
			Context{"foo": int16(42)},
			"foo",
			42,
			false,
		},
		{
			"int32",
			Context{"foo": int32(42)},
			"foo",
			42,
			false,
		},
		{
			"int64",
			Context{"foo": int64(42)},
			"foo",
			42,
			false,
		},
		{
			"uint",
			Context{"foo": uint(42)},
			"foo",
			42,
			false,
		},
		{
			"uint8",
			Context{"foo": uint8(42)},
			"foo",
			42,
			false,
		},
		{
			"uint16",
			Context{"foo": uint16(42)},
			"foo",
			42,
			false,
		},
		{
			"uint32",
			Context{"foo": uint32(42)},
			"foo",
			42,
			false,
		},
		{
			"max uint32",
			Context{"foo": uint32(math.MaxUint32)},
			"foo",
			math.MaxUint32,
			false,
		},
		{
			"uint64",
			Context{"foo": uint64(42)},
			"foo",
			42,
			false,
		},
		{
			"max uint64",
			Context{"foo": uint64(math.MaxUint64)},
			"foo",
			0,
			true,
		},
		{
			"float32",
			Context{"foo": float32(42.42)},
			"foo",
			42,
			false,
		},
		{
			"float64",
			Context{"foo": float64(42.42)},
			"foo",
			42,
			false,
		},
		{
			"string float",
			Context{"foo": "42.42"},
			"foo",
			42,
			false,
		},
		{
			"string int",
			Context{"foo": "42"},
			"foo",
			42,
			false,
		},
		{
			"invalid string",
			Context{"foo": "bar"},
			"foo",
			0,
			true,
		},
		{
			"unknown type",
			Context{"foo": []byte("bar")},
			"foo",
			0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.Int64(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Context.Int64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Context.Int64() = %v, want %v", got, tt.want)
			}
		})
	}
}
