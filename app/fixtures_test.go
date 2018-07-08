package app_test

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/vbogretsov/go-validation"

	"github.com/vbogretsov/maild/app"
	"github.com/vbogretsov/maild/model"
)

const (
	errorStrCannotBeBlank     = "cannot be blank"
	errorStrInvalidEmail      = "invalid email"
	errorMissingRecipients    = "missing recipients"
	errorInvalidBodyTypeEmpty = "invalid body type '', allowed values [text/plain text/html]"
	errorInvalidBodyTypeXxx   = "invalid body type 'text/xxx', allowed values [text/plain text/html]"
)

const expectedBody = `Hello SuperUser!
This is test body!
`

const validMsg = `
From:
  Email: user@mail.com
  Name: Sender
Subject: Subject
BodyType: "text/plain"
Body: |
  Hello {{.Username}}!
  This is test body!
`

const invalidMsgInvalidBodyType = `
From:
  Email: user@mail.com
  Name: LevelUp
Subject: Subject
BodyType: "text/xxx"
Body: |
  Hello {{.Username}}!
  This is test body!
`

const invalidMsgMissingBodyType = `
From:
  Email: user@mail.com
  Name: LevelUp
Subject: Subject
Body: |
  Hello {{.Username}}!
  This is test body!
`

const invalidMsgMissingBody = `
From:
  Email: user@mail.com
  Name: LevelUp
Subject: Subject
BodyType: "text/plain"
`

const invalidMsgMissingSubject = `
From:
  Email: user@mail.com
  Name: LevelUp
BodyType: "text/plain"
Body: |
  Hello {{.Username}}!
  This is test body!
`

const invalidMsgMissingFrom = `
Subject: Subject
BodyType: "text/plain"
Body: |
  Hello {{.Username}}!
  This is test body!
`

const invalidMsgMissingFromEmail = `
From:
  Name: LevelUp
Subject: Subject
BodyType: "text/plain"
Body: |
  Hello {{.Username}}!
  This is test body!
`

const invalidMsgInvalidFromEmail = `
From:
  Email: user.mail.com
  Name: LevelUp
Subject: Subject
BodyType: "text/plain"
Body: |
  Hello {{.Username}}!
  This is test body!
`

type fixture struct {
	name    string
	request model.Request
	result  error
}

var fixtures = []fixture{
	{
		name: "ErrorIfMissingTemplateLang",
		request: model.Request{
			TemplateName: "valid",
			TemplateArgs: map[string]interface{}{
				"Username": "user@mail.com",
			},
			To: []model.Address{
				{
					Email: "to1@mail.com", Name: "",
				},
			},
		},
		result: app.ArgumentError{
			Err: validation.Errors([]error{
				validation.StructError{
					Field:  "templateLang",
					Errors: validation.Error(errorStrCannotBeBlank),
				},
			}),
		},
	},
	{
		name: "ErrorIfMissingTemplateName",
		request: model.Request{
			TemplateLang: "en",
			TemplateArgs: map[string]interface{}{
				"Username": "user@mail.com",
			},
			To: []model.Address{
				{
					Email: "to1@mail.com",
					Name:  "",
				},
			},
		},
		result: app.ArgumentError{
			Err: validation.Errors([]error{
				validation.StructError{
					Field:  "templateName",
					Errors: validation.Error(errorStrCannotBeBlank),
				},
			}),
		},
	},
	{
		name: "ErrorIfMissingRecipients",
		request: model.Request{
			TemplateLang: "en",
			TemplateName: "valid",
			TemplateArgs: map[string]interface{}{
				"Username": "user@mail.com",
			},
		},
		result: app.ArgumentError{
			Err: validation.Errors([]error{
				validation.StructError{
					Field:  "",
					Errors: validation.Error(errorMissingRecipients),
				},
			}),
		},
	},
	{
		name: "ErrorIfRecipientMissingEmail",
		request: model.Request{
			TemplateLang: "en",
			TemplateName: "valid",
			TemplateArgs: map[string]interface{}{
				"Username": "user@mail.com",
			},
			To: []model.Address{
				{
					Email: "",
					Name:  "",
				},
			},
		},
		result: app.ArgumentError{
			Err: validation.Errors([]error{
				validation.StructError{
					Field: "to",
					Errors: []error{
						validation.SliceError{
							Index: 0,
							Errors: []error{
								validation.StructError{
									Field:  "Email",
									Errors: validation.Error(errorStrInvalidEmail),
								},
							},
						},
					},
				},
			}),
		},
	},
	{
		name: "ErrorIfRecipientEmailInvalid",
		request: model.Request{
			TemplateLang: "en",
			TemplateName: "valid",
			TemplateArgs: map[string]interface{}{
				"Username": "user@mail.com",
			},
			To: []model.Address{
				{
					Email: "to1.mail.com",
					Name:  "",
				},
			},
		},
		result: app.ArgumentError{
			Err: validation.Errors([]error{
				validation.StructError{
					Field: "to",
					Errors: []error{
						validation.SliceError{
							Index: 0,
							Errors: []error{
								validation.StructError{
									Field:  "Email",
									Errors: validation.Error(errorStrInvalidEmail),
								},
							},
						},
					},
				},
			}),
		},
	},
	{
		name: "ErrorIfTemplateNotFound",
		request: model.Request{
			TemplateLang: "en",
			TemplateName: "xxx",
			TemplateArgs: map[string]interface{}{
				"Username": "user@mail.com",
			},
			To: []model.Address{
				{
					Email: "to1@mail.com",
					Name:  "",
				},
			},
		},
		result: app.ArgumentError{
			Err: validation.Error("template en-xxx not found"),
		},
	},
	{
		name: "ErrorIfMissingBodyType",
		request: model.Request{
			TemplateLang: "en",
			TemplateName: "invalid-missing-body-type",
			TemplateArgs: map[string]interface{}{
				"Username": "user@mail.com",
			},
			To: []model.Address{
				{
					Email: "to1@mail.com", Name: "",
				},
			},
		},
		result: validation.Errors([]error{
			validation.StructError{
				Field: "BodyType",
				Errors: []error{
					errors.New(errorStrCannotBeBlank),
					errors.New(errorInvalidBodyTypeEmpty),
				},
			},
		}),
	},
	{
		name: "ErrorIfInvalidBodyType",
		request: model.Request{
			TemplateLang: "en",
			TemplateName: "invalid-body-type",
			TemplateArgs: map[string]interface{}{
				"Username": "user@mail.com",
			},
			To: []model.Address{
				{
					Email: "to1@mail.com", Name: "",
				},
			},
		},
		result: validation.Errors([]error{
			validation.StructError{
				Field:  "BodyType",
				Errors: validation.Error(errorInvalidBodyTypeXxx),
			},
		}),
	},
	{
		name: "ErrorIfMissingBody",
		request: model.Request{
			TemplateLang: "en",
			TemplateName: "invalid-missing-body",
			TemplateArgs: map[string]interface{}{
				"Username": "user@mail.com",
			},
			To: []model.Address{
				{
					Email: "to1@mail.com", Name: "",
				},
			},
		},
		result: validation.Errors([]error{
			validation.StructError{
				Field:  "Body",
				Errors: validation.Error(errorStrCannotBeBlank),
			},
		}),
	},
	{
		name: "ErrorIfMissingSubject",
		request: model.Request{
			TemplateLang: "en",
			TemplateName: "invalid-missing-subject",
			TemplateArgs: map[string]interface{}{
				"Username": "user@mail.com",
			},
			To: []model.Address{
				{
					Email: "to1@mail.com", Name: "",
				},
			},
		},
		result: validation.Errors([]error{
			validation.StructError{
				Field:  "Subject",
				Errors: validation.Error(errorStrCannotBeBlank),
			},
		}),
	},
	{
		name: "ErrorIfMissingFrom",
		request: model.Request{
			TemplateLang: "en",
			TemplateName: "invalid-missing-from",
			TemplateArgs: map[string]interface{}{
				"Username": "user@mail.com",
			},
			To: []model.Address{
				{
					Email: "to1@mail.com", Name: "",
				},
			},
		},
		result: validation.Errors([]error{
			validation.StructError{
				Field: "From",
				Errors: []error{
					validation.StructError{
						Field:  "Email",
						Errors: validation.Error(errorStrInvalidEmail),
					},
				},
			},
		}),
	},
	{
		name: "ErrorIfMissingFromEmail",
		request: model.Request{
			TemplateLang: "en",
			TemplateName: "invalid-missing-from-email",
			TemplateArgs: map[string]interface{}{
				"Username": "user@mail.com",
			},
			To: []model.Address{
				{
					Email: "to1@mail.com", Name: "",
				},
			},
		},
		result: validation.Errors([]error{
			validation.StructError{
				Field: "From",
				Errors: []error{
					validation.StructError{
						Field:  "Email",
						Errors: validation.Error(errorStrInvalidEmail),
					},
				},
			},
		}),
	},
	{
		name: "ErrorIfInvalidFromEmail",
		request: model.Request{
			TemplateLang: "en",
			TemplateName: "invalid-from-email",
			TemplateArgs: map[string]interface{}{
				"Username": "user@mail.com",
			},
			To: []model.Address{
				{
					Email: "to1@mail.com", Name: "",
				},
			},
		},
		result: validation.Errors([]error{
			validation.StructError{
				Field: "From",
				Errors: []error{
					validation.StructError{
						Field:  "Email",
						Errors: validation.Error(errorStrInvalidEmail),
					},
				},
			},
		}),
	},
}

type Loader struct {
	templates map[string]string
}

func NewLoader() *Loader {
	return &Loader{templates: map[string]string{
		"en-valid":                      validMsg,
		"en-invalid-missing-body-type":  invalidMsgMissingBodyType,
		"en-invalid-body-type":          invalidMsgInvalidBodyType,
		"en-invalid-missing-body":       invalidMsgMissingBody,
		"en-invalid-missing-subject":    invalidMsgMissingSubject,
		"en-invalid-missing-from":       invalidMsgMissingFrom,
		"en-invalid-missing-from-email": invalidMsgMissingFromEmail,
		"en-invalid-from-email":         invalidMsgInvalidFromEmail,
	}}
}

func (self *Loader) Load(lang, name string) (io.Reader, error) {
	id := fmt.Sprintf("%s-%s", lang, name)
	text, ok := self.templates[id]
	if !ok {
		e := validation.Errorf("template %s not found", id)
		return nil, app.ArgumentError{Err: e}
	}
	return strings.NewReader(text), nil
}

type Sender struct {
	inbox []model.Message
}

func NewSender() *Sender {
	return &Sender{inbox: []model.Message{}}
}

func (self *Sender) Send(msg model.Message) error {
	self.inbox = append(self.inbox, msg)
	return nil
}
