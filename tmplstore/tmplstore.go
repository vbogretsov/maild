package tmplstore

import (
	"github.com/vbogretsov/maild"
)

// A TemplateStore implements loading/storing tempaltes.
type TemplateStore interface {
	maild.TemplateLoader
	Store(lang string, id string, content []byte) error
}
