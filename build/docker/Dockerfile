FROM golang:1.19-alpine AS builder

ENV GO111MODULE=on

RUN apk update && apk add bash ca-certificates git gcc g++ libc-dev

#RUN mkdir /aliens

WORKDIR /aliens

#COPY . .

#COPY data/map.txt ./
COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build ./cmd/main.go -o /alien_invade

ENTRYPOINT ["./alien_invade", "-A", "1000", "-M", "/aliens/map.txt"]