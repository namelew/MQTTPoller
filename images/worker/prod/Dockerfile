FROM golang:1.18-alpine as base

WORKDIR /app
COPY . /app/

RUN apk update
RUN go build -o bin/worker cmd/worker/main.go

FROM alpine:3.17 as binary

ENV LTIMEOUT=30
ENV LTHRESHOUT=-1
ENV TOOL=./tools/mqttloader/bin/mqttloader
ENV BROKER=tcp://localhost:1883

COPY --from=base /app/bin/worker ./app/worker
COPY --from=base /app/tools ./app/tools

RUN apk add openjdk8

ENTRYPOINT cd app; ./worker --broker ${BROKER} --tool ${TOOL} --login_t ${LTIMEOUT} --login_th ${LTHRESHOUT}