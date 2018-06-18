package main

const (
	name              = `maild`
	usage             = `Notification service for micro service architecture`
	providerURLHelp   = `URL of SMTP service provider`
	providerKeyHelp   = `SMTP service provider security key`
	providerNameHelp  = `SMTP provider name, allowed valus: [sendgrid, log]`
	amqpURLHelp       = `AMQP broker url`
	amqpQNameHelp     = `AMQP queue name`
	templatesPathHelp = `Email templates location`
	logLevelHelp      = `Log level, allowed values: %v`
	logFormatHelp     = `Log output format, allowed values %v`
)
