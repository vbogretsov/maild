#!/bin/sh

term()
{
    kill -15 "$child"
    wait "$child"
}

trap term SIGTERM

wait_rabbit()
{
    local BROKER_HOST=$(echo $BROKER_URL | sed 's/\///g' | sed s/amqp://g | sed s/.*:.*@//g | cut -d':' -f 1)
    local BROKER_PORT=$(echo $BROKER_URL | sed 's/\///g' | sed s/amqp://g | sed s/.*:.*@//g | cut -d':' -f 2)

    for i in $(seq 1 10)
    do
            if nc -z $BROKER_HOST $BROKER_PORT
            then
                    break;
            else
                    echo unable to connect broker $BROKER_URL;
                    sleep 1;
            fi
    done
}

wait_rabbit

/bin/maild	--provider-url $PROVIDER_URL \
			--provider-key $PROVIDER_KEY \
			--provider-name $PROVIDER_NAME \
			--service-name $SERVICE_NAME \
			--broker-url $BROKER_URL \
			--log-level $LOG_LEVEL \
			--template-dir /var/lib/maild/templates &

child=$!
wait "$child"