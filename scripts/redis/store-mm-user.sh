#!/bin/bash

if ! command -v redis-cli &> /dev/null
then
    echo "redis-cli could not be found. Please install it."
    exit 1
fi


if [ -z "${REDIS_HOST}" ]; then
    echo "Please set the REDIS_HOST environment variable."
    exit 1
else
    HOST="${REDIS_HOST}"
fi

if [ -z "${REDIS_PORT}" ]; then
    echo "Please set the REDIS_PORT environment variable."
    exit 1
else
    PORT="${REDIS_PORT}"
fi

if [ -z "${REDIS_PASSWORD}" ]; then
    echo "Please set the REDIS_PASSWORD environment variable."
    exit 1
else
    PASSWORD="${REDIS_PASSWORD}"
fi

if [ -z "${API_KEY}" ]; then
    echo "Please set the API_KEY environment variable."
    exit 1
else
    API_KEY="${API_KEY}"
fi

if ! redis-cli -h $HOST -p $PORT -a $PASSWORD PING &> /dev/null
then
    echo "Unable to reach Redis at $HOST:$PORT. Please check your connection."
    exit 1
fi

ID="00000000-0000-0000-0000-000000000001"
KEY_USER_API_KEY="userApiKey:$API_KEY:user"
KEY_USER_ID="userId:$ID:user"
TYPE="MARKET_MAKER"
PUB_KEY="0x6a04ab98d9e4774ad806e302dddeb63bea16b5cb5f223ee77478e861bb583eb336b6fbcb60b5b3d4f1551ac45e5ffc4936466e7d98f6c7c0ec736539f74691a6"

redis-cli -h $HOST -p $PORT -a $PASSWORD HSET "$KEY_USER_API_KEY" id "$ID" type "$TYPE" pubKey "$PUB_KEY" apiKey "$API_KEY"

redis-cli -h $HOST -p $PORT -a $PASSWORD HSET "$KEY_USER_ID" id "$ID" type "$TYPE" pubKey "$PUB_KEY" apiKey "$API_KEY"

echo "Market maker user stored successfully."