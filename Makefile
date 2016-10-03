# Make file for go project

NAME=nubis-user-management
GO15VENDOREXPERIMENT := 1

build:
	rm -rf build && mkdir -p build
	mkdir -p build/linux && GOOS=linux go build -o build/linux/$(NAME) ./*.go
	mkdir -p build/darwin && GOOS=darwin go build -o build/darwin/$(NAME) ./*.go

fmt:
	gofmt -w=true $$(find . -type f -name '*.go')

generate:
	go generate

clean:
	rm -rf build

all: clean build

.PHONY: build fmt
