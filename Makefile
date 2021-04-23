# Build variables
VERSION = $(shell git describe --tags --exact-match 2>/dev/null || git symbolic-ref -q --short HEAD)
COMMIT_HASH = $(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_DATE = $(shell date "+%FT%T%z")
LDFLAGS = -X main.version=${VERSION} -X main.commitHash=${COMMIT_HASH} -X main.buildDate=${BUILD_DATE}

all: test build run

test_cover:
	go test -v -coverprofile cover.out ./...
	go tool cover -html cover.out

test:
	go test -v ./...

build:
	go build -ldflags "${LDFLAGS}" -o ./bin/bitburst ./cmd/bitburst/*.go

run:
	./bin/bitburst --log-beautify