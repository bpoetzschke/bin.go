.PHONY: build run test run-infra

SLACK_TOKEN ?=

default: build

build:
	go build

run:
	go run main.go

test:
	go test ./... -race

run-infra:
	docker run -d -p 5432:5432 \
        -e POSTGRES_PASSWORD=postgres \
        -e POSTGRES_USER=postgres \
        -e POSTGRES_DB=bin.go \
        -v pgdata:/var/lib/postgresql/data \
        postgres