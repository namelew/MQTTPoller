FROM golang:1.17.13-bullseye

ENV tolerance=5
ENV broker=localhost

WORKDIR /app
COPY dump/orquestrator/ /app/

RUN apt-get update; apt-get upgrade -y
RUN apt-get install git make -y
RUN make

EXPOSE 8080

ENTRYPOINT cd bin; ./orquestrator --broker tcp://${broker}:1883 --tl ${tolerance}
