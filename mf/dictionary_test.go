package mf

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDictionaryGet(t *testing.T) {
	d := NewDictionary(`
one:
    two: msg1-2
foo: bar
multyline: |-
    a
    b
`)

	res, err := d.Get("one.two")
	require.NoError(t, err)
	assert.Equal(t, "msg1-2", res)

	res, err = d.Get("foo")
	require.NoError(t, err)
	assert.Equal(t, "bar", res)

	res, err = d.Get("multyline")
	require.NoError(t, err)
	assert.Equal(t, "a\nb", res)

	_, err = d.Get("...invalid path$.$")
	require.Error(t, err)

	_, err = d.Get("...invalid path$.$")
	require.Error(t, err)

	_, err = d.Get("no.path")
	require.Error(t, err)
}
