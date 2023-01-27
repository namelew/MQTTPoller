FROM golang:1.17.13-bullseye

ENV timeout=5
ENV login_t=30
ENV tool="source/tools/mqttloader/bin/mqttloader"
ENV broker="tcp://localhost:1883"

WORKDIR /app
COPY dump/orquestrator/ /app/

RUN apt-get update; apt-get upgrade -y
RUN apt-get install git make -y
RUN make

ENTRYPOINT cd bin; ./orquestrator --broker ${broker} --timeout ${timeout} --tool ${tool} --login_t ${login_t}