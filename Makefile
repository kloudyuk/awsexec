SHELL = /bin/bash

.PHONY: test cover

test:
	go test -coverprofile cover.out -race -v ./...

cover: test
	go tool cover -html cover.out
