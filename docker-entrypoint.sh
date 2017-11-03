#!/bin/sh

term()
{
    kill -15 "$child"
    wait "$child"
}

trap term SIGTERM

wait_rabbit()
{
    local BROKER_HOST=$(echo $MAILD_BROKER_URL | sed 's/\///g' | sed s/amqp://g | sed s/.*:.*@//g | cut -d':' -f 1)
    local BROKER_PORT=$(echo $MAILD_BROKER_URL | sed 's/\///g' | sed s/amqp://g | sed s/.*:.*@//g | cut -d':' -f 2)

    for i in $(seq 1 10)
    do
            if nc -z $BROKER_HOST $BROKER_PORT
            then
                    break;
            else
                    echo unable to connect broker $MAILD_BROKER_URL;
                    sleep 1;
            fi
    done
}

wait_rabbit

/bin/maild	--provider-url $MAILD_PROVIDER_URL \
			--provider-key $MAILD_PROVIDER_KEY \
			--provider-name $MAILD_PROVIDER_NAME \
			--service-name $MAILD_SERVICE_NAME \
			--broker-url $MAILD_BROKER_URL \
			--log-level $MAILD_LOG_LEVEL \
			--template-dir /var/lib/maild/templates &

child=$!
wait "$child"