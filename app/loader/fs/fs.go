package fs

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/vbogretsov/go-validation"

	"github.com/vbogretsov/maild/app"
)

// ErrorTemplateNotFound defines error message if template not found.
var ErrorTemplateNotFound = "template not found"

var (
	// ParamTemplateLang defines key name for error parameter template lang.
	ParamTemplateLang = "lang"
	// ParamTemplateName defines key name for error parameter template name.
	ParamTemplateName = "name"
)

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
			return nil, app.ArgumentError{
				Err: errorTemplateNotFound(lang, name),
			}
		}
		return nil, err
	}

	return bufio.NewReader(file), nil
}

func errorTemplateNotFound(lang, name string) error {
	e := validation.Error{
		Message: ErrorTemplateNotFound,
		Params: validation.Params{
			ParamTemplateLang: lang,
			ParamTemplateName: name,
		}}
	return validation.Errors([]error{e})
}
