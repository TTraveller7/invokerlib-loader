FROM golang:1.19.0-alpine3.16 AS builder

WORKDIR /loader

COPY go.mod ./
COPY go.sum ./

RUN mkdir bin

RUN go mod download

COPY *.go ./

RUN go build -o ./bin/loader ./


# Copy it the app to an alpine
FROM alpine:3.11.3

WORKDIR /app

COPY --from=builder /loader/bin/loader ./
COPY order-line.csv ./

CMD ["./loader"]
