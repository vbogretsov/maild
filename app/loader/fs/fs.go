package fs

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/vbogretsov/maild/app"
)

const errorTemplateNotFound = "template '%s/%s.msg' not found"

type fsloader struct {
	root string
}

// New creates a new loader. The exact loader type is determined by root prefix.
func New(root string) (app.Loader, error) {
	return fsloader{root: root}, nil
}

// Load loads a template with the language and name provided from the local
// file system.
func (ld fsloader) Load(lang, name string) (io.Reader, error) {
	fname := path.Join(ld.root, lang, fmt.Sprintf("%s.msg", name))

	file, err := os.Open(fname)

	if err != nil {
		if os.IsNotExist(err) {
			e := fmt.Errorf(errorTemplateNotFound, lang, name)
			return nil, app.ArgumentError{Err: e}
		}

		return nil, err
	}

	return bufio.NewReader(file), nil
}
