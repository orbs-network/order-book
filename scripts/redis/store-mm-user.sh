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

if ! redis-cli -h $HOST -p $PORT -a $PASSWORD PING &> /dev/null
then
    echo "Unable to reach Redis at $HOST:$PORT. Please check your connection."
    exit 1
fi

KEY="user:MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAEhqhj8rWPzkghzOZTUCOo/sdkE53sU1coVhaYskKGKrgiUF7lsSmxy46i3j8w7E7KMTfYBpCGAFYiWWARa0KQwg==:publicKey"
ID="00000000-0000-0000-0000-000000000001"
TYPE="MARKET_MAKER"
PUBKEY="MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAEhqhj8rWPzkghzOZTUCOo/sdkE53sU1coVhaYskKGKrgiUF7lsSmxy46i3j8w7E7KMTfYBpCGAFYiWWARa0KQwg=="

redis-cli -h $HOST -p $PORT -a $PASSWORD HSET "$KEY" id "$ID" type "$TYPE" pubKey "$PUBKEY"

echo "Market maker user stored successfully."