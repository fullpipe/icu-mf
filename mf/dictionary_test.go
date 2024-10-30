package mf

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDictionaryGet(t *testing.T) {
	d, err := NewYamlDictionary([]byte(`
one:
    two: msg1-2
foo: bar
multyline: |-
    a
    b
oneline: >-
    c
    d
deep:
    multyline: |-
        e
        f
`))
	require.NoError(t, err)

	res, err := d.Get("one.two")
	require.NoError(t, err)
	assert.Equal(t, "msg1-2", res)

	res, err = d.Get("foo")
	require.NoError(t, err)
	assert.Equal(t, "bar", res)

	res, err = d.Get("multyline")
	require.NoError(t, err)
	assert.Equal(t, "a\nb", res)

	res, err = d.Get("oneline")
	require.NoError(t, err)
	assert.Equal(t, "c d", res)

	res, err = d.Get("deep.multyline")
	require.NoError(t, err)
	assert.Equal(t, "e\nf", res)

	_, err = d.Get("...invalid path$.$")
	require.Error(t, err)

	_, err = d.Get("...invalid path$.$")
	require.Error(t, err)

	_, err = d.Get("no.path")
	require.Error(t, err)
}

func TestNewDictionary(t *testing.T) {
	tests := []struct {
		name    string
		yaml    []byte
		want    map[string]string
		wantErr bool
	}{
		{
			"valid yaml",
			[]byte("foo: bar\nfizz: buzz"),
			map[string]string{"foo": "bar", "fizz": "buzz"},
			false,
		},
		{
			"invalid yaml",
			[]byte(":"),

			map[string]string{},
			true,
		},
		{
			"empty yaml",
			[]byte(""),
			map[string]string{},
			false,
		},
		{
			"valid yaml, but sequence",
			[]byte(`
foo:
  - one
  - two
`),
			map[string]string{"foo": ""},
			false,
		},
		{
			"sequence only",
			[]byte("- one\n- two"),
			map[string]string{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewYamlDictionary(tt.yaml)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDictionary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for k, v := range tt.want {
				msg, _ := d.Get(k)
				assert.Equal(t, v, msg)
			}
		})
	}
}
