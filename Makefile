.PHONY: build run test

SLACK_TOKEN ?=

default: build

build:
	go build

run:
	go run main.go

test:
	go test ./...