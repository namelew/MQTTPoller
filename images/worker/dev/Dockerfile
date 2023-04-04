FROM golang:1.18-alpine as base

ENV broker=tcp://localhost:1883

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest

COPY . .

RUN go mod download github.com/namelew/mqtt-bm-latency
RUN apk add openjdk8

ENTRYPOINT ["air", "-c", "worker.air.toml"]