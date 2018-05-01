package client

type Client interface {
	Send(templateID string, args map[string]interface{}) error
}

// TODO(vbogretsov): add implementation
