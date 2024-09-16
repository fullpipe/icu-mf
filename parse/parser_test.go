package parse

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser_Escaping(t *testing.T) {
	parser := NewParser()

	msg, err := parser.Parse("", strings.NewReader("foo '{ '' ' foo"))

	require.NoError(t, err)
	assert.Len(t, msg.Fragments, 7)

	assert.Equal(t, &Fragment{Text: "foo "}, msg.Fragments[0])
	assert.Equal(t, &Fragment{Escaped: "'{"}, msg.Fragments[1])
	assert.Equal(t, &Fragment{Text: " "}, msg.Fragments[2])
	assert.Equal(t, &Fragment{Escaped: "''"}, msg.Fragments[3])
	assert.Equal(t, &Fragment{Text: " "}, msg.Fragments[4])
	assert.Equal(t, &Fragment{Text: "'"}, msg.Fragments[5])
	assert.Equal(t, &Fragment{Text: " foo"}, msg.Fragments[6])
}
