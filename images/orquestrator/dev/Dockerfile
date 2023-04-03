FROM golang:1.18-alpine as base

WORKDIR /app
COPY . .

ENV tolerance=5
ENV hkInterval=1
ENV broker=tcp://localhost:1883
ENV port=8000

RUN apk update; apk add build-base
RUN go install github.com/cosmtrek/air@latest
RUN go mod download github.com/namelew/mqtt-bm-latency

EXPOSE ${port}

ENTRYPOINT ["air", "-c", "orquestrator.air.toml"]
