FROM alpine:3.6

ENV BROKER_URL amqp://guest:guest@rabbitmq:5672/
ENV PROVIDER_URL none
ENV PROVIDER_KEY none
ENV PROVIDER_NAME log
ENV SERVICE_NAME maild
ENV LOG_LEVEL INFO

ADD ./bin/maild /bin/maild
ADD ./docker-entrypoint.sh /bin/docker-entrypoint.sh

RUN addgroup -S maild && \
	adduser -S -g maild maild && \
	mkdir -p /var/lib/maild/templates && \
	chown -R maild:maild /var/lib/maild && \
	chown maild:maild /bin/maild && \
	chown maild:maild /bin/docker-entrypoint.sh && \
	chmod u+x /bin/maild && \
	chmod u+x /bin/docker-entrypoint.sh

USER maild

ENTRYPOINT ["docker-entrypoint.sh"]