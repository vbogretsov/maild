FROM alpine:3.6

ENV	MAILD_BROKER_URL=amqp://guest:guest@rabbitmq:5672/ \
	MAILD_PROVIDER_URL=none \
	MAILD_PROVIDER_KEY=none \
	MAILD_PROVIDER_NAME=log \
	MAILD_SERVICE_NAME=maild \
	MAILD_LOG_LEVEL=INFO

ADD ./bin/maild /bin/maild
ADD ./docker-entrypoint.sh /bin/docker-entrypoint.sh

RUN	addgroup -S maild && \
	adduser -S -g maild maild && \
	mkdir -p /var/lib/maild/templates && \
	chown -R maild:maild /var/lib/maild && \
	chown maild:maild /bin/maild && \
	chown maild:maild /bin/docker-entrypoint.sh && \
	chmod u+x /bin/maild && \
	chmod u+x /bin/docker-entrypoint.sh

USER maild

ENTRYPOINT ["docker-entrypoint.sh"]