#!/bin/sh

term()
{
    kill -15 "$child"
    wait "$child"
}

trap term SIGTERM

/bin/maild	--provider-url $MAILD_PROVIDER_URL \
			--provider-key $MAILD_PROVIDER_KEY \
			--provider-name $MAILD_PROVIDER_NAME \
			--service-name $MAILD_SERVICE_NAME \
			--broker-url $MAILD_BROKER_URL \
			--log-level $MAILD_LOG_LEVEL \
			--template-dir /var/lib/maild/templates &

child=$!
wait "$child"