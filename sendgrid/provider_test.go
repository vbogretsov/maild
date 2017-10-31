package sendgrid

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/gernest/alien"
	"github.com/kr/pretty"

	"github.com/vbogretsov/maild/model"
	"github.com/vbogretsov/maild/sendgrid"
)

func startSendGridMock(addr string, response *[]byte) *http.Server {
	m := alien.New()
	m.Post("/v3/mail/send", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			*response = body
		}
	})

	server := http.Server{Addr: addr, Handler: m}
	go server.ListenAndServe()

	return &server
}

func toSendGridMessage(m model.Message) sgMessage {
	return sgMessage{
		From: m.From,
		Personalization: sgPersonalization{
			Subject: m.Subject,
			To:      m.To,
			Cc:      m.Cc,
		},
		Content: sgContent{
			Type:  m.BodyType,
			Value: m.Body,
		},
	}
}

func testSendSucceed(p *sendgrid.Provider, m model.Message, r *[]byte) func(*testing.T) {
	return func(t *testing.T) {
		err := p.SendMail(&m)
		if err != nil {
			t.Error(err)
		}

		exp := toSendGridMessage(m)

		act := sgMessage{}
		if err := json.Unmarshal(*r, &act); err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(exp, act) {
			t.Errorf(
				"error: actual response do not match expected reponse: %v",
				pretty.Diff(exp, act),
			)
		}
	}
}

func TestSendGridProvider(t *testing.T) {
	var response []byte
	mockaddr := ":3000"

	server := startSendGridMock(mockaddr, &response)
	defer server.Close()

	sgurl := fmt.Sprintf("http://localhost%s", mockaddr)
	provider := sendgrid.NewProvider(sgurl, "none")

	msgs := []model.Message{
		{
			From: model.Address{
				Name:  "sender",
				Email: "sender@mail.com",
			},
			To: []model.Address{
				{Email: "to1@mail.com"},
				{Email: "to2@mail.com"},
			},
			Cc: []model.Address{
				{Email: "cc1@mail.com"},
				{Email: "cc2@mail.com"},
			},
			Subject:  "Test Subj",
			BodyType: "text/plain",
			Body:     "Hello!",
		},
	}

	for _, msg := range msgs {
		t.Run("TestSendSucceed", testSendSucceed(provider, msg, &response))
	}
}
