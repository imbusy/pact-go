all: compile

dependencies:
	go get ./...

vet: dependencies
	go vet -v ./...

test: vet
	go test -v -cover ./...

compile: test
	go build -v

run: all

.PHONY: all compile test vet dependencies
