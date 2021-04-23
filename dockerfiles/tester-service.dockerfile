# syntax=docker/dockerfile:experimental
FROM golang:1.16.3-alpine

WORKDIR /go/src/app

COPY ./testdata/tester_service.go .

EXPOSE 9010

CMD ["go", "run", "tester_service.go", "-service-addr", "server:9090"]