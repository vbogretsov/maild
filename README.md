# maild
Simple notification service for micro service architecture

## Installation

### go tool

```{bash}
$ go install github.com/vbogretsov/maild/cmd/maild
```

Command line options:

* --broker-url value     URL of the broker which holds the queue of requests (default: "amqp://localhost:5672")
* --service-name value   name of the service is used for routing requests (default: "maild")
* --provider-url value   URL of SMTP service provider
* --provider-key value   SMTP service provider security key
* --provider-name value  SMTP provider name, allowed valus: [sendgrid, log] (default: "log")
* --template-dir value   email templates location
* --log-level value      log level, allowed values: [DEBUG, INFO, NOTICE, WARNING, ERROR, CRITICAL] (default: "INFO")
* --help, -h             show help
* --version, -v          print the version

### Docker

```{bash}
$ docker pull vbogretsov/maild:1
```

Create a Dockerfile

```{Dockerfile}
FROM vbogretsov/maild:1

COPY ./templates/*.msg /var/lib/maild/templates

ENTRYPOINT ["docker-entrypoint.sh"]
```

Available environment variables:

* MAILD_BROKER_URL - RabbitMQ broker URL (default: amqp://guest:guest@localhost:5672)
* MAILD_PROVIDER_URL - SMTP provider URL
* MAILD_PROVIDER_KEY - SMTP provider security key
* MAILD_PROVIDER_NAME - SMTP provider name
* MAILD_SERVICE_NAME - service name (default: maild)
* MAILD_LOG_LEVEL - log level (default: INFO)

## Usage

Email tempaltes should be a golang templates named according to the pattern: lang-template_name.msg

Available clients:

* [GOLANG](https://github.com/vbogretsov/go-mailcd)

## Licence

See the LICENCE file.
