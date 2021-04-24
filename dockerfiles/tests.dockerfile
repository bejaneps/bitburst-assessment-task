FROM golang:alpine

LABEL MAINTAINER=bejanhtc@gmail.com

WORKDIR /go/src/app

COPY . .

ENV CGO_ENABLED=0

CMD ["/usr/local/go/bin/go", "test", "-v", "./..."]