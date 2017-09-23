package provider

import (
	"fmt"

	"github.com/vbogretsov/maild"
)

type consoleProvider struct{}

func (cp *consoleProvider) Send(msg maild.Message) error {
	fmt.Printf("sent: %v\n", msg)
	return nil
}

// NewConsolePorovider creates new console provider.
func NewConsolePorovider() maild.Provider {
	return &consoleProvider{}
}
