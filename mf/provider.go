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

type Provider interface {
	Get(lang language.Tag, path string) (string, error)
}

type FSProvider struct {
	dictionaries map[language.Tag]Dictionary
}

func (p *FSProvider) Get(lang language.Tag, path string) (string, error) {
	d, hasDictionary := p.dictionaries[lang]
	if !hasDictionary {
		return "", errors.Errorf("no dictionary for lang %s", lang)
	}

	return d.Get(path)
}

func NewFSProvider(dir fs.FS) (*FSProvider, error) {
	provider := FSProvider{
		dictionaries: map[language.Tag]Dictionary{},
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

func (p *FSProvider) loadMessages(rd fs.FS, path string, lang language.Tag) error {
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

	p.dictionaries[lang], err = NewDictionary(yamlData)
	if err != nil {
		return errors.Wrap(err, "unable to create dictionary")
	}

	return nil
}
