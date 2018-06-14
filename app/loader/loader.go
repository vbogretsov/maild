package loader

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/vbogretsov/go-validation"
)

const errTemplateNotFound = "template '%s/%s.msg' not found"

type Loader interface {
	Load(lang, name string) (io.Reader, error)
}

type fsloader struct {
	root string
}

func New(root string) (Loader, error) {
	// TODO: add loaders from network storages: S3, GlusterFs, etc.
	return fsloader{root: root}, nil
}

func (self fsloader) Load(lang, name string) (io.Reader, error) {
	fname := path.Join(self.root, lang, fmt.Sprintf("%s.msg", name))

	file, err := os.Open(fname)

	if err != nil {
		if os.IsNotExist(err) {
			return nil, validation.Error(
				fmt.Errorf(errTemplateNotFound, lang, name))
		}

		return nil, err
	}

	return bufio.NewReader(file), nil
}
