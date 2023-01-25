FROM golang:1.19-alpine AS builder

ENV GO111MODULE=on

RUN apk update && apk add bash ca-certificates git gcc g++ libc-dev

RUN go clean -modcache

RUN mkdir /aliens

WORKDIR /aliens

COPY . .

COPY data/map.txt ./
COPY go.mod ./
COPY go.sum ./

RUN go mod download

#COPY *.go ./

#RUN cd cmd

RUN go build -ldflags="-d -s -w" -tags timetzdata -trimpath -o app ./cmd

CMD [ "./app", "-A", "1000", "-M", "/aliens/map.txt" ]

#ENTRYPOINT ["./app", "-A", "1000", "-M", "/aliens/map.txt"]