package model

// Address represents an email address.
type Address struct {
	Email string `yaml:"Email" json:"email,omitempty" validate:"required"`
	Name  string `yaml:"Name" json:"name,omitempty"`
}

// Message represent an email message.
type Message struct {
	From     Address   `yaml:"From"`
	To       []Address `yaml:"To"`
	Cc       []Address `yaml:"Cc"`
	Bcc      []Address `yaml:Bcc`
	Subject  string    `yaml:"Subject"`
	BodyType string    `yaml:"BodyType"`
	Body     string    `yaml:"Body"`
}

// Request represents request parameters required to build and send an email.
type Request struct {
	TemplateLang string                 `validate:"required" json:"templateLang"`
	TemplateName string                 `validate:"required" json:"templateName"`
	TemplateArgs map[string]interface{} `validate:"required" json:"templateArgs"`
	To           []Address              `json:"to"`
	Cc           []Address              `json:"cc"`
	Bcc          []Address              `json:"bcc"`
}
