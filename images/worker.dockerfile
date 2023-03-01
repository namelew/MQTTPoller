FROM golang:1.17.13-alpine

ENV login_t=30
ENV tool=./tools/mqttloader/bin/mqttloader
ENV broker=localhost

WORKDIR /app
COPY dump/worker/ /app/

RUN apk update
RUN apk add git make
RUN make

ENTRYPOINT cd bin; ./worker --broker tcp://${broker}:1883 --tool ${tool} --login_t ${login_t}