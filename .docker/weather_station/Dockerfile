FROM golang:alpine

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

RUN go get git.apache.org/thrift.git/lib/go/thrift
RUN go get github.com/eclipse/paho.mqtt.golang

WORKDIR /src

COPY ./WeatherStation /src

RUN cd /src && go build -o goapp
