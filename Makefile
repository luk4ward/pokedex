.PHONY: build run test

LINTER_VERSION=v1.33.0

default: build

build: test
	go build -o ./app ./cmd/server

run: build
	./app

lint: get-linter
	golangci-lint run --timeout=5m

get-linter:
	command -v golangci-lint || curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b ${GOPATH}/bin $(LINTER_VERSION)

test: generate lint
	go fmt ./...
	go test -vet all ./...

cover:
	echo unit tests only
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

generate: get-generator
	go generate -x ./...

get-generator:
	go install github.com/golang/mock/mockgen

mod:
	go mod vendor -v

build-docker: test
	docker-compose build

run-docker: build-docker
	docker-compose up

start-docker: build-docker
	docker-compose up -d

stop-docker:
	docker-compose stop

clean-docker:
	docker-compose rm -s -f