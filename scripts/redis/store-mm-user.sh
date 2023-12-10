#!/bin/bash

if ! command -v redis-cli &> /dev/null
then
    echo "redis-cli could not be found. Please install it."
    exit 1
fi


HOST="ec2-52-30-0-190.eu-west-1.compute.amazonaws.com"
PORT=20759

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
KEY_USER_API_KEY="user:userApiKey:$API_KEY"
KEY_USER_ID="user:userId:$ID"
TYPE="MARKET_MAKER"
PUB_KEY="MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAEhqhj8rWPzkghzOZTUCOo/sdkE53sU1coVhaYskKGKrgiUF7lsSmxy46i3j8w7E7KMTfYBpCGAFYiWWARa0KQwg=="

redis-cli -h $HOST -p $PORT -a $PASSWORD HSET "$KEY_USER_API_KEY" id "$ID" type "$TYPE" pubKey "$PUB_KEY" apiKey "$API_KEY"

redis-cli -h $HOST -p $PORT -a $PASSWORD HSET "$KEY_USER_ID" id "$ID" type "$TYPE" pubKey "$PUB_KEY" apiKey "$API_KEY"

echo "Market maker user stored successfully."