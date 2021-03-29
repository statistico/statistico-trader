FROM golang:1.13.9-alpine

RUN apk update && apk add git gcc g++
ENV GO111MODULE=on
WORKDIR /opt
COPY go.mod .
COPY go.sum .
RUN go get gotest.tools/gotestsum
RUN go mod download

CMD ["/bin/sh"]