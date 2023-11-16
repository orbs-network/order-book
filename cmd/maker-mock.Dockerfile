FROM golang:1.21.3-alpine3.17

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN find . -type f ! -name "*.go" ! -name "go.mod" ! -name "go.sum" -delete

RUN CGO_ENABLED=0 GOOS=linux go build -o /maker-mock ./cmd/maker-mock

CMD [ "/maker-mock" ]