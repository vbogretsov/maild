package model

type Template struct {
	Lang string
	Name string
	Text string
}

type Request struct {
	TemplateLang string `validate:"required"`
	TemplateName string `validate:"required"`
	TemplateArgs string `validate:"required"`
}

type Address struct {
	Email string `yaml:"Email" json:"email,omitempty"`
	Name  string `yaml:"Name" json:"name,omitempty"`
}

type Message struct {
	From     Address   `yaml:"From"`
	To       []Address `yaml:"To"`
	Cc       []Address `yaml:"Cc"`
	Subject  string    `yaml:"Subject"`
	BodyType string    `yaml:"BodyType"`
	Body     string    `yaml:"Body"`
}
