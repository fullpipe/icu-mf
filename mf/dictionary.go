package mf

import (
	"fmt"

	y3 "gopkg.in/yaml.v3"
)

type Dictionary interface {
	Get(path string) (string, error)
}

type dummyDictionary struct{}

func (*dummyDictionary) Get(id string) (string, error) {
	return "", fmt.Errorf("no message with id %s", id)
}

func NewDictionary(yaml []byte) (Dictionary, error) {
	d := &dictionary{
		flatMap: map[string]string{},
	}

	var document y3.Node
	if err := y3.Unmarshal(yaml, &document); err != nil {
		return nil, err
	}

	if len(document.Content) > 0 {
		d.buildFlatMap("", document.Content[0])
	}

	return d, nil
}

type dictionary struct {
	flatMap map[string]string
}

func (d *dictionary) Get(id string) (string, error) {
	msg, ok := d.flatMap[id]
	if !ok {
		return "", fmt.Errorf("no message with id %s", id)
	}

	return msg, nil
}

func (d *dictionary) buildFlatMap(prefix string, yn *y3.Node) {
	for i := 0; i < len(yn.Content); i++ {
		n := yn.Content[i]

		if n.Kind == y3.MappingNode {
			d.buildFlatMap(prefix+n.Value+".", n)

			continue
		}

		if n.Kind == y3.ScalarNode {
			if yn.Content[i+1].Kind == y3.ScalarNode {
				d.flatMap[prefix+n.Value] = yn.Content[i+1].Value
			} else if yn.Content[i+1].Kind == y3.MappingNode {
				d.buildFlatMap(prefix+n.Value+".", yn.Content[i+1])
			}

			i++
			continue
		}
	}
}
