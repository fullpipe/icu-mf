package mf

import (
	"fmt"

	y3 "gopkg.in/yaml.v3"
)

type Dictionary interface {
	Get(path string) (string, error)
}

type DummyDictionary struct{}

func (*DummyDictionary) Get(id string) (string, error) {
	return "", fmt.Errorf("no message with id %s", id)
}

func NewYamlDictionary(yaml []byte) (*YamlDictionary, error) {
	d := &YamlDictionary{
		flatMap: make(map[string]string),
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

type YamlDictionary struct {
	flatMap map[string]string
}

func (d *YamlDictionary) Get(id string) (string, error) {
	if msg, ok := d.flatMap[id]; ok {
		return msg, nil
	}

	return "", fmt.Errorf("no message with id %s", id)
}

func (d *YamlDictionary) buildFlatMap(prefix string, yn *y3.Node) {
	for i := 0; i < len(yn.Content); i += 2 {
		keyNode := yn.Content[i]
		valueNode := yn.Content[i+1]

		key := prefix + keyNode.Value

		switch valueNode.Kind {
		case y3.ScalarNode:
			d.flatMap[key] = valueNode.Value
		case y3.MappingNode:
			d.buildFlatMap(key+".", valueNode)
		}
	}
}
