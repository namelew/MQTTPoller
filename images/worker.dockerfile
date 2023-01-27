FROM golang:1.17.13-bullseye

ENV timeout=5
ENV login_t=30
ENV tool=source/tools/mqttloader/bin/mqttloader
ENV broker=localhost

WORKDIR /app
COPY dump/worker/ /app/

RUN apt-get update; apt-get upgrade -y
RUN apt-get install git make -y
RUN make

ENTRYPOINT cd bin; ./worker --broker tcp://${broker}:1883 --timeout ${timeout} --tool ${tool} --login_t ${login_t}