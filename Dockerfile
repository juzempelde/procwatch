FROM golang:1.10.3-alpine3.7
RUN mkdir -p /go/src/github.com/juzempelde/procwatch
WORKDIR /go/src/github.com/juzempelde/procwatch

RUN apk add --update-cache \
    git

RUN go get -u golang.org/x/vgo
COPY . .
WORKDIR /go/src/github.com/juzempelde/procwatch/backend
RUN vgo build -o procwatch ./cmd/procwatch
