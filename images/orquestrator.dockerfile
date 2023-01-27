FROM golang:1.17.13-bullseye

ENV tolerance=5
ENV broker="tcp://localhost:1883"

WORKDIR /app
COPY dump/orquestrator/ /app/

RUN apt-get update; apt-get upgrade -y
RUN apt-get install git make -y
RUN make

ENTRYPOINT cd bin; ./orquestrator --broker ${broker} --tl ${tolerance}
