# syntax=docker/dockerfile:experimental
FROM golang:1.16.3 AS build

LABEL MAINTAINER=bejanhtc@gmail.com

WORKDIR /go/src/app

# download dependencies specified in go.mod and go.sum
COPY ./go.mod .
COPY ./go.sum .
RUN --mount=type=cache,target=/go/pkg/mod go mod download -x

# copy all source files to container
COPY . .

# build executable
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 \
    go build -o ./bin/bitburst ./cmd/bitburst/main.go

# copy executable to new container
FROM alpine:latest
COPY --from=build /go/src/app/bin /go/src/app/bin

WORKDIR /go/src/app

EXPOSE 9090

CMD ["./bin/bitburst", "--log-beautify", "--log-level", "0"]