.PHONY: build
build:
	go build -v ./cmd/server
	GOOS=linux GOARCH=amd64 go build -o server-linux-amd64 -v ./cmd/server
.DEFAULT_GOAL := build