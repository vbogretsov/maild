package model

// Address represents an email address.
type Address struct {
	Email string `yaml:"Email" json:"email,omitempty" validate:"required"`
	Name  string `yaml:"Name" json:"name,omitempty"`
}

// Request represents request parameters required to build and send an email.
type Request struct {
	TemplateLang string      `validate:"required"`
	TemplateName string      `validate:"required"`
	TemplateArgs interface{} `validate:"required"`
	To           []Address   `validate:"required"`
	Cc           []Address
}

// Message represent an email message.
type Message struct {
	From     Address   `yaml:"From"`
	To       []Address `yaml:"To"`
	Cc       []Address `yaml:"Cc"`
	Subject  string    `yaml:"Subject"`
	BodyType string    `yaml:"BodyType"`
	Body     string    `yaml:"Body"`
}
