package mf

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

func TestNewFSProvider(t *testing.T) {
	p, err := NewYamlMessageProvider(fstest.MapFS{"dir.en": {Mode: fs.ModeDir}})
	require.NoError(t, err, "no error on empty fs")
	assert.Empty(t, p.dictionaries, "does not loads dirs")

	p, err = NewYamlMessageProvider(fstest.MapFS{"messages.en.toml": {Data: []byte("foo=bar")}})
	require.NoError(t, err, "no error on non yaml files")
	assert.Empty(t, p.dictionaries, "but does not loads them")

	_, err = NewYamlMessageProvider(fstest.MapFS{".yaml": {Data: []byte("foo: bar")}})
	require.Error(t, err, "error on invalid filename")

	_, err = NewYamlMessageProvider(fstest.MapFS{"messages.FOO.yaml": {Data: []byte("foo: bar")}})
	require.Error(t, err, "error on invalid lang in filename")

	p, err = NewYamlMessageProvider(fstest.MapFS{
		"messages.en.yaml": {Data: []byte("foo: bar")},
		"messages.es.yaml": {Data: []byte("foo: bar")},
		"messages.ru.yml":  {Data: []byte("foo: bar")},
	})
	require.NoError(t, err, "no error on normal yaml files")
	assert.NotEmpty(t, p.dictionaries, "has dictionaries")
	assert.NotNil(t, p.dictionaries[language.English])
	assert.NotNil(t, p.dictionaries[language.Spanish])
	assert.NotNil(t, p.dictionaries[language.Russian])
}

func TestFSProvider_loadMessages(t *testing.T) {
	p := &YamlMessageProvider{
		dictionaries: map[language.Tag]*YamlDictionary{},
	}
	fs := fstest.MapFS{
		"non_readable": {
			Mode: fs.ModeDir,
		},
		"foo.en.yaml": {
			Data: []byte("foo: bar"),
		},
	}

	require.Error(t, p.loadMessages(fs, "file_not_exists.yaml", language.English))
	assert.Nil(t, p.dictionaries[language.English])

	require.Error(t, p.loadMessages(fs, "non_readable", language.English))
	assert.Nil(t, p.dictionaries[language.English])

	require.NoError(t, p.loadMessages(fs, "foo.en.yaml", language.English))
	assert.NotNil(t, p.dictionaries[language.English])

	require.Error(t, p.loadMessages(fs, "foo.en.yaml", language.English), "error when lang already loaded")
	assert.NotNil(t, p.dictionaries[language.English])
}
