# maild

#### 1.0.0

Simple notification service for micro service architecture

## Installation

### go tool

```{bash}
$ go install github.com/vbogretsov/maild/cmd/maild
```

Command line options:

* --provider-url value   URL of SMTP service provider
* --provider-key value   SMTP service provider security key
* --provider-name value  SMTP provider name, allowed valus: [sendgrid, log] (default: log)
* --templates-path value email templates location
* --log-level value      log level, allowed values: [DEBUG, INFO, NOTICE, WARNING, ERROR, CRITICAL] (default: INFO)
* --help, -h             show help
* --version, -v          print the version

### Docker

```{bash}
$ docker pull vbogretsov/maild:1
```

To add custom templates either use volume mounted into `/var/lib/maild/templates` or create a new image based on `vbogretsov/maild`. Example Dockerfile:

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

* [Go](https://github.com/vbogretsov/go-mailcd)

## Licence

See the LICENCE file.
