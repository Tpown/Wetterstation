FROM golang:alpine as build-env

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

RUN go get github.com/eclipse/paho.mqtt.golang

WORKDIR /src

COPY ./Sensor /src

RUN cd /src && go build -o goapp

