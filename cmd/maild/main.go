package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"

	"github.com/vbogretsov/maild/api"
	"github.com/vbogretsov/maild/app"
	"github.com/vbogretsov/maild/app/loader"
	"github.com/vbogretsov/maild/app/sender"
)

type argT struct {
	port     *int
	provider struct {
		Name *string
		URL  *string
		Key  *string
	}
	template struct {
		Path *string
	}
	log struct {
		Level *string
	}
}

var (
	args   = argT{}
	parser = argparse.NewParser(name, usage)
)

func init() {
	args.provider.Name = parser.String(
		"",
		"provider-name",
		&argparse.Options{
			Required: true,
			Help:     providerNameHelp,
		})
	args.provider.URL = parser.String(
		"",
		"provider-url",
		&argparse.Options{
			Required: true,
			Help:     providerURLHelp,
		})
	args.provider.Key = parser.String(
		"",
		"provider-key",
		&argparse.Options{
			Required: true,
			Help:     providerKeyHelp,
		})
	args.template.Path = parser.String(
		"",
		"templates-path",
		&argparse.Options{
			Required: true,
			Help:     templatesPathHelp,
		})
	args.port = parser.Int(
		"",
		"port",
		&argparse.Options{
			Required: false,
			Help:     portHelp,
			Default:  8000,
		})
}

func run() error {
	if err := parser.Parse(os.Args); err != nil {
		return err
	}

	lr, err := loader.New(*args.template.Path)
	if err != nil {
		return err
	}

	sr, err := sender.New(
		*args.provider.Name,
		*args.provider.URL,
		*args.provider.Key)

	rt, err := api.New(app.New(lr, sr))
	if err != nil {
		return err
	}

	rt.Logger.Fatal(rt.Start(fmt.Sprintf(":%d", *args.port)))

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
