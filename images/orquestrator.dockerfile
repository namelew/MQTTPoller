FROM golang:1.17.13-bullseye

ENV tolerance=5
ENV broker=localhost
ENV adress=""
ENV port=8000

WORKDIR /app
COPY dump/orquestrator/ /app/

RUN apt-get update; apt-get upgrade -y
RUN apt-get install git make -y
RUN make

EXPOSE ${port}

ENTRYPOINT cd bin; ./orquestrator --broker tcp://${broker}:1883 --tl ${tolerance} --adress ${adress} --port ${port}
