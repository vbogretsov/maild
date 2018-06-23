#! /bin/sh

function term()
{
    kill -15 $child
    wait $child
}

trap term SIGTERM

exec "`maild \
	--provider-url ${MAILD_PROVIDER_URL} \
	--provider-key ${MAILD_PROVIDER_KEY} \
	--provider-name ${MAILD_PROVIDER_NAME} \
	--templates-path ${MAILD_TEMPLATES_PATH} \
	--amqp-url ${MAILD_AMQP_URL} \
	--log-level ${MAILD_LOG_LEVEL} \
	--log-format ${MAILD_LOG_FORMAT} \
	`" &

child=$!
wait $child
