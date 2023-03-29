FROM golang:1.17.13-alpine as base

WORKDIR /app
COPY dump /app/

RUN apk update
RUN go build -o bin/worker cmd/worker/main.go

FROM alpine:3.17 as binary

ENV login_t=30
ENV tool=./tools/mqttloader/bin/mqttloader
ENV broker=tcp://localhost:1883

COPY --from=base /app/bin/worker ./app/worker
COPY --from=base /app/tools ./app/tools

ENTRYPOINT cd app; ./worker --broker ${broker} --tool ${tool} --login_t ${login_t}