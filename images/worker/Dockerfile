FROM golang:1.21.0-alpine as base

WORKDIR /app
COPY . .

RUN apk update; apk add --no-cache curl unzip
RUN go build -o bin/worker src/core/cmd/worker/main.go
RUN mkdir tools; curl -OL https://github.com/dist-sys/mqttloader/releases/download/v0.8.6/mqttloader-0.8.6.zip; unzip ./mqttloader-0.8.6.zip -d /app/tools

FROM alpine:3.17 as binary

ENV LTIMEOUT=30
ENV LTHRESHOUT=-1
ENV TOOL=./tools/mqttloader/bin/mqttloader
ENV BROKER=tcp://localhost:1883

COPY --from=base /app/bin/worker ./app/worker
COPY --from=base /app/tools ./app/tools

RUN apk add openjdk8

ENTRYPOINT cd app; ./worker --broker ${BROKER} --tool ${TOOL} --login_t ${LTIMEOUT} --login_th ${LTHRESHOUT}