# Build variables
VERSION = $(shell git describe --tags --exact-match 2>/dev/null || git symbolic-ref -q --short HEAD)
COMMIT_HASH = $(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_DATE = $(shell date "+%FT%T%z")
LDFLAGS = -X main.version=${VERSION} -X main.commitHash=${COMMIT_HASH} -X main.buildDate=${BUILD_DATE}

all: test install generate build run

test-cover:
	go test -v -coverprofile cover.out ./...
	go tool cover -html cover.out

test:
	go test -v ./...

install: install-go-bindata install-sqlc

install-go-bindata:
	go get -v github.com/tmthrgd/go-bindata/go-bindata
	go install github.com/tmthrgd/go-bindata/go-bindata

install-sqlc:
	go get -v github.com/kyleconroy/sqlc/cmd/sqlc
	go install github.com/kyleconroy/sqlc/cmd/sqlc

generate: generate-go-bindata generate-sqlc

generate-go-bindata:
	go-bindata -pkg migrations -ignore bindata -nometadata -prefix internal/db/migrations/ -o ./internal/db/migrations/bindata.go ./internal/db/migrations

generate-sqlc:
	sqlc generate

build:
	go build -ldflags "${LDFLAGS}" -o ./bin/bitburst ./cmd/bitburst/*.go

run:
	./bin/bitburst --log-beautify