package model

import (
	"text/template"
)

type Address struct {
	Email string `yaml:"Email" json:"email,omitempty" validate:"required"`
	Name  string `yaml:"Name" json:"name,omitempty"`
}

type Request struct {
	TemplateLang string      `validate:"required"`
	TemplateName string      `validate:"required"`
	TemplateArgs interface{} `validate:"required"`
	To           []Address   `validate:"required"`
	Cc           []Address
}

type Message struct {
	From     Address   `yaml:"From"`
	To       []Address `yaml:"To"`
	Cc       []Address `yaml:"Cc"`
	Subject  string    `yaml:"Subject"`
	BodyType string    `yaml:"BodyType"`
	Body     string    `yaml:"Body"`
}

type TemplateLoader func(*TemplateID) (*template.Template, error)

type TemplateID struct {
	Lang string
	Name string
}

type Provider interface {
	SendMail(*Message) error
}
