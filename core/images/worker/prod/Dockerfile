FROM golang:1.18-alpine as base

WORKDIR /app
COPY . .

RUN apk update; apk add --no-cache git
RUN go build -o bin/worker cmd/worker/main.go
RUN git clone https://github.com/namelew/mqttloader.git ./tools/mqttloader

FROM alpine:3.17 as binary

ENV LTIMEOUT=30
ENV LTHRESHOUT=-1
ENV TOOL=./tools/mqttloader/bin/mqttloader
ENV BROKER=tcp://localhost:1883

COPY --from=base /app/bin/worker ./app/worker
COPY --from=base /app/tools ./app/tools

RUN apk add openjdk8

ENTRYPOINT cd app; ./worker --broker ${BROKER} --tool ${TOOL} --login_t ${LTIMEOUT} --login_th ${LTHRESHOUT}