FROM golang:1.21.0-alpine as base

WORKDIR /app
COPY . /app/

RUN apk update; apk add build-base
RUN go build -o bin/orquestrator cmd/orquestrator/main.go

FROM alpine:3.17 as binary

ENV TOLERANCE=5
ENV BROKER=tcp://localhost:1883
ENV PORT=8000

COPY --from=base /app/bin/orquestrator ./app/orquestrator

EXPOSE ${port}

ENTRYPOINT cd /app; ./orquestrator --broker ${BROKER} --tl ${TOLERANCE} --port ${PORT}
