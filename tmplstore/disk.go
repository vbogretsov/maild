package tmplstore

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"text/template"
)

type diskStore struct {
	root string
}

func (ds *diskStore) Load(lang string, id string) (*template.Template, error) {
	tempaltePath := path.Join(ds.root, lang, id+".tmpl")
	if _, err := os.Stat(tempaltePath); err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadFile(tempaltePath)
	if err != nil {
		return nil, err
	}

	return template.New(id).Parse(string(bytes))
}

func (ds *diskStore) Store(lang string, id string, content []byte) error {
	langDir := path.Join(ds.root, lang)
	if _, err := os.Stat(langDir); os.IsNotExist(err) {
		err = os.Mkdir(langDir, os.ModeDir)
		if err != nil {
			return err
		}
	}

	tempaltePath := path.Join(langDir, id+".tmpl")
	return ioutil.WriteFile(tempaltePath, content, 0644)
}

// NewDiskStore creates a new on disk templates store.
func NewDiskStore(root string) (TemplateStore, error) {
	fi, err := os.Stat(root)
	if err != nil {
		return nil, err
	}

	if !fi.IsDir() {
		return nil, errors.New("path is not a directory")
	}

	return &diskStore{root}, nil
}
