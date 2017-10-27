package main

const (
	name    = "maild"
	usage   = "notification service for micro service architecture"
	version = "0.0.0"
	logfmt  = `%{color}#%{id:03x} [%{pid}] %{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{message}%{color:reset}`

	brokerURLUsage   = `URL of the broker which holds the queue of requests`
	serviceNameUsage = `name of the service is used for routing requests`
	providerURLUsage = `URL of SMTP service provider`
	providerKeyUsage = `SMTP service provider security key`
	dbTypeUsage      = `type of the database where templates are stored, allowed values: [postgresql, mysql]`
	dbDSNUsage       = `DSN of the database where templates are stored`
	logLevelUsage    = `log level, allowed values: [DEBUG, INFO, NOTICE, WARNING, ERROR, CRITICAL]`
)
