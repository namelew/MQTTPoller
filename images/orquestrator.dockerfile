FROM golang:1.17.13-alpine

ENV tolerance=5
ENV broker=localhost
ENV port=8000

WORKDIR /app
COPY dump/orquestrator/ /app/

RUN apk update
RUN apk add git make
RUN make

EXPOSE ${port}

ENTRYPOINT cd bin; ./orquestrator --broker tcp://${broker}:1883 --tl ${tolerance} --port ${port}
