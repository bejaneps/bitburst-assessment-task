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
	go install github.com/tmthrgd/go-bindata/go-bindata

install-sqlc:
	go install github.com/kyleconroy/sqlc/cmd/sqlc

generate: generate-go-bindata generate-sqlc

generate-go-bindata:
	go-bindata -pkg migrations -ignore bindata -nometadata -prefix internal/db/migrations/ -o ./internal/db/migrations/bindata.go ./internal/db/migrations

generate-sqlc:
	sqlc generate

build:
	go build -trimpath -ldflags "${LDFLAGS}" -o ./bin/bitburst ./cmd/bitburst/main.go

build-darwin-linux:
	GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags "${LDFLAGS}" -o ./bin/bitburst_darwin_amd64 ./cmd/bitburst/main.go
	GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "${LDFLAGS}" -o ./bin/bitburst_linux_amd64 ./cmd/bitburst/main.go

run:
	./bin/bitburst --log-beautify