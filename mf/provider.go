package mf

import (
	"fmt"
	"io"
	"io/fs"
	"path"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/text/language"
)

type MessageProvider interface {
	Get(lang language.Tag, id string) (string, error)
}

type LanguagesProvider interface {
	MessageProvider
	Languages() []language.Tag
}

type YamlMessageProvider struct {
	dictionaries map[language.Tag]*YamlDictionary
	languages    []language.Tag
}

func (p *YamlMessageProvider) Get(lang language.Tag, path string) (string, error) {
	d, hasDictionary := p.dictionaries[lang]
	if !hasDictionary {
		return "", errors.Errorf("no dictionary for lang %s", lang)
	}

	return d.Get(path)
}

func (p *YamlMessageProvider) Languages() []language.Tag {
	return p.languages
}

func NewYamlMessageProvider(dir fs.FS) (*YamlMessageProvider, error) {
	provider := YamlMessageProvider{
		dictionaries: map[language.Tag]*YamlDictionary{},
	}

	err := fs.WalkDir(dir, ".", func(p string, f fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if f.IsDir() {
			return nil
		}

		if path.Ext(f.Name()) != ".yaml" && path.Ext(f.Name()) != ".yml" {
			return nil
		}

		nameParts := strings.Split(f.Name(), ".")
		if len(nameParts) < 2 {
			return fmt.Errorf("no lang in file %s", f.Name())
		}

		tag, err := language.Parse(nameParts[len(nameParts)-2])
		if err != nil {
			return errors.Wrap(err, "unable to parse language from filename")
		}

		return provider.loadMessages(dir, p, tag)
	})

	return &provider, err
}

func (p *YamlMessageProvider) loadMessages(rd fs.FS, path string, lang language.Tag) error {
	yamlFile, err := rd.Open(path)
	if err != nil {
		return errors.Wrap(err, "unable to open file")
	}

	yamlData, err := io.ReadAll(yamlFile)
	if err != nil {
		return errors.Wrap(err, "unable to read file")
	}

	_, hasDictionary := p.dictionaries[lang]
	if hasDictionary {
		return fmt.Errorf("unable to load %s: language %s already has messages loaded", path, lang)
	}

	p.dictionaries[lang], err = NewYamlDictionary(yamlData)
	if err != nil {
		return errors.Wrap(err, "unable to create dictionary")
	}
	p.languages = append(p.languages, lang)
	return nil
}
