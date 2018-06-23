# maild (work in progress)

#### 0.1.0

Simple notification service for micro service architecture

## Installation

### go tool

```{bash}
$ go install github.com/vbogretsov/maild/cmd/maild
```

Command line options:

 * --help            Print help information
 * --provider-name   SMTP provider name, allowed valus: [sendgrid, log]
 * --provider-url    URL of SMTP service provider
 * --provider-key    SMTP service provider security key
 * --templates-path  Email templates location
 * --amqp-url        AMQP broker url. Default: amqp://guest:guest@localhost
 * --amqp-qname      AMQP queue name. Default: maild
 * --log-level       Log level, allowed values: [panic fatal error warn info
                     debug]. Default: info
 * --log-format      Log output format, allowed values [kubernetes json].
                     Default: json

### Docker

```{bash}
$ docker pull vbogretsov/maild:1
```

To add custom templates either use volume mounted into `/var/lib/maild/templates` or create a new image based on `vbogretsov/maild`. Example Dockerfile:

```{Dockerfile}
FROM vbogretsov/maild:0.1.0

COPY ./templates/*.msg /var/lib/maild/templates

ENTRYPOINT ["docker-entrypoint.sh"]
```

Available environment variables:

* MAILD_PROVIDER_URL - URL of SMTP service provider
* MAILD_PROVIDER_KEY - SMTP service provider security key
* MAILD_PROVIDER_NAME - SMTP provider name, allowed valus: [sendgrid, log]
* MAILD_AMQP_URL - AMQP broker url. Default: amqp://guest:guest@localhost
* MAILD_AMQP_QNAME - AMQP queue name. Default: maild
* MAILD_LOG_LEVEL - Log level, allowed values: [panic fatal error warn info debug]. Default: info
* MAILD_LOG_FORMAT - Log output format, allowed values [kubernetes json]. Default: json

## Usage

Email tempaltes should be a golang templates named according to the pattern: lang-template_name.msg

Available clients:

* [Go](https://github.com/vbogretsov/go-mailcd)

## Licence

See the LICENCE file.
