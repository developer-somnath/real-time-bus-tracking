FROM golang:alpine

WORKDIR /app

COPY ../../go.mod ../../go.sum ./
RUN go mod download

COPY ../../cmd/api-gateway ./cmd/api-gateway
COPY ../../internal ./internal
COPY ../../pkg ./pkg
COPY ../../migrations ./migrations

RUN go build -o api-gateway ./cmd/api-gateway/main.go

CMD ["./api-gateway"]