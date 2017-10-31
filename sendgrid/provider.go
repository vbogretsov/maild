package sendgrid

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"

	api "github.com/sendgrid/sendgrid-go"
	"gopkg.in/go-playground/validator.v9"

	"github.com/vbogretsov/maild/model"
)

var (
	valid = validator.New()
)

type sgContent struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type sgPersonalization struct {
	Subject string          `json:"subject"`
	To      []model.Address `json:"to"`
	Cc      []model.Address `json:"cc,omitempty"`
}

type sgMessage struct {
	From            model.Address     `json:"from"`
	Personalization sgPersonalization `json:"personalization"`
	Content         sgContent         `json:"content"`
}

type Provider struct {
	url string
	key string
}

func (p *Provider) SendMail(message *model.Message) error {
	if err := valid.Struct(message); err != nil {
		return err
	}

	sgdata := sgMessage{
		From: message.From,
		Personalization: sgPersonalization{
			Subject: message.Subject,
			To:      message.To,
			Cc:      message.Cc,
		},
		Content: sgContent{
			Type:  message.BodyType,
			Value: message.Body,
		},
	}

	data, err := json.Marshal(&sgdata)
	if err != nil {
		return err
	}

	request := api.GetRequest(p.key, "/v3/mail/send", p.url)
	request.Method = "POST"
	request.Body = data

	resp, err := api.API(request)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("SendGrid API call failed: %d", resp.StatusCode)
	}

	return nil
}

func NewProvider(url, key string) *Provider {
	return &Provider{url: url, key: key}
}
