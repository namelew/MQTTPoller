FROM golang:1.18-alpine as base

WORKDIR /app
COPY . /app/

RUN apk update; apk add build-base
RUN go build -o bin/orquestrator cmd/orquestrator/main.go

FROM alpine:3.17 as binary

ENV tolerance=5
ENV hkInterval=1
ENV broker=tcp://localhost:1883
ENV port=8000

COPY --from=base /app/bin/orquestrator ./app/orquestrator

EXPOSE ${port}

ENTRYPOINT cd /app; ./orquestrator --broker ${broker} --tl ${tolerance} --port ${port} --hk-interval ${hkInterval}
