FROM golang:alpine as build-env

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

RUN go get git.apache.org/thrift.git/lib/go/thrift
RUN go get github.com/lib/pq

WORKDIR /src

COPY ./WeatherService /src

RUN cd /src && go build -o goapp
