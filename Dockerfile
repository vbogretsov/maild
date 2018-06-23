FROM alpine:3.7

ENV \
	MAILD_AMQP_URL=amqp://guest:guest@rabbitmq:5672/ \
	MAILD_PROVIDER_URL= \
	MAILD_PROVIDER_KEY= \
	MAILD_PROVIDER_NAME= \
	MAILD_TEMPLATES_PATH=fs:///var/lib/maild/templates \
	MAILD_LOG_LEVEL=info \
	MAILD_LOG_FORMAT=json

ADD ./maild /bin/maild
ADD ./docker-entrypoint.sh /bin/docker-entrypoint.sh

RUN	adduser -D maild && \
	mkdir -p /var/lib/maild/templates && \
	# chown -R maild:maild /var/lib/maild && \
	# chown maild:maild /bin/maild && \
	# chown maild:maild /bin/docker-entrypoint.sh && \
	chmod +x /bin/maild && \
	chmod +x /bin/docker-entrypoint.sh

USER maild

ENTRYPOINT ["docker-entrypoint.sh"]