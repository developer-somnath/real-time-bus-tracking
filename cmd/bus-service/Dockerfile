FROM golang:alpine

WORKDIR /app

COPY ../../go.mod ../../go.sum ./
RUN go mod download

COPY ../../cmd/bus-service ./cmd/bus-service
COPY ../../internal ./internal
COPY ../../pkg ./pkg
COPY ../../migrations ./migrations

RUN go build -o bus-service ./cmd/bus-service/main.go

CMD ["./bus-service"]