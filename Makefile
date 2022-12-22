SHELL = /bin/bash

.PHONY: mod test

mod:
	go mod tidy

cover.out:
	go test -race -coverprofile cover.out -v ./...

test: cover.out

cover: cover.out
	go tool cover -html=cover.out
