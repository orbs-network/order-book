curl --location --request GET 'http://localhost:8080/taker/v1/quote' \
--header 'Content-Type: application/json' \
--data '{
    "inAmount": "100",
    "inToken": "MATIC",
    "outToken": "USDC"
}'

curl --location --request POST 'http://localhost:8080/taker/v1/swap' \
--header 'Content-Type: application/json' \
--data '{
    "inAmount": "1",
    "inToken": "USDC",
    "outToken": "MATIC"
}'


curl --location --request POST 'http://localhost:8080/taker/v1/abort/8c6f17dc-5c5c-44c2-bd84-ccc2d5f0c321'


