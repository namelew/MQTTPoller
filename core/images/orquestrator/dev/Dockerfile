FROM golang:1.18-alpine as base

WORKDIR /app
COPY . .

ENV TOLERANCE=5
ENV HKINTERVAL=1
ENV BROKER=tcp://localhost:1883
ENV PORT=8000

RUN apk update; apk add build-base
RUN go install github.com/cosmtrek/air@latest
RUN go mod download github.com/namelew/mqtt-bm-latency

EXPOSE ${port}

ENTRYPOINT ["air", "-c", "orquestrator.air.toml"]
