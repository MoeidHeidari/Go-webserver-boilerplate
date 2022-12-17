# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download


COPY *.go ./
RUN go get

EXPOSE 4000

CMD [ "go","run","main.go" ]