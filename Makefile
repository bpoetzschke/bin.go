.PHONY: build run

SLACK_TOKEN ?=

default: build

build:
	go build

run:
	go run main.go