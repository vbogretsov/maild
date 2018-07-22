package app_test

import (
	"fmt"
	"io"
	"strings"

	"github.com/vbogretsov/go-validation"

	"github.com/vbogretsov/maild/app"
	"github.com/vbogretsov/maild/model"
)

const (
	errorStrCannotBeBlank  = "cannot be blank"
	errorStrInvalidEmail   = "invalid email"
	errorMissingRecipients = "missing recipients"
	errorInvalidBodyType   = "invalid body type"
)

var bodyTypes = []interface{}{"text/plain", "text/html"}

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
					Field: "templateLang",
					Errors: []error{
						validation.Error{Message: errorStrCannotBeBlank},
					},
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
					Field: "templateName",
					Errors: []error{
						validation.Error{Message: errorStrCannotBeBlank},
					},
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
					Field: "",
					Errors: []error{
						validation.Error{Message: errorMissingRecipients},
					},
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
									Field: "Email",
									Errors: []error{
										validation.Error{
											Message: errorStrInvalidEmail,
										},
									},
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
									Field: "Email",
									Errors: []error{validation.Error{
										Message: errorStrInvalidEmail,
									}},
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
			Err: validation.Errors([]error{
				validation.Error{
					Message: "template not found",
					Params: validation.Params{
						"lang": "en",
						"name": "xxx",
					},
				},
			}),
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
					validation.Error{
						Message: errorStrCannotBeBlank,
					},
					validation.Error{
						Message: errorInvalidBodyType,
						Params: validation.Params{
							"unsupported": "",
							"supported":   bodyTypes,
						},
					},
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
				Field: "BodyType",
				Errors: []error{
					validation.Error{
						Message: errorInvalidBodyType,
						Params: validation.Params{
							"unsupported": "text/xxx",
							"supported":   bodyTypes,
						},
					},
				},
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
				Field: "Body",
				Errors: []error{validation.Error{
					Message: errorStrCannotBeBlank,
				}},
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
				Field: "Subject",
				Errors: []error{validation.Error{
					Message: errorStrCannotBeBlank,
				}},
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
						Field: "Email",
						Errors: []error{validation.Error{
							Message: errorStrInvalidEmail,
						}},
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
						Field: "Email",
						Errors: []error{validation.Error{
							Message: errorStrInvalidEmail,
						}},
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
						Field: "Email",
						Errors: []error{validation.Error{
							Message: errorStrInvalidEmail,
						}},
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
		e := validation.Error{
			Message: "template not found",
			Params: validation.Params{
				"lang": lang,
				"name": name,
			},
		}
		return nil, app.ArgumentError{
			Err: validation.Errors([]error{e}),
		}
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
