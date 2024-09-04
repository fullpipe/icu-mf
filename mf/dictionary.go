package mf

import (
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
)

type Dictionary interface {
	Get(path string) (string, error)
}

type dummyDictionary struct{}

func (*dummyDictionary) Get(id string) (string, error) {
	return "", fmt.Errorf("no message with id %s", id)
}

func NewDictionary(yaml string) Dictionary {
	return &dictionary{
		yamlContent: yaml,
		cache: map[string]struct {
			data string
			err  error
		}{},
	}
}

type dictionary struct {
	yamlContent string
	cache       map[string]struct {
		data string
		err  error
	}
}

func (d *dictionary) Get(id string) (string, error) {
	data, isCached := d.cache[id]
	if isCached {
		return data.data, data.err
	}

	var message string

	path, err := yaml.PathString("$." + id)
	if err != nil {
		return d.cacheAndReturn(id, message, err)
	}

	err = path.Read(strings.NewReader(d.yamlContent), &message)
	if err != nil {
		return d.cacheAndReturn(id, message, err)
	}

	return d.cacheAndReturn(id, message, err)
}

func (d *dictionary) cacheAndReturn(key string, data string, err error) (string, error) {
	d.cache[key] = struct {
		data string
		err  error
	}{
		data: data,
		err:  err,
	}

	return data, err
}
