package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
	log "github.com/sirupsen/logrus"

	amqpapi "github.com/vbogretsov/maild/api/amqp"
	"github.com/vbogretsov/maild/app"
	"github.com/vbogretsov/maild/app/loader"
	"github.com/vbogretsov/maild/app/sender"
)

const (
	logFormatKubernetes = "kubernetes"
	logFormatJSON       = "json"
)

type argT struct {
	provider struct {
		Name *string
		URL  *string
		Key  *string
	}
	amqp struct {
		URL   *string
		QName *string
	}
	template struct {
		Path *string
	}
	log struct {
		Level  *string
		Format *string
	}
}

var (
	args       = argT{}
	parser     = argparse.NewParser(fmt.Sprintf("%s %s", name, version), usage)
	logLevels  = []string{"panic", "fatal", "error", "warn", "info", "debug"}
	logFormats = []string{logFormatKubernetes, logFormatJSON}
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
	args.amqp.URL = parser.String(
		"",
		"amqp-url",
		&argparse.Options{
			Required: false,
			Default:  "amqp://guest:guest@localhost",
			Help:     amqpURLHelp,
		})
	args.amqp.QName = parser.String("", "amqp-qname", &argparse.Options{
		Required: false,
		Default:  name,
		Help:     amqpQNameHelp,
	})
	args.log.Level = parser.Selector(
		"",
		"log-level",
		logLevels,
		&argparse.Options{
			Required: false,
			Default:  "info",
			Help:     fmt.Sprintf(logLevelHelp, logLevels),
		})
	args.log.Format = parser.Selector(
		"",
		"log-format",
		logFormats,
		&argparse.Options{
			Required: false,
			Default:  logFormatJSON,
			Help:     fmt.Sprintf(logFormatHelp, logFormats),
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
	if err != nil {
		return err
	}

	ap := app.New(lr, sr)

	lv, err := log.ParseLevel(*args.log.Level)
	if err != nil {
		return err
	}
	log.SetLevel(lv)

	return amqpapi.Run(ap, *args.amqp.URL, *args.amqp.QName)
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
