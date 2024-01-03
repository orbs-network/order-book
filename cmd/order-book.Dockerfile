FROM golang:1.21.3-alpine3.17

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN find . -type f ! -name "*.go" ! -name "go.mod" ! -name "go.sum" ! -name "supportedTokens.json" -delete

ARG APP_PATH

RUN CGO_ENABLED=0 GOOS=linux go build -o /order-book $APP_PATH

CMD [ "/order-book" ]