package server

import (
	"bytes"
	"text/template"

	"github.com/bluele/gcache"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v1"

	"github.com/vbogretsov/maild/model"
)

var (
	valid = validator.New()
)

// TODO(vbogretsov): rename sender to Maild
type Maild struct {
	cache    gcache.Cache
	provider model.Provider
}

func (m *Maild) Send(r *model.Request, out *struct{}) error {
	if err := valid.Struct(r); err != nil {
		return errors.Wrap(err, "invalid request")
	}

	tid := model.TemplateID{
		Lang: r.TemplateLang,
		Name: r.TemplateName,
	}

	val, err := m.cache.Get(&tid)
	if err != nil {
		return errors.Wrap(err, "template {%s, %s} not found")
	}

	tml, _ := val.(*template.Template)

	buf := bytes.Buffer{}
	if err := tml.Execute(&buf, r.TemplateArgs); err != nil {
		return err
	}

	msg := model.Message{}

	if err := yaml.Unmarshal(buf.Bytes(), &msg); err != nil {
		return err
	}

	msg.To = append(msg.To, r.To...)
	msg.Cc = append(msg.Cc, r.Cc...)

	return m.provider.SendMail(&msg)
}

func NewMaild(provider model.Provider, loader model.TemplateLoader, cacheSize int) *Maild {
	cache := gcache.New(cacheSize).
		LRU().
		LoaderFunc(func(key interface{}) (interface{}, error) {
			return loader(key.(*model.TemplateID))
		}).
		Build()
	return &Maild{cache: cache, provider: provider}
}
